import { createSelector } from "@reduxjs/toolkit";
import { formatSlots } from "../utils/format";
import type { RootState } from ".";

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
