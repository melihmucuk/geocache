package geocache

import (
	"errors"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// Range specifies range of cache
type Range int32

const (
	// WithIn11KM The first decimal place is worth up to 11.1 km
	// eg: 41.3, 29.6
	WithIn11KM Range = 1 + iota

	// WithIn1KM The second decimal place is worth up to 1.1 km
	// eg: 41.36, 29.63
	WithIn1KM

	// WithIn110M The third decimal place is worth up to 110 m
	// eg: 41.367, 29.631
	WithIn110M

	// WithIn11M The fourth decimal place is worth up to 11 m
	// eg: 41.3674, 29.6316
	WithIn11M

	// WithIn1M The fifth decimal place is worth up to 1.1 m
	// eg: 41.36742, 29.63168
	WithIn1M

	// WithIn11CM The sixth decimal place is worth up to 0.11 m
	// eg: 41.367421, 29.631689
	WithIn11CM

	// WithIn11MM The seventh decimal place is worth up to 11 mm
	// eg: 41.3674211, 29.6316893
	WithIn11MM

	// WithIn1MM The eighth decimal place is worth up to 1.1 mm
	// eg: 41.36742115, 29.63168932
	WithIn1MM
)

// Item struct keeps cache value and expiration time of object
type Item struct {
	Object     interface{}
	Expiration int64
}

// Cache struct manages items, expirations and clean ups
type Cache struct {
	items           map[GeoPoint]Item
	m               sync.RWMutex
	expiration      time.Duration
	cleanUpInterval time.Duration
	precision       int32
	stopCleanUp     chan bool
}

// GeoPoint specifies point that used as key of cache
type GeoPoint struct {
	Latitude  float64
	Longitude float64
}

// NewCache creates new Cache with params and returns pointer of Cache and error
// cleanUpInterval used for deleting expired objects from cache.
func NewCache(expiration, cleanUpInterval time.Duration, withInRange Range) (*Cache, error) {
	if withInRange < 1 || withInRange > 8 {
		return nil, errors.New("Range must be within 1-8!")
	}

	c := &Cache{items: make(map[GeoPoint]Item), expiration: expiration, cleanUpInterval: cleanUpInterval, precision: int32(withInRange), stopCleanUp: make(chan bool)}

	ticker := time.NewTicker(cleanUpInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.cleanUp()
			case <-c.stopCleanUp:
				ticker.Stop()
				return
			}
		}
	}()

	return c, nil
}

// Set adds object to cache with given geopoint
func (c *Cache) Set(position GeoPoint, value interface{}, expiration time.Duration) {
	var exp int64
	if expiration == 0 {
		expiration = c.expiration
	}

	if expiration > 0 {
		exp = time.Now().Add(expiration).UnixNano()
	}
	c.m.Lock()
	defer c.m.Unlock()
	c.items[position.truncate(c.precision)] = Item{Object: value, Expiration: exp}
}

// Get gets object from cache with given geopoint
func (c *Cache) Get(position GeoPoint) (interface{}, bool) {
	c.m.RLock()
	defer c.m.RUnlock()
	item, found := c.items[position.truncate(c.precision)]
	return item.Object, found
}

// Items returns cached items
func (c *Cache) Items() map[GeoPoint]Item {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.items
}

// ItemCount returns cached items count
func (c *Cache) ItemCount() int {
	c.m.RLock()
	defer c.m.RUnlock()
	n := len(c.items)
	return n
}

// Flush deletes all cached items
func (c *Cache) Flush() {
	c.m.Lock()
	defer c.m.Unlock()
	c.items = map[GeoPoint]Item{}
}

// StopCleanUp stops clean up process.
func (c *Cache) StopCleanUp() {
	c.stopCleanUp <- true
}

func (c *Cache) cleanUp() {
	c.m.Lock()
	defer c.m.Unlock()
	for k, v := range c.items {
		if v.Expiration < time.Now().UnixNano() {
			delete(c.items, k)
		}
	}
}

func (g GeoPoint) truncate(precision int32) GeoPoint {
	dLat := decimal.NewFromFloat(g.Latitude)
	dLong := decimal.NewFromFloat(g.Longitude)
	lat, _ := dLat.Truncate(precision).Float64()
	long, _ := dLong.Truncate(precision).Float64()

	return GeoPoint{
		Latitude:  lat,
		Longitude: long,
	}
}
