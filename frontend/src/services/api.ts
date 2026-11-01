import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { apiBase } from "../utils/runtime";

// Artes de cada mapa pelo nome: a faixa leve para o carrossel e o splash para o resultado.
export type MapImages = Record<string, { card: string; splash: string }>;
export interface MapsResponse {
  maps: string[];
  images: MapImages;
}

export const authApi = createApi({
  reducerPath: "authApi",
  baseQuery: fetchBaseQuery({ baseUrl: apiBase }),
  endpoints: (build) => ({
    getConfig: build.query<{ clientId: string }, void>({
      query: () => "/config",
    }),
    getMaps: build.query<MapsResponse, void>({
      query: () => "/maps",
    }),
    exchangeCode: build.mutation<{ access_token: string }, { code: string }>({
      query: (body) => ({ url: "/token", method: "POST", body }),
    }),
  }),
});

export const { useGetConfigQuery, useGetMapsQuery, useExchangeCodeMutation } = authApi;
