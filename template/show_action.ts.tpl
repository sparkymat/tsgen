    show: builder.query<{{ .ResponseType }}, {{ .RequestType }}>({
      query: id => `{{ .ResourceURL }}/${id}`,
      providesTags: (_result, _error, arg) => [{ type: '{{ .Resource }}', id: arg }],
    }),
