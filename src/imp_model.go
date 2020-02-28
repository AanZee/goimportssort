package main

// ImpModel is used for storing import information
type ImpModel struct {
	path           string
	localReference string
}

// String is used to get a string representation of an import
func (m ImpModel) String() string {
	if m.localReference == "" {
		return m.path
	}

	return m.localReference + " " + m.path
}
