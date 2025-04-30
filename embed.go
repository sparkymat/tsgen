package tsgen

import _ "embed"

//go:embed model.ts.tpl
var modelTS string

//go:embed interface.ts.tpl
var interfaceTS string

//go:embed slice.ts.tpl
var sliceTS string

//go:embed create_action.ts.tpl
var createActionTS string

//go:embed show_action.ts.tpl
var showActionTS string

//go:embed custom_query_action.ts.tpl
var customQueryActionTS string

//go:embed list_action.ts.tpl
var listActionTS string

//go:embed destroy_action.ts.tpl
var destroyActionTS string

//go:embed update_action.ts.tpl
var updateActionTS string

//go:embed custom_member_action.ts.tpl
var customMemberActionTS string

//go:embed custom_member_multipart_action.ts.tpl
var customMemberMultipartAction string
