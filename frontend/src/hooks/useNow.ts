import { useEffect, useState } from "react";

// Hora atual que se atualiza sozinha enquanto active for true (contagens regressivas).
export function useNow(active: boolean, intervalMs = 1000) {
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    if (!active) return;
    setNow(Date.now());
    const id = setInterval(() => setNow(Date.now()), intervalMs);
    return () => clearInterval(id);
  }, [active, intervalMs]);

  return now;
}
