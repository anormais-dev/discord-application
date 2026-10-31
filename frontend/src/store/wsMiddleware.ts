import type { Middleware } from "@reduxjs/toolkit";
import { errorReceived, stateReceived, statusChanged, wsConnect, wsSend } from "./roomSlice";
import type { ServerMessage } from "../types/room";
import { apiBase } from "../utils/runtime";

const RECONNECT_MIN_MS = 1000;
const RECONNECT_MAX_MS = 10000;

type ConnectParams = ReturnType<typeof wsConnect>["payload"];

// Abre wss://<host>/.proxy/api/ws (ou /api/ws fora do Discord), despacha eventos do servidor e envia comandos.
export const wsMiddleware: Middleware = (store) => {
  let socket: WebSocket | null = null;
  let params: ConnectParams | null = null;
  let retryTimer: ReturnType<typeof setTimeout> | undefined;
  let retryDelay = RECONNECT_MIN_MS;

  const open = ({ accessToken, instanceId, guildId }: ConnectParams) => {
    clearTimeout(retryTimer);
    const protocol = location.protocol === "https:" ? "wss" : "ws";
    const url = `${protocol}://${location.host}${apiBase}/ws?instance=${encodeURIComponent(instanceId)}`;
    store.dispatch(statusChanged("connecting"));

    const ws = new WebSocket(url);
    socket = ws;
    ws.onopen = () => ws.send(JSON.stringify({ type: "auth", accessToken, guildId }));
    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data) as ServerMessage;
      if (msg.type === "state") {
        retryDelay = RECONNECT_MIN_MS;
        store.dispatch(stateReceived({ room: msg.room, you: msg.you, formats: msg.formats }));
      } else if (msg.type === "error") {
        store.dispatch(errorReceived(msg.message));
      }
    };
    ws.onclose = () => {
      if (socket !== ws) return;
      socket = null;
      store.dispatch(statusChanged("disconnected"));
      retryTimer = setTimeout(reconnect, retryDelay);
      retryDelay = Math.min(retryDelay * 2, RECONNECT_MAX_MS);
    };
  };

  const reconnect = () => {
    if (params && !socket) open(params);
  };

  // Com a aba em segundo plano o navegador atrasa os timers; ao voltar, reconecta na hora em vez de esperar o backoff.
  document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "visible") reconnect();
  });
  window.addEventListener("online", reconnect);

  return (next) => (action) => {
    if (wsConnect.match(action)) {
      params = action.payload;
      open(action.payload);
    } else if (wsSend.match(action)) {
      if (socket?.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(action.payload));
      }
    }
    return next(action);
  };
};
