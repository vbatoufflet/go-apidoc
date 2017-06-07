package main

type endpoint struct {
	Path    string  `yaml:"path"`
	Methods methods `yaml:"methods"`

	methodsMap map[string]*method
}

func (e *endpoint) SetMethod(method *method) {
	_, ok := e.methodsMap[method.Method]
	e.methodsMap[method.Method] = method

	if !ok {
		e.Methods = append(e.Methods, method)
	}
}

type endpoints []*endpoint

func (e endpoints) Len() int {
	return len(e)
}

func (e endpoints) Less(i, j int) bool {
	return e[i].Path < e[j].Path
}

func (e endpoints) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
