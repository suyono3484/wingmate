package testconfig

var One = `service:
  one:
    command: "mycommand -o output"
  two:
    command: ["cmd", "now"]
    workdir: /
`
