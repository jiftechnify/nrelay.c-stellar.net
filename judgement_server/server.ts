import { Hono } from "https://deno.land/x/hono@v3.12.10/mod.ts";
import { isSyncInProgress, syncFollowerList } from "./followers_updater.ts";
import { AppContext } from "./types.ts";

import * as log from "https://deno.land/std@0.208.0/log/mod.ts";

export const launchServer = (ctx: AppContext, signal: AbortSignal) => {
  log.info("launching judgement server...");

  const app = new Hono();

  app.get("/is_follower", async (c) => {
    const pubkey = c.req.query("pubkey");
    if (pubkey === undefined) {
      return c.text("param 'pubkey' is required", 400);
    }
    const { value } = await ctx.kv.get(["followers", pubkey]);
    return c.json({ isFollower: value != null });
  });

  app.put("/sync", (c) => {
    if (isSyncInProgress()) {
      return c.text("sync in progress", 400);
    }
    syncFollowerList(ctx, { force: true });
    return c.text("sync triggered", 200);
  });

  app.get("/health", (c) => {
    return c.text("ok");
  });

  Deno.serve({
    port: 8080,
    signal,
    onListen: ({ hostname, port }) => {
      log.info(`judgement server listening on ${hostname}:${port}`);
    },
  }, app.fetch);
};
