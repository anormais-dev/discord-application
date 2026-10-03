import { useGetMapsQuery } from "../services/api";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { wsSend } from "../store/roomSlice";
import { selectPermissions, selectRoom } from "../store/selectors";
import { Wheel } from "./Wheel";

export function MapPanel() {
  const dispatch = useAppDispatch();
  const room = useAppSelector(selectRoom);
  const perms = useAppSelector(selectPermissions);
  const { data: maps = [], isError } = useGetMapsQuery();
  if (!room) return null;

  const spin = room.mapSpin;
  const ids = spin ? spin.poolSnapshot : maps;

  return (
    <section className="panel map-panel">
      <h2>Mapa</h2>
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
      <div className="controls">
        {perms.isAdmin ? (
          <button
            className="primary"
            disabled={!!spin || maps.length === 0}
            onClick={() => dispatch(wsSend({ type: "spin_map" }))}
          >
            {spin ? "Girando..." : "Sortear mapa"}
          </button>
        ) : (
          <p className="muted">Só o admin sorteia o mapa</p>
        )}
      </div>
    </section>
  );
}
