/** Explicit local Pages handoff checker. Never loads CMS settings or uploads/extracts archive members. */
import fs from "node:fs/promises";
import { verifyPagesPayload, MAX_PAYLOAD } from "./pages-payload.mjs";
/** Read only a bounded regular file; input paths are operator-selected, never supplied by public content. */
async function read(name, max) {
  const s = await fs.lstat(name);
  if (!s.isFile() || s.isSymbolicLink() || s.size < 1 || s.size > max)
    throw Error("PAGES_PAYLOAD_UNVERIFIED");
  const f = await fs.open(name, "r");
  try {
    const b = Buffer.alloc(s.size + 1);
    let n = 0;
    while (n < b.length) {
      const r = await f.read(b, n, b.length - n, null);
      if (!r.bytesRead) break;
      n += r.bytesRead;
    }
    if (n !== s.size) throw Error("PAGES_PAYLOAD_UNVERIFIED");
    return b.subarray(0, n);
  } finally {
    await f.close();
  }
}
try {
  if (process.argv.length !== 4) throw Error();
  const raw = await read(process.argv[2], MAX_PAYLOAD),
    manifest = JSON.parse(
      (await read(process.argv[3], 1024 * 1024)).toString("utf8"),
    );
  verifyPagesPayload(raw, manifest);
  console.log(
    "Exact Pages payload and complete file inventory verified locally. No deployment or CMS receipt.",
  );
} catch {
  console.error("PAGES_PAYLOAD_UNVERIFIED");
  process.exitCode = 1;
}
