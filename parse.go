package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	tagSection = "api:section"
	tagMethod  = "api:method"
)

func parse(path string) (*sectionGroup, error) {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	sections := &sectionGroup{
		Sections:    []*section{},
		sectionsMap: map[string]*section{},
	}

	// Loop through comments groups for documentation tags
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, c := range file.Comments {
				if err := parseGroup(fset, c, sections); err != nil {
					return nil, err
				}
			}
		}
	}

	return sections, nil
}

func parseGroup(fset *token.FileSet, comments *ast.CommentGroup, sections *sectionGroup) error {
	for idx, comment := range comments.List {
		text := trimPrefix(comment.Text)

		if strings.HasPrefix(text, tagSection+" ") {
			if err := parseSection(fset, comments.List[idx:], sections); err != nil {
				return err
			}
		} else if strings.HasPrefix(text, tagMethod+" ") {
			if err := parseEndpoint(fset, comments.List[idx:], sections); err != nil {
				return err
			}
		}
	}

	return nil
}

func parseSection(fset *token.FileSet, comments []*ast.Comment, sections *sectionGroup) error {
	var s *section

	for _, comment := range comments {
		text := trimPrefix(comment.Text)

		if strings.HasPrefix(text, tagSection+" ") {
			parts := strings.SplitN(text, " ", 3)
			if len(parts) < 3 {
				return fmt.Errorf("%s: malformed section tag", fset.PositionFor(comment.Pos(), true))
			}

			s = sections.Get(parts[1])
			s.DisplayName, _ = strconv.Unquote(parts[2])
		} else if len(s.Description) > 0 {
			s.Description += "\n" + text
		} else {
			s.Description += text
		}
	}

	return nil
}

func parseEndpoint(fset *token.FileSet, comments []*ast.Comment, sections *sectionGroup) error {
	var (
		endpointPath string
		dataPos      token.Pos
	)

	data := bytes.NewBuffer(nil)
	m := &method{}

	for _, comment := range comments {
		text := trimPrefix(comment.Text)

		if strings.HasPrefix(text, tagMethod+" ") {
			parts := strings.SplitN(text, " ", 4)
			if len(parts) < 3 {
				return fmt.Errorf("%s: malformed endpoint tag", fset.PositionFor(comment.Pos(), true))
			}

			m.Method = parts[1]
			endpointPath = parts[2]
			m.Summary, _ = strconv.Unquote(parts[3])
		} else if text == "---" {
			dataPos = comment.Pos()
		} else if dataPos.IsValid() {
			data.WriteString(text + "\n")
		} else if len(m.Description) > 0 {
			m.Description += "\n" + text
		} else {
			m.Description += text
		}
	}

	if dataPos.IsValid() {
		if err := yaml.Unmarshal(data.Bytes(), m); err != nil {
			return fmt.Errorf("%s: failed to parse data: %s", fset.PositionFor(dataPos, true), err)
		}

		m.Description = strings.Trim(m.Description, "\n")

		for _, p := range m.Parameters {
			if p.In == "" {
				p.In = "query"
			}
		}
	}

	sections.GetEndpoint(m.Section, endpointPath).SetMethod(m)

	return nil
}

func trimPrefix(input string) string {
	if input == "//" {
		return ""
	}

	return strings.TrimPrefix(input, "// ")
}
