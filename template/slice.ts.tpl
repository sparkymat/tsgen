import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

{{ .RenderedImports }}
{{ .RenderedInterfaceDefinitions }}export const api = createApi({
  reducerPath: '{{ .ReducerPath }}',
  baseQuery: fetchBaseQuery({ baseUrl: '/api' }),
  tagTypes: ['{{ .Name }}'],
  endpoints: builder => ({
{{ .RenderedEndpoints }}  })
});

export const {
{{ .RenderedExports }}} = api;
