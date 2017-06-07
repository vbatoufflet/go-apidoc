package main

import (
	"io"
	"net/http"
	"path/filepath"
	"sort"
	"text/template"
)

type sectionGroup struct {
	Sections []*section

	sectionsMap map[string]*section
}

func (sg *sectionGroup) Get(name string) *section {
	if _, ok := sg.sectionsMap[name]; !ok {
		sg.sectionsMap[name] = &section{
			Name:         name,
			Endpoints:    endpoints{},
			endpointsMap: map[string]*endpoint{},
		}

		sg.Sections = append(sg.Sections, sg.sectionsMap[name])
	}

	return sg.sectionsMap[name]
}

func (sg *sectionGroup) GetEndpoint(name, path string) *endpoint {
	return sg.Get(name).GetEndpoint(path)
}

type section struct {
	Name        string    `yaml:"name"`
	DisplayName string    `yaml:"display_name"`
	Description string    `yaml:"description"`
	Endpoints   endpoints `yaml:"endpoints"`

	endpointsMap map[string]*endpoint
}

func (s *section) GetEndpoint(path string) *endpoint {
	if _, ok := s.endpointsMap[path]; !ok {
		s.endpointsMap[path] = &endpoint{
			Path:       path,
			Methods:    []*method{},
			methodsMap: map[string]*method{},
		}

		s.Endpoints = append(s.Endpoints, s.endpointsMap[path])
	}

	return s.endpointsMap[path]
}

func (s *section) Render(path string, r io.Writer) error {
	// Sort endpoints and methods
	sort.Sort(s.Endpoints)
	for _, e := range s.Endpoints {
		sort.Sort(e.Methods)
	}

	// Render documentation
	tmpl, err := template.New(filepath.Base(path)).Funcs(template.FuncMap{
		"status_text": http.StatusText,
	}).ParseFiles(path)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(r, s); err != nil {
		return err
	}

	return nil
}
