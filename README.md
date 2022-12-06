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
- https://localhost:1313/gw/<CID>
- https://localhost:1313/gw/file/<CID>
- https://localhost:1313/gw/dir/<CID>

# Live Demo
- https://whypfs-gateway.onrender.com/gw 

# Serve files
![image](https://user-images.githubusercontent.com/4479171/205086971-5b3a67ae-3ac3-42f9-961a-0ef22fae5f32.png)

# Serve Directories
![image](https://user-images.githubusercontent.com/4479171/205777409-45045f54-1c4a-4373-a7ee-d38440dd9a3e.png)
