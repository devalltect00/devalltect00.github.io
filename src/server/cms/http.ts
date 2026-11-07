/** Anonymous fixed-origin read for trusted preparation tooling; no cookies, keys, proxies, redirects or content URLs. */
import http from "node:http";
import https from "node:https";
import { CmsError, packageLimit } from "./snapshot.ts";
/** Loopback HTTP is an explicit development capability, not a production default or browser endpoint. */
export function publicSource(
  origin: string,
  allowLoopback = false,
  timeoutMs = 10000,
): () => Promise<Uint8Array> {
  let u: URL;
  try {
    u = new URL(origin);
  } catch {
    throw new CmsError();
  }
  const loopback = ["127.0.0.1", "localhost", "[::1]"].includes(u.hostname);
  if (
    (origin !== u.origin && origin !== u.origin + "/") ||
    u.username ||
    u.password ||
    u.search ||
    u.hash ||
    u.pathname !== "/" ||
    !Number.isInteger(timeoutMs) ||
    timeoutMs < 1 ||
    timeoutMs > 10000 ||
    (u.protocol !== "https:" &&
      !(allowLoopback && loopback && u.protocol === "http:")) ||
    (loopback && !allowLoopback)
  )
    throw new CmsError();
  return () =>
    new Promise((resolve, reject) => {
      let settled = false;
      const client = u.protocol === "https:" ? https : http;
      const fail = () => {
        if (!settled) {
          settled = true;
          reject(new CmsError());
        }
        request.destroy();
      };
      const request = client.get(
        new URL("/public/v1/content/en", u),
        {
          agent: false,
          headers: { Accept: "application/json", "Cache-Control": "no-cache" },
        },
        (response) => {
          const length = response.headers["content-length"],
            mime = response.headers["content-type"]?.split(";")[0]?.trim();
          if (
            response.statusCode !== 200 ||
            mime !== "application/json" ||
            (length && (!/^\d+$/.test(length) || Number(length) > packageLimit))
          ) {
            response.destroy();
            fail();
            return;
          }
          const buffer = Buffer.allocUnsafe(packageLimit);
          let size = 0;
          response.on("data", (chunk: Buffer) => {
            size += chunk.length;
            if (size > packageLimit) {
              response.destroy();
              fail();
              return;
            }
            chunk.copy(buffer, size - chunk.length);
          });
          response.on("aborted", fail);
          response.on("error", fail);
          response.on("end", () => {
            if (!settled) {
              settled = true;
              clearTimeout(timer);
              if (!size) reject(new CmsError());
              else resolve(Buffer.from(buffer.subarray(0, size)));
            }
          });
        },
      );
      const timer = setTimeout(fail, timeoutMs);
      timer.unref();
      request.on("error", () => {
        clearTimeout(timer);
        fail();
      });
      request.on("close", () => clearTimeout(timer));
    });
}
