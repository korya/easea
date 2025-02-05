package main

type Match struct {
	File   string
	Line   string
	Column string
}

type Parser interface {
	Process(text string) []Match
}

type ParserFunc func(text string) []Match

func (f ParserFunc) Process(text string) []Match {
	return f(text)
}

type Handler interface {
	Handle(match Match)
}
