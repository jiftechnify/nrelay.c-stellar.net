export type AppContext = {
  kv: Deno.Kv;
  pubkey: string;
  relays: string[];
};
