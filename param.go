package router

import "fmt"

// Param describes a named parameter and its value. As the key the parameter
// name is used (without a trailing parameter) and the value a string of the
// path corresponding to the given position.
//
// I did not use the settings dictionary, because this method allows to keep the
// order and use the parameters with the same name.
type Param struct {
	Key, Value string
}

// String returns the string representation of the name and value parameter.
func (p *Param) String() string {
	return fmt.Sprintf("%v: %v", p.Key, p.Value)
}

// Params describes a list of named parameters.
type Params []Param

// Get returns the value of the first parameter in the list with the specified
// name. If such a parameter is not listed, it returns the empty string.
func (p Params) Get(name string) string {
	for _, param := range p {
		if param.Key == name {
			return param.Value
		}
	}
	return ""
}
