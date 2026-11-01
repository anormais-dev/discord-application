import { useEffect, useRef, useState } from "react";
import type { SpinState } from "../types/room";
import { playConfetti, playTick } from "../utils/sound";
import { POINTER_ANGLE, TAU, WHEEL_FONT, drawWheel, easeOutQuart } from "../utils/wheel";
import { Confetti } from "./Confetti";

const SPIRAL_START = 0.15;
const SPIRAL_TURNS = 4;

interface Animation {
  key: string;
  winnerId: string;
  from: number;
  to: number;
  startedAt: number;
  durationMs: number;
}

interface WheelProps {
  spin: SpinState | null;
  ids: string[];
  label: (id: string) => string;
  empty: boolean;
  celebrate: boolean;
  resultPrefix: string;
  // Resultado já decidido antes de a página abrir, mostrado sem festa.
  settledId?: string;
  emptyText?: string;
}

export function Wheel({ spin, ids, label, empty, celebrate, resultPrefix, settledId, emptyText }: WheelProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const monkeyRef = useRef<HTMLDivElement>(null);
  const rotationRef = useRef(0);
  const animationRef = useRef<Animation | null>(null);
  const dizzyRef = useRef(false);
  const [dizzy, setDizzy] = useState(false);
  const [landed, setLanded] = useState<string | null>(null);
  // Largura do canvas; muda quando o Discord redimensiona a janela e força redesenhar.
  const [size, setSize] = useState(0);
  // O canvas não espera a fonte carregar; redesenha quando ela chega.
  const [fontReady, setFontReady] = useState(false);
  // Fatia que estava no ponteiro no quadro anterior: quando muda, toca o "tic".
  const sliceUnderPointerRef = useRef<number | null>(null);
  // Giro que já tocou o som de confete, para não repetir ao redesenhar.
  const celebratedRef = useRef<string | null>(null);
  const celebrateRef = useRef(celebrate);
  celebrateRef.current = celebrate;

  const celebrationActive = !!landed && celebrate;
  const resultId = celebrationActive ? landed : spin ? null : settledId;
  const names = ids.map(label);

  // Começa uma animação nova quando chega um giro que ainda não foi animado.
  useEffect(() => {
    if (!spin) {
      const finishedAnimation = animationRef.current;
      if (finishedAnimation) {
        rotationRef.current = finishedAnimation.to;
        setLanded(finishedAnimation.winnerId);
      }
      animationRef.current = null;
      dizzyRef.current = false;
      setDizzy(false);
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

    animationRef.current = { key, winnerId: spin.winnerId, from, to, startedAt: performance.now(), durationMs: spin.durationMs };
    setLanded(null);
  }, [spin]);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const observer = new ResizeObserver(() => setSize(canvas.clientWidth));
    observer.observe(canvas);
    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    let alive = true;
    document.fonts
      ?.load(`700 16px ${WHEEL_FONT}`)
      .then(() => alive && setFontReady(true))
      .catch(() => {});
    return () => {
      alive = false;
    };
  }, []);

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
        if (progress < 1 && names.length > 1) {
          const slice = TAU / names.length;
          const atPointer = Math.floor(((((POINTER_ANGLE - rotationRef.current) % TAU) + TAU) % TAU) / slice);
          if (sliceUnderPointerRef.current !== null && atPointer !== sliceUnderPointerRef.current) playTick();
          sliceUnderPointerRef.current = atPointer;
        }
        if (progress >= 1 && spin) {
          setLanded(spin.winnerId);
          if (celebratedRef.current !== anim.key) {
            celebratedRef.current = anim.key;
            sliceUnderPointerRef.current = null;
            if (celebrateRef.current) playConfetti();
          }
        }
      }

      const shouldBeDizzy = progress >= SPIRAL_START && progress < 1;
      if (shouldBeDizzy !== dizzyRef.current) {
        dizzyRef.current = shouldBeDizzy;
        setDizzy(shouldBeDizzy);
      }

      const angle = progress >= 1 ? POINTER_ANGLE : POINTER_ANGLE + rotationRef.current;
      const pupilX = Math.cos(angle) * 6;
      const pupilY = Math.sin(angle) * 6;
      const spiralProgress = Math.max(0, Math.min(1, (progress - SPIRAL_START) / (1 - SPIRAL_START)));
      const spiralRotation = spiralProgress * TAU * SPIRAL_TURNS;

      if (monkeyRef.current) {
        monkeyRef.current.style.setProperty("--pupil-x", `${pupilX}px`);
        monkeyRef.current.style.setProperty("--pupil-y", `${pupilY}px`);
        monkeyRef.current.style.setProperty("--spiral-rotation", `${spiralRotation}rad`);
      }

      drawWheel(canvas, names, rotationRef.current, emptyText);
      if (anim && performance.now() - anim.startedAt < anim.durationMs) {
        frame = requestAnimationFrame(render);
      }
    };
    render();
    return () => cancelAnimationFrame(frame);
  }, [names.join("\u0000"), spin, emptyText, size, fontReady]);

  return (
    <div className="wheel">
      <div className="wheel-stage">
        <canvas ref={canvasRef} className="wheel-canvas" />
        <div
          ref={monkeyRef}
          className={`wheel-monkey${empty ? " is-empty" : ""}${dizzy ? " is-dizzy" : ""}${celebrationActive ? " is-winner" : ""}`}
          aria-hidden="true"
        >
          <img className="monkey-head monkey-head-normal" src="/images/monkey-head.png" alt="" />
          <img className="monkey-head monkey-head-stand" src="/images/stand-monkey-head.png" alt="" />
          <img className="monkey-head monkey-head-dizzy" src="/images/dizzy-monkey-head.png" alt="" />
          <img className="monkey-eyes" src="/images/monkey-eye.png" alt="" />
          <img className="monkey-spiral-eye spiral-eye-left" src="/images/espiral-eye.png" alt="" />
          <img className="monkey-spiral-eye spiral-eye-right" src="/images/espiral-eye.png" alt="" />
          <img className="monkey-head monkey-head-happy" src="/images/happy-monkey-head.png" alt="" />
          <div className="celebration-cone-group">
            <img className="celebration-cone" src="/images/cone.png" alt="" />
            <img className="celebration-confetti" src="/images/confetti.png" alt="" />
          </div>
        </div>
        {celebrationActive && <Confetti />}
      </div>
      <p className="wheel-result" aria-live="polite">
        {resultId ? `${resultPrefix}${label(resultId)}` : spin ? "Girando..." : " "}
      </p>
    </div>
  );
}
