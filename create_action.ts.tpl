    create: builder.mutation<{{ .ResponseType }}, {{ .RequestType }}>({
      query: {{ .RequestInput }} => ({
        url: `{{ .ResourceURL }}`,
        method: 'POST',
        body: request,
        headers: {
          'X-CSRF-Token': (
            document.querySelector('meta[name="csrf-token"]') as any
          ).content,
        },
      }),
      invalidatesTags: [
        { type: '{{ .Resource }}', id: 'LIST' },
        {{ .CustomInvalidates }}
      ],
    }),
