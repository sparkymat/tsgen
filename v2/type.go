package v2

import "github.com/sparkymat/tsgen/tstype"

type ObjectType string

const (
	ObjectTypeClass     ObjectType = "class"
	ObjectTypeInterface ObjectType = "interface"
)

type Action string

const (
	ActionCreate                      Action = "create"
	ActionShow                        Action = "show"
	ActionDestroy                     Action = "destroy"
	ActionList                        Action = "list"
	ActionUpdate                      Action = "update"
	ActionCustomAction                Action = "customAction"
	ActionCustomQuery                 Action = "customQuery"
	ActionCustomMemberAction          Action = "customMemberAction"
	ActionCustomMemberMultipartAction Action = "customMemberMultipartAction"
	ActionCustomMemberQuery           Action = "customMemberQuery"
)

type Object struct {
	ObjectType ObjectType
	Name       string
}

type Endpoint struct {
	Action         Action
	MethodName     string
	ParentResource string
	RequestType    string
	RequestTSType  *tstype.TSType
	ResponseType   string
	ResponseTSType *tstype.TSType
	FileField      string
}

type Slice struct {
	Resource  string
	Types     []Object
	Endpoints []Endpoint
}

type Collection struct {
	Types  []Object
	Slices []Slice
}
