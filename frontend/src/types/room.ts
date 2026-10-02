export type Phase = "lobby" | "drafting" | "finished";

export interface Format {
  teams: number;
  size: number;
}

export interface Participant {
  id: string;
  name: string;
  avatar: string;
  online: boolean;
  joinedAt: string;
}

// Giro em andamento. Todos os clientes animam a partir destes dados.
export interface SpinState {
  winnerId: string;
  poolSnapshot: string[];
  targetOffset: number;
  rotations: number;
  startedAt: string;
  durationMs: number;
}

export interface RoomState {
  participants: Participant[];
  adminId: string;
  spinners: Record<string, boolean>;
  format: Format;
  pool: string[];
  teams: Participant[][];
  phase: Phase;
  spin: SpinState | null;
  picks: number;
}

export type ClientMessage =
  | { type: "auth"; accessToken: string }
  | { type: "join" }
  | { type: "leave" }
  | { type: "set_format"; format: Format }
  | { type: "set_spinner"; userId: string; enabled: boolean }
  | { type: "spin" }
  | { type: "reset" };

export type ServerMessage =
  | { type: "state"; room: RoomState; you: string; formats: Format[] }
  | { type: "error"; message: string };
