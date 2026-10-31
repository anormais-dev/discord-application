// Sons 8-bit gerados na hora com a Web Audio API: sem arquivos de áudio para baixar ou licenciar.
// O navegador só libera áudio depois de um clique na página; antes disso os sons são ignorados.

let ctx: AudioContext | null = null;
let muted = false;
let lastTick = 0;

function audio() {
  if (muted) return null;
  if (!ctx) {
    try {
      ctx = new AudioContext();
    } catch {
      return null;
    }
  }
  if (ctx.state === "suspended") void ctx.resume();
  return ctx.state === "running" ? ctx : null;
}

// Destrava o áudio no primeiro gesto do usuário (regra de autoplay dos navegadores).
export function unlockAudioOnFirstGesture() {
  const unlock = () => {
    if (!ctx) {
      try {
        ctx = new AudioContext();
      } catch {
        return;
      }
    }
    void ctx.resume();
  };
  document.addEventListener("pointerdown", unlock, { once: true, capture: true });
  document.addEventListener("keydown", unlock, { once: true, capture: true });
}

export function setMuted(value: boolean) {
  muted = value;
}

type Wave = OscillatorType;

// Uma nota com envelope curto (ataque rápido, queda exponencial), como num chip de videogame.
function note(ac: AudioContext, freq: number, start: number, duration: number, { type = "square" as Wave, volume = 0.08, slideTo = 0 } = {}) {
  const osc = ac.createOscillator();
  const gain = ac.createGain();
  osc.type = type;
  osc.frequency.setValueAtTime(freq, start);
  if (slideTo) osc.frequency.exponentialRampToValueAtTime(slideTo, start + duration);
  gain.gain.setValueAtTime(0.0001, start);
  gain.gain.exponentialRampToValueAtTime(volume, start + 0.008);
  gain.gain.exponentialRampToValueAtTime(0.0001, start + duration);
  osc.connect(gain).connect(ac.destination);
  osc.start(start);
  osc.stop(start + duration + 0.02);
}

// Estalo de ruído curto, para o "pop" do confete.
function pop(ac: AudioContext, start: number, volume = 0.12) {
  const length = Math.floor(ac.sampleRate * 0.06);
  const buffer = ac.createBuffer(1, length, ac.sampleRate);
  const data = buffer.getChannelData(0);
  for (let i = 0; i < length; i++) data[i] = (Math.random() * 2 - 1) * (1 - i / length) ** 2;
  const src = ac.createBufferSource();
  const gain = ac.createGain();
  src.buffer = buffer;
  gain.gain.value = volume;
  src.connect(gain).connect(ac.destination);
  src.start(start);
}

// "Tic" de uma fatia passando pelo ponteiro. Limitado para não virar um zumbido no começo do giro.
export function playTick() {
  const ac = audio();
  if (!ac) return;
  const now = ac.currentTime;
  if (now - lastTick < 0.035) return;
  lastTick = now;
  note(ac, 1400 + Math.random() * 300, now, 0.03, { volume: 0.04 });
}

// Pop e arpejo subindo quando a roda para no sorteado.
export function playConfetti() {
  const ac = audio();
  if (!ac) return;
  const t = ac.currentTime;
  pop(ac, t);
  [523, 659, 784, 1047].forEach((f, i) => note(ac, f, t + 0.05 + i * 0.06, 0.12, { volume: 0.06 }));
}

// Fanfarrinha de quem entrou num time.
export function playFanfare() {
  const ac = audio();
  if (!ac) return;
  const t = ac.currentTime;
  const melody: [number, number, number][] = [
    [523, 0, 0.12],
    [523, 0.13, 0.12],
    [523, 0.26, 0.12],
    [659, 0.39, 0.3],
    [587, 0.7, 0.12],
    [659, 0.83, 0.12],
    [784, 0.96, 0.45],
  ];
  for (const [f, at, d] of melody) {
    note(ac, f, t + at, d, { volume: 0.07 });
    note(ac, f / 2, t + at, d, { type: "triangle", volume: 0.06 });
  }
}

// Trombone triste de quem ficou de fora: três notas descendo e a última escorregando.
export function playSadTrombone() {
  const ac = audio();
  if (!ac) return;
  const t = ac.currentTime;
  const notes: [number, number, number, number][] = [
    [392, 0, 0.32, 0],
    [370, 0.36, 0.32, 0],
    [349, 0.72, 0.32, 0],
    [330, 1.08, 0.9, 262],
  ];
  for (const [f, at, d, slide] of notes) note(ac, f, t + at, d, { type: "sawtooth", volume: 0.05, slideTo: slide });
}
