# A Go library for geospatial operations

The library is basically written as a wrapper around **liblwgeom** which comes with PostGIS. Liblwgeom uses following libraries in order to perform geospatial calculations.

- [GEOS](https://geos.osgeo.org/)
- [liblwgeom](https://github.com/postgis/postgis/tree/svn-trunk/liblwgeom)
- [proj.4](https://proj4.org/)

The library containes some extra functions which uses **libgeos** directly.

[![Build Status](https://travis-ci.org/srimaln91/go-geom.svg?branch=master)](https://travis-ci.org/srimaln91/go-geom)
[![codecov](https://codecov.io/gh/srimaln91/go-geom/branch/master/graph/badge.svg)](https://codecov.io/gh/srimaln91/go-geom)
[![Go Report Card](https://goreportcard.com/badge/github.com/srimaln91/go-geom)](https://goreportcard.com/report/github.com/srimaln91/go-geom)
[![GoDoc](https://godoc.org/github.com/srimaln91/go-geom/geos?status.svg)](https://godoc.org/github.com/srimaln91/go-geom/geom)

## Install required libraries on Ubuntu

### Using package manager

```bash
# Add ubuntugis repository
sudo apt-add-repository -y ppa:ubuntugis/ubuntugis-unstable
sudo apt-get install libgeos-dev libgeos-3.7.0 liblwgeom-2.5-0 liblwgeom-dev libproj-dev
```

### Compile GEOS from source

You can download the source code from [https://trac.osgeo.org/geos/](https://trac.osgeo.org/geos/).

More on: [https://trac.osgeo.org/geos#BuildandInstall](https://trac.osgeo.org/geos#BuildandInstall)

Building GEOS requires a C++11 compiler

```bash
# Install gcc and cmake
sudo apt install gcc
sudo apt install cmake
```

Below are the compilation steps

```bash
# Commpilation
wget http://download.osgeo.org/geos/geos-3.7.1.tar.bz2
tar xvfj geos-3.7.1.tar.bz2
cd geos-3.7.1
mkdir build && cd build
cmake ..
make
make check

# Installation. Should be executed with superuser privileges
sudo make install

# Create necessary links
sudo ldconfig
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/srimaln91/go-geom/geom"
)

func main() {

	jsonLineString := `{
    "type": "LineString",
    "coordinates": [
            [79.86064,6.933669],
            [79.87326,6.923529],
            [79.87017,6.910237],
            [79.88305,6.899671],
            [79.89008,6.890724],
            [79.89051,6.882118],
            [79.88854,6.869762],
            [79.87532,6.863541]
        ]
    }`

    lwgeom, err := geom.FromGeoJSON(jsonLineString)
    if err != nil {
        fmt.Println(err)
        return
    }

	lwgeom.SetSRID(4326)
	defer lwgeom.Free()

	lwgeom.LineSubstring(0.5, 0.9)

	fromSRS := geom.SRS["EPSG:4326"]
	toSRS := geom.SRS["EPSG:3857"]

	// Transform the geometry to SRS EPSG:3757 so we can measure from metres.
	lwgeom.Project(fromSRS, toSRS)
	lwgeom.Buffer(200)

	// Reset SRS to EPSG:4326
	lwgeom.Project(toSRS, fromSRS)

	bufJSON, err := lwgeom.ToGeoJSON(4, 0)
    if err != nil {
        fmt.Println(err)
        return
    }
    
	fmt.Println(bufJSON)
}
```
