package crestmarket

import (
	"strconv"
	"time"
)

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

// TODO: Build a map indexed by something useful
func (r *Regions) ById(id int) *Region {
	for _, region := range r.AllRegions {
		if region.Id == id {
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

func (r *MarketTypes) ById(id int) *MarketType {
	for _, mtype := range r.Types {
		if mtype.Id == id {
			return mtype
		}
	}
	return nil
}

// Build a new MarketTypes structure
func newMarketTypes() *MarketTypes {
	return &MarketTypes{make([]*MarketType, 0)}
}

// A station
type Station struct {
	Name string
	Href string
	Id   int
}

type MarketOrder struct {
	Buy       bool
	Duration  int
	Href      string
	Id        int
	Issued    time.Time
	Station   Station
	MinVolume int
	Price     float64
	Range     string
	Type      MarketType
	Volume    int
}

// Numericrange returns the classical numeric range key
// based on the string input/
func (order *MarketOrder) NumericRange() int {
	orderRange := 0
	if order.Range == "solarsystem" {
		orderRange = 0
	} else if order.Range == "region" {
		orderRange = 65535
	} else if order.Range == "station" {
		orderRange = -1
	} else {
		or, _ := strconv.ParseInt(order.Range, 10, 64)
		orderRange = int(or)
	}
	return orderRange

}

type MarketOrders struct {
	Region  *Region
	Type    *MarketType
	Orders  []*MarketOrder
	Fetched time.Time
}

// Make a new MarketOrders structure
func NewMarketOrders() *MarketOrders {
	return &MarketOrders{Orders: make([]*MarketOrder, 0)}
}

// Defines a Root of all possible resources on this API
type Root struct {
	// Provides a list of canonical resources and their URL root
	Resources map[string]string
}
