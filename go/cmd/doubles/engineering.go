package main

type Engineering interface {
	DoEngineering(*Ship)
}

type NoEngineering struct{}

func (_ NoEngineering) DoEngineering(_ *Ship) {}
