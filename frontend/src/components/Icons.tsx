// Engrenagem em pixel art (grade 16x16): dentes nos quatro lados e nas diagonais, furo no meio.
const GEAR = [
  "......####......",
  "..##..####..##..",
  "..############..",
  "...##########...",
  "...##########...",
  "..############..",
  "######....######",
  "######....######",
  "######....######",
  "######....######",
  "..############..",
  "...##########...",
  "...##########...",
  "..############..",
  "..##..####..##..",
  "......####......",
];
const GEAR_PATH = GEAR.flatMap((row, y) => [...row].map((v, x) => (v === "#" ? `M${x} ${y}h1v1h-1z` : ""))).join("");

export function GearIcon({ size = 20 }: { size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 16 16" aria-hidden="true" shapeRendering="crispEdges">
      <path d={GEAR_PATH} fill="currentColor" />
    </svg>
  );
}

// Símbolo do Valorant (Simple Icons 13.21.0), com o viewBox recortado no desenho para centralizar.
const VALORANT_PATH =
  "M23.792 2.152a.252.252 0 0 0-.098.083c-3.384 4.23-6.769 8.46-10.15 12.69-.107.093-.025.288.119.265 2.439.003 4.877 0 7.316.001a.66.66 0 0 0 .552-.25c.774-.967 1.55-1.934 2.324-2.903a.72.72 0 0 0 .144-.49c-.002-3.077 0-6.153-.003-9.23.016-.11-.1-.206-.204-.167zM.077 2.166c-.077.038-.074.132-.076.205.002 3.074.001 6.15.001 9.225a.679.679 0 0 0 .158.463l7.64 9.55c.12.152.308.25.505.247 2.455 0 4.91.003 7.365 0 .142.02.222-.174.116-.265C10.661 15.176 5.526 8.766.4 2.35c-.08-.094-.174-.272-.322-.184z";

export function ValorantIcon({ size = 24 }: { size?: number }) {
  return (
    <svg width={size} height={(size * 20) / 24} viewBox="0 2 24 20" aria-hidden="true">
      <path d={VALORANT_PATH} fill="var(--red)" />
    </svg>
  );
}

// Alto-falante em pixel art: com ondas quando o som está ligado, com um X quando está mudo.
const SPEAKER = [
  "................",
  "................",
  "......#.........",
  ".....##.........",
  "....###.........",
  "#######.........",
  "#######.........",
  "#######.........",
  "#######.........",
  "#######.........",
  "#######.........",
  "....###.........",
  ".....##.........",
  "......#.........",
  "................",
  "................",
];
const WAVES = [
  "................",
  "................",
  "............#...",
  ".............#..",
  ".............#..",
  ".........#....#.",
  "..........#...#.",
  "..........#...#.",
  "..........#...#.",
  "..........#...#.",
  ".........#....#.",
  ".............#..",
  ".............#..",
  "............#...",
  "................",
  "................",
];
const CROSS = [
  "................",
  "................",
  "................",
  "................",
  "................",
  ".........#....#.",
  "..........#..#..",
  "...........##...",
  "...........##...",
  "..........#..#..",
  ".........#....#.",
  "................",
  "................",
  "................",
  "................",
  "................",
];
const toPath = (grid: string[]) => grid.flatMap((row, y) => [...row].map((v, x) => (v === "#" ? `M${x} ${y}h1v1h-1z` : ""))).join("");
const SPEAKER_ON = toPath(SPEAKER) + toPath(WAVES);
const SPEAKER_OFF = toPath(SPEAKER) + toPath(CROSS);

export function SoundIcon({ muted, size = 20 }: { muted: boolean; size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 16 16" aria-hidden="true" shapeRendering="crispEdges">
      <path d={muted ? SPEAKER_OFF : SPEAKER_ON} fill="currentColor" />
    </svg>
  );
}
