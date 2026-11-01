import { createSelector } from "@reduxjs/toolkit";
import { formatSlots } from "../utils/format";
import type { RootState } from ".";

// Mínimo para girar, igual ao backend (MinPlayers).
export const MIN_PLAYERS = 2;

export const selectStatus = (s: RootState) => s.room.status;
export const selectRoom = (s: RootState) => s.room.room;
export const selectMe = (s: RootState) => s.room.me;
export const selectFormats = (s: RootState) => s.room.formats;
export const selectError = (s: RootState) => s.room.error;
export const selectFatal = (s: RootState) => s.room.fatal;
export const selectCompact = (s: RootState) => s.room.compact;
export const selectMuted = (s: RootState) => s.room.muted;

export const selectParticipantsById = createSelector([selectRoom], (room) =>
  Object.fromEntries((room?.participants ?? []).map((p) => [p.id, p])),
);

export const selectPermissions = createSelector([selectRoom, selectMe], (room, me) => {
  const isParticipant = !!room && !!me && room.participants.some((p) => p.id === me);
  const isAdmin = !!room && !!me && room.adminId === me;
  const slots = room ? formatSlots(room.format) : 0;
  const notReady = room ? room.participants.filter((p) => !room.ready[p.id]).length : 0;
  return {
    isParticipant,
    isAdmin,
    canSpin: !!room && !!me && (isAdmin || !!room.spinners[me]),
    canJoin: !!room && !isParticipant && room.phase === "lobby",
    canLeave: isParticipant && room?.phase !== "drafting",
    slots,
    // Com uma pessoa só não há o que sortear, mesmo no 1x1: precisa de pelo menos MIN_PLAYERS.
    enoughPlayers: !!room && room.participants.length >= Math.max(slots, MIN_PLAYERS),
    isReady: !!room && !!me && !!room.ready[me],
    notReady,
    allReady: !!room && notReady === 0,
  };
});
