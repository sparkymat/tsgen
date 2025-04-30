package template

import _ "embed"

//go:embed model.ts.tpl
var ModelTS string

//go:embed interface.ts.tpl
var InterfaceTS string

//go:embed slice.ts.tpl
var SliceTS string

//go:embed create_action.ts.tpl
var CreateActionTS string

//go:embed show_action.ts.tpl
var ShowActionTS string

//go:embed custom_query_action.ts.tpl
var CustomQueryActionTS string

//go:embed list_action.ts.tpl
var ListActionTS string

//go:embed destroy_action.ts.tpl
var DestroyActionTS string

//go:embed update_action.ts.tpl
var UpdateActionTS string

//go:embed custom_member_action.ts.tpl
var CustomMemberActionTS string

//go:embed custom_member_multipart_action.ts.tpl
var CustomMemberMultipartAction string
