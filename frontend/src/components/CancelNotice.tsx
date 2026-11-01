import { useEffect, useState } from "react";
import { useAppSelector } from "../store/hooks";
import { selectMe, selectRoom } from "../store/selectors";

const NOTICE_MS = 4000;
// Último cancelamento já visto; a página remonta ao voltar do PiP e não deve repetir o aviso.
// undefined enquanto não chegou nenhum estado: quem abre depois do cancelamento não vê aviso velho.
let seenAt: string | null | undefined;

// Avisa a sala quem cancelou a roleta.
export function CancelNotice() {
  const room = useAppSelector(selectRoom);
  const me = useAppSelector(selectMe);
  const [text, setText] = useState<string | null>(null);
  const cancelled = room?.cancelled ?? null;
  const at = cancelled?.at ?? null;

  useEffect(() => {
    if (!room) return;
    const isNew = seenAt !== undefined && at !== null && at !== seenAt;
    seenAt = at;
    if (isNew && cancelled) setText(cancelled.byId === me ? "Você cancelou a roleta" : `${cancelled.byName} cancelou a roleta`);
  }, [at, !!room]);

  useEffect(() => {
    if (!text) return;
    const id = setTimeout(() => setText(null), NOTICE_MS);
    return () => clearTimeout(id);
  }, [text]);

  if (!text) return null;
  return (
    <div className="toast toast-info" role="status">
      {text}
    </div>
  );
}
