    {{ .MethodName }}: builder.mutation<void, {{ .RequestType }}>({
      query: ({ {{ .FieldNames }} }: {{ .RequestType }}) => ({
        url: `{{ .ResourceURL }}/${id}/{{ .MethodName }}`,
        method: 'POST',
        body: {
          {{ .FieldAssignments }}
        },
        headers: {
          'X-CSRF-Token': (
            document.querySelector('meta[name="csrf-token"]') as any
          ).content,
        },
      }),
      invalidatesTags: (_result, _error, arg) => [
        { type: '{{ .Resource }}', id: arg.{{ .InvalidateIdField }} },
        { type: '{{ .Resource }}', id: 'LIST' },
        {{ .CustomInvalidates }}
      ],
    }),
