// Dentro do Discord a página recebe frame_id na URL e as chamadas passam pelo proxy /.proxy.
export const inDiscord = new URLSearchParams(location.search).has("frame_id");

export const devAuth = process.env.DEV_AUTH === "true";

export const apiBase = inDiscord ? "/.proxy/api" : "/api";
