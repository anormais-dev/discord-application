import { DiscordSDK } from "@discord/embedded-app-sdk";

export const discordSdk = new DiscordSDK(process.env.DISCORD_CLIENT_ID);

// Fluxo: ready -> authorize -> /.proxy/api/token -> authenticate.
export async function setupDiscord(): Promise<{ accessToken: string; instanceId: string }> {
  throw new Error("não implementado");
}
