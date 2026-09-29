import { configureStore } from "@reduxjs/toolkit";
import { authApi } from "../features/auth/authApi";
import { roomReducer } from "../features/room/roomSlice";
import { wsMiddleware } from "../features/room/wsMiddleware";

export const store = configureStore({
  reducer: {
    room: roomReducer,
    [authApi.reducerPath]: authApi.reducer,
  },
  middleware: (getDefault) => getDefault().concat(authApi.middleware, wsMiddleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
