# beacon
Random beacon simulator

## Build

`docker pull dfinity/build-env`

`git clone git@github.com:timohanke/rb.git`

`cd rb`

```docker run --rm -it -v `pwd`/go:/go dfinity/build-env go run sim.go```

### Run linter(s)

```docker run --rm -it -v `pwd`/go:/go dfinity/build-env gometalinter ./...```

As of Jan 22 (EOD)
```gometalinter --deadline=10s --disable=gas --disable=gocyclo --disable=gotype --disable=golint --disable=interfacer ./...```
runs without warnings.

## Run

### Simulation

`go run sim.go`

### Cgo test and benchmark

`go run test.go`
 
## Dependencies

All dependencies below are taken care off in the docker image `dfinity/build-env`.

### go-ethereum

Code currently depends on `github.com/ethereum/go-ethereum/common` being present in the `src` directory.
This will be removed.

### bls command line tool

Code currently calls the command line tool `bls_tool.exe` compiled from Shigeo's code. This is to be replaced by Shigeo's Cgo code.

The build instructions are here: `https://github.com/herumi/bls` -> readme.md -> Installation Requirements

That also requires `apt-get install libssl-dev`

### cgo bindings

For cgo, which is transitioning in, we need the environment variables set:

`export LIBRARY_PATH=/build/herumi/bls/lib:/build/herumi/mcl/lib:$LIBRARY_PATH`

`export CPATH=/build/herumi/bls/include:$CPATH`
 
