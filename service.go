package tsgen

func New() *Service {
	return &Service{
		models: map[string]TSType{},
		slices: map[string]Slice{},
	}
}

type Service struct {
	models map[string]TSType
	slices map[string]Slice
}
