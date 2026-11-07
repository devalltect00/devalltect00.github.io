/** Real tar/gzip adversarial fixtures test the independent checker without private Go code or a hosted provider. */
import { test } from "node:test";
import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { gzipSync } from "node:zlib";
import {
  verifyPagesPayload,
  verifyPagesTar,
  MAX_TAR,
} from "../../content/pages-payload.mjs";
/** Public synthetic file hash, never a credential. */
const hash = (b) => createHash("sha256").update(b).digest("hex");
/** Build canonical USTAR numeric fields for tests, independently of Go's tar implementation. */
function number(n, width) {
  return n.toString(8).padStart(width - 1, "0") + "\0";
}
/** Create exact tar bytes around synthetic complete files; the callback injects deliberately hostile headers. */
function tar(files, change = () => {}) {
  const blocks = [];
  for (const [name, body] of files) {
    const h = Buffer.alloc(512);
    let leaf = name,
      prefix = "";
    if (name.length > 100) {
      const i = name.lastIndexOf("/");
      prefix = name.slice(0, i);
      leaf = name.slice(i + 1);
    }
    h.write(leaf);
    h.write(number(0o644, 8), 100);
    h.write(number(0, 8), 108);
    h.write(number(0, 8), 116);
    h.write(number(body.length, 12), 124);
    h.write(number(0, 12), 136);
    h[156] = 48;
    h.write("ustar\0" + "00", 257);
    h.write(number(0, 8), 329);
    h.write(number(0, 8), 337);
    h.write(prefix, 345);
    change(h, name);
    h.fill(32, 148, 156);
    const sum = h.reduce((n, b) => n + b, 0);
    h.write(sum.toString(8).padStart(6, "0") + "\0 ", 148);
    blocks.push(h, body, Buffer.alloc((512 - (body.length % 512)) % 512));
  }
  blocks.push(Buffer.alloc(1024));
  return Buffer.concat(blocks);
}
/** Gzip metadata matches the documented portable Go header. */
function gzip(raw) {
  const b = gzipSync(raw);
  b[9] = 255;
  return b;
}
/** Complete structural manifest; actual frozen HTML semantics are checked by Go before this boundary. */
function fixture(change) {
  const files = [
    [".well-known/homepage-content.json", Buffer.from("{}")],
    ["_astro/runtime.js", Buffer.from([0, 255, 1, 2])],
    ["id/index.html", Buffer.from("id")],
    ["index.html", Buffer.from("en")],
    ["zero.txt", Buffer.alloc(0)],
  ];
  const payload = gzip(tar(files, change)),
    r = {
      schema_version: 1,
      scope: "astro-project-profile-build-v1",
      release_id: "11111111-1111-4111-8111-111111111111",
      sequence: "1",
      snapshot_digest: "a".repeat(64),
      projection_digest: "b".repeat(64),
    };
  const evidence = {
    schema_version: 1,
    scope: "astro-pages-payload-v1",
    source: {
      content: {
        request: r,
        en_sha256: hash(files[3][1]),
        id_sha256: hash(files[2][1]),
        marker_sha256: hash(files[0][1]),
        projects: 0,
        pages: 2,
      },
      zip_sha256: "c".repeat(64),
      files: files.length,
      expanded_bytes: files.reduce((n, [, b]) => n + b.length, 0),
    },
    payload_sha256: hash(payload),
    payload_bytes: payload.length,
    entries: files.map(([name, b]) => ({
      name,
      size: b.length,
      sha256: hash(b),
    })),
  };
  return { files, payload, evidence };
}
/** Recompute only the outer public digest so structure/member checks, rather than the first digest guard, must reject corruption. */
function repin(evidence, payload) {
  evidence.payload_sha256 = hash(payload);
  evidence.payload_bytes = payload.length;
  return evidence;
}

test("Pages payload checks all assets and empty files without extraction", () => {
  const { payload, evidence } = fixture();
  const out = verifyPagesPayload(payload, evidence);
  assert.equal(out.files, 5);
  assert.equal(out.payloadSHA256, hash(payload));
  assert.ok(Object.isFrozen(out));
});
test("Pages rejects invalid or coercible manifest fields", () => {
  for (const change of [
    (e) => (e.extra = "destination"),
    (e) => (e.scope = "deployment-receipt"),
    (e) => (e.payload_sha256 = [e.payload_sha256]),
    (e) => (e.source.zip_sha256 = [e.source.zip_sha256]),
    (e) => e.entries.pop(),
    (e) => e.entries.reverse(),
    (e) => (e.entries[0].name = "../private"),
    (e) => (e.entries[4].name = "INDEX.HTML"),
    (e) => (e.source.content.request.sequence = "9223372036854775808"),
    (e) => (e.source.content.en_sha256 = "0".repeat(64)),
  ]) {
    const { payload, evidence } = fixture();
    change(evidence);
    assert.throws(
      () => verifyPagesPayload(payload, evidence),
      /^Error: PAGES_PAYLOAD_UNVERIFIED$/,
    );
  }
});
test("Pages rejects link, PAX, owner, mode and timestamp header substitutions", () => {
  for (const change of [
    (h) => (h[156] = 50),
    (h) => (h[156] = 49),
    (h) => (h[156] = 120),
    (h) => h.write("../private", 0),
    (h) => h.write(number(1, 8), 108),
    (h) => h.write(number(0o755, 8), 100),
    (h) => h.write(number(1, 12), 136),
    (h) => (h[500] = 1),
  ]) {
    const { payload, evidence } = fixture((h) => change(h));
    assert.throws(
      () => verifyPagesPayload(payload, evidence),
      /PAGES_PAYLOAD_UNVERIFIED/,
    );
  }
});
test("Pages checks runtime file bytes, tar padding, CRC and complete endings", () => {
  for (const kind of [
    "runtime",
    "padding",
    "checksum",
    "tail",
    "CRC",
    "truncated",
    "gzip-metadata",
    "multigzip",
  ]) {
    const { files, evidence } = fixture();
    let raw = tar(files),
      payload;
    if (kind === "runtime") raw[1536] ^= 1;
    if (kind === "padding") raw[515] = 1;
    if (kind === "checksum") raw[148] ^= 1;
    if (kind === "tail") raw = Buffer.concat([raw, Buffer.alloc(512)]);
    payload = gzip(raw);
    if (kind === "CRC") payload[payload.length - 8] ^= 1;
    if (kind === "truncated") payload = payload.subarray(0, payload.length - 1);
    if (kind === "gzip-metadata") payload[4] = 1;
    if (kind === "multigzip")
      payload = Buffer.concat([payload, gzip(Buffer.alloc(0))]);
    assert.throws(
      () => verifyPagesPayload(payload, repin(evidence, payload)),
      /PAGES_PAYLOAD_UNVERIFIED/,
      kind,
    );
  }
});
test("Pages bounds decompression before reading a purported inventory", () => {
  const { evidence } = fixture();
  const payload = gzip(Buffer.alloc(MAX_TAR + 1));
  assert.throws(
    () => verifyPagesPayload(payload, repin(evidence, payload)),
    /PAGES_PAYLOAD_UNVERIFIED/,
  );
});
test("Pages supports canonical USTAR prefixes without permitting private hidden roots", () => {
  const { files, evidence } = fixture();
  const name = "folder/".repeat(14) + "asset.bin";
  files.splice(2, 0, [name, Buffer.from([12, 34])]);
  files.sort(([a], [b]) => (a < b ? -1 : a > b ? 1 : 0));
  const payload = gzip(tar(files));
  evidence.entries = files.map(([name, b]) => ({
    name,
    size: b.length,
    sha256: hash(b),
  }));
  evidence.source.files = files.length;
  evidence.source.expanded_bytes += 2;
  assert.equal(verifyPagesPayload(payload, repin(evidence, payload)).files, 6);
  evidence.entries[0].name = ".private/file";
  assert.throws(
    () => verifyPagesPayload(payload, evidence),
    /PAGES_PAYLOAD_UNVERIFIED/,
  );
});

test("Exact upload TAR must equal the independently verified payload bytes", () => {
  const { payload, evidence, files } = fixture();
  const raw = tar(files);
  assert.equal(verifyPagesTar(payload, evidence, raw), true);
  raw[512] ^= 1;
  assert.throws(
    () => verifyPagesTar(payload, evidence, raw),
    /PAGES_PAYLOAD_UNVERIFIED/,
  );
});
test("Reserved COM0 names follow the same strict Go and Node portability policy", () => {
  const { files, evidence } = fixture();
  files[4][0] = "zero/COM0.txt";
  const payload = gzip(tar(files));
  evidence.entries[4].name = files[4][0];
  assert.throws(
    () => verifyPagesPayload(payload, repin(evidence, payload)),
    /PAGES_PAYLOAD_UNVERIFIED/,
  );
});
