package gorpheus

import (
	"strings"
)

func (m Migrations) Len() int {
	return len(m)
}

func (m Migrations) Less(i, j int) bool {
	// let's answer the question - should left be applied before right ?
	// if so, let's return True here.
	left := m[i]
	right := m[j]

	for _, dep := range right.Dependencies() {
		if left.Revision() == dep {
			return true
		}
	}

	// if we made it here this means that there weren't any dependencies.
	if left.GetNamespace() == right.GetNamespace() {
		return strings.Compare(left.Revision(), right.Revision()) == -1
	}

	return false
}

func (m Migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
