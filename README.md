# A Go wrapper for [libgeos](https://geos.osgeo.org/)

## Install geos library on Ubuntu

### Using package manager

```bash
# Add ubuntugis repository
sudo apt-add-repository -y ppa:ubuntugis/ubuntugis-unstable
sudo apt-get install libgeos-dev libgeos-3.7.0
```

### Compile from source

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
