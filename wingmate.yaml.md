YAML Configuration
---

Table of content
- [Service](#service)
  - [Command](#command)
  - [Environ](#environ)
  - [User and Group](#user-and-group)
  - [Working Directory](#working-directory)
  - [setsid](#setsid)
  - [PID File](#pid-file)
- [Cron](#cron)
  - [Schedule](#schedule)

Example
```yaml
service:
  spawner:
    command: [ wmspawner ]
    user: "1200"
    working_dir: "/var/run/test"

  bgtest:
    command:
      - "wmstarter"
      - "--no-wait"
      - "--"
      - "wmexec"
      - "--setsid"
      - "--"
      - "wmbg"
      - "--name"
      - "test-run"
      - "--pause"
      - "10"
      - "--log-path"
      - "/var/log/wmbg.log"
      - "--pid-file"
      - "/var/run/wmbg.pid"
    pidfile: "/var/run/wmbg.pid"

cron:
  cron1:
    command: ["wmoneshot", "--", "sleep", "5"]
    schedule: "*/5 * * * *"
    working_dir: "/var/run/cron"
    environ:
      - "WINGMATE_LOG=/var/log/cron1.log"
      - "WINGMATE_LOG_MESSAGE=cron executed in minute 5,10,15,20,25,30,35,40,45,50,55"
  cron2:
    command: ["wmoneshot", "--", "sleep", "5"]
    schedule: "17,42 */2 * * *"
    environ:
      - "WINGMATE_LOG=/var/log/cron2.log"
      - "WINGMATE_LOG_MESSAGE=cron scheduled using 17,42 */2 * * *"
  cron3:
    command:
      - "wmoneshot"
      - "--"
      - "sleep"
      - "5"
    schedule: "7,19,23,47 22 * * *"
    environ:
      - "WINGMATE_LOG=/var/log/cron3.log"
      - "WINGMATE_LOG_MESSAGE=cron scheduled using 7,19,23,47 22 * * *"

```

At the top-level, there are two possible entries: Service and Cron.

## Service

`service` is a top-level element that hosts the definition of services to be started by `wingmate`.

Example
```yaml
service:
  svc1:
    command: [ some_executable ]
    user: "1200"
    working_dir: "/var/run/test"
```

In the example above, we declare a service called `svc1`. `wingmate` will start a process based on all
elements defined under `svc1`. To learn more about elements for a service, read below.

### Command

`command` element is an array of strings consists of an executable name (optionally with path) and
its arguments (if any). `wingmate` will start the service as its child process by executing
the executable with its arguments.

Example

```yaml
command: [ executable1, argument1, argument2 ]
```

Based on YAML standard, the above example can also be written like

```yaml
command:
  - executable1
  - argument1
  - argument2
```

### Environ

`environ` element is an array of strings. It is a list of environment variables `wingmate` will pass to
the child process or service. The format of each environment variable is a pair of key and value separated
by `=` sign. By default, the child process or service will inherit all environment variables of its parent.

Example

```yaml
environ:
  - "S3_BUCKET=YOURS3BUCKET"
  - "SECRET_KEY=YOUR_SECRET_KEY"
```

Note: don't worry if an environment variable value has one or more `=` character(s) in it. `wingmate` will
separate key and value using the first `=` character only.

### Working Directory

`working_dir` is a string contains the path where the child process will be running in. By default, the child
process will run in the `wingmate` current directory.

### User and Group

Both `user` and `group` take string value. `user` and `group` refer to the operating system's user and group.
They can be in the form of name, like username or groupname, or in the form of id, like uid or gid.
If they are set, the child process will run as the specified user and group. By default, the child process
will run as the same user and group as the `wingmate` process. The `user` and `group` are only effective
when the `wingmate` running as privileged user, such as `root`. The `user` and `group` configuration depends
on the [wmexec](README.md#wingmate-exec-binary).

### setsid

`setsid` takes a boolean value, `true` or `false`. This feature is operating system dependant. If set to `true`,
the child process will run in a new session. Read `man setsid` on Linux/UNIX. The `setsid` configuration depends
on the [wmexec](README.md#wingmate-exec-binary).

### PID File

This feature is designated to handle service that run in the background. This kind of service usually forks a
new process, terminate the parent process, and continue running in the background child process. It writes its
background process PID in a file. This file is referred as PID file. Put the path of the PID file to this
`pidfile` element. It will help `wingmate` to restart the service if its process exited / terminated. The `pidfile`
configuration depends on the [wmpidproxy](README.md#wingmate-pid-proxy-binary).

## Cron

`cron` is a top-level element that hosts the definition of crones to run by `wingmate` on the specified schedule.
Cron shares almost all configuration elements with Service, except `schedule` and `pidfile`. For the following
elements, please refer to the [Service](#service) section

- [Command](#command)
- [Environ](#environ)
- [Working Directory](#working-directory)
- [setsid](#setsid)
- [User and Group](#user-and-group)

`pidfile` is an invalid config parameter for cron because `wingmate` cannot start cron in background mode. This
limitation is intentionally built into `wingmate` because it doesn't make any sense to run a periodic cron process
in background.

### Schedule

The schedule configuration field uses a format similar to the one described in the [README.md](README.md).

```shell
 ┌───────────── minute (0–59)
 │ ┌───────────── hour (0–23)
 │ │ ┌───────────── day of the month (1–31)
 │ │ │ ┌───────────── month (1–12)
 │ │ │ │ ┌───────────── day of the week (0–6) (Sunday to Saturday)
 │ │ │ │ │ 
 │ │ │ │ │
 │ │ │ │ │
 * * * * *
```
