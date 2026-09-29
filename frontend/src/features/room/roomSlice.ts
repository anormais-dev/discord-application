import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { Format, RoomState, SpinState } from "./types";

type Status = "idle" | "connecting" | "connected" | "error";

interface RoomSliceState {
  status: Status;
  me: string | null;
  room: RoomState | null;
  formats: Format[];
  activeSpin: SpinState | null;
  error: string | null;
}

const initialState: RoomSliceState = {
  status: "idle",
  me: null,
  room: null,
  formats: [],
  activeSpin: null,
  error: null,
};

const roomSlice = createSlice({
  name: "room",
  initialState,
  reducers: {
    stateReceived(state, action: PayloadAction<{ room: RoomState; you: string; formats: Format[] }>) {
      state.status = "connected";
      state.room = action.payload.room;
      state.me = action.payload.you;
      state.formats = action.payload.formats;
    },
    spinStarted(state, action: PayloadAction<SpinState>) {
      state.activeSpin = action.payload;
    },
    errorReceived(state, action: PayloadAction<string>) {
      state.error = action.payload;
    },
  },
});

export const { stateReceived, spinStarted, errorReceived } = roomSlice.actions;
export const roomReducer = roomSlice.reducer;
