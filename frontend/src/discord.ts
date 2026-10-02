import { DiscordSDK } from "@discord/embedded-app-sdk";
import { createAsyncThunk } from "@reduxjs/toolkit";
import { authApi } from "./features/auth/authApi";
import { sessionFailed, wsConnect } from "./features/room/roomSlice";
import { describe } from "./utils/errors";

// Fluxo: ready -> authorize -> /.proxy/api/token -> authenticate -> WebSocket.
export const startSession = createAsyncThunk("session/start", async (_, { dispatch }) => {
  let step = "ready";
  try {
    // Criado aqui porque o construtor falha fora do Discord.
    const discordSdk = new DiscordSDK(process.env.DISCORD_CLIENT_ID);
    await discordSdk.ready();
    step = "authorize";
    const { code } = await discordSdk.commands.authorize({
      client_id: process.env.DISCORD_CLIENT_ID,
      response_type: "code",
      state: "",
      prompt: "none",
      scope: ["identify"],
    });
    step = "token";
    const { access_token } = await dispatch(authApi.endpoints.exchangeCode.initiate({ code })).unwrap();
    step = "authenticate";
    await discordSdk.commands.authenticate({ access_token });
    dispatch(wsConnect({ accessToken: access_token, instanceId: discordSdk.instanceId }));
  } catch (err) {
    console.error(`falha na etapa ${step}`, err);
    dispatch(sessionFailed(`Etapa ${step}: ${describe(err)}`));
    throw err;
  }
});
