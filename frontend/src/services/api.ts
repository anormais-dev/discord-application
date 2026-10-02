import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { apiBase } from "../utils/runtime";

export const authApi = createApi({
  reducerPath: "authApi",
  baseQuery: fetchBaseQuery({ baseUrl: apiBase }),
  endpoints: (build) => ({
    getConfig: build.query<{ clientId: string }, void>({
      query: () => "/config",
    }),
    exchangeCode: build.mutation<{ access_token: string }, { code: string }>({
      query: (body) => ({ url: "/token", method: "POST", body }),
    }),
  }),
});

export const { useGetConfigQuery, useExchangeCodeMutation } = authApi;
