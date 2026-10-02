import { configureStore } from "@reduxjs/toolkit";
import { authApi } from "../services/api";
import { roomReducer } from "./roomSlice";
import { wsMiddleware } from "./wsMiddleware";

export const store = configureStore({
  reducer: {
    room: roomReducer,
    [authApi.reducerPath]: authApi.reducer,
  },
  middleware: (getDefault) => getDefault().concat(authApi.middleware, wsMiddleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
