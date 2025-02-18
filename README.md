# Network Exporter

## Metrics
* 

## Networks
* 

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

