package generator

// IdGenerator interface defines a method for generating unique identifiers.
type IdGenerator interface {
	// Generate generates a new unique identifier and returns it as a string.
	Generate() string
}
