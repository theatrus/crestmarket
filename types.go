package crestmarket

import (
	"strconv"
	"sync"
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
	byId       map[int]*Region
	byName     map[string]*Region
	m          sync.Mutex
}

// Build a new Regions structure
func newRegions() *Regions {
	return &Regions{AllRegions: make([]*Region, 0)}
}

func (r *Regions) ByName(name string) *Region {
	if r.byName == nil {
		r.m.Lock()
		defer r.m.Unlock()

		r.byName = make(map[string]*Region)

		for _, region := range r.AllRegions {
			r.byName[region.Name] = region
		}
	}

	return r.byName[name]
}

func (r *Regions) ById(id int) *Region {
	if r.byId == nil {
		r.m.Lock()
		defer r.m.Unlock()

		r.byId = make(map[int]*Region)
		for _, region := range r.AllRegions {
			r.byId[region.Id] = region
		}
	}
	return r.byId[id]
}

// An inventory type
type MarketType struct {
	Name string
	Href string
	Id   int
}

// A collection of inventory types
type MarketTypes struct {
	Types  []*MarketType
	byName map[string]*MarketType
	byId   map[int]*MarketType
	m      sync.Mutex
}

func (r *MarketTypes) ByName(name string) *MarketType {
	if r.byName == nil {
		r.m.Lock()
		defer r.m.Unlock()

		r.byName = make(map[string]*MarketType)
		for _, mtype := range r.Types {
			r.byName[mtype.Name] = mtype
		}
	}
	return r.byName[name]
}

func (r *MarketTypes) ById(id int) *MarketType {
	if r.byId == nil {
		r.m.Lock()
		defer r.m.Unlock()

		r.byId = make(map[int]*MarketType)

		for _, mtype := range r.Types {
			r.byId[mtype.Id] = mtype
		}
	}
	return r.byId[id]
}

// Build a new MarketTypes structure
func newMarketTypes() *MarketTypes {
	return &MarketTypes{Types: make([]*MarketType, 0)}
}

// A station
type Station struct {
	Name string
	Href string
	Id   int
}

type MarketOrder struct {
	Bid           bool
	Duration      int
	Href          string
	Id            int
	Issued        time.Time
	Station       Station
	MinVolume     int
	VolumeEntered int
	Price         float64
	Range         string
	Type          MarketType
	Volume        int
}

const (
	RangeSolarsystem = 0
	RangeRegion      = 65535
	RangeStation     = -1
)

// Numericrange returns the classical numeric range key
// based on the string input/
func (order *MarketOrder) NumericRange() int {
	orderRange := 0
	if order.Range == "solarsystem" {
		orderRange = RangeSolarsystem
	} else if order.Range == "region" {
		orderRange = RangeRegion
	} else if order.Range == "station" {
		orderRange = RangeStation
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
