import { useEffect, useRef, useState } from "react";
import { useAppSelector } from "../store/hooks";
import { selectParticipantsById, selectRoom } from "../store/selectors";
import { POINTER_ANGLE, TAU, drawWheel, easeOutQuart } from "../utils/wheel";

const SPIRAL_START = 0.15;
const SPIRAL_TURNS = 4;

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
  const monkeyRef = useRef<HTMLDivElement>(null);
  const rotationRef = useRef(0);
  const animationRef = useRef<Animation | null>(null);
  const dizzyRef = useRef(false);
  const [dizzy, setDizzy] = useState(false);
  const [landed, setLanded] = useState<string | null>(null);

  const spin = room?.spin ?? null;
  const ids = spin ? spin.poolSnapshot : (room?.pool ?? []);
  const names = ids.map((id) => byId[id]?.name ?? "?");

  // Começa uma animação nova quando chega um giro que ainda não foi animado.
  useEffect(() => {
    if (!spin) {
      animationRef.current = null;
      dizzyRef.current = false;
      setDizzy(false);
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
      let progress = 1;
      if (anim) {
        progress = Math.min(1, (performance.now() - anim.startedAt) / anim.durationMs);
        rotationRef.current = anim.from + (anim.to - anim.from) * easeOutQuart(progress);
        if (progress >= 1 && spin) setLanded(spin.winnerId);
      }

      const shouldBeDizzy = progress >= SPIRAL_START && progress < 1;
      if (shouldBeDizzy !== dizzyRef.current) {
        dizzyRef.current = shouldBeDizzy;
        setDizzy(shouldBeDizzy);
      }

      const angle = progress >= 1 ? POINTER_ANGLE : POINTER_ANGLE - rotationRef.current;
      const pupilX = Math.cos(angle) * 6;
      const pupilY = Math.sin(angle) * 6;
      const spiralProgress = Math.max(0, Math.min(1, (progress - SPIRAL_START) / (1 - SPIRAL_START)));
      const spiralRotation = spiralProgress * TAU * SPIRAL_TURNS;

      if (monkeyRef.current) {
        monkeyRef.current.style.setProperty("--pupil-x", `${pupilX}px`);
        monkeyRef.current.style.setProperty("--pupil-y", `${pupilY}px`);
        monkeyRef.current.style.setProperty("--spiral-rotation", `${spiralRotation}rad`);
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
      <div className="wheel-stage">
        <canvas ref={canvasRef} className="wheel-canvas" />
        <div ref={monkeyRef} className={`wheel-monkey${dizzy ? " is-dizzy" : ""}`} aria-hidden="true">
          <img className="monkey-head monkey-head-normal" src="/images/monkey-head.png" alt="" />
          <img className="monkey-head monkey-head-dizzy" src="/images/dizzy-monkey-head.png" alt="" />
          <img className="monkey-eyes" src="/images/monkey-eye.png" alt="" />
          <img className="monkey-spiral-eye spiral-eye-left" src="/images/espiral-eye.png" alt="" />
          <img className="monkey-spiral-eye spiral-eye-right" src="/images/espiral-eye.png" alt="" />
        </div>
      </div>
      <p className="wheel-result" aria-live="polite">
        {landed ? `Sorteado: ${byId[landed]?.name ?? "?"}` : spin ? "Girando..." : " "}
      </p>
    </div>
  );
}
