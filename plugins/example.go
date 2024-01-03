package main

// the empty struct used for the plugin
type plugin struct {
}

// returns the version of plugin in major, minor, and patch
// levels as separate integers
func (p *plugin) Version() (int, int, int) {
	return 0, 1, 0
}

// returns the short name of the module as a string
func (p *plugin) Name() string {
	return "example"
}

// returns a short description of the module as a string
func (p *plugin) Description() string {
	return "example plugin for code template"
}

var PlugIn plugin
