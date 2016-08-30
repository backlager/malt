package ingredient

// Ingredient is a core component used in the making of a good
// base malt.
type Ingredient struct {
	FilePath string
	Name     string
	Grains   map[string]Grain
}
