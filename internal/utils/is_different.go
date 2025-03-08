package utils

// this method is used by dbState to compare column revisions
// against each other.
func IsDifferent(a, b map[string]interface{}) bool {
	// some sanity checks - if the number of keys is different then
	// obviously the maps are different
	if len(a) != len(b) {
		return true
	}
	// now check each key
	for keyA, valueA := range a {
		// nothing fancy really
		valueB, exists := b[keyA]
		if !exists || valueB != valueA {
			return true
		}
	}
	// if we haven't returned yet then the maps are equal
	return false
}
