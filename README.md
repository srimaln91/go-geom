A Go wrapper for [libgeos](https://geos.osgeo.org/)

## Install geos library on Ubuntu

```bash
sudo apt-get install libgeos-dev libgeos-3.7.0
```

## Quick Start

```go
import "github.com/srimaln91/geos-go/geos"

geom := geos.FromWKT("POINT (0 0)")
defer geom.Destroy()

//set SRID
geom.SetSRID(4326)

//Create a buffer around geometry (2 radians)
geom.Buffer(2)
wktString := geom.ToWKT()

```
