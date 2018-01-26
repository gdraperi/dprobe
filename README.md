# dprobe
Docker security auditing tool

## Security features
- TODO: Compare running version with current stable
- TODO: Check live restore
- Audit running containers:
    - Check if they are `Privileged`
    - Check if they have `Capabilities`
    - Check if they are running outdated image
    - Ensure memory limits are set
- Restrict host access:
    - Do not allow containers access to hosts process namespace
    - Do not allow containers access to hosts IPC namespace
    - Do not allow containers access to hosts UTS namespace
    - Do not allow containers access to hosts devices
- Docker daemon audit:
    - Various directories and configuration
    - Confirm cgroups configuration
    - Disable swarm mode if it's not in use
- TODO: Container/Image sprawl (audit the amount of abandoned/orphaned containers/images on the host)
- ECS agent:
    - Version check
- Image audit:
    - Ensure they do not have unnecessary packages
    - Ensure they are not running SSH
    - Use COPY instead of ADD
    - Use HEALTHCHECK
    - Do not store secrets in Dockerfiles
    - Ensure images are rebuilt with security patches
        - We use ubuntu:trusty (and others); these images need to be rebuilt when there is a security update related to the base image
