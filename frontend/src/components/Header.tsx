import { useState } from "react";
import { useAppSelector } from "../store/hooks";
import { selectPermissions, selectRoom } from "../store/selectors";
import type { Phase } from "../types/room";
import { GearIcon } from "./Icons";
import { SettingsPanel } from "./SettingsPanel";
import { SoundToggle } from "./SoundToggle";

const STEPS: { phase: Phase; label: string }[] = [
  { phase: "lobby", label: "Lobby" },
  { phase: "drafting", label: "Sorteio" },
  { phase: "finished", label: "Times" },
];

// Etapas à esquerda; à direita, o som (para todos) e as configurações (só o admin abre).
// O nome da call ficaria aqui também: ver o bloco comentado em services/discord.ts.
export function Header() {
  const room = useAppSelector(selectRoom);
  const perms = useAppSelector(selectPermissions);
  const [settingsOpen, setSettingsOpen] = useState(false);
  if (!room) return null;

  const current = STEPS.findIndex((s) => s.phase === room.phase);

  return (
    <header className="top">
      <nav className="steps" aria-label="Etapas">
        {STEPS.map((s, i) => (
          <div key={s.phase} className={`step ${i === current ? "on" : i < current ? "done" : ""}`}>
            <span className="step-n px">{i + 1}</span>
            <span className="step-label">{s.label}</span>
          </div>
        ))}
      </nav>
      <div className="top-r">
        <SoundToggle />
        {perms.isAdmin && (
          <>
            <button
              className={`icon-btn gear ${settingsOpen ? "on" : ""}`}
              aria-label="Configurações"
              aria-expanded={settingsOpen}
              onClick={() => setSettingsOpen(!settingsOpen)}
            >
              <GearIcon />
            </button>
            {settingsOpen && <SettingsPanel onClose={() => setSettingsOpen(false)} />}
          </>
        )}
      </div>
    </header>
  );
}
