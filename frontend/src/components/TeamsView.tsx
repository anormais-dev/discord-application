import { useAppSelector } from "../store/hooks";
import { selectParticipantsById, selectRoom } from "../store/selectors";
import { teamName } from "../utils/format";

export function TeamsView() {
  const room = useAppSelector(selectRoom);
  const byId = useAppSelector(selectParticipantsById);
  if (!room || room.phase === "lobby") return null;

  const finished = room.phase === "finished";
  const leftOut = finished ? room.pool.map((id) => byId[id]).filter(Boolean) : [];

  return (
    <section className="panel">
      <h2>{finished ? "Times formados!" : "Montando os times"}</h2>
      <div className="teams">
        {room.teams.map((team, i) => (
          <div key={i} className="team">
            <h3>{teamName(i)}</h3>
            <ol>
              {Array.from({ length: room.format.size }, (_, slot) => {
                const p = team[slot];
                return (
                  <li key={slot} className={p ? "" : "empty"}>
                    {p ? (
                      <>
                        <img src={p.avatar} alt="" width={24} height={24} />
                        {p.name}
                      </>
                    ) : (
                      "Aguardando sorteio"
                    )}
                  </li>
                );
              })}
            </ol>
          </div>
        ))}
      </div>
      {leftOut.length > 0 && (
        <p className="muted">Ficaram de fora: {leftOut.map((p) => p.name).join(", ")}</p>
      )}
    </section>
  );
}
