    destroy: builder.mutation<void, {{ .RequestType }}>({
      query: ({{ .QueryParams }}) => ({
        url: `{{ .ResourceURL }}/${id}`,
        method: 'DELETE',
        headers: {
          'X-CSRF-Token': (
            document.querySelector('meta[name="csrf-token"]') as any
          ).content,
        },
        invalidatesTags: (_result, _error, arg) => [
          { type: '{{ .Resource }}', id: arg },
          { type: '{{ .Resource }}', id: 'LIST' },
          {{ .CustomInvalidates }}        ],
      }),
    }),
