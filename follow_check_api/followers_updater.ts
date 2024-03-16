import * as log from "https://deno.land/std@0.220.1/log/mod.ts";
import { NostrFetcher } from "npm:nostr-fetch@0.15.0";
import { createRxForwardReq, createRxNostr, uniq } from "npm:rx-nostr@2.7.0";
import { AppContext } from "./types.ts";

/**
 * Launch followers updater
 * - start subscription of follow list changes
 * - synchronize follower list
 * - schedule periodic synchronization
 */
export const launchFollowersUpdater = async (
  appCtx: AppContext,
  signal: AbortSignal,
) => {
  log.info("launching followers updater...");

  subNewFollowers(appCtx, signal);
  await syncFollowerList(appCtx);

  // sync follower list every 6 hours
  Deno.cron("periodic sync", { hour: { every: 6 } }, async () => {
    await syncFollowerList(appCtx);
  });
};

export const subNewFollowers = (
  { kv, pubkey, relays }: AppContext,
  signal: AbortSignal,
) => {
  log.info(`start: subscribe to new followers`);

  const rxn = createRxNostr();
  rxn.addDefaultRelays(
    relays.map((r) => ({ url: r, read: true, write: false })),
  );

  rxn.createConnectionStateObservable().subscribe(({ from, state }) => {
    log.info(`[${from}] ${state}`);
  });

  const rxReq = createRxForwardReq();
  rxn
    .use(rxReq)
    .pipe(uniq())
    .subscribe(async ({ event }) => {
      const { pubkey: newFollower } = event;
      await storeSingleFollower(kv, newFollower, currUnixtime());
    });
  rxReq.emit({ kinds: [3], "#p": [pubkey], limit: 0 });

  signal.addEventListener("abort", () => {
    log.info("aborting new followers subscription");
    rxn.dispose();
  }, { once: true });
};

let _syncInProgress = false;
export const isSyncInProgress = () => _syncInProgress;

export const syncFollowerList = async (
  { kv, pubkey, relays }: AppContext,
  { force } = { force: false },
) => {
  _syncInProgress = true;
  try {
    log.info("start to synchronize follower list");

    const now = currUnixtime();

    // skip if recently synced(within past 3 hrs)
    const lastSycnedAt = (await kv.get<number>(["last_synced_at"])).value;
    if (
      !force && lastSycnedAt != null &&
      now - lastSycnedAt < 3 * 60 * 60
    ) {
      log.info("recently synced, skipping sync");
      return;
    }

    const followers = await fetchAllFollowers(pubkey, relays);
    await storeFollowers(kv, followers, now);
    await evictFollowersRemovedMe(kv, now);

    await kv.set(["last_synced_at"], now);
    log.info("finish synchronizing follower list");
  } catch (err) {
    log.error(`failed to synchronize follower list: ${err}`);
  } finally {
    _syncInProgress = false;
  }
};

const storeFollowers = async (
  kv: Deno.Kv,
  followers: string[],
  unixtime: number,
): Promise<void> => {
  const res = await Promise.allSettled(
    followers.map((pk) => kv.set(["followers", pk], unixtime)),
  );

  const stat = res.reduce(
    (acc, r) => {
      switch (r.status) {
        case "fulfilled":
          return { ok: acc.ok + 1, err: acc.err };
        case "rejected":
          return { ok: acc.ok, err: acc.err + 1 };
      }
    },
    { ok: 0, err: 0 },
  );
  log.info(`ok: ${stat.ok}, err: ${stat.err}`);
};

const storeSingleFollower = async (
  kv: Deno.Kv,
  pk: string,
  unixtime: number,
): Promise<void> => {
  const followed = (await kv.get(["followers", pk])).value != null;
  if (followed) {
    log.debug(`already followed by ${pk}`);
    return;
  }

  log.info(`newly followed by ${pk}`);
  const res = await kv.set(["followers", pk], unixtime);
  if (res.ok) {
    log.info(`stored: ${pk}`);
  } else {
    log.error(`failed to store: ${pk}`);
  }
};

const currUnixtime = () => Math.floor(Date.now() / 1000);

const removeThresholdSec = 24 * 60 * 60; // 24 hours
const removedMe = (lastSeenSec: number, nowSec: number): boolean => {
  return nowSec - lastSeenSec > removeThresholdSec;
};

const evictFollowersRemovedMe = async (
  kv: Deno.Kv,
  now: number,
): Promise<void> => {
  const iter = kv.list<number>({ prefix: ["followers"] });

  for await (const { key, value: lastSeen } of iter) {
    if (removedMe(lastSeen, now)) {
      await kv.delete(key);
      log.info(`removed by: ${key}`);
    }
  }
};

const fetchAllFollowers = async (
  pubkey: string,
  relays: string[],
): Promise<string[]> => {
  const fetcher = NostrFetcher.init();

  const followingMe = await fetcher.fetchAllEvents(
    relays,
    { kinds: [3], "#p": [pubkey] },
    {},
    {
      abortSignal: AbortSignal.timeout(5 * 60 * 1000),
      abortSubBeforeEoseTimeoutMs: 5000,
    },
  );
  fetcher.shutdown();

  return followingMe.map((e) => e.pubkey);
};
