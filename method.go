package main

type method struct {
	Method      string            `yaml:"method"`
	Section     string            `yaml:"section"`
	Summary     string            `yaml:"summary"`
	Description string            `yaml:"description"`
	Parameters  []*parameter      `yaml:"parameters"`
	Responses   map[int]*response `yaml:"responses"`
}

type methods []*method

func (m methods) Len() int {
	return len(m)
}

func (m methods) Less(i, j int) bool {
	return m[i].Method < m[j].Method
}

func (m methods) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type parameter struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required,omitempty"`
	In          string `yaml:"in"`
	Default     string `yaml:"default,omitempty"`
}

type response struct {
	Type    string            `yaml:"type,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
	Example *example          `yaml:"example,omitempty"`
}

type example struct {
	Headers map[string]string `yaml:"headers,omitempty"`
	Format  string            `yaml:"format"`
	Body    string            `yaml:"body,omitempty"`
}
