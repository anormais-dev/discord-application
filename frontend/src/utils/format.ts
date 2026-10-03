import type { Format } from "../types/room";

export const formatLabel = (f: Format) => (f.teams === 1 ? `Time de ${f.size}` : `${f.size}x${f.size}`);

export const teamsLabel = (teams: number) => (teams === 1 ? "1 time" : `${teams} times`);

export const formatSlots = (f: Format) => f.teams * f.size;

export const teamName = (i: number) => `Time ${String.fromCharCode(65 + i)}`;
