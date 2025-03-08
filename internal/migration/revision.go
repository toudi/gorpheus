package migration

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"
)

// Revision is a structure that will allow us to traverse a graph. It just so happens that our
// nodes are namespace revisions
type Revision struct {
	Namespace string
	Name      string
	Version   int
}

func (r Revision) String() string {
	return r.Namespace + "/" + r.Name
}

func LookupID(namespace string, version int) string {
	return fmt.Sprintf("%s/%d", namespace, version)
}

// We need to know if revision `a` should be moved before `b` inside the graph.
func (r Revision) Compare(other Revision) int {
	// let's see what happens here.
	// This function is called during building of the graph stage. Each of our migrations
	// has a `Revision` and can have 1 .. n dependencies.
	// When we're traversing the nodes, we inspect the dependencies array and check if we
	// need to reorder them
	if r.Namespace == other.Namespace {
		// the easiest of cases.
		return cmp.Compare(r.Version, other.Version)
	}
	// if these two are from different namespaces then we can at least compare them by
	// namespaces themselves and sort them alphabetically.
	return strings.Compare(r.Namespace, other.Namespace)
}

func ParseId(id string) (namespace string, name string, version int, err error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 2 {
		return "", "", -1, ErrInvalidId
	}
	namespace = idParts[0]
	name = idParts[1]

	nameParts := strings.SplitN(name, "_", 2)
	if nameParts == nil {
		return "", "", -1, ErrInvalidId
	}
	version, err = strconv.Atoi(nameParts[0])
	return namespace, name, version, err
}

// this function must parse `Revision` from the provided `input` string.
// example input: "users/0001_initial"
func RevisionFromString(input string) (Revision, error) {
	namespace, name, version, err := ParseId(input)
	return Revision{
		Namespace: namespace,
		Name:      name,
		Version:   version,
	}, err
}
