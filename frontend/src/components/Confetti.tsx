// Chuva de confete em pixel: peças coloridas caindo por cima de quem a usa (posição relativa).
const PIECES = Array.from({ length: 40 }, (_, index) => ({
  left: `${(index * 37) % 101}%`,
  delay: `-${(index % 8) * 0.13}s`,
  duration: `${1.2 + (index % 5) * 0.12}s`,
  color: index % 5,
}));

export function Confetti({ pieces = 24, className = "" }: { pieces?: number; className?: string }) {
  return (
    <div className={`generated-confetti ${className}`} aria-hidden="true">
      {PIECES.slice(0, pieces).map((piece, index) => (
        <span
          key={index}
          className={`generated-confetti-piece confetti-color-${piece.color}`}
          style={{ left: piece.left, animationDelay: piece.delay, animationDuration: piece.duration }}
        />
      ))}
    </div>
  );
}
