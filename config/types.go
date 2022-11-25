package config

// File contains types that is used to export data from package

// Link represents symbolic link from user config
type Link struct {
	TargetPath string
	LinkPath   string
}

// Command represents command from user config
type Command struct {
	InputPath  string `yaml:"input"`
	OutputPath string `yaml:"output"`
	Command    string `yaml:"command"`
}

// Command represents template from user config
type Template struct {
	InputPath  string            `yaml:"input"`
	OutputPath string            `yaml:"output"`
	Data       map[string]string `yaml:"data"`
}

// Config represents parsed user config
type Config struct {
	Links     map[string]Link     `yaml:"links"`
	Commands  map[string]Command  `yaml:"commands"`
	Templates map[string]Template `yaml:"templates"`
}
