package ingredient

// Grain represents a single file and its associated properties
type Grain struct {
	FilePath   string
	Title      string
	Properties map[string]interface{}
}
