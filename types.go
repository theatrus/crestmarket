package crestmarket

// An EVE-Online Region
type Region struct {
	Name string
	Href string
	Id   int
}

// A set of Regions, with pre-filtered views
type Regions struct {
	AllRegions []*Region
}

// Build a new Regions structure
func newRegions() *Regions {
	return &Regions{make([]*Region, 0)}
}

// An inventory type
type InventoryType struct {
	Name string
	Href string
	Id   int
}

// A collection of inventory types
type InventoryTypes struct {
	Types []*InventoryType
}

// Build a new InventoryTypes structure
func newInventoryTypes() *InventoryTypes {
	return &InventoryTypes{make([]*InventoryType, 0)}
}

// Defines a Root of all possible resources on this API
type Root struct {
	// Provides a list of canonical resources and their URL root
	Resources map[string]string
}
