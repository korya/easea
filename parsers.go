package main

type ParserRegistry map[string]Parser

func (r ParserRegistry) Process(format string, text string) []Match {
	for name, parser := range r {
		if name != format {
			continue
		}

		matches := parser.Process(text)
		if len(matches) > 0 {
			return matches
		}
	}

	return nil
}
