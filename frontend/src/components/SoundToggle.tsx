import { useEffect } from "react";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { MUTED_KEY, mutedToggled } from "../store/roomSlice";
import { selectMuted } from "../store/selectors";
import { setMuted } from "../utils/sound";
import { SoundIcon } from "./Icons";

export function SoundToggle() {
  const dispatch = useAppDispatch();
  const muted = useAppSelector(selectMuted);

  // Mantém o módulo de som e o navegador em dia com a escolha.
  useEffect(() => {
    setMuted(muted);
    try {
      localStorage.setItem(MUTED_KEY, muted ? "1" : "0");
    } catch {
      // sem localStorage, a escolha vale só até recarregar
    }
  }, [muted]);

  return (
    <button
      className={`icon-btn ${muted ? "" : "on"}`}
      aria-label={muted ? "Ligar o som" : "Desligar o som"}
      aria-pressed={!muted}
      title={muted ? "Ligar o som" : "Desligar o som"}
      onClick={() => dispatch(mutedToggled())}
    >
      <SoundIcon muted={muted} />
    </button>
  );
}
