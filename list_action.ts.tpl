    list: builder.query<{{ .ResponseType }}, {{ .RequestType }}>({
      query: ({ {{ .FieldNames }} }: {{.RequestType }}) =>
        `{{ .ResourceURL }}{{ .ResourceQuery }}`,
      providesTags: [{ type: '{{ .Resource }}', id: 'LIST' }],
    }),
