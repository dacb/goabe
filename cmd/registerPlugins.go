package cmd

import (
	"github.com/dacb/goabe/example"
	"github.com/dacb/goabe/life"
)

func registerPlugins() {
	example.Register()
	life.Register()
}
