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

// TODO: Build a map indexed by something useful
func (r *Regions) ByName(name string) *Region {
	for _, region := range r.AllRegions {
		if region.Name == name {
			return region
		}
	}
	return nil
}

// An inventory type
type MarketType struct {
	Name string
	Href string
	Id   int
}

// A collection of inventory types
type MarketTypes struct {
	Types []*MarketType
}

// TODO: Build a map indexed by something useful
func (r *MarketTypes) ByName(name string) *MarketType {
	for _, mtype := range r.Types {
		if mtype.Name == name {
			return mtype
		}
	}
	return nil
}

// Build a new MarketTypes structure
func newMarketTypes() *MarketTypes {
	return &MarketTypes{make([]*MarketType, 0)}
}

// Defines a Root of all possible resources on this API
type Root struct {
	// Provides a list of canonical resources and their URL root
	Resources map[string]string
}
