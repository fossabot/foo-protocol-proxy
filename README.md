Foo Protocol Proxy
==================

A Golang implementation for proxy receiver over Foo protocol.
It also has some extra features like:

- Reporting stats when sending `SIGUSR2` signal to the process and over HTTP.
- Health check over HTTP.
- Data recovery after failure by saving to file, if the failure period does not exceed the `10s` time frame.

Installation
------------

### Prerequisites

* [Golang][1] installation, having [$GOPATH][2] properly set.

To install [**Foo Protocol Proxy**](https://github.com/ahmedkamals/foo-protocol-proxy)

```bash
$ go get github.com/ahmedkamals/foo-protocol-proxy
```

Test Driver
-----------

You can use the following steps as a testing procedure

  * **Server**
    ```bash
    $ bin/server-linux -listen ":8001"
    ```

  * **Proxy**
    ```bash
    $ make build
    $ bin/foo-protocol-proxy-{OS}-{ARCH} -forward "{FORWARDING_PORT}" -listen "{LISTENING_PORT}" -http "{HTTP_ADDRESS}" -recovery-path "{RECOVERY_PATH}"
    ```
    
    **`Environment`**:
       * `OS` - the current operating system, e.g. (linux, darwin, ...etc.)
       * `ARCH` - the current system architecture, e.g. (386, amd64)
        
    **`Params`**:
       * `FORWARDING_PORT` - e.g. `":8081"`
       * `LISTENING_PORT` - e.g. `":8082"`
       * `HTTP_ADDRESS` - e.g. `"0.0.0.0:8088"`
       * `RECOVERY_PATH` - e.g. `"data/recovery.json"`
             
  * **Multiple Client Connections**
    ```bash
    $ for i in {0..1000..1}
      do 
         bin/client-linux -connect "localhost:8002";
      done
    ```

  * **Sending `SIGUSR2` Signal**
      
    ```bash
    $ kill -SIGUSR2 $(pidof foo-protocol-proxy) > /dev/null 2>&1
    ```

## Tests
    
Not all items covered, just made one example.
    
```bash
$ make unit
```

## Coding - __Structure & Design__
* `Dispatcher` - parses configuration, builds configuration object,
and passes it to the proxy.
* `Proxy` - orchestrates the interactions between the components. 
* `Listner` - awaits for client connections, and on every new connection, 
a `BridgeConnection` instance is created.
* `BridgeConnection` - acts as Bi-directional communication object, that
passes data forward and backward to the server.
* `Analyzer` - perform analysis by sniffing all the data read and written from the server.
* `Stats` - wraps stats data the would be flushed to stdout upon request.
* `TimeTable` - contains snapshot of aggregated number of requests/responses at specific timestamp.
* `Saver` - handles reading/writing data. 
* `Recovery` - handles storing/retrieval of backup data. 
* `HTTPServer` - reports metrics/stats over HTTP using the path `/metrics` or `/stats`,
also used for health check using `/health` or `/status`.

## Todo
   - Resource pooling for connections, to enable reuse of a limited number of open connections with the server,
     and to requeue unused ones.
   - Performance and memory optimization.
   - More unit tests coverage
   - Dockerization
   - Refactoring

Enjoy!

[1]: https://golang.org/dl/
[2]: https://golang.org/doc/install
