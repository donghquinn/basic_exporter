# Basic Exporter
* It's very simple but basic exporter for monitoring your own machine.
* My purpose of this exporter is building fast, ligth, and stable exporter.
* It provides System and Network metrics so far. Planning to add more functions and give users more options for choosing what they want.

## Metrics
* It's available to read metrics by :9468/metrics
    * The port is fixed with 9468

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
