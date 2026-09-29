import { useAppDispatch, useAppSelector } from "../app/hooks";
import { selectFormats, selectPermissions, selectRoom, wsSend } from "../features/room/roomSlice";
import type { Format } from "../features/room/types";

const formatLabel = (f: Format) => `${f.size}x${f.size}`;

export function AdminControls() {
  const dispatch = useAppDispatch();
  const room = useAppSelector(selectRoom);
  const formats = useAppSelector(selectFormats);
  const perms = useAppSelector(selectPermissions);
  if (!room) return null;

  const spinning = !!room.spin;
  const missing = perms.slots - room.participants.length;
  const showSpin = perms.canSpin && room.phase !== "finished";

  return (
    <section className="controls">
      {room.phase === "lobby" && (
        <div className="formats">
          <span className="muted">Formato</span>
          {perms.isAdmin ? (
            formats.map((f) => (
              <button
                key={formatLabel(f)}
                className={`chip ${f.teams === room.format.teams && f.size === room.format.size ? "active" : ""}`}
                onClick={() => dispatch(wsSend({ type: "set_format", format: f }))}
              >
                {formatLabel(f)}
              </button>
            ))
          ) : (
            <span className="chip active">{formatLabel(room.format)}</span>
          )}
        </div>
      )}

      <div className="actions">
        {perms.canJoin && (
          <button className="primary" onClick={() => dispatch(wsSend({ type: "join" }))}>
            Participar
          </button>
        )}
        {showSpin && (
          <button
            className="primary"
            disabled={spinning || (room.phase === "lobby" && !perms.enoughPlayers)}
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
        {perms.canLeave && (
          <button className="secondary" onClick={() => dispatch(wsSend({ type: "leave" }))}>
            Sair
          </button>
        )}
      </div>

      {room.phase === "lobby" && missing > 0 && (
        <p className="muted">
          {missing === 1 ? "Falta 1 participante" : `Faltam ${missing} participantes`} para o {formatLabel(room.format)}
        </p>
      )}
      {room.phase === "drafting" && !perms.canSpin && <p className="muted">Aguardando o admin girar</p>}
      {room.phase === "finished" && !perms.isAdmin && <p className="muted">Aguardando o admin iniciar uma nova girada</p>}
    </section>
  );
}
