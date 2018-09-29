package gorpheus

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (m Migrations) Len() int {
	return len(m)
}

func (m Migrations) Less(i, j int) bool {
	// let's answer the question - should left be applied before right ?
	// if so, let's return True here.
	left := m[i]
	right := m[j]

	fmt.Printf("Comparing %s with %s\n", left.Revision(), right.Revision())
	fmt.Printf("right dependencies: %+v\n", right.Dependencies())

	for _, dep := range right.Dependencies() {
		if left.Revision() == dep {
			log.Debugf("Found %v in %v dependencies")
			return true
		}
	}

	// empirical bool table:
	// v: version[i] < version[j]
	// p: prev.v [i] < prev.v [j]
	// out: 0 => don't swap, 1 => swap

	// v | p | out
	// 0 | 0 | 0
	// 1 | 0 | 1
	// 1 | 1 | 0
	// 0 | 1 | 0
	// if (m[i].Version < m[j].Version) && (m[i].PreviousVersion > m[j].PreviousVersion) {
	// 	return true
	// }
	// fmt.Printf("Dependencies: %+v with %+v\n", m[i].Depends(), m[j].Depends())

	// leftDep := m[i].Depends()
	// rightDep := m[j].Depends()

	// leftNamespace := m[i].ID().Namespace

	// if leftDep == nil && rightDep != nil {
	// 	if rightDep.Namespace == leftNamespace {
	// 		if left.Version < rightDep.Version {
	// 			return true
	// 		}
	// 	}
	// }

	return false
}

func (m Migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
