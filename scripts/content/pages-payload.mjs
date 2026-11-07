/** Independent Node validation of the complete Go Pages handoff. No files are extracted, executed, uploaded or fetched. */
import { createHash } from "node:crypto";
import { crc32, inflateRawSync } from "node:zlib";

export const MAX_TAR = 32 * 1024 * 1024 + 256 * 1024 + 1024;
export const MAX_PAYLOAD = MAX_TAR + 1024 * 1024;
const digest = /^[a-f0-9]{64}$/;
/** Hash fields are strings, never coercible arrays or objects. */
function validHash(value) {
  return typeof value === "string" && digest.test(value);
}
/** Constant diagnostics never contain paths, file content or caller/provider text. */
function invalid() {
  throw new Error("PAGES_PAYLOAD_UNVERIFIED");
}
/** SHA-256 covers actual binary data, including images, fonts and empty files. */
function hash(raw) {
  return createHash("sha256").update(raw).digest("hex");
}
/** Reject unknown/missing public fields instead of accepting credentials or caller-chosen destinations. */
function exact(value, keys) {
  if (
    !value ||
    Array.isArray(value) ||
    typeof value !== "object" ||
    Object.keys(value).sort().join("|") !== [...keys].sort().join("|")
  )
    invalid();
}
/** Archive names have one portable, non-aliased interpretation on supported development hosts. */
function portable(name) {
  if (
    typeof name !== "string" ||
    name.length < 1 ||
    name.length > 1024 ||
    !/^[A-Za-z0-9_. /-]+$/.test(name)
  )
    return false;
  return name
    .split("/")
    .every(
      (p, i) =>
        p.length > 0 &&
        p.length <= 240 &&
        p !== "." &&
        p !== ".." &&
        !/[. ]$/.test(p) &&
        (!p.startsWith(".") || (i === 0 && p === ".well-known")) &&
        !/^(CON|PRN|AUX|NUL|COM[0-9]|LPT[0-9])$/i.test(p.split(".")[0]),
    );
}
/** NUL padding must be canonical; high bytes and hidden suffixes cannot alias another file. */
function text(raw) {
  const n = raw.indexOf(0),
    end = n < 0 ? raw.length : n;
  if (
    raw.subarray(end).some((v) => v !== 0) ||
    raw.subarray(0, end).some((v) => v < 32 || v > 126)
  )
    invalid();
  return raw.subarray(0, end).toString("ascii");
}
/** Only ordinary zero-padded octal numeric fields are supported, never base-256 or negative encodings. */
function octal(raw) {
  if (!/^[0-7]+\0$/.test(raw.toString("latin1"))) invalid();
  const n = Number.parseInt(raw.toString("ascii"), 8);
  if (!Number.isSafeInteger(n) || n < 0) invalid();
  return n;
}
/** Validate one bounded manifest before allocating decompressed bytes or trusting any member identity. */
function manifest(value) {
  exact(value, [
    "schema_version",
    "scope",
    "source",
    "payload_sha256",
    "payload_bytes",
    "entries",
  ]);
  if (
    value.schema_version !== 1 ||
    value.scope !== "astro-pages-payload-v1" ||
    !validHash(value.payload_sha256) ||
    !Number.isInteger(value.payload_bytes) ||
    value.payload_bytes < 18 ||
    value.payload_bytes > MAX_PAYLOAD
  )
    invalid();
  const source = value.source;
  exact(source, ["content", "zip_sha256", "files", "expanded_bytes"]);
  if (
    !validHash(source.zip_sha256) ||
    !Number.isInteger(source.files) ||
    source.files < 3 ||
    source.files > 256 ||
    !Number.isInteger(source.expanded_bytes) ||
    source.expanded_bytes < 1 ||
    source.expanded_bytes > 32 * 1024 * 1024
  )
    invalid();
  const c = source.content;
  exact(c, [
    "request",
    "en_sha256",
    "id_sha256",
    "marker_sha256",
    "projects",
    "pages",
  ]);
  if (
    ![c.en_sha256, c.id_sha256, c.marker_sha256].every(
      (s) => typeof s === "string" && validHash(s),
    ) ||
    c.pages !== 2 ||
    !Number.isInteger(c.projects) ||
    c.projects < 0
  )
    invalid();
  const r = c.request;
  exact(r, [
    "schema_version",
    "scope",
    "release_id",
    "sequence",
    "snapshot_digest",
    "projection_digest",
  ]);
  if (
    r.schema_version !== 1 ||
    r.scope !== "astro-project-profile-build-v1" ||
    typeof r.release_id !== "string" ||
    !/^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$/.test(
      r.release_id,
    ) ||
    typeof r.sequence !== "string" ||
    !/^[1-9][0-9]{0,18}$/.test(r.sequence) ||
    BigInt(r.sequence) > 9223372036854775807n ||
    !validHash(r.snapshot_digest) ||
    !validHash(r.projection_digest)
  )
    invalid();
  if (!Array.isArray(value.entries) || value.entries.length !== source.files)
    invalid();
  let prior = "",
    total = 0;
  const names = new Set();
  for (const f of value.entries) {
    exact(f, ["name", "size", "sha256"]);
    if (
      !portable(f.name) ||
      f.name <= prior ||
      !Number.isInteger(f.size) ||
      f.size < 0 ||
      f.size > 32 * 1024 * 1024 ||
      !validHash(f.sha256)
    )
      invalid();
    const key = f.name.toLowerCase();
    if (
      names.has(key) ||
      [...names].some((n) => n.startsWith(key + "/") || key.startsWith(n + "/"))
    )
      invalid();
    names.add(key);
    prior = f.name;
    total += f.size;
  }
  if (total !== source.expanded_bytes) invalid();
  for (const [name, expected] of [
    ["index.html", c.en_sha256],
    ["id/index.html", c.id_sha256],
    [".well-known/homepage-content.json", c.marker_sha256],
  ])
    if (value.entries.find((f) => f.name === name)?.sha256 !== expected)
      invalid();
}
/** Verify exact gzip bytes and all canonical USTAR members, preserving all public assets without filesystem extraction.
 * The manifest must come from a trusted Go check; public hashes alone are not approval/deployment authority.
 */
export function verifyPagesPayload(payload, evidence) {
  try {
    manifest(evidence);
    if (
      !Buffer.isBuffer(payload) ||
      payload.length !== evidence.payload_bytes ||
      payload.length > MAX_PAYLOAD ||
      hash(payload) !== evidence.payload_sha256 ||
      !payload
        .subarray(0, 10)
        .equals(Buffer.from([31, 139, 8, 0, 0, 0, 0, 0, 0, 255]))
    )
      invalid();
    const inflated = inflateRawSync(payload.subarray(10), {
        maxOutputLength: MAX_TAR,
        info: true,
      }),
      raw = inflated.buffer,
      end = 10 + inflated.engine.bytesWritten;
    if (
      end + 8 !== payload.length ||
      payload.readUInt32LE(end) !== crc32(raw) ||
      payload.readUInt32LE(end + 4) !== raw.length
    )
      invalid();
    let offset = 0;
    for (const file of evidence.entries) {
      if (offset + 512 > raw.length) invalid();
      const h = raw.subarray(offset, offset + 512);
      let checksum = 0;
      for (let i = 0; i < 512; i++) checksum += i >= 148 && i < 156 ? 32 : h[i];
      if (
        !/^[0-7]{6}\0 $/.test(h.subarray(148, 156).toString("latin1")) ||
        Number.parseInt(h.subarray(148, 154).toString("ascii"), 8) !== checksum
      )
        invalid();
      const prefix = text(h.subarray(345, 500)),
        name = (prefix ? prefix + "/" : "") + text(h.subarray(0, 100));
      if (
        name !== file.name ||
        h[156] !== 48 ||
        !h.subarray(257, 265).equals(Buffer.from("ustar\0" + "00")) ||
        h.subarray(157, 257).some((v) => v !== 0) ||
        h.subarray(265, 329).some((v) => v !== 0) ||
        h.subarray(500).some((v) => v !== 0)
      )
        invalid();
      if (
        octal(h.subarray(100, 108)) !== 0o644 ||
        octal(h.subarray(108, 116)) !== 0 ||
        octal(h.subarray(116, 124)) !== 0 ||
        octal(h.subarray(124, 136)) !== file.size ||
        octal(h.subarray(136, 148)) !== 0 ||
        octal(h.subarray(329, 337)) !== 0 ||
        octal(h.subarray(337, 345)) !== 0
      )
        invalid();
      offset += 512;
      const next = offset + Math.ceil(file.size / 512) * 512;
      if (
        next > raw.length ||
        hash(raw.subarray(offset, offset + file.size)) !== file.sha256 ||
        raw.subarray(offset + file.size, next).some((v) => v !== 0)
      )
        invalid();
      offset = next;
    }
    if (
      raw.length !== offset + 1024 ||
      raw.subarray(offset).some((v) => v !== 0)
    )
      invalid();
    return Object.freeze({
      scope: evidence.scope,
      files: evidence.source.files,
      bytes: evidence.source.expanded_bytes,
      payloadSHA256: evidence.payload_sha256,
    });
  } catch {
    invalid();
  }
}

/** Check the exact uncompressed TAR that upload-artifact will receive, without extracting any member. */
export function verifyPagesTar(payload, evidence, tar) {
  try {
    verifyPagesPayload(payload, evidence);
    if (!Buffer.isBuffer(tar) || tar.length > MAX_TAR) invalid();
    const raw = inflateRawSync(payload.subarray(10), {
      maxOutputLength: MAX_TAR,
    });
    if (!raw.equals(tar)) invalid();
    return true;
  } catch {
    invalid();
  }
}
