import { DiscordSDK } from "@discord/embedded-app-sdk";
import { createAsyncThunk } from "@reduxjs/toolkit";
import { authApi } from "./api";
import { sessionFailed, wsConnect } from "../store/roomSlice";
import { describe } from "../utils/errors";

// Fluxo: config -> ready -> authorize -> /.proxy/api/token -> authenticate -> WebSocket.
export const startSession = createAsyncThunk("session/start", async (_, { dispatch }) => {
  let step = "config";
  try {
    const { clientId } = await dispatch(authApi.endpoints.getConfig.initiate()).unwrap();
    step = "ready";
    // Criado aqui porque o construtor falha fora do Discord.
    const discordSdk = new DiscordSDK(clientId);
    await discordSdk.ready();
    step = "authorize";
    const { code } = await discordSdk.commands.authorize({
      client_id: clientId,
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
