package geocache

import (
	"testing"
	"time"
)

type TestStruct struct {
	ID    int
	Prop1 string
	Prop2 float64
}

func TestCache(t *testing.T) {
	cache, err := NewCache(5*time.Minute, 30*time.Second, WithIn1KM)

	if err != nil {
		t.Error("An error occured while creating cache: ", err.Error())
	}

	point := GeoPoint{Latitude: 41.234, Longitude: 29.432}

	value, found := cache.Get(point)

	if found {
		t.Errorf("Point: %v should be nil and found should be false", value)
	}

	item := TestStruct{ID: 1, Prop1: "hello", Prop2: 3.14}

	cache.Set(point, item, 10*time.Minute)

	cachedItem, found := cache.Get(point)

	if !found || cachedItem == nil {
		t.Errorf("Item: %v shouln't be nil", item)
		t.Errorf("Found expected: %t got: %t", true, found)
	}
}

func TestItemCount(t *testing.T) {
	cache, err := NewCache(5*time.Minute, 30*time.Second, WithIn1KM)

	if err != nil {
		t.Error("An error occured while creating cache: ", err.Error())
	}

	point := GeoPoint{Latitude: 41.234, Longitude: 29.432}
	item := TestStruct{ID: 1, Prop1: "hello", Prop2: 3.14}
	cache.Set(point, item, 10*time.Minute)

	itemCount := cache.ItemCount()

	if itemCount != 1 {
		t.Errorf("Item count should be: %d , got: %d", 1, itemCount)
	}
}

func TestFlush(t *testing.T) {
	cache, err := NewCache(5*time.Minute, 30*time.Second, WithIn1KM)

	if err != nil {
		t.Error("An error occured while creating cache: ", err.Error())
	}

	point := GeoPoint{Latitude: 41.234, Longitude: 29.432}
	item := TestStruct{ID: 1, Prop1: "hello", Prop2: 3.14}
	cache.Set(point, item, 10*time.Minute)

	itemCount := cache.ItemCount()

	if itemCount != 1 {
		t.Errorf("Item count should be: %d , got: %d", 1, itemCount)
	}

	cache.Flush()

	itemCountAfterFlush := cache.ItemCount()

	if itemCountAfterFlush != 0 {
		t.Errorf("Item count should be: %d , got: %d", 0, itemCountAfterFlush)
	}
}

func TestItems(t *testing.T) {
	cache, err := NewCache(5*time.Minute, 30*time.Second, WithIn1KM)

	if err != nil {
		t.Error("An error occured while creating cache: ", err.Error())
	}

	point1 := GeoPoint{Latitude: 41.234, Longitude: 29.432}
	item1 := TestStruct{ID: 1, Prop1: "ping", Prop2: 3.141}

	point2 := GeoPoint{Latitude: 41.25221, Longitude: 29.44421}
	item2 := TestStruct{ID: 2, Prop1: "pong", Prop2: 3.142}
	cache.Set(point1, item1, 10*time.Minute)
	cache.Set(point2, item2, 10*time.Minute)

	items := cache.Items()

	if len(items) != 2 {
		t.Errorf("Item count should be: %d, got: %d", 2, len(items))
	}
}
