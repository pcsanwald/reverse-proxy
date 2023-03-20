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

1. Clone and `cd` into repo: `git clone git@github.com:pcsanwald/reverse-proxy.git; cd reverse-proxy` 
2. run the tests: `go test`

Next, we'll startup a backend HTTP server, start up our reverse proxy, and pass a `curl` command
to the proxy, to see a request being proxied and then served.

1. start `nginx` on a given port, we'll use 8080 for this example.
2. edit `config.json` to point to the server and port nginx is running on, localhost:8080 is default.
3. start reverse proxy, listening on port 9090: `go run listener.go`. Logs will appear in this window.
4. in another window, `curl -XGET localhost:9090`

### Assumptions, notes and areas for improvement

* installation instructions assume a sophisticated user: the install process could be simplified quite a lot by providing an image, for example.
* Request blocking is done by checking for existence of Header or Query String parameter, ignoring the value.
* To keep the scope reasonable, I choose to define "sensitive information" as email or phone, and used a [3rd party validation library](https://github.com/nyaruka/phonenumbers) for phone numbers. To detect sensitive information, I only considered the value of the parameter, but detection could be improved by considering the parameter name as well.
* testing a broader range of inputs, using matrix style testing, would be a nice improvement.
* separate log files, and logging configuration, is another area for improvement. 
* I've tried to keep logging in the proxy relatively close to the spec, but I've added logging where I feel it makes the functionality clearer
* I'm intentionally logging sensitive info at the proxy level, prior to masking, but this can be easily removed.
