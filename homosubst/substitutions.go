package homosubst

func (s *Substitutor) SubstituteByte(b byte) rune {
	if b >= 'A' && b <= 'Z' {
		return s.substituteByte(b)
	} else {
		return rune(b)
	}
}

func (s *Substitutor) SubstituteRune(r rune) rune {
	return s.substituteByte(byte(r))
}

func (s *Substitutor) substituteByte(b byte) rune {
	bi := b - 'A'
	index := s.substitutionIndex[bi]
	substitutionList := s.substitutions[bi]
	substitutionListSize := uint16(len(substitutionList))
	result := substitutionList[index]

	if substitutionListSize > 1 {
		index++
		if index >= substitutionListSize {
			index = 0
		}
		s.substitutionIndex[bi] = index
	}

	return result
}
