# dprobe
Docker security auditing tool

## Overview
A tool to audit underlying docker host and containers. The audit information comes from CIS benchmarks for Docker.

## Features
- Output to slack
- **New** All output is now sent as JSON
- **New** When sending to slack, the audit data is sent as a snippet

### Docker Daemon/Host
- Compare running version with current stable
- Check live restore
- Docker file/directory ownership
- Docker file/directory mode
- Container sprawl
- Image sprawl
- If running ECS, detect the version and cluster
- Gathers basic host info (IPs, Instance ID (aws), and hostname)

### Docker Containers
- Check if they are `Privileged`
- Check if they have `Capabilities`
- Ensure memory limits are set per container
- Each container should have HEALTHCHECK
- Mount propagation
- Check if container is using ports < 1024
- Do not allow containers access to hosts process namespace
- Do not allow containers access to hosts IPC namespace
- Do not allow containers access to hosts UTS namespace
- Do not allow containers access to hosts devices

## Usage
1. Download a release binary.
2. Before running figure out the Docker API version you're running: `docker version` will tell you.
3. To run the audit: `env DOCKER_API_VERSION="x.xx" ./dprobe` (note: you most likely need to `sudo` to run at all [to access the docker socket and to stat files, etc.])

### Output
- `--output`/`-o` supports `slack` or `stdout`; you must configure `dprobe.json` to output to slack.
- `--csprawl`/`-c` to set container sprawl amount.
- `--isprawl`/`-i` to set image sprawl amount.

## Todo
- Check if containers are running outdated image
- Compare ECS version to stable
- Disable swarm mode if it's not in use
- stdout colorizing
- Link CIS benchmark numbers?
