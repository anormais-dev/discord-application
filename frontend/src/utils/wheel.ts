export const TAU = Math.PI * 2;
// O ponteiro fica no topo da roda.
export const POINTER_ANGLE = -Math.PI / 2;
const COLORS = ["#5865f2", "#eb459e", "#57f287", "#fee75c", "#ed4245", "#3ba55c", "#faa61a", "#9b59b6"];

export const easeOutQuart = (t: number) => 1 - (1 - t) ** 4;

export function drawWheel(canvas: HTMLCanvasElement, names: string[], rotation: number) {
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
    const mid = start + slice / 2;
    // Na metade esquerda o texto giraria de ponta-cabeça; vira meia volta e ancora pelo outro lado.
    const flipped = Math.cos(mid) < 0;
    ctx.rotate(flipped ? mid + Math.PI : mid);
    ctx.fillStyle = "#111214";
    ctx.font = `600 ${fontSize}px system-ui, sans-serif`;
    ctx.textAlign = flipped ? "left" : "right";
    ctx.textBaseline = "middle";
    ctx.fillText(truncate(ctx, name, radius - 40), flipped ? -(radius - 14) : radius - 14, 0);
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
