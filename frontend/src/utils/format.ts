import type { Format } from "../types/room";

export const formatLabel = (f: Format) => (f.teams === 1 ? `Time de ${f.size}` : `${f.size}x${f.size}`);

export const teamsLabel = (teams: number) => (teams === 1 ? "1 time" : `${teams} times`);

export const formatSlots = (f: Format) => f.teams * f.size;

export const teamName = (i: number) => `Time ${String.fromCharCode(65 + i)}`;

// Rótulo curto do redesign: quantidade de times x jogadores por time ("2x5", "1x3").
export const shortFormat = (f: Format) => `${f.teams}x${f.size}`;

// "1 jogador" / "3 jogadores".
export const plural = (n: number, one: string, many: string) => `${n} ${n === 1 ? one : many}`;
