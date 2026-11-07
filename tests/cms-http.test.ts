// @vitest-environment node
/** Actual local HTTP servers prove no redirects/ambient credentials, strict MIME/size and bounded complete reads. */
import { it, expect } from "vitest";
import http from "node:http";
import { once } from "node:events";
import { publicSource } from "../src/server/cms/http.ts";
import { packageLimit } from "../src/server/cms/snapshot.ts";
/** Every test server is loopback-only and closed even when the read is rejected. */
async function serve(
  handler: http.RequestListener,
  run: (origin: string) => Promise<void>,
) {
  const server = http.createServer(handler);
  server.listen(0, "127.0.0.1");
  await once(server, "listening");
  const address = server.address();
  if (!address || typeof address === "string") throw Error();
  try {
    await run("http://127.0.0.1:" + address.port);
  } finally {
    server.closeAllConnections();
    await new Promise<void>((resolve) => server.close(() => resolve()));
  }
}
it("allows only explicit root origins and opt-in loopback HTTP", () => {
  for (const origin of [
    "http://example.com",
    "http://127.0.0.1:1234",
    "https://user:secret@example.com",
    "https://example.com/path",
    "https://example.com?",
    "https://example.com#x",
    "https://localhost/",
  ])
    expect(() => publicSource(origin)).toThrow("CMS_CONTENT_INVALID");
  expect(() => publicSource("https://example.com")).not.toThrow();
});
it("reads only the fixed anonymous public route with no ambient authority", async () => {
  await serve(
    (req, res) => {
      expect(req.url).toBe("/public/v1/content/en");
      expect(req.headers.authorization).toBeUndefined();
      expect(req.headers.cookie).toBeUndefined();
      res.writeHead(200, { "content-type": "application/json" });
      res.end("{}");
    },
    async (origin) => {
      expect(Buffer.from(await publicSource(origin, true)()).toString()).toBe(
        "{}",
      );
    },
  );
});
it("rejects redirects, wrong MIME, empty/oversized bodies and incomplete streams", async () => {
  for (const mode of [
    "redirect",
    "mime",
    "empty",
    "length",
    "actual",
    "partial",
  ]) {
    await serve(
      (_req, res) => {
        if (mode === "redirect") {
          res.writeHead(302, { location: "http://127.0.0.1:1/" });
          res.end();
          return;
        }
        res.writeHead(200, {
          "content-type": mode === "mime" ? "text/html" : "application/json",
          ...(mode === "length"
            ? { "content-length": String(packageLimit + 1) }
            : {}),
        });
        if (mode === "partial") {
          res.write("{");
          res.destroy();
        } else
          res.end(
            mode === "empty"
              ? ""
              : mode === "actual"
                ? Buffer.alloc(packageLimit + 1)
                : "{}",
          );
      },
      async (origin) => {
        await expect(publicSource(origin, true)()).rejects.toThrow(
          /^CMS_CONTENT_INVALID$/,
        );
      },
    );
  }
});
it("times out stalled body reads without returning partial bytes or raw errors", async () => {
  await serve(
    (_req, res) => {
      res.writeHead(200, { "content-type": "application/json" });
      res.write("{");
    },
    async (origin) => {
      await expect(publicSource(origin, true, 80)()).rejects.toThrow(
        /^CMS_CONTENT_INVALID$/,
      );
    },
  );
});
