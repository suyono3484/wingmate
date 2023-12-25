# wingmate

Wingmate is a process manager for services. It works like init. It starts and restarts services.
It also has cron feature. It is designed to run in a container/docker. 
The Wingmate binary do not need any external dependency. 
Just copy the binary, and exec from the entry point script.

# Getting Started

## Binaries

There are three binaries in this project: `wingmate`, `wmpidproxy`, and `wmexec`.

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
You can find some examples for shell script in [alpine docker](docker/alpine) and 
[bookworm docker](docker/bookworm).

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
## Wingmate PID Proxy binary

`wingmate` works by monitoring its direct children process. When it sees one of its children
exited, it will start the child process again.

Sometimes you find some services work by running in the background. It means it forks a new 
process, disconnect the new child from terminal, exit the parent process, and continue 
running in the child process. This kind of service usually write its background process
PID in a pid file.

To monitor the background services, `wingmate` utilizes `wmpidproxy`. `wmpidproxy` runs
in foreground in-place of the background service. It also periodically check whether the
background service is still running, in current implementation it checks every second.

```shell
wmpidproxy --pid-file <path to pid file> -- <background service binary/start script>
```
#### Example
Running sshd background with `wingmate` and `wmpidproxy`: [here](example/ssh-docker)

## Wingmate Exec binary


## Wingmate core binary
### Service
### Cron
