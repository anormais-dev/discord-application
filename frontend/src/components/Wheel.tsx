import { useEffect, useRef, useState } from "react";
import { useAppSelector } from "../store/hooks";
import { selectParticipantsById, selectRoom } from "../store/selectors";
import { POINTER_ANGLE, TAU, drawWheel, easeOutQuart } from "../utils/wheel";

interface Animation {
  key: string;
  from: number;
  to: number;
  startedAt: number;
  durationMs: number;
}

export function Wheel() {
  const room = useAppSelector(selectRoom);
  const byId = useAppSelector(selectParticipantsById);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const rotationRef = useRef(0);
  const animationRef = useRef<Animation | null>(null);
  const [landed, setLanded] = useState<string | null>(null);

  const spin = room?.spin ?? null;
  const ids = spin ? spin.poolSnapshot : (room?.pool ?? []);
  const names = ids.map((id) => byId[id]?.name ?? "?");

  // Começa uma animação nova quando chega um giro que ainda não foi animado.
  useEffect(() => {
    if (!spin) {
      animationRef.current = null;
      setLanded(null);
      return;
    }
    const key = spin.startedAt;
    if (animationRef.current?.key === key) return;

    const slice = TAU / spin.poolSnapshot.length;
    const winnerIndex = spin.poolSnapshot.indexOf(spin.winnerId);
    const target = POINTER_ANGLE - (winnerIndex + spin.targetOffset) * slice;
    const from = rotationRef.current;
    // Menor ângulo equivalente ao alvo que dá pelo menos `rotations` voltas completas.
    const minimum = from + spin.rotations * TAU;
    const to = target + Math.ceil((minimum - target) / TAU) * TAU;

    animationRef.current = { key, from, to, startedAt: performance.now(), durationMs: spin.durationMs };
    setLanded(null);
  }, [spin]);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    let frame = 0;

    const render = () => {
      const anim = animationRef.current;
      if (anim) {
        const t = Math.min(1, (performance.now() - anim.startedAt) / anim.durationMs);
        rotationRef.current = anim.from + (anim.to - anim.from) * easeOutQuart(t);
        if (t >= 1 && spin) setLanded(spin.winnerId);
      }
      drawWheel(canvas, names, rotationRef.current);
      if (anim && performance.now() - anim.startedAt < anim.durationMs) {
        frame = requestAnimationFrame(render);
      }
    };
    render();
    return () => cancelAnimationFrame(frame);
  }, [names.join("\u0000"), spin]);

  return (
    <div className="wheel">
      <canvas ref={canvasRef} className="wheel-canvas" />
      <p className="wheel-result" aria-live="polite">
        {landed ? `Sorteado: ${byId[landed]?.name ?? "?"}` : spin ? "Girando..." : " "}
      </p>
    </div>
  );
}
