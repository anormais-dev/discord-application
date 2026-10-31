export const TAU = Math.PI * 2;
// O ponteiro fica no topo da roda.
export const POINTER_ANGLE = -Math.PI / 2;
// Fatias claras e escuras alternadas, com o texto na cor que contrasta.
const SLICES = [
  ["#ff4655", "#0a0c10"],
  ["#1f2430", "#ece8e1"],
  ["#19d3c5", "#0a0c10"],
  ["#2a1e2a", "#ece8e1"],
  ["#ffd447", "#0a0c10"],
  ["#1f2430", "#ece8e1"],
  ["#8b5cf6", "#ece8e1"],
  ["#162a2a", "#ece8e1"],
];
const DARK = "#0a0c10";
const LINE = "#2c3444";
const LIGHTS = ["#ff4655", "#ffd447"];
export const WHEEL_FONT = '"Chakra Petch", system-ui, sans-serif';
// Tamanho de referência do desenho; menor que isso, tudo encolhe junto.
const BASE_SIZE = 430;

export const easeOutQuart = (t: number) => 1 - (1 - t) ** 4;

export function drawWheel(canvas: HTMLCanvasElement, names: string[], rotation: number, emptyText = "Ninguém na roleta") {
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
  const scale = Math.min(1, size / BASE_SIZE);
  const ring = Math.max(8, 18 * scale);
  const radius = c - ring;
  // Na transição do PiP o Discord encolhe a janela a poucos pixels; com raio negativo o arc lança erro.
  if (radius <= 0) return;

  // Aro escuro com as luzes de fliperama.
  ctx.fillStyle = DARK;
  ctx.strokeStyle = LINE;
  ctx.lineWidth = Math.max(2, 4 * scale);
  ctx.beginPath();
  ctx.arc(c, c, c - ctx.lineWidth / 2, 0, TAU);
  ctx.fill();
  ctx.stroke();
  const light = Math.max(3, 6 * scale);
  for (let i = 0; i < 24; i++) {
    const a = (i / 24) * TAU;
    ctx.fillStyle = LIGHTS[i % 2];
    ctx.fillRect(c + (c - ring / 2) * Math.cos(a) - light / 2, c + (c - ring / 2) * Math.sin(a) - light / 2, light, light);
  }

  const hub = size * 0.2;
  if (names.length === 0) {
    ctx.fillStyle = "#1f2430";
    ctx.beginPath();
    ctx.arc(c, c, radius, 0, TAU);
    ctx.fill();
    ctx.fillStyle = "#7c8698";
    ctx.font = `700 ${Math.max(11, 15 * scale)}px ${WHEEL_FONT}`;
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.fillText(emptyText, c, c + hub + 18 * scale);
    drawHub(ctx, c, hub, scale);
    drawPointer(ctx, c, scale);
    return;
  }

  const slice = TAU / names.length;
  names.forEach((name, i) => {
    const start = rotation + i * slice;
    let [fill, text] = SLICES[i % SLICES.length];
    // Evita duas fatias vizinhas da mesma cor quando fecha a volta.
    if (i === names.length - 1 && i % SLICES.length === 0 && i > 0) [fill, text] = SLICES[1];
    ctx.fillStyle = fill;
    ctx.beginPath();
    // Com uma pessoa só a fatia é o círculo inteiro; partir do centro desenharia um raio solto.
    if (names.length > 1) ctx.moveTo(c, c);
    ctx.arc(c, c, radius, start, start + slice);
    ctx.closePath();
    ctx.fill();
    ctx.strokeStyle = DARK;
    ctx.lineWidth = Math.max(1.5, 3 * scale);
    ctx.stroke();

    ctx.save();
    ctx.translate(c, c);
    const mid = start + slice / 2;
    // Na metade esquerda o texto giraria de ponta-cabeça; vira meia volta e ancora pelo outro lado.
    const flipped = Math.cos(mid) < 0;
    ctx.rotate(flipped ? mid + Math.PI : mid);
    const label = name.toUpperCase();
    // Nomes longos ficam menores antes de serem cortados.
    const fontSize = Math.max(9, Math.min(15, 220 / names.length) * scale * Math.min(1, 7 / Math.max(7, label.length)));
    ctx.fillStyle = text;
    ctx.font = `700 ${fontSize}px ${WHEEL_FONT}`;
    ctx.textAlign = flipped ? "left" : "right";
    ctx.textBaseline = "middle";
    const edge = radius - 12 * scale;
    ctx.fillText(truncate(ctx, label, edge - hub - 8 * scale), flipped ? -edge : edge, 0);
    ctx.restore();
  });

  drawHub(ctx, c, hub, scale);
  drawPointer(ctx, c, scale);
}

// Miolo escuro com borda amarela, onde fica o macaco.
function drawHub(ctx: CanvasRenderingContext2D, c: number, hub: number, scale: number) {
  ctx.fillStyle = DARK;
  ctx.strokeStyle = "#ffd447";
  ctx.lineWidth = Math.max(2, 3 * scale);
  ctx.beginPath();
  ctx.arc(c, c, hub, 0, TAU);
  ctx.fill();
  ctx.stroke();
}

// Ponteiro em degraus, como pixel art.
function drawPointer(ctx: CanvasRenderingContext2D, c: number, scale: number) {
  const u = Math.max(2, 4 * scale);
  ctx.fillStyle = "#ece8e1";
  ctx.strokeStyle = DARK;
  ctx.lineWidth = Math.max(1, 2 * scale);
  ctx.beginPath();
  ctx.moveTo(c - 3.5 * u, 0);
  ctx.lineTo(c + 3.5 * u, 0);
  ctx.lineTo(c + 3.5 * u, 2 * u);
  ctx.lineTo(c + 2 * u, 2 * u);
  ctx.lineTo(c + 2 * u, 4 * u);
  ctx.lineTo(c + u, 4 * u);
  ctx.lineTo(c + u, 6 * u);
  ctx.lineTo(c - u, 6 * u);
  ctx.lineTo(c - u, 4 * u);
  ctx.lineTo(c - 2 * u, 4 * u);
  ctx.lineTo(c - 2 * u, 2 * u);
  ctx.lineTo(c - 3.5 * u, 2 * u);
  ctx.closePath();
  ctx.fill();
  ctx.stroke();
}

function truncate(ctx: CanvasRenderingContext2D, text: string, max: number) {
  if (ctx.measureText(text).width <= max) return text;
  let t = text;
  while (t.length > 1 && ctx.measureText(t + "…").width > max) t = t.slice(0, -1);
  return t + "…";
}
