package common

import (
	"basics/internal/runtime"
	"strings"
)

func VarType(name string) string {
	if strings.HasSuffix(name, "%") {
		return "int"
	}
	if strings.HasSuffix(name, "$") {
		return "string"
	}
	return "float"
}

func VarTypeAsInt(name string) runtime.ValueType {
	if strings.HasSuffix(name, "%") {
		return runtime.INTEGER
	}
	if strings.HasSuffix(name, "$") {
		return runtime.STRING
	}
	return runtime.NUMBER
}
