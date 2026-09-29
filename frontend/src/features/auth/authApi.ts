import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { apiBase } from "../../runtime";

export const authApi = createApi({
  reducerPath: "authApi",
  baseQuery: fetchBaseQuery({ baseUrl: apiBase }),
  endpoints: (build) => ({
    exchangeCode: build.mutation<{ access_token: string }, { code: string }>({
      query: (body) => ({ url: "/token", method: "POST", body }),
    }),
  }),
});

export const { useExchangeCodeMutation } = authApi;
