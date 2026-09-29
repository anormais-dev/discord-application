import { DiscordSDK } from "@discord/embedded-app-sdk";
import { createAsyncThunk } from "@reduxjs/toolkit";
import { authApi } from "./features/auth/authApi";
import { statusChanged, wsConnect } from "./features/room/roomSlice";

export const discordSdk = new DiscordSDK(process.env.DISCORD_CLIENT_ID);

// Fluxo: ready -> authorize -> /.proxy/api/token -> authenticate -> WebSocket.
export const startSession = createAsyncThunk("session/start", async (_, { dispatch }) => {
  try {
    await discordSdk.ready();
    const { code } = await discordSdk.commands.authorize({
      client_id: process.env.DISCORD_CLIENT_ID,
      response_type: "code",
      state: "",
      prompt: "none",
      scope: ["identify"],
    });
    const { access_token } = await dispatch(authApi.endpoints.exchangeCode.initiate({ code })).unwrap();
    await discordSdk.commands.authenticate({ access_token });
    dispatch(wsConnect({ accessToken: access_token, instanceId: discordSdk.instanceId }));
  } catch (err) {
    dispatch(statusChanged("error"));
    throw err;
  }
});
