import { useGetMapsQuery } from "../services/api";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { wsSend } from "../store/roomSlice";
import { selectPermissions, selectRoom } from "../store/selectors";
import { Wheel } from "./Wheel";

export function MapModal() {
  const dispatch = useAppDispatch();
  const room = useAppSelector(selectRoom);
  const perms = useAppSelector(selectPermissions);
  const { data: maps = [], isError } = useGetMapsQuery();
  if (!room?.mapOpen) return null;

  const spin = room.mapSpin;
  const ids = spin ? spin.poolSnapshot : maps;

  return (
    <div className="modal-backdrop">
      <section className="panel modal" role="dialog" aria-modal="true" aria-labelledby="map-modal-title">
        <h2 id="map-modal-title">Mapa do Valorant</h2>
        <Wheel
          spin={spin}
          ids={ids}
          label={(id) => id}
          empty={ids.length === 0}
          celebrate={!!room.map || !!spin}
          resultPrefix="Mapa: "
          settledId={room.map || undefined}
          emptyText={isError ? "Mapas indisponíveis" : "Carregando mapas..."}
        />
        <div className="actions">
          {perms.isAdmin ? (
            <>
              <button
                className="primary"
                disabled={!!spin || maps.length === 0}
                onClick={() => dispatch(wsSend({ type: "spin_map" }))}
              >
                {spin ? "Girando..." : "Girar"}
              </button>
              <button
                className="secondary"
                disabled={!!spin}
                onClick={() => dispatch(wsSend({ type: "set_map_open", enabled: false }))}
              >
                Fechar
              </button>
            </>
          ) : (
            <p className="muted">Só o admin gira o mapa</p>
          )}
        </div>
      </section>
    </div>
  );
}
