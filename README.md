Foo Protocol Proxy
==================

[![GoDoc](https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy?status.svg)](https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy "API Documentation")

A Golang implementation for proxy receiver over Foo protocol.
It also has some extra features like:

- Reporting stats when sending `SIGUSR2` signal to the process and over HTTP.
- Health check over HTTP.
- Data recovery after failure by saving to file, if the failure period does not exceed the `10s` time frame.

Installation
------------

### Prerequisites

* [Golang][1] installation, having [$GOPATH][2] properly set.

To install [**Foo Protocol Proxy**][3]

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
    - Normal Approach
    
        Running proxy on the host directly.
        
        ```bash
        $ make build
        $ bin/foo-protocol-proxy-{OS}-{ARCH} -forward "{FORWARDING_PORT}" -listen "{LISTENING_PORT}" -http "{HTTP_ADDRESS}" -recovery-path "{RECOVERY_PATH}"
        ```
        
        **`Environment`**:
        + `OS` - the current operating system, e.g. (linux, darwin, ...etc.)
        + `ARCH` - the current system architecture, e.g. (amd64, 386)
            
        **`Params`**:           
        + `FORWARDING_PORT` - e.g. `":8081"`
        + `LISTENING_PORT` - e.g. `":8082"`
        + `HTTP_ADDRESS` - e.g. `"0.0.0.0:8088"`
        + `RECOVERY_PATH` - e.g. `"data/recovery.json"`
        
        **Sending `SIGUSR2` Signal**
                  
        ```bash
        # Process name: foo-protocol-proxy-{OS}-{ARCH}
        # For linux, and amd64, it would be as following:
        $ kill -SIGUSR2 $(pidof foo-protocol-proxy-linux-amd64) > /dev/null 2>&1
        ```
                   
    - Docker Approach
       
       Running proxy through docker container.
       
       ```bash
       $ bash deploy.sh -f "{FORWARDING_PORT}" -l "{LISTENING_PORT}" -h "{HTTP_ADDRESS}" -r "{RECOVERY_PATH}"
       ```
        
       **`Params`**:
       + `f` - e.g. `8081`
       + `l` - e.g. `8082`
       + `h` - e.g. `8088`
       + `r` - e.g. `"data/recovery.json"`
       
       **Sending `SIGUSR2` Signal**
         
       ```bash
       $ docker exec -it foo-proxy-0.0.1 pkill -SIGUSR2 foo-protocol-proxy > /dev/null 2>&1
       $ docker logs -f foo-proxy-0.0.1
       ```
       
  * **Multiple Client Connections**
    ```bash
    $ for i in {0..1000..1}
      do 
         bin/client-linux -connect "localhost:8002";
      done
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
* `Analyzer` - perform [`analysis`][4] by sniffing all the data read and written from the server.
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
   - Adding Go package documentation.
   - More unit tests coverage.
   - Refactoring.

Enjoy!

[1]: https://golang.org/dl/
[2]: https://golang.org/doc/install
[3]: https://github.com/ahmedkamals/foo-protocol-proxy
[4]: https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy/analysis "API Documentation"
