# geocache [![GoDoc](https://godoc.org/github.com/melihmucuk/geocache?status.svg)](https://godoc.org/github.com/melihmucuk/geocache) [![Go Report Card](https://goreportcard.com/badge/melihmucuk/geocache)](https://goreportcard.com/report/melihmucuk/geocache) [![Build Status](http://img.shields.io/travis/melihmucuk/geocache.svg?style=flat)](https://travis-ci.org/melihmucuk/geocache)

geocache is an in-memory cache that is suitable for geolocation based applications. It uses geolocation as a key for storing items. You can specify range on initialization and thats it! You can store any object, it uses interface.

### Installation

`go get github.com/melihmucuk/geocache`

### Usage

![geolocation cache](http://i.imgur.com/O6UzVEW.png "Geolocation Cache")

```go

import (
	"fmt"
	"time"

	"github.com/melihmucuk/geocache"
)

func main() {
	c, err := geocache.NewCache(5*time.Minute, 30*time.Second, geocache.WithIn1KM)
	geoPoint := geocache.GeoPoint{Latitude: 40.9887, Longitude: 28.7817}
	if err != nil {
		fmt.Println("Error: ", err.Error())
	} else {
		c.Set(geoPoint, "helloooo", 2*time.Minute)
		v1, ok1 := c.Get(geocache.GeoPoint{Latitude: 41.2, Longitude: 29.3})
		v2, ok2 := c.Get(geocache.GeoPoint{Latitude: 41.2142, Longitude: 29.4234})
		v3, ok3 := c.Get(geocache.GeoPoint{Latitude: 40.9858, Longitude: 28.7852})
		v4, ok4 := c.Get(geocache.GeoPoint{Latitude: 40.9827, Longitude: 28.7883})
		fmt.Println(v1, ok1)
		fmt.Println(v2, ok2)
		fmt.Println(v3, ok3)
		fmt.Println(v4, ok4)
	}
}

```

outputs:
```
<nil>, false
<nil>, false
helloooo, true
helloooo, true
```

### Information

You can specify 8 different range. More info can be found [here](http://gis.stackexchange.com/questions/8650/how-to-measure-the-accuracy-of-latitude-and-longitude).

* `WithIn11KM`

The first decimal place is worth up to 11.1 km `eg: 41.3, 29.6`

* `WithIn1KM`

The second decimal place is worth up to 1.1 km `eg: 41.36, 29.63`

* `WithIn110M`

The third decimal place is worth up to 110 m `eg: 41.367, 29.631`

* `WithIn11M`

The fourth decimal place is worth up to 11 m `eg: 41.3674, 29.6316`

* `WithIn1M`

The fifth decimal place is worth up to 1.1 m `eg: 41.36742, 29.63168`

* `WithIn11CM`

The sixth decimal place is worth up to 0.11 m `eg: 41.367421, 29.631689`

* `WithIn11MM`

The seventh decimal place is worth up to 11 mm `eg: 41.3674211, 29.6316893`

* `WithIn1MM`

The eighth decimal place is worth up to 1.1 mm `eg: 41.36742115, 29.63168932`
