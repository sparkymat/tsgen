    {{ .MethodName }}: builder.mutation<void, {{ .RequestType }}>({
      query: ({ id, {{ .FieldNames }} }: {{ .RequestType }}) => ({
        url: `{{ .ResourceURL }}/${id}/{{ .MethodName }}`,
        method: 'POST',
        body: function() {
          const formData = new FormData();

          {{ .FormDataAssignments }}

          return formData;
        }(),
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
