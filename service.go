package tsgen

import (
	"github.com/sparkymat/tsgen/tstype"
)

func New() *Service {
	return &Service{
		models: map[string]tstype.TSType{},
		slices: map[string]Slice{},
	}
}

type Service struct {
	models map[string]tstype.TSType
	slices map[string]Slice
}
