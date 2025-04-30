package tsgen

import "github.com/sparkymat/tsgen/tstype"

func (s *Service) AddModel(v any) error {
	m, err := tstype.StructToTSType(v, false)
	if err != nil {
		return err
	}

	s.models[m.Name()] = m

	return nil
}
