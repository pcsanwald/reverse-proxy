# reverse-proxy

## Installation instructions

These instructions describe how to set up the reverse proxy to run on your local machine.
The reverse proxy can sit in front of any web server.

### Prerequisites

1. golang should be installed for your system.
2. nginx or equivalent http server should be installed for your system. 
3. `git` should be installed on your system.

### Install and Run

First, we'll get a copy of the project and ensure the tests pass.

1. `git clone git@github.com:pcsanwald/reverse-proxy.git; cd reverse-proxy` the reverse-proxy repo
2. `cd reverse-proxy`
3. run the tests: `go test`

Next, we'll startup a backend HTTP server, start up our reverse proxy, and pass a `curl` command
to the proxy, to see a request being proxied and then served.

1. start `nginx` on a given port, we'll use 8080 for this example.
2. edit `config.json` to point to the server and port nginx is running on, localhost:8080 is default.
3. start reverse proxy, listening on port 9090: `go run`. Logs will appear in this window.
4. in another window, `curl -XGET localhost:9090`


