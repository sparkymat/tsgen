{{ .Imports }}export class {{.Name }} {
{{ .RenderedFieldsForClass }}  constructor(json: {{.Name }}) {
{{ .RenderedFieldAssignments }}  }
}
