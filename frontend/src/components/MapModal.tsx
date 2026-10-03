import { useEffect, useLayoutEffect, useMemo, useRef, useState } from "react";
import { useGetMapsQuery } from "../services/api";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { wsSend } from "../store/roomSlice";
import { selectPermissions, selectRoom } from "../store/selectors";
import type { SpinState } from "../types/room";

function mapImagePath(map: string) {
  const filename = map
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "");
  return `/images/maps/${filename}.jpg`;
}

function MapCard({
  map,
  featured = false,
}: {
  map: string;
  featured?: boolean;
}) {
  const [imageMissing, setImageMissing] = useState(false);

  return (
    <article className={`map-card${featured ? " map-card-featured" : ""}`}>
      <div className="map-card-image">
        {!imageMissing ? (
          <img src={mapImagePath(map)} alt="" onError={() => setImageMissing(true)} />
        ) : (
          <span className="map-image-missing">Imagem ausente: {mapImagePath(map)}</span>
        )}
      </div>
      <strong>{map}</strong>
    </article>
  );
}

function makeCarouselItems(spin: SpinState, maps: string[]) {
  const uniqueMaps = [...new Set(maps)];
  if (!uniqueMaps.includes(spin.winnerId)) uniqueMaps.push(spin.winnerId);
  if (uniqueMaps.length === 1) return { items: uniqueMaps, targetIndex: 0 };

  const targetCycle = spin.rotations;
  const trailingCycles = 2;
  const itemCount = (targetCycle + trailingCycles + 1) * uniqueMaps.length;
  const offset = Math.max(0, Math.min(0.999999, spin.targetOffset));
  const targetIndex = targetCycle * uniqueMaps.length + Math.floor(offset * uniqueMaps.length);
  const winnerIndex = uniqueMaps.indexOf(spin.winnerId);
  const startIndex = (winnerIndex - (targetIndex % uniqueMaps.length) + uniqueMaps.length) % uniqueMaps.length;
  const items = Array.from(
    { length: itemCount },
    (_, index) => uniqueMaps[(index + startIndex) % uniqueMaps.length],
  );

  return { items, targetIndex };
}

export function MapModal() {
  const dispatch = useAppDispatch();
  const room = useAppSelector(selectRoom);
  const perms = useAppSelector(selectPermissions);
  const { data: maps = [], isError } = useGetMapsQuery();
  const viewportRef = useRef<HTMLDivElement>(null);
  const activeSpinRef = useRef<string | null>(null);
  const [offset, setOffset] = useState(0);
  const [isRolling, setIsRolling] = useState(false);
  const [animationDuration, setAnimationDuration] = useState(0);
  const [revealedMap, setRevealedMap] = useState<string | null>(null);

  const spin = room?.mapSpin ?? null;
  const ids = spin ? spin.poolSnapshot : maps;
  const carousel = useMemo(
    () => (spin ? makeCarouselItems(spin, ids) : { items: ids, targetIndex: Math.floor(ids.length / 2) }),
    [spin?.startedAt, ids],
  );
  const carouselItems = carousel.items;
  const spinKey = spin?.startedAt;

  useEffect(() => {
    if (!room?.mapOpen) {
      activeSpinRef.current = null;
      setRevealedMap(null);
      return;
    }
    if (spin) {
      activeSpinRef.current = spin.startedAt;
      setRevealedMap(null);
      return;
    }
    if (activeSpinRef.current) {
      setRevealedMap(room.map || null);
      activeSpinRef.current = null;
    }
  }, [room?.mapOpen, room?.map, spinKey]);

  useLayoutEffect(() => {
    const viewport = viewportRef.current;
    const track = viewport?.querySelector<HTMLElement>(".map-track");
    const cards = track?.querySelectorAll<HTMLElement>(".map-card");
    if (!viewport || !cards?.length) return;

    const alignOffset = (card: HTMLElement, cardPosition = 0.5) =>
      viewport.clientWidth / 2 - (card.offsetLeft + card.offsetWidth * cardPosition);

    setIsRolling(false);
    if (!spin) {
      setAnimationDuration(0);
      setOffset(alignOffset(cards[carousel.targetIndex]));
      return;
    }

    const targetPosition = Math.max(0.15, Math.min(0.85, spin.targetOffset));
    const elapsed = Math.max(0, Date.now() - new Date(spin.startedAt).getTime());
    const duration = Math.max(0, spin.durationMs - elapsed - 32);
    setAnimationDuration(duration);
    setOffset(alignOffset(cards[0]));
    let secondFrame = 0;
    const frame = requestAnimationFrame(() => {
      secondFrame = requestAnimationFrame(() => {
        setIsRolling(true);
        setOffset(alignOffset(cards[carousel.targetIndex], targetPosition));
      });
    });
    return () => {
      cancelAnimationFrame(frame);
      cancelAnimationFrame(secondFrame);
    };
  }, [spinKey, carouselItems.length, carousel.targetIndex, !!spin, room?.mapOpen]);

  if (!room?.mapOpen) return null;

  const hasResult = !spin && !!revealedMap;

  return (
    <div className="modal-backdrop">
      <section className="panel modal map-modal" role="dialog" aria-modal="true" aria-labelledby="map-modal-title">
        <h2 id="map-modal-title">Sorteio do mapa</h2>
        {hasResult ? (
          <>
            <MapCard map={revealedMap} featured />
            <p className="map-status" aria-live="polite">Mapa escolhido</p>
          </>
        ) : ids.length > 0 ? (
          <>
            <div className={`map-carousel${spin ? " is-spinning" : " is-idle"}`} ref={viewportRef}>
              <div
                className={`map-track${isRolling ? " is-rolling" : ""}`}
                style={{
                  transform: `translateX(${offset}px)`,
                  transitionDuration: `${animationDuration}ms`,
                }}
              >
                {carouselItems.map((map, index) => (
                  <MapCard key={`${map}-${index}`} map={map} />
                ))}
              </div>
              <span className="map-selection-line" aria-hidden="true" />
            </div>
            <p className="map-status" aria-live="polite">
              {spin ? "Sorteando mapa..." : "Pronto para sortear"}
            </p>
          </>
        ) : (
          <p className="map-status muted">{isError ? "Mapas indisponíveis" : "Carregando mapas..."}</p>
        )}
        <div className="actions">
          {perms.isAdmin ? (
            <>
              {!hasResult && (
                <button
                  className="primary"
                  disabled={!!spin || maps.length === 0}
                  onClick={() => dispatch(wsSend({ type: "spin_map" }))}
                >
                  {spin ? "Sorteando..." : "Sortear mapa"}
                </button>
              )}
              <button
                className="secondary"
                disabled={!!spin}
                onClick={() => dispatch(wsSend({ type: "set_map_open", enabled: false }))}
              >
                Fechar
              </button>
              {hasResult && (
                <button className="primary" onClick={() => dispatch(wsSend({ type: "spin_map" }))}>
                  Girar de novo
                </button>
              )}
            </>
          ) : (
            <p className="muted">Só o admin pode sortear ou fechar</p>
          )}
        </div>
      </section>
    </div>
  );
}
