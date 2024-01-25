import { open, writeFile } from "fs/promises";
import { NostrFetcher } from "nostr-fetch";
import "websocket-polyfill";

const relays = [
  "wss://directory.yabu.me",
  "wss://purplepag.es",
  "wss://relay.nostr.band",
  "wss://nrelay.c-stellar.net",
];
const pubkey =
  "d1d1747115d16751a97c239f46ec1703292c3b7e9988b9ebdd4ec4705b15ed44";

const pubkeyDbPath = "pubkey_db.txt";
const removeThresholdMs = 12 * 60 * 60 * 1000; // 12 Hrs

/**
 *
 * @param {import('fs/promises').FileHandle} file
 * @returns {Promise<Map<string, string>>}
 */
const readFollowersDb = async (file) => {
  const txt = await file.readFile({ encoding: "utf8" });
  return new Map(
    txt
      .split("\n")
      .filter((l) => l.trim() !== "")
      .map((l) => l.split(","))
  );
};
/**
 *
 * @param {import('fs/promises').FileHandle} file
 * @param {Map<string, string>} data
 * @returns {Promise<void>}
 */
const writeFollowersDb = async (file, data) => {
  await file.truncate();

  const txt = [...data.entries()]
    .map(([pubkey, time]) => `${pubkey},${time}`)
    .join("\n");
  await file.writeFile(txt);
};

/**
 *
 * @param {Set<string>} followers
 * @returns {Priomise<string[]>}
 */
const updateFollowersDb = async (followers) => {
  const file = await open(pubkeyDbPath, "a+");
  const db = await readFollowersDb(file);

  const now = new Date().getTime();

  // remove pubkeys that are not in followers list recently
  [...db.entries()]
    .filter(
      ([_, lastSeenAt]) =>
        lastSeenAt !== "manual" && now - Number(lastSeenAt) > removeThresholdMs
    )
    .forEach(([pk]) => {
      db.delete(pk);
    });
  // add/update entries for current followers
  for (const f of followers) {
    db.set(f, String(now));
  }

  await writeFollowersDb(file, db);
  await file.close();

  // return pubkeys in db
  return [...db.keys()];
};

const main = async () => {
  const fetcher = NostrFetcher.init({ minLogLevel: "info" });

  const ac = new AbortController();
  setTimeout(() => {
    ac.abort();
  }, 10 * 60 * 1000);
  
  const followingMe = await fetcher.fetchAllEvents(
    relays,
    { kinds: [3], "#p": [pubkey] },
    {},
    { abortSubBeforeEoseTimeoutMs: 5000, abortSignal: ac.signal }
  );
  fetcher.shutdown();

  const whitelist = await updateFollowersDb(
    new Set(followingMe.map((ev) => ev.pubkey))
  );

  await writeFile("whitelist.txt", whitelist.join("\n"), { encoding: "utf8" });
};

main()
  .then(() => process.exit(0))
  .catch((err) => {
    console.error(err);
    process.exit(1);
  });
