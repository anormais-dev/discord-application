import type { CSSProperties } from "react";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { wsSend } from "../store/roomSlice";
import { selectParticipantsById, selectPermissions, selectRoom } from "../store/selectors";
import type { Participant } from "../types/room";
import { teamName } from "../utils/format";
import { PlayerWheel } from "./PlayerWheel";

// Cores por time; com um time só, os selecionados ficam em amarelo.
export const TEAM_COLORS = ["var(--red)", "var(--cyan)"];
export const SOLO_COLOR = "var(--yellow)";

// Os picks alternam entre os times: o n-ésimo de cada time saiu no giro slot * times + time + 1.
export const pickNumber = (team: number, slot: number, teams: number) => slot * teams + team + 1;

// Sorteio em andamento: com 2 times, um de cada lado da roda; com 1, selecionados de um lado e a roda do outro.
export function DraftView() {
  const dispatch = useAppDispatch();
  const room = useAppSelector(selectRoom);
  const byId = useAppSelector(selectParticipantsById);
  const perms = useAppSelector(selectPermissions);
  if (!room) return null;

  const teams = room.format.teams;
  const total = perms.slots;
  const nextTeam = room.picks % teams;
  const pool = room.pool.map((id) => byId[id]).filter(Boolean);

  const column = (i: number) => (
    <TeamColumn
      key={i}
      title={teams === 1 ? "Selecionados" : teamName(i)}
      color={teams === 1 ? SOLO_COLOR : TEAM_COLORS[i]}
      members={room.teams[i] ?? []}
      size={room.format.size}
      team={i}
      teams={teams}
      next={i === nextTeam}
    />
  );

  return (
    <main className="draft">
      {column(0)}
      <div className="draft-center">
        <div className="round">
          <span className="px">
            Pick {Math.min(room.picks + 1, total)}/{total}
          </span>
          <div className="round-bar" aria-hidden="true">
            {Array.from({ length: total }, (_, i) => (
              <span key={i} className={i < room.picks ? "on" : i === room.picks ? "now" : ""} />
            ))}
          </div>
        </div>
        <PlayerWheel />
        <div className="draft-actions">
          {perms.canSpin && !room.autoSpin && (
            <button className="btn btn-red px" disabled={!!room.spin} onClick={() => dispatch(wsSend({ type: "spin" }))}>
              Girar
            </button>
          )}
          {(perms.canSpin || room.autoSpin) && (
            <button
              className={`switch ${room.autoSpin ? "on" : ""}`}
              disabled={!perms.canSpin}
              aria-pressed={room.autoSpin}
              onClick={() => dispatch(wsSend({ type: "set_auto_spin", enabled: !room.autoSpin }))}
            >
              <span className="switch-box" />
              Giro auto a cada 2s
            </button>
          )}
          {!perms.canSpin && !room.autoSpin && <p className="hint">Esperando o admin girar</p>}
        </div>
      </div>
      {teams === 1 ? <PoolColumn people={pool} /> : column(1)}
    </main>
  );
}

function TeamColumn({
  title,
  color,
  members,
  size,
  team,
  teams,
  next,
}: {
  title: string;
  color: string;
  members: Participant[];
  size: number;
  team: number;
  teams: number;
  next: boolean;
}) {
  return (
    <section className="panel team-col" style={{ "--tc": color } as CSSProperties}>
      <div className="team-h">
        <span className="px">{title}</span>
        <span className="count">
          {members.length}
          <small>/{size}</small>
        </span>
      </div>
      <ol className="team-list">
        {Array.from({ length: size }, (_, slot) => {
          const p = members[slot];
          if (p) {
            return (
              <li key={slot}>
                <img className="avatar" src={p.avatar} alt="" width={32} height={32} />
                <span className="nm">{p.name}</span>
                <span className="pick px">#{pickNumber(team, slot, teams)}</span>
              </li>
            );
          }
          const isNext = next && slot === members.length;
          return (
            <li key={slot} className={`empty ${isNext ? "next" : ""}`}>
              <span className="ph" />
              <span>{isNext ? "Próximo pick" : "Vaga"}</span>
            </li>
          );
        })}
      </ol>
    </section>
  );
}

function PoolColumn({ people }: { people: Participant[] }) {
  return (
    <section className="panel team-col pool" style={{ "--tc": "var(--muted)" } as CSSProperties}>
      <div className="team-h">
        <span className="px">Na roda</span>
        <span className="count">{people.length}</span>
      </div>
      <ol className="team-list">
        {people.map((p) => (
          <li key={p.id}>
            <img className="avatar" src={p.avatar} alt="" width={32} height={32} />
            <span className="nm">{p.name}</span>
          </li>
        ))}
      </ol>
    </section>
  );
}
