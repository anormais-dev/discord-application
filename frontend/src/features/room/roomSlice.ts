import { createAction, createSelector, createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { RootState } from "../../app/store";
import type { ClientMessage, Format, RoomState } from "../../types/room";
import { formatSlots } from "../../utils/format";

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
export const wsConnect = createAction<{ accessToken: string; instanceId: string }>("room/wsConnect");
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

export const selectStatus = (s: RootState) => s.room.status;
export const selectRoom = (s: RootState) => s.room.room;
export const selectMe = (s: RootState) => s.room.me;
export const selectFormats = (s: RootState) => s.room.formats;
export const selectError = (s: RootState) => s.room.error;
export const selectFatal = (s: RootState) => s.room.fatal;

export const selectParticipantsById = createSelector([selectRoom], (room) =>
  Object.fromEntries((room?.participants ?? []).map((p) => [p.id, p])),
);

export const selectPermissions = createSelector([selectRoom, selectMe], (room, me) => {
  const isParticipant = !!room && !!me && room.participants.some((p) => p.id === me);
  const isAdmin = !!room && !!me && room.adminId === me;
  const slots = room ? formatSlots(room.format) : 0;
  return {
    isParticipant,
    isAdmin,
    canSpin: !!room && !!me && (isAdmin || !!room.spinners[me]),
    canJoin: !!room && !isParticipant && room.phase === "lobby",
    canLeave: isParticipant && room?.phase !== "drafting",
    slots,
    enoughPlayers: !!room && room.participants.length >= slots,
  };
});
