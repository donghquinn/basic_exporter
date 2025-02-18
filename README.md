# Network Exporter

## System
* Memory
    * Total Physical Memroy in bytes
    * Used memory in bytes
    * Available Memory in bytes
    * Memory usage in bytes
* Load Average
    * 1, 5, 15
* CPU Usage
    * Total CPU Usage Percentage

## Networks
* Network Total Received & Transmitted bytes
* Current Network Received & Transmitted Bandwidth in bytes per second.

## Build and Deploy
### Docker and Docker compose
* use docker-compose.yml with building using Dockerfile

### Compile binary executable file
* Use build.sh. This script is the example build command script for ubuntu. customize yourself.

## Dependencies
### gopsutil

* Ported library from python psutil for retrieving information on running proccesses and system utilizations.
    * CPU, memory, disks, networks, sensors, etc...
    * Cross Platform available

```zsh
go get -u github.com/shirou/gopsutil/v4
```

### prometheus client
* Go client for Prometheus.

```zsh
go get -u github.com/prometheus/client_golang
```

### Go Env Loader

```zsh
go get -u github.com/joho/godotenv
```

