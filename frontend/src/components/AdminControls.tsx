import { useAppDispatch, useAppSelector } from "../store/hooks";
import { wsSend } from "../store/roomSlice";
import { selectFormats, selectPermissions, selectRoom } from "../store/selectors";
import type { Format } from "../types/room";
import { formatLabel, teamsLabel } from "../utils/format";

export function AdminControls() {
  const dispatch = useAppDispatch();
  const room = useAppSelector(selectRoom);
  const formats = useAppSelector(selectFormats);
  const perms = useAppSelector(selectPermissions);
  if (!room) return null;

  const spinning = !!room.spin;
  const missing = perms.slots - room.participants.length;
  const showSpin = perms.canSpin && room.phase !== "finished";
  const teamCounts = [...new Set(formats.map((f) => f.teams))].sort((a, b) => a - b);

  const setFormat = (format: Format) => dispatch(wsSend({ type: "set_format", format }));
  // Trocar a quantidade de times mantém o tamanho quando ele existe no outro modo.
  const selectTeams = (teams: number) => {
    const options = formats.filter((f) => f.teams === teams);
    const next = options.find((f) => f.size === room.format.size) ?? options[0];
    if (next) setFormat(next);
  };

  return (
    <section className="controls">
      {room.phase === "lobby" &&
        (perms.isAdmin ? (
          <>
            <div className="formats">
              <span className="muted">Times</span>
              {teamCounts.map((teams) => (
                <button
                  key={teams}
                  className={`chip ${teams === room.format.teams ? "active" : ""}`}
                  onClick={() => selectTeams(teams)}
                >
                  {teamsLabel(teams)}
                </button>
              ))}
            </div>
            <div className="formats">
              <span className="muted">Formato</span>
              {formats
                .filter((f) => f.teams === room.format.teams)
                .map((f) => (
                  <button
                    key={formatLabel(f)}
                    className={`chip ${f.size === room.format.size ? "active" : ""}`}
                    onClick={() => setFormat(f)}
                  >
                    {formatLabel(f)}
                  </button>
                ))}
            </div>
          </>
        ) : (
          <div className="formats">
            <span className="muted">Formato</span>
            <span className="chip active">{formatLabel(room.format)}</span>
          </div>
        ))}

      {room.phase !== "finished" && (perms.canSpin || room.autoSpin) && (
        <div className="formats">
          <span className="muted">Giro automático</span>
          {perms.canSpin ? (
            <button
              className={`chip ${room.autoSpin ? "active" : ""}`}
              onClick={() => dispatch(wsSend({ type: "set_auto_spin", enabled: !room.autoSpin }))}
            >
              {room.autoSpin ? "Ligado" : "Desligado"}
            </button>
          ) : (
            <span className="chip active">Ligado</span>
          )}
        </div>
      )}

      <div className="actions">
        {perms.canJoin && (
          <button className="primary" onClick={() => dispatch(wsSend({ type: "join" }))}>
            Participar
          </button>
        )}
        {perms.isParticipant && room.phase === "lobby" && (
          <button
            className={perms.isReady ? "secondary" : "primary"}
            onClick={() => dispatch(wsSend({ type: "set_ready", enabled: !perms.isReady }))}
          >
            {perms.isReady ? "Cancelar pronto" : "Pronto"}
          </button>
        )}
        {showSpin && (
          <button
            className="primary"
            disabled={spinning || (room.phase === "lobby" && (!perms.enoughPlayers || !perms.allReady))}
            onClick={() => dispatch(wsSend({ type: "spin" }))}
          >
            {spinning ? "Girando..." : "Girar"}
          </button>
        )}
        {perms.isAdmin && room.phase === "finished" && (
          <button className="primary" onClick={() => dispatch(wsSend({ type: "reset" }))}>
            Nova girada
          </button>
        )}
        {perms.isAdmin && (
          <button className="secondary" onClick={() => dispatch(wsSend({ type: "set_map_open", enabled: true }))}>
            Escolher mapa
          </button>
        )}
        {perms.canLeave && (
          <button className="secondary" onClick={() => dispatch(wsSend({ type: "leave" }))}>
            Sair
          </button>
        )}
      </div>

      {room.map && <p className="muted">Mapa: {room.map}</p>}
      {room.phase === "lobby" && missing > 0 && (
        <p className="muted">
          {missing === 1 ? "Falta 1 participante" : `Faltam ${missing} participantes`} para o {formatLabel(room.format)}
        </p>
      )}
      {room.phase === "lobby" && missing <= 0 && perms.notReady > 0 && (
        <p className="muted">
          {perms.notReady === 1 ? "Aguardando 1 participante ficar pronto" : `Aguardando ${perms.notReady} participantes ficarem prontos`}
        </p>
      )}
      {room.phase === "lobby" && room.autoSpin && <p className="muted">Depois do primeiro giro, a roleta segue girando sozinha</p>}
      {room.phase === "drafting" && !perms.canSpin && !room.autoSpin && <p className="muted">Aguardando o admin girar</p>}
      {room.phase === "finished" && !perms.isAdmin && <p className="muted">Aguardando o admin iniciar uma nova girada</p>}
    </section>
  );
}
