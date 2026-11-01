import { useAppDispatch, useAppSelector } from "../store/hooks";
import { wsSend } from "../store/roomSlice";
import { selectFormats, selectPermissions, selectRoom } from "../store/selectors";
import type { Format } from "../types/room";
import { plural, shortFormat } from "../utils/format";

// Controles ao lado da roda no lobby. Quem ainda não entrou vê só o convite para participar.
export function LobbyControls() {
  const dispatch = useAppDispatch();
  const room = useAppSelector(selectRoom);
  const formats = useAppSelector(selectFormats);
  const perms = useAppSelector(selectPermissions);
  if (!room) return null;

  const missing = perms.slots - room.participants.length;

  if (perms.canJoin) {
    return (
      <div className="ctrl">
        <div className="join-box">
          <div className="px join-k">Você está assistindo</div>
          <p className="join-d">
            Entre na roleta para concorrer a uma vaga.{" "}
            {missing > 0 ? (
              <>
                Ainda tem <b>{plural(missing, "vaga", "vagas")}</b> no {shortFormat(room.format)}.
              </>
            ) : (
              "As vagas já estão cheias, mas quem entrar também concorre."
            )}
          </p>
          <button className="btn btn-red btn-big px" onClick={() => dispatch(wsSend({ type: "join" }))}>
            Participar
          </button>
          <FormatReadOnly format={room.format} />
        </div>
      </div>
    );
  }

  const setFormat = (format: Format) => dispatch(wsSend({ type: "set_format", format }));
  // Trocar a quantidade de times mantém o tamanho quando ele existe no outro modo.
  const selectTeams = (teams: number) => {
    const options = formats.filter((f) => f.teams === teams);
    const next = options.find((f) => f.size === room.format.size) ?? options[0];
    if (next) setFormat(next);
  };
  const teamCounts = [...new Set(formats.map((f) => f.teams))].sort((a, b) => a - b);
  const sizes = formats.filter((f) => f.teams === room.format.teams);

  return (
    <div className="ctrl">
      {perms.isAdmin ? (
        <div className="fmt">
          <div className="fmt-row">
            <span className="lbl px">Times</span>
            <div className="seg">
              {teamCounts.map((teams) => (
                <button key={teams} className={teams === room.format.teams ? "on" : ""} onClick={() => selectTeams(teams)}>
                  {teams}
                </button>
              ))}
            </div>
          </div>
          <div className="fmt-row">
            <span className="lbl px">Por time</span>
            <div className="seg">
              {sizes.map((f) => (
                <button key={f.size} className={f.size === room.format.size ? "on" : ""} onClick={() => setFormat(f)}>
                  {f.size}
                </button>
              ))}
            </div>
          </div>
          <div className="fmt-total">
            <span className="px">{shortFormat(room.format)}</span>
            <span>{plural(perms.slots, "vaga", "vagas")}</span>
          </div>
        </div>
      ) : (
        <FormatReadOnly format={room.format} />
      )}

      {(perms.canSpin || room.autoSpin) && (
        <button
          className={`switch ${room.autoSpin ? "on" : ""}`}
          disabled={!perms.canSpin}
          aria-pressed={room.autoSpin}
          onClick={() => dispatch(wsSend({ type: "set_auto_spin", enabled: !room.autoSpin }))}
        >
          <span className="switch-box" />
          Giro auto
        </button>
      )}

      {perms.canSpin && (
        <button
          className="btn btn-red btn-big px"
          disabled={!!room.spin || !perms.enoughPlayers || !perms.allReady}
          onClick={() => dispatch(wsSend({ type: "spin" }))}
        >
          Girar
        </button>
      )}

      <p className="hint">
        {missing > 0 ? (
          <>
            {missing === 1 ? "Falta" : "Faltam"} <b>{plural(missing, "jogador", "jogadores")}</b> para o {shortFormat(room.format)} ·{" "}
            {room.participants.length}/{perms.slots} vagas
          </>
        ) : perms.notReady > 0 ? (
          <>
            {perms.notReady === 1 ? "Falta" : "Faltam"} <b>{plural(perms.notReady, "jogador", "jogadores")}</b>{" "}
            {perms.notReady === 1 ? "ficar pronto" : "ficarem prontos"}
          </>
        ) : perms.canSpin ? (
          "Tudo pronto, pode girar!"
        ) : (
          "Tudo pronto, esperando o admin girar"
        )}
      </p>
      {room.autoSpin && <p className="hint small">Depois do primeiro giro, a roleta segue girando sozinha</p>}
    </div>
  );
}

function FormatReadOnly({ format }: { format: Format }) {
  return (
    <div className="fmt-ro">
      <span className="lbl px">Formato</span>
      <span className="px fmt-v">{shortFormat(format)}</span>
      <span>definido pelo admin</span>
    </div>
  );
}
