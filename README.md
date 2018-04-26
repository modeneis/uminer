# UMiner

This Miner is intended to be a GPU miner that works well for a number of handpicked currencies and pools.

All available opencl capable GPU's are detected and used in parallel.

## Binary releases


## Installation from source

### Prerequisites
* go 1.9
* dep


## Build
```
    make build
```

## Start
```
    make start
```

## Test
```
    make test
```

Usage:
```
./uminer --help
Usage:
  uminer [OPTIONS]

Application Options:
  -t, --treads=    CPU threads count for specified currency
  -i, --intensity= GPU mining intensity (NVidia only) (values range: 1..4. Recommended: 2) (default: 28)
  -v, --version    Display version and exit

Help Options:
  -h, --help       Show this help message
```

See what intensity gives you the best hashrate, increasing the intensity also increases the stale rate though.
