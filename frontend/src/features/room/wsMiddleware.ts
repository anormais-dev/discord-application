import type { Middleware } from "@reduxjs/toolkit";

// Abre wss://<host>/.proxy/api/ws, despacha eventos do servidor e envia comandos.
export const wsMiddleware: Middleware = () => (next) => (action) => next(action);
