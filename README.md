# WhyPFS Gateway 

Quick HTTP gateway using whypfs-core.  

This is experimental and not recommended for production use. 

## Installation
```
go mod tidy
go mod download
```
## Setup
```
go build -tags netgo -ldflags '-s -w' -o whypfs-gw
./whypfs-gw
```

# Test
https://localhost:1313/gw/ipfs/<CID>

