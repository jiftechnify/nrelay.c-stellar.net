import { z } from "https://deno.land/x/zod@v3.22.4/mod.ts";

const configSchema = z.object({
  pubkey: z.string().length(64),
  relays: z.array(z.string().startsWith("wss://")).min(1),
  kvPath: z.string().min(1),
});

export const parseConfig = (rawConfig: unknown) => {
  return configSchema.safeParse(rawConfig);
}
