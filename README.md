Foo Protocol Proxy
===========

A Golang implementation for proxy receiver over Foo protocol.

Installation
-------------

### Prerequisites

* [Golang][1] installation, having [$GOPATH][2] properly set.

To install [**foo-protocol-proxy**](https://github.com/ahmedkamals/foo-protocol-proxy)

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
    $ bin/foo-protocol-proxy -forward "{FORWARDING_PORT}" -listen "{LISTENING_PORT}"
    ```
    
**`Params`**:
   * FORWARDING_PORT
   * LISTENING_PORT
          
  * **Multiple Client Connections**
    ```bash
    $ for i in {0..1000..1}
      do 
         bin/client-linux -connect "localhost:8002";
      done
    ```

  * **Sending SIGUSR2 Signal**
      
    ```bash
    $ kill -SIGUSR2 $(pidof foo-protocol-proxy) > /dev/null 2>&1
    ```

## Tests
TBD

## Coding - __Structure & Design__
* `Dispatcher` - parses configuration, builds configuration object,
and passes it to the proxy.
* `Proxy` - Orchestrates the interactions between the components. 
* `Listner` - awaits for client connections, and on every new connection, 
a `BridgeConnection` instance is created.
* `BridgeConnection` - acts as Bi-directional communication object, that
passes data forward and backward to the server.
* `Analyzer` - perform analysis by sniffing all the data read and written from the server.
* `Stats` - wraps stats data the would be flushed to stdout upon request.
* `TimeTable` - contains snapshot of aggregated number of requests/responses at specific timestamp.
  
## Todo
   - Health check through http url.
   - Resource pooling for connections, to enable reuse of a limited number of open connections with the server,
     and to requeue unused ones.
   - Performance and memory optimization.
   - Data recovery after failure by saving to file.
   - Unit tests
   - Refactoring

Enjoy!

[1]: https://golang.org/dl/
[2]: https://golang.org/doc/install
