import { useEffect, useRef, useState } from "react";
import { useAppSelector } from "../app/hooks";
import { selectParticipantsById, selectRoom } from "../features/room/roomSlice";

const TAU = Math.PI * 2;
// O ponteiro fica no topo da roda.
const POINTER_ANGLE = -Math.PI / 2;
const COLORS = ["#5865f2", "#eb459e", "#57f287", "#fee75c", "#ed4245", "#3ba55c", "#faa61a", "#9b59b6"];

const easeOutQuart = (t: number) => 1 - (1 - t) ** 4;

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
      draw(canvas, names, rotationRef.current);
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

function draw(canvas: HTMLCanvasElement, names: string[], rotation: number) {
  const dpr = window.devicePixelRatio || 1;
  const size = canvas.clientWidth;
  if (canvas.width !== size * dpr) {
    canvas.width = size * dpr;
    canvas.height = size * dpr;
  }
  const ctx = canvas.getContext("2d")!;
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  ctx.clearRect(0, 0, size, size);

  const c = size / 2;
  const radius = c - 12;

  if (names.length === 0) {
    ctx.fillStyle = "#2b2d31";
    ctx.beginPath();
    ctx.arc(c, c, radius, 0, TAU);
    ctx.fill();
    ctx.fillStyle = "#949ba4";
    ctx.font = "16px system-ui, sans-serif";
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.fillText("Ninguém na roleta", c, c);
    drawPointer(ctx, c);
    return;
  }

  const slice = TAU / names.length;
  const fontSize = Math.max(11, Math.min(18, 220 / names.length));
  names.forEach((name, i) => {
    const start = rotation + i * slice;
    ctx.fillStyle = COLORS[i % COLORS.length];
    // Evita duas fatias vizinhas da mesma cor quando fecha a volta.
    if (i === names.length - 1 && i % COLORS.length === 0 && i > 0) ctx.fillStyle = COLORS[1];
    ctx.beginPath();
    ctx.moveTo(c, c);
    ctx.arc(c, c, radius, start, start + slice);
    ctx.closePath();
    ctx.fill();
    ctx.strokeStyle = "#1e1f22";
    ctx.lineWidth = 2;
    ctx.stroke();

    ctx.save();
    ctx.translate(c, c);
    ctx.rotate(start + slice / 2);
    ctx.fillStyle = "#111214";
    ctx.font = `600 ${fontSize}px system-ui, sans-serif`;
    ctx.textAlign = "right";
    ctx.textBaseline = "middle";
    ctx.fillText(truncate(ctx, name, radius - 40), radius - 14, 0);
    ctx.restore();
  });

  ctx.fillStyle = "#1e1f22";
  ctx.beginPath();
  ctx.arc(c, c, 22, 0, TAU);
  ctx.fill();
  drawPointer(ctx, c);
}

function drawPointer(ctx: CanvasRenderingContext2D, c: number) {
  ctx.fillStyle = "#f2f3f5";
  ctx.beginPath();
  ctx.moveTo(c - 12, 2);
  ctx.lineTo(c + 12, 2);
  ctx.lineTo(c, 30);
  ctx.closePath();
  ctx.fill();
}

function truncate(ctx: CanvasRenderingContext2D, text: string, max: number) {
  if (ctx.measureText(text).width <= max) return text;
  let t = text;
  while (t.length > 1 && ctx.measureText(t + "…").width > max) t = t.slice(0, -1);
  return t + "…";
}
