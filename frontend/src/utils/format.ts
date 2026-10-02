import type { Format } from "../types/room";

export const formatLabel = (f: Format) => `${f.size}x${f.size}`;

export const formatSlots = (f: Format) => f.teams * f.size;

export const teamName = (i: number) => `Time ${String.fromCharCode(65 + i)}`;
