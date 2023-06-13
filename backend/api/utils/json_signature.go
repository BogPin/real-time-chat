package utils

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
)

var specialKinds = []string{"struct", "array", "slice"}

func isSpecialKind(t reflect.Type) bool {
	return slices.Contains(specialKinds, t.Kind().String())
}

func GetJSONSignature(t reflect.Type) string {
	switch t.Kind().String() {
	case "array":
		if isSpecialKind(t.Elem()) {
			return GetJSONSignature(t.Elem()) + "[]"
		} else {
			return t.Elem().String() + "[]"
		}
	case "slice":
		if isSpecialKind(t.Elem()) {
			return GetJSONSignature(t.Elem()) + "[]"
		} else {
			return t.Elem().String() + "[]"
		}
	case "struct":
		fields := reflect.VisibleFields(t)
		fieldDefs := make([]string, 0)
		for _, field := range fields {
			strType := GetJSONSignature(field.Type)
			strName := field.Name
			if jsonTagVal := field.Tag.Get("json"); jsonTagVal != "" {
				strName = jsonTagVal
			}
			fieldDefs = append(fieldDefs, fmt.Sprintf("\t%s: %s", strName, strType))
		}
		return "{\n" + strings.Join(fieldDefs, ",\n") + "\n}"
	default:
		return t.String()
	}
}
