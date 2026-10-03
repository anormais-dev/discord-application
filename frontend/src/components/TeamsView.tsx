import { useAppSelector } from "../store/hooks";
import { selectMe, selectParticipantsById, selectRoom } from "../store/selectors";
import { teamName } from "../utils/format";

export function TeamsView() {
  const room = useAppSelector(selectRoom);
  const byId = useAppSelector(selectParticipantsById);
  const me = useAppSelector(selectMe);
  if (!room || room.phase === "lobby") return null;

  const finished = room.phase === "finished";
  const single = room.teams.length === 1;
  const heading = finished ? (single ? "Time formado!" : "Times formados!") : single ? "Montando o time" : "Montando os times";
  const leftOut = finished ? room.pool.map((id) => byId[id]).filter(Boolean) : [];
  const meLeftOut = !!me && leftOut.some((p) => p.id === me);

  return (
    <section className="panel">
      <h2>{heading}</h2>
      {meLeftOut && (
        <p className="notice" role="status">
          Você não foi selecionado desta vez.
        </p>
      )}
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
