// Macaco fora da roda (resultado, PiP). Montado como o da roda: a cabeça normal recebe as
// pupilas por cima, a tonta recebe as espirais; feliz e "parado" já têm os olhos desenhados.
type Mood = "normal" | "happy" | "dizzy" | "stand";

const HEADS: Record<Mood, string> = {
  normal: "/images/monkey-head.png",
  happy: "/images/happy-monkey-head.png",
  dizzy: "/images/dizzy-monkey-head.png",
  stand: "/images/stand-monkey-head.png",
};

export function Monkey({ mood, size, party = false }: { mood: Mood; size: number; party?: boolean }) {
  return (
    <div className="mk" style={{ width: size, height: size * 0.9 }} aria-hidden="true">
      {party && (
        <div className="mk-party">
          <img src="/images/cone.png" alt="" />
          <img className="mk-confetti" src="/images/confetti.png" alt="" />
        </div>
      )}
      <img className="mk-head" src={HEADS[mood]} alt="" />
      {mood === "normal" && <img className="mk-eyes" src="/images/monkey-eye.png" alt="" />}
      {mood === "dizzy" && (
        <>
          <img className="mk-spiral left" src="/images/espiral-eye.png" alt="" />
          <img className="mk-spiral right" src="/images/espiral-eye.png" alt="" />
        </>
      )}
    </div>
  );
}
