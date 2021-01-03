Golang SSH Tarpit
===================

A simple SSH tarpit written in Go.


## Background

There wasn't really a good reason to write yet-another-tarpit other than that I felt like doing it.

I'd originally planned to do it in Python, but then changed my mind and figured I'd give Go a go instead.

The basic concept, is you run this tarpit on `TCP 22` (the default for SSH) in order to inconvenience bad bots and tie up their resources (rather than moving onto spamming other boxes) and run your actual SSH on another port.

----

## Build/Usage

### Docker

You can run with docker, the image uses the default port `2222` so when running the image, just map it across to `22`

    docker run -d -p 22:2222 bentasker12/go_ssh_tarpit

This will fetch it from [Docker Hub](https://hub.docker.com/repository/docker/bentasker12/go_ssh_tarpit)

### Manual

If you'd rather not use docker, you just need to build it with `Go`

    go build ssh_tarpit.go

And then run it

    ./ssh_tarpit.go

Or, of course, you can use `go run`

    go run ssh_tarpit.go

However, by default the script binds to port 2222 - this is so that it could easily run as a non-privileged user within the docker container. If you're running directly, you have 2 options

* Edit the constant to bind to `22` and run as root (*very, very, very* bad idea)
* Run as unprivileged user and use IPTables to NAT `22` to `2222`

The latter can be achieved with

    iptables -t nat -A PREROUTING -p tcp --dport 22 -j REDIRECT --to-port 2222

----

### Copyright

Golang SSH Tarpit is Copyright (C) 2021 B Tasker. All Rights Reserved. 

Released Under [GNU GPL V3 License](http://www.gnu.org/licenses/gpl-3.0.txt).
