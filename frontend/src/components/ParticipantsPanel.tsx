import { useAppDispatch, useAppSelector } from "../store/hooks";
import { wsSend } from "../store/roomSlice";
import { selectMe, selectPermissions, selectRoom } from "../store/selectors";

export function ParticipantsPanel() {
  const dispatch = useAppDispatch();
  const room = useAppSelector(selectRoom);
  const me = useAppSelector(selectMe);
  const perms = useAppSelector(selectPermissions);
  if (!room) return null;

  return (
    <section className="panel">
      <h2>
        Participantes <span className="muted">{room.participants.length}/{perms.slots}</span>
      </h2>
      {room.participants.length === 0 && <p className="muted">Clique em Participar para entrar na roleta.</p>}
      <ul className="people">
        {room.participants.map((p) => {
          const isAdmin = p.id === room.adminId;
          const isSpinner = !!room.spinners[p.id];
          return (
            <li key={p.id} className={p.online ? "" : "offline"}>
              <img src={p.avatar} alt="" width={28} height={28} />
              <span className="name">
                {p.name}
                {p.id === me && <span className="muted"> (você)</span>}
              </span>
              {isAdmin && <span className="badge admin">Admin</span>}
              {!isAdmin && isSpinner && <span className="badge">Gira</span>}
              {!p.online && <span className="badge">Offline</span>}
              {perms.isAdmin && !isAdmin && (
                <button
                  className="link"
                  onClick={() => dispatch(wsSend({ type: "set_spinner", userId: p.id, enabled: !isSpinner }))}
                >
                  {isSpinner ? "Remover giro" : "Pode girar"}
                </button>
              )}
            </li>
          );
        })}
      </ul>
    </section>
  );
}
