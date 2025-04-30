    {{ .MethodName }}: builder.query<{{ .ResponseType }}, {{ .RequestType }}>({
      query: ({ {{ .FieldNames }} }: {{.RequestType }}) =>
        `{{ .ResourceURL }}/{{ .MethodName }}{{ .ResourceQuery }}`,
      providesTags: [{ type: '{{ .Resource }}', id: '{{ .MethodName }}' }],
    }),
