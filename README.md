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

# Live Demo
- https://whypfs-gateway.onrender.com/gw 

# Serve files
https://gateway.estuary.tech/gw/ipfs/bafybeibpkuvcuatbkt4s6pvr46uc7flbwp53bmryypssqsuob55oznt5fu
![image](https://user-images.githubusercontent.com/4479171/206327573-0d2bdf75-723c-4d15-a52a-522f04fb0991.png)

# Serve Dirs 
https://gateway.estuary.tech/gw/ipfs/QmPBHAjRLZqvJwcBUTiVxNtvugToAnTyJxpzTCgKZVHsvw
![image](https://user-images.githubusercontent.com/4479171/206327483-0a939510-ac5a-408c-8773-cbb9ae72d7ff.png)
