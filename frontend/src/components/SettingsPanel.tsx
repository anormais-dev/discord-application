import { useAppDispatch } from "../store/hooks";
import { wsSend } from "../store/roomSlice";
import { ValorantIcon } from "./Icons";

// A roleta de mapas abre em qualquer fase e formato, com ou sem times sorteados.
export function SettingsPanel({ onClose }: { onClose: () => void }) {
  const dispatch = useAppDispatch();

  const openMaps = () => {
    dispatch(wsSend({ type: "set_map_open", enabled: true }));
    onClose();
  };

  return (
    <section className="panel settings" aria-label="Configurações">
      <div className="panel-h">
        <span className="px">Configurações</span>
        <span className="tag tag-red">Admin</span>
      </div>
      <div className="opt">
        <div className="opt-icon">
          <ValorantIcon size={30} />
        </div>
        <div className="opt-txt">
          <div className="opt-name">Roleta de mapas</div>
          <div className="opt-desc">Sorteia um mapa competitivo do Valorant. Abre a qualquer momento, com ou sem times sorteados.</div>
        </div>
        <button className="btn btn-red px btn-small" onClick={openMaps}>
          Abrir
        </button>
      </div>
    </section>
  );
}
