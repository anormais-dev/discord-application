import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { apiBase } from "../utils/runtime";

export const authApi = createApi({
  reducerPath: "authApi",
  baseQuery: fetchBaseQuery({ baseUrl: apiBase }),
  endpoints: (build) => ({
    getConfig: build.query<{ clientId: string }, void>({
      query: () => "/config",
    }),
    getMaps: build.query<string[], void>({
      query: () => "/maps",
      transformResponse: (res: { maps: string[] }) => res.maps,
    }),
    exchangeCode: build.mutation<{ access_token: string }, { code: string }>({
      query: (body) => ({ url: "/token", method: "POST", body }),
    }),
  }),
});

export const { useGetConfigQuery, useGetMapsQuery, useExchangeCodeMutation } = authApi;
