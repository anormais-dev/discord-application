import { useNow } from "../hooks/useNow";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { wsSend } from "../store/roomSlice";
import { selectMe, selectPermissions, selectRoom } from "../store/selectors";
import type { Participant } from "../types/room";

// Lista de quem está na roleta no lobby, com o estado de cada um e as ações do próprio usuário.
export function SquadPanel() {
  const dispatch = useAppDispatch();
  const room = useAppSelector(selectRoom);
  const me = useAppSelector(selectMe);
  const perms = useAppSelector(selectPermissions);
  const reconnecting = !!room?.participants.some((p) => !p.online && p.reconnectUntil);
  const now = useNow(reconnecting);
  if (!room) return null;

  const count = room.participants.length;

  return (
    <aside className="panel squad">
      <div className="panel-h">
        <span className="px">Esquadrão</span>
        <span className="count">
          {count}
          <small>/{perms.slots}</small>
        </span>
      </div>
      <div className="slots" aria-hidden="true">
        {Array.from({ length: perms.slots }, (_, i) => (
          <span key={i} className={i < count ? "on" : ""} />
        ))}
      </div>

      {count === 0 && <p className="muted squad-empty">Ninguém entrou ainda.</p>}
      <ul className="people">
        {room.participants.map((p) => {
          const isAdmin = p.id === room.adminId;
          const isSpinner = !!room.spinners[p.id];
          return (
            <li key={p.id} className={p.online ? "" : "off"}>
              <img className="avatar" src={p.avatar} alt="" width={30} height={30} />
              <span className="nm">
                {p.name}
                {p.id === me && <span className="you"> (você)</span>}
              </span>
              <span className="tags">
                {isAdmin && <span className="tag tag-red">Admin</span>}
                <StatusTag participant={p} ready={!!room.ready[p.id]} now={now} />
              </span>
              {!isAdmin &&
                (perms.isAdmin ? (
                  <button
                    className={`star ${isSpinner ? "on" : ""}`}
                    title={isSpinner ? "Tirar a permissão de girar" : "Deixar girar a roleta"}
                    aria-label={isSpinner ? `Tirar a permissão de girar de ${p.name}` : `Deixar ${p.name} girar a roleta`}
                    aria-pressed={isSpinner}
                    onClick={() => dispatch(wsSend({ type: "set_spinner", userId: p.id, enabled: !isSpinner }))}
                  >
                    ★
                  </button>
                ) : (
                  isSpinner && (
                    <span className="star on" title="Pode girar a roleta">
                      ★
                    </span>
                  )
                ))}
            </li>
          );
        })}
      </ul>

      <div className="squad-foot">
        {perms.isParticipant ? (
          <>
            <button
              className={`btn px ${perms.isReady ? "btn-ghost" : "btn-cyan"}`}
              onClick={() => dispatch(wsSend({ type: "set_ready", enabled: !perms.isReady }))}
            >
              {perms.isReady ? "Cancelar pronto" : "✓ Pronto"}
            </button>
            {perms.canLeave && (
              <button className="btn btn-ghost" onClick={() => dispatch(wsSend({ type: "leave" }))}>
                Sair
              </button>
            )}
          </>
        ) : (
          <span className="muted">Você não está na roleta</span>
        )}
      </div>
    </aside>
  );
}

// Quem caiu mostra quanto tempo ainda tem para voltar sem perder o lugar.
function StatusTag({ participant, ready, now }: { participant: Participant; ready: boolean; now: number }) {
  if (!participant.online) {
    const left = participant.reconnectUntil ? Math.ceil((Date.parse(participant.reconnectUntil) - now) / 1000) : 0;
    return left > 0 ? <span className="tag tag-yellow">Voltando {left}s</span> : <span className="tag tag-muted">Offline</span>;
  }
  return ready ? <span className="tag tag-cyan">Pronto</span> : <span className="tag tag-muted">Esperando</span>;
}
