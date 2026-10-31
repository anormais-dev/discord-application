import type { Middleware } from "@reduxjs/toolkit";
import { errorReceived, stateReceived, statusChanged, wsConnect, wsSend } from "./roomSlice";
import type { ServerMessage } from "../types/room";
import { apiBase } from "../utils/runtime";

const RECONNECT_DELAY_MS = 2000;

// Abre wss://<host>/.proxy/api/ws (ou /api/ws fora do Discord), despacha eventos do servidor e envia comandos.
export const wsMiddleware: Middleware = (store) => {
  let socket: WebSocket | null = null;

  const open = (accessToken: string, instanceId: string, guildId?: string | null, channelId?: string | null) => {
    const protocol = location.protocol === "https:" ? "wss" : "ws";
    // guild e channel separam a sala por call; sem eles o backend usa a instance.
    const query = new URLSearchParams({ instance: instanceId });
    if (guildId) query.set("guild", guildId);
    if (channelId) query.set("channel", channelId);
    const url = `${protocol}://${location.host}${apiBase}/ws?${query}`;
    store.dispatch(statusChanged("connecting"));

    const ws = new WebSocket(url);
    socket = ws;
    ws.onopen = () => ws.send(JSON.stringify({ type: "auth", accessToken, guildId }));
    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data) as ServerMessage;
      if (msg.type === "state") {
        store.dispatch(stateReceived({ room: msg.room, you: msg.you, formats: msg.formats }));
      } else if (msg.type === "error") {
        store.dispatch(errorReceived(msg.message));
      }
    };
    ws.onclose = () => {
      if (socket !== ws) return;
      socket = null;
      store.dispatch(statusChanged("disconnected"));
      setTimeout(() => open(accessToken, instanceId, guildId, channelId), RECONNECT_DELAY_MS);
    };
  };

  return (next) => (action) => {
    if (wsConnect.match(action)) {
      open(action.payload.accessToken, action.payload.instanceId, action.payload.guildId, action.payload.channelId);
    } else if (wsSend.match(action)) {
      if (socket?.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(action.payload));
      }
    }
    return next(action);
  };
};
