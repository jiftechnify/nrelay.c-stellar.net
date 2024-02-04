import * as log from "https://deno.land/std@0.208.0/log/mod.ts";

import rawConfig from "./config.json" with { type: "json" };
import { parseConfig } from "./config.ts";
import { launchFollowersUpdater } from "./followers_updater.ts";
import { launchServer } from "./server.ts";

if (import.meta.main) {
  log.setup({
    handlers: {
      console: new log.handlers.ConsoleHandler("DEBUG", {
        formatter: ({ levelName, datetime, msg }) => {
          return `${datetime.toLocaleString()} [${levelName.padEnd(8)}] ${msg}`;
        },
      }),
    },
    loggers: {
      default: {
        handlers: ["console"],
      },
    },
  });

  const parseRes = parseConfig(rawConfig);
  if (!parseRes.success) {
    console.error(parseRes.error);
    Deno.exit(1);
  }
  const { pubkey, relays, kvPath } = parseRes.data;
  console.log(`config: ${JSON.stringify(parseRes.data)}`);

  const kv = await Deno.openKv(kvPath);
  const appCtx = { kv, pubkey, relays };

  // handle SIGTERM
  const ac = new AbortController();
  Deno.addSignalListener("SIGTERM", () => {
    console.log("received SIGTERM");
    ac.abort();
  });

  await launchFollowersUpdater(appCtx);
  launchServer(appCtx, ac.signal);
}
