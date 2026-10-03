import { createAction, createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { ClientMessage, Format, RoomState } from "../types/room";

type Status = "idle" | "connecting" | "connected" | "disconnected" | "error";

interface RoomSliceState {
  status: Status;
  me: string | null;
  room: RoomState | null;
  formats: Format[];
  error: string | null;
  fatal: string | null;
}

const initialState: RoomSliceState = {
  status: "idle",
  me: null,
  room: null,
  formats: [],
  error: null,
  fatal: null,
};

// Tratadas pelo wsMiddleware.
export const wsConnect = createAction<{ accessToken: string; instanceId: string; guildId?: string | null }>("room/wsConnect");
export const wsSend = createAction<ClientMessage>("room/wsSend");

const roomSlice = createSlice({
  name: "room",
  initialState,
  reducers: {
    statusChanged(state, action: PayloadAction<Status>) {
      state.status = action.payload;
    },
    stateReceived(state, action: PayloadAction<{ room: RoomState; you: string; formats: Format[] }>) {
      state.status = "connected";
      state.room = action.payload.room;
      state.me = action.payload.you;
      state.formats = action.payload.formats;
    },
    sessionFailed(state, action: PayloadAction<string>) {
      state.status = "error";
      state.fatal = action.payload;
    },
    errorReceived(state, action: PayloadAction<string>) {
      state.error = action.payload;
    },
    errorCleared(state) {
      state.error = null;
    },
  },
});

export const { statusChanged, stateReceived, sessionFailed, errorReceived, errorCleared } = roomSlice.actions;
export const roomReducer = roomSlice.reducer;
