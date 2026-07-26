package messages

func mustJSON(msg interface{ MarshalJSON() ([]byte, error) }) string {
	b, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(b)
}

func addJSONSpaces(jsonData []byte) []byte {
	spaced := make([]byte, 0, len(jsonData))
	inString := false
	isEscaped := false

	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	for _, b := range jsonData {
		if inString {
			spaced = append(spaced, b)
			if isEscaped {
				isEscaped = false
				continue
			}
			if b == '\\' {
				isEscaped = true
				continue
			}
			if b == '"' {
				inString = false
			}
			continue
		}

		switch b {
		case '"':
			inString = true
			spaced = append(spaced, b)
		case ':':
			spaced = append(spaced, ':', ' ')
		case ',':
			spaced = append(spaced, ',', ' ')
		default:
			spaced = append(spaced, b)
		}
	}

	return spaced
}
