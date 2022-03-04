package lox

import "fmt"

// Environment holds the values of identifiers for a particular scope
type Environment struct {
	Enclosing *Environment
	Values    map[string]*Node
}

func (env *Environment) printScope() {
	fmt.Print("\n")
	scopes := []*Environment{}
	scope := env
	for scope != nil {
		scopes = append([]*Environment{scope}, scopes...)
		scope = scope.Enclosing
	}

	for depth := 0; depth < len(scopes); depth++ {
		fmt.Printf("Scope %d:\n", depth)
		for k, v := range scopes[depth].Values {
			fmt.Printf("\t%s: %s\n", k, v.ToString())
		}
	}
}
