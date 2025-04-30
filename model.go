package tsgen

func (s *Service) AddModel(v any) error {
	m, err := StructToTSType(v, false)
	if err != nil {
		return err
	}

	s.models[m.Name] = m

	return nil
}
