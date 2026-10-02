import { wsConnect } from "./store/roomSlice";

// Modo de desenvolvimento fora do Discord. Só funciona com o backend rodando com DEV_AUTH=true.
// Todas as abas entram na mesma roleta; o nome fica salvo por aba.
const STORAGE_KEY = "dev-user";
const DEV_INSTANCE = "local";

export function savedDevName(): string | null {
  const fromUrl = new URLSearchParams(location.search).get("user");
  if (fromUrl) return fromUrl;
  try {
    return sessionStorage.getItem(STORAGE_KEY);
  } catch {
    return null;
  }
}

export function startDevSession(name: string) {
  try {
    sessionStorage.setItem(STORAGE_KEY, name);
  } catch {
    // sem sessionStorage, o nome só vale até recarregar
  }
  return wsConnect({ accessToken: `dev:${name}`, instanceId: DEV_INSTANCE });
}
