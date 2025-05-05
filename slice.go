package tsgen

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"github.com/sparkymat/tsgen/template"
	"github.com/sparkymat/tsgen/tstype"
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

type SliceEntry struct {
	ParentResourceName string
	MethodName         string
	Action             Action
	RequestType        string
	RequestTSType      *tstype.TSType
	ResponseType       string
	ResponseTSType     *tstype.TSType
	RequestFields      []string
	FileField          string
}

type Slice struct {
	Name           string
	Entries        []SliceEntry
	Interfaces     map[string]tstype.TSType
	ImportedModels []string
}

func (s Slice) ReducerPath() string {
	pl := pluralize.NewClient()

	return strcase.ToSnake(pl.Plural(s.Name))
}

func (s Slice) ResourceURL() string {
	pl := pluralize.NewClient()

	return strcase.ToSnake(pl.Plural(s.Name))
}

func (s Slice) RenderedInterfaceDefinitions() (string, error) {
	v := ""

	for _, interfaceEntry := range s.Interfaces {
		renderedInterface, err := renderTemplateToString(template.InterfaceTS, &interfaceEntry)
		if err != nil {
			return "", err
		}

		v += (renderedInterface + "\n")
	}

	return v, nil
}

func (s Slice) RenderedEndpoints() (string, error) {
	v := ""

	customCollectionEntries := lo.Filter(s.Entries, func(e SliceEntry, _ int) bool {
		return e.Action == ActionCustomQuery
	})

	customInvalidates := lo.Map(customCollectionEntries, func(e SliceEntry, _ int) string {
		return fmt.Sprintf("{type: '%s', id: '%s'},", s.Name, e.MethodName)
	})

	customInvalidatesString := strings.Join(customInvalidates, "\n")

	pl := pluralize.NewClient()

	for _, entry := range s.Entries {
		resourceURL := s.ResourceURL()

		if entry.ParentResourceName != "" {
			parentURL := strcase.ToSnake(pl.Plural(entry.ParentResourceName)) + "/${parentId}/"

			resourceURL = parentURL + resourceURL
		}

		switch entry.Action {
		case ActionCreate:
			requestType := entry.RequestType
			if entry.ParentResourceName != "" {
				requestType += "WithParent"
			}

			requestInput := "request"
			if entry.ParentResourceName != "" {
				requestInput = "({ parentId, request })"
			}

			renderedEntry, err := renderTemplateToString(template.CreateActionTS, map[string]string{
				"RequestInput":      requestInput,
				"ResponseType":      entry.ResponseType,
				"RequestType":       requestType,
				"ResourceURL":       resourceURL,
				"Resource":          s.Name,
				"CustomInvalidates": customInvalidatesString,
			})
			if err != nil {
				return "", err
			}

			v += renderedEntry
		case ActionShow:
			renderedEntry, err := renderTemplateToString(template.ShowActionTS, map[string]string{
				"ResponseType": entry.ResponseType,
				"RequestType":  entry.RequestType,
				"ResourceURL":  resourceURL,
				"Resource":     s.Name,
			})
			if err != nil {
				return "", err
			}

			v += renderedEntry
		case ActionCustomQuery:
			resourceQuery := strings.Join(lo.Map(entry.RequestFields, func(v string, _ int) string {
				return v + "=${encodeURIComponent(" + v + ")}"
			}), "&")

			fieldNames := strings.Join(entry.RequestFields, ", ")

			if entry.ParentResourceName != "" {
				fieldNames = "parentId, request: {" + fieldNames + "}"
			}

			requestType := entry.RequestType
			if entry.ParentResourceName != "" {
				requestType += "WithParent"
			}

			values := map[string]string{
				"MethodName":    entry.MethodName,
				"ResponseType":  entry.ResponseType,
				"RequestType":   requestType,
				"ResourceURL":   resourceURL,
				"Resource":      s.Name,
				"ResourceQuery": "?" + resourceQuery,
				"FieldNames":    fieldNames,
			}

			renderedEntry, err := renderTemplateToString(template.CustomQueryActionTS, values)
			if err != nil {
				return "", err
			}

			v += renderedEntry
		case ActionList:
			resourceQuery := strings.Join(lo.Map(entry.RequestFields, func(v string, _ int) string {
				return v + "=${encodeURIComponent(" + v + ")}"
			}), "&")

			fieldNames := strings.Join(entry.RequestFields, ", ")

			if entry.ParentResourceName != "" {
				if len(entry.RequestFields) > 0 {
					fieldNames = "parentId, request: {" + fieldNames + "}"
				} else {
					fieldNames = "parentId"
				}
			}

			requestType := entry.RequestType
			if entry.ParentResourceName != "" {
				requestType += "WithParent"
			}

			values := map[string]string{
				"ResponseType":  entry.ResponseType,
				"RequestType":   requestType,
				"ResourceURL":   resourceURL,
				"Resource":      s.Name,
				"ResourceQuery": "?" + resourceQuery,
				"FieldNames":    fieldNames,
			}

			renderedEntry, err := renderTemplateToString(template.ListActionTS, values)
			if err != nil {
				return "", err
			}

			v += renderedEntry
		case ActionDestroy:
			queryParams := "id"
			requestType := entry.RequestType

			// Assume request type is string since you can't pass info to DELETE
			if entry.ParentResourceName != "" {
				requestType = s.Name + "DestroyRequest"
				queryParams = "{ id, parentId } : " + s.Name + "DestroyRequest"
			}

			renderedEntry, err := renderTemplateToString(template.DestroyActionTS, map[string]string{
				"ResourceURL":       resourceURL,
				"Resource":          s.Name,
				"RequestType":       requestType,
				"QueryParams":       queryParams,
				"CustomInvalidates": customInvalidatesString,
			})
			if err != nil {
				return "", err
			}

			v += renderedEntry
		case ActionUpdate:
			otherFields := lo.Filter(entry.RequestFields, func(f string, _ int) bool { return f != "id" })
			fieldAssignments := strings.Join(otherFields, ",\n")
			fieldNames := strings.Join(otherFields, ", ")

			if entry.ParentResourceName != "" {
				fieldNames = "parentId"
				innerFieldNames := strings.Join(entry.RequestFields, ", ")
				fieldNames += ", request: {" + innerFieldNames + "}"
			}

			requestType := entry.RequestType
			if entry.ParentResourceName != "" {
				requestType += "WithParent"
			}

			invalidateIdField := "id"
			if entry.ParentResourceName != "" {
				invalidateIdField = "request.id"
			}

			renderedEntry, err := renderTemplateToString(template.UpdateActionTS, map[string]string{
				"ResourceURL":       resourceURL,
				"Resource":          s.Name,
				"FieldNames":        fieldNames,
				"ResponseType":      entry.ResponseType,
				"RequestType":       requestType,
				"FieldAssignments":  fieldAssignments,
				"CustomInvalidates": customInvalidatesString,
				"InvalidateIdField": invalidateIdField,
			})
			if err != nil {
				return "", err
			}

			v += renderedEntry
		case ActionCustomMemberAction:
			otherFields := lo.Filter(entry.RequestFields, func(f string, _ int) bool { return f != "id" })
			fieldAssignments := strings.Join(otherFields, ",\n")
			fieldNames := strings.Join(otherFields, ", ")

			if entry.ParentResourceName != "" {
				fieldNames = "parentId"
				innerFieldNames := strings.Join(entry.RequestFields, ", ")
				fieldNames += ", request: {" + innerFieldNames + "}"
			}

			requestType := entry.RequestType
			if entry.ParentResourceName != "" {
				requestType += "WithParent"
			}

			invalidateIdField := "id"
			if entry.ParentResourceName != "" {
				invalidateIdField = "request.id"
			}

			renderedEntry, err := renderTemplateToString(template.CustomMemberActionTS, map[string]string{
				"MethodName":        entry.MethodName,
				"ResourceURL":       resourceURL,
				"Resource":          s.Name,
				"FieldNames":        fieldNames,
				"ResponseType":      entry.ResponseType,
				"RequestType":       requestType,
				"FieldAssignments":  fieldAssignments,
				"CustomInvalidates": customInvalidatesString,
				"InvalidateIdField": invalidateIdField,
			})
			if err != nil {
				return "", err
			}

			v += renderedEntry
		case ActionCustomMemberMultipartAction:
			otherFields := lo.Filter(entry.RequestFields, func(f string, _ int) bool { return f != "id" })
			formDataAssignments := strings.Join(lo.Map(otherFields, func(fName string, _ int) string {
				return fmt.Sprintf("formData.append('%s', %s);", fName, fName)
			}), "\n")
			fieldNames := strings.Join(otherFields, ", ")

			if entry.ParentResourceName != "" {
				fieldNames = "parentId"
				innerFieldNames := strings.Join(entry.RequestFields, ", ")
				fieldNames += ", request: {" + innerFieldNames + "}"
			}

			requestType := entry.RequestType
			if entry.ParentResourceName != "" {
				requestType += "WithParent"
			}

			invalidateIdField := "id"
			if entry.ParentResourceName != "" {
				invalidateIdField = "request.id"
			}

			renderedEntry, err := renderTemplateToString(template.CustomMemberMultipartAction, map[string]string{
				"MethodName":          entry.MethodName,
				"ResourceURL":         resourceURL,
				"Resource":            s.Name,
				"FieldNames":          fieldNames,
				"ResponseType":        entry.ResponseType,
				"RequestType":         requestType,
				"FormDataAssignments": formDataAssignments,
				"FileField":           entry.FileField,
				"CustomInvalidates":   customInvalidatesString,
				"InvalidateIdField":   invalidateIdField,
			})
			if err != nil {
				return "", err
			}

			v += renderedEntry

		case ActionCustomAction:
		case ActionCustomMemberQuery:
		default:
		}
	}

	return v, nil
}

func (s Slice) RenderedExports() string {
	v := ""

	for _, entry := range s.Entries {
		switch entry.Action {
		case ActionCreate:
			v += "  useCreateMutation,\n"
		case ActionShow:
			v += "  useShowQuery,\n"
		case ActionDestroy:
			v += "  useDestroyMutation,\n"
		case ActionCustomQuery:
			v += fmt.Sprintf("  use%sQuery,\n", strcase.ToCamel(entry.MethodName))
		case ActionList:
			v += "  useListQuery,\n"
		case ActionUpdate:
			v += "  useUpdateMutation,\n"
		case ActionCustomMemberAction:
			v += fmt.Sprintf("  use%sMutation,\n", strcase.ToCamel(entry.MethodName))
		case ActionCustomMemberMultipartAction:
			v += fmt.Sprintf("  use%sMutation,\n", strcase.ToCamel(entry.MethodName))
		case ActionCustomAction:
		case ActionCustomMemberQuery:
		default:
		}
	}

	return v
}

func (s Slice) RenderedImports() string {
	v := ""

	interfaceTypes := []string{}
	usedTypes := []string{}

	for _, entry := range s.Entries {
		if entry.RequestTSType != nil {
			for _, fieldType := range entry.RequestTSType.Fields() {
				usedTypes = append(usedTypes, fieldType)
			}
		}

		if entry.ResponseTSType != nil {
			for _, fieldType := range entry.ResponseTSType.Fields() {
				// strip array bits
				cleanedupFieldType := strings.TrimSuffix(fieldType, "[]")
				usedTypes = append(usedTypes, cleanedupFieldType)
			}
		}
	}

	missingTypes := lo.Filter(usedTypes, func(tt string, _ int) bool {
		return !lo.Contains(interfaceTypes, tt)
	})

	// Remove basic types and sort it
	missingTypes = lo.Filter(missingTypes, func(tt string, _ int) bool {
		return !lo.Contains([]string{"string", "number", "boolean"}, tt)
	})
	slices.Sort(missingTypes)
	missingTypes = lo.Uniq(missingTypes)

	for _, m := range missingTypes {
		v += fmt.Sprintf("import { %s } from '../models/%s';\n", m, m)
	}

	return v
}

func (s *Service) AddSliceEntry(
	resourceName string,
	parentResourceName string,
	methodName string,
	action Action,
	in any,
	out any,
	addIDToIn bool,
	fileField string,
) error {
	thisSlice, found := s.slices[resourceName]

	if !found {
		thisSlice = Slice{
			Name:       resourceName,
			Entries:    []SliceEntry{},
			Interfaces: map[string]tstype.TSType{},
		}
	}

	entry := SliceEntry{
		ParentResourceName: parentResourceName,
		MethodName:         methodName,
		Action:             action,
		FileField:          fileField,
	}

	if in != nil {
		if inString, isString := in.(string); isString {
			entry.RequestType = inString
			entry.RequestFields = []string{}

			if parentResourceName != "" {
				wrapperTypeName := resourceName + strcase.ToCamel(methodName) + "Request"
				wrapperType := tstype.New(wrapperTypeName)
				wrapperType.AddField("parentId", "string")
				wrapperType.AddField("id", "string")
				thisSlice.Interfaces[wrapperTypeName] = wrapperType
			}
		} else {
			inType, err := tstype.StructToTSType(in, addIDToIn)
			if err != nil {
				return err
			}

			if fileField != "" {
				inType.AddField(fileField, "File")
			}

			if _, found := s.models[inType.Name()]; found {
				if !lo.Contains(thisSlice.ImportedModels, inType.Name()) {
					thisSlice.ImportedModels = append(thisSlice.ImportedModels, inType.Name())
				}
			} else {
				thisSlice.Interfaces[inType.Name()] = inType

				if parentResourceName != "" {
					wrapperType := tstype.New(inType.Name() + "WithParent")
					wrapperType.AddField("parentId", "string")
					wrapperType.AddField("request", inType.Name())
					thisSlice.Interfaces[inType.Name()+"WithParent"] = wrapperType
				}
			}

			entry.RequestType = inType.Name()
			entry.RequestTSType = &inType
			entry.RequestFields = lo.Keys(inType.Fields())
		}
	}

	if out != nil {
		outType, err := tstype.StructToTSType(out, false)
		if err != nil {
			return err
		}

		if _, found := s.models[outType.Name()]; found {
			if !lo.Contains(thisSlice.ImportedModels, outType.Name()) {
				thisSlice.ImportedModels = append(thisSlice.ImportedModels, outType.Name())
			}
		} else {
			thisSlice.Interfaces[outType.Name()] = outType
		}

		entry.ResponseType = outType.Name()
		entry.ResponseTSType = &outType
	}

	thisSlice.Entries = append(thisSlice.Entries, entry)

	s.slices[resourceName] = thisSlice

	return nil
}
