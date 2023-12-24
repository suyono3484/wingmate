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

## Building a container image based on wingmate image in Docker Hub

Wingmate has no dependency other than `alpine` base image, so you just need to copy
the binaries directly. If you have built your application into an `alpine` based image,
all you need to do is copy whichever binary you need, crontab file (if you use cron)
and add some shell script to glue them together. Here is a Dockerfile example.

```Dockerfile
# Dockerfile
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
You can find some examples for shell script in this [directory](../alpine/).

## Configuration

When `wingmate` binary starts, it will look for some files. By default, it will
try to read the content of `/etc/wingmate` directory. You can change the directory
where it reads by setting `WINGMATE_CONFIG_PATH` environment variable. The structure
inside the config path should look like this.

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

First, `wingmate` will try to read the content of `service` directory. The content of
this directory should be executables (either shell scripts or binaries). The wingmate
will run every executable in `service` directory without going into any subdirectory.

Next, `wingmate` will read the `crontab` file. `wingmate` expects the `crontab` file using
common UNIX crontab file format. Something like this

```shell
 ┌───────────── minute (0–59)
 │ ┌───────────── hour (0–23)
 │ │ ┌───────────── day of the month (1–31)
 │ │ │ ┌───────────── month (1–12)
 │ │ │ │ ┌───────────── day of the week (0–6) (Sunday to Saturday)
 │ │ │ │ │ 
 │ │ │ │ │
 │ │ │ │ │
 * * * * * <commad or shell script or binary>
```

The command part only support simple command and arguments. Shell expression is not supported
yet. It is recommended to write a shell script and put the path to shell script in
the command part.

# Appendix
## Wingmate core binary
### Service
### Cron
## Wingmate PID Proxy binary
## Wingmate Exec binary