package operation

// Some of the operations (like CreateTable / AddField) can produce more than a single
// operation. For instance, when you create a table and the fields contain indexes,
// CreateTable will be followed by CreateIndex.
// Similarly, when having some foreign keys, the CreteTable may be folllowed by
// Set constraint DDL's.
func (o Operation) Unpack() ([]Operation, error) {
	// regardless of the outcome, add the "root" operation to the result queue
	var result []Operation

	// here are known examples of DDL operations that can yield side effects:
	if o.CreateTable != nil {
		result = append(result, o)
		if sideEffects, err := o.CreateTable.SideEffects(); err != nil {
			return nil, err
		} else {
			result = append(result, sideEffects...)
		}
		return result, nil
	}
	if o.AddField != nil {
		// addField supports multiple fields, but the syntax to do this is not the same across
		// database engines. Therefore let's first convert them to an array of single operations.
		for _, f := range o.AddField.SingleFieldAdds() {
			result = append(result, f)
			// each addField can have their own side effects (like creating an index or a foreign key)
			addFieldSideEffects, err := f.AddField.SideEffects()
			if err != nil {
				return nil, err
			}
			result = append(result, addFieldSideEffects...)
		}

		return result, nil
	}

	// if we are here then it means that no specific branch was executed
	// and we can safely return the source operation as is.
	return []Operation{o}, nil
}

// for the internal logic of the code, it is important for the dbState to know when a "source"
// operation was applied. However, because we're unpacking the operations it is possible for
// a single "source" operation to yield more than one operations. In order to keep everything
// at sync, we're just going to insert a fake operation after unpacking operations.
func UnpackOperations(operations []Operation, adjustTableVersion bool) ([]Operation, error) {
	var result []Operation

	for _, operation := range operations {
		unpacked, err := operation.Unpack()
		if err != nil {
			return nil, err
		}
		result = append(result, unpacked...)
		if adjustTableVersion {
			result = append(
				result,
				Operation{
					AdjustTableVersion: &AdjustTableVersion{
						TableName: operations[0].TableName(),
						Delta:     1,
					},
				},
			)
		}
	}

	return result, nil
}
