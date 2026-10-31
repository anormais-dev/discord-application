import type { CSSProperties } from "react";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { wsSend } from "../store/roomSlice";
import { selectMe, selectParticipantsById, selectPermissions, selectRoom } from "../store/selectors";
import type { Participant } from "../types/room";
import { plural, shortFormat, teamName } from "../utils/format";
import { SOLO_COLOR, TEAM_COLORS, pickNumber } from "./DraftView";
import { Monkey } from "./Monkey";

// Fim do sorteio. Cada pessoa vê o próprio resultado em destaque: entrou num time ou ficou de fora.
// Quem não participou (inclusive o admin, se só organizou) vê o resultado neutro.
export function ResultView() {
  const dispatch = useAppDispatch();
  const room = useAppSelector(selectRoom);
  const me = useAppSelector(selectMe);
  const byId = useAppSelector(selectParticipantsById);
  const perms = useAppSelector(selectPermissions);
  if (!room) return null;

  const teams = room.format.teams;
  const leftOut = room.pool.map((id) => byId[id]).filter(Boolean);
  const myTeam = room.teams.findIndex((t) => t.some((p) => p.id === me));
  const mySlot = myTeam >= 0 ? room.teams[myTeam].findIndex((p) => p.id === me) : -1;
  const meLeftOut = leftOut.some((p) => p.id === me);
  const color = (i: number) => (teams === 1 ? SOLO_COLOR : TEAM_COLORS[i]);

  return (
    <main className="result">
      {myTeam >= 0 ? (
        <section className="panel hero" style={{ "--tc": color(myTeam) } as CSSProperties}>
          <Monkey mood="happy" size={120} party />
          <div>
            <div className="px hero-k">{teams === 1 ? "Você foi" : "Você está no"}</div>
            <div className="px hero-t">{teams === 1 ? "Selecionado" : teamName(myTeam)}</div>
            <p className="hero-d">
              Você saiu no <b>pick #{pickNumber(myTeam, mySlot, teams)}</b>.
              {teams === 2 && ` Seu time joga contra o ${teamName(1 - myTeam)}, chama a galera na call.`}
            </p>
          </div>
        </section>
      ) : meLeftOut ? (
        <section className="panel hero" style={{ "--tc": "var(--yellow)" } as CSSProperties}>
          <Monkey mood="stand" size={120} />
          <div>
            <div className="px hero-k">Não foi dessa vez</div>
            <div className="px hero-t">Você ficou de fora</div>
            <p className="hero-d">
              O formato era <b>{shortFormat(room.format)}</b>: eram {room.participants.length} na roda para{" "}
              {plural(perms.slots, "vaga", "vagas")}. Na próxima rodada todo mundo volta para a roda.
            </p>
          </div>
        </section>
      ) : (
        <section className="result-title">
          <Monkey mood="happy" size={72} />
          <div>
            <div className="px big-title">{teams === 1 ? "Time definido" : "Times definidos"}</div>
            <div className="muted">
              {shortFormat(room.format)} · {plural(room.participants.length, "jogador", "jogadores")} na roda
            </div>
          </div>
        </section>
      )}

      <div className="result-teams">
        <TeamCard title={teams === 1 ? "Selecionados" : teamName(0)} color={color(0)} members={room.teams[0] ?? []} me={me} team={0} teams={teams} />
        {teams === 2 && <div className="vs px">VS</div>}
        {teams === 2 && <TeamCard title={teamName(1)} color={color(1)} members={room.teams[1] ?? []} me={me} team={1} teams={teams} />}
        {teams === 1 && leftOut.length > 0 && (
          <TeamCard title="Ficaram de fora" color="var(--muted)" members={leftOut} me={me} />
        )}
      </div>
      {teams === 2 && leftOut.length > 0 && (
        <p className="muted result-out">Ficaram de fora: {leftOut.map((p) => p.name).join(", ")}</p>
      )}

      <div className="result-bar">
        {teams === 2 && (
          <div className="map-strip">
            <span className="lbl px">Mapa</span>
            <span className={`map-name ${room.map ? "set" : ""}`}>{room.map || "Ainda não sorteado"}</span>
            {perms.isAdmin && (
              <button className="btn btn-yellow px" onClick={() => dispatch(wsSend({ type: "set_map_open", enabled: true }))}>
                Sortear mapa
              </button>
            )}
          </div>
        )}
        {perms.isAdmin ? (
          <button className="btn btn-ghost" onClick={() => dispatch(wsSend({ type: "reset" }))}>
            Nova rodada
          </button>
        ) : (
          <p className="muted wait">
            {teams === 2 ? "Esperando o admin sortear o mapa ou começar uma nova rodada" : "Esperando o admin começar uma nova rodada"}
          </p>
        )}
      </div>
    </main>
  );
}

function TeamCard({
  title,
  color,
  members,
  me,
  team,
  teams,
}: {
  title: string;
  color: string;
  members: Participant[];
  me: string | null;
  team?: number;
  teams?: number;
}) {
  return (
    <section className="panel result-team" style={{ "--tc": color } as CSSProperties}>
      <div className="team-h">
        <span className="px">{title}</span>
      </div>
      <ol className="team-list">
        {members.map((p, slot) => (
          <li key={p.id} className={p.id === me ? "me" : ""}>
            <img className="avatar" src={p.avatar} alt="" width={36} height={36} />
            <span className="nm">{p.name}</span>
            {p.id === me && <span className="tag tag-yellow">Você</span>}
            {team !== undefined && teams !== undefined && <span className="pick">#{pickNumber(team, slot, teams)}</span>}
          </li>
        ))}
      </ol>
    </section>
  );
}
