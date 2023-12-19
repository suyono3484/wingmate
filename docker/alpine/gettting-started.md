# Getting Started

_Note: This document is for Alpine base image user, if you are using Debian/Ubuntu read [this one](../bookworm/gettting-started.md) instead._

## Binaries

There are three binaries in in this project: `wingmate`, `wmpidproxy`, and `wmexec`.

`wingmate` is the core binary. It reads config, starts, restarts services. It also
runs cron. Read [here](#wingmate-core-binary) for further details about `wingmate`.

`wmpidproxy` is a helper binary for monitoring _legacy style_ service (fork, exit 
initial proces, and continue in background). Read [here](#wingmate-pid-proxy-binary)
for further details about `wmpidproxy`.

`wmexec` is a helper binary for running process in different `user` or `group`.
It also useful for setting the process as process group leader.
Read [here](#wingmate-exec-binary) for further details about `wmexec`.

## Building your image based on wingmate image in Docker Hub

```Dockerfile
FROM suyono/wingmate:alpine as source

FROM alpine:latest
ADD --chmod=755 wingmate/ /etc/wingmate/
ADD --chmod=755 entry.sh /usr/local/bin/entry.sh
COPY --from=source /usr/local/bin/wingmate /usr/local/bin/wingmate
COPY --from=source /usr/local/bin/wmpidproxy /usr/local/bin/wmpidproxy
COPY --from=source /usr/local/bin/wmexec /usr/local/bin/wmexec
ENTRYPOINT [ "/usr/local/bin/entry.sh" ]
CMD [ "/usr/local/bin/wingmate" ]
```


## Configuration

```shell
/etc
 └── wingmate
     ├── crontab
     ├── crontab.d
     │   ├── cron1.sh
     │   ├── cron2.sh
     │   └── cron3.sh
     └── service
         ├── one.sh
         └── spawner.sh
```

## Appendix
### Wingmate core binary
### Wingmate PID Proxy binary
### Wingmate Exec binary