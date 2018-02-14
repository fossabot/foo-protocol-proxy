Foo Protocol Proxy [![CircleCI](https://circleci.com/gh/ahmedkamals/foo-protocol-proxy.svg?style=svg)](https://circleci.com/gh/ahmedkamals/foo-protocol-proxy "Build Status")
==================

[![GitHub tag](https://img.shields.io/github/tag/ahmedkamals/foo-protocol-proxy.svg?style=flat)](https://github.com/ahmedkamals/foo-protocol-proxy/releases  "Version Tag")
[![Travis CI](https://travis-ci.org/ahmedkamals/foo-protocol-proxy.svg)](https://travis-ci.org/ahmedkamals/foo-protocol-proxy "Cross Build Status [Linux, OSx]") 
[![Coverage Status](https://coveralls.io/repos/github/ahmedkamals/foo-protocol-proxy/badge.svg?branch=master)](https://coveralls.io/github/ahmedkamals/foo-protocol-proxy?branch=master  "Code Coverage")
[![Go Report Card](https://goreportcard.com/badge/github.com/ahmedkamals/foo-protocol-proxy)](https://goreportcard.com/report/github.com/ahmedkamals/foo-protocol-proxy  "Go Report Card")
[![GoDoc](https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy?status.svg)](https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy "API Documentation")
[![Docker Pulls](https://img.shields.io/docker/pulls/ahmedkamal/foo-protocol-proxy.svg?maxAge=604800)](https://hub.docker.com/r/ahmedkamal/foo-protocol-proxy/ "Docker Pulls")
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE  "License")
[![Join the chat at https://gitter.im/ahmedkamals/foo-protocol-proxy](https://badges.gitter.im/ahmedkamals/foo-protocol-proxy.svg)](https://gitter.im/ahmedkamals/foo-protocol-proxy?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge "Let's discuss")

```
  _____             ____            _                  _   ____                      
 |  ___|__   ___   |  _ \ _ __ ___ | |_ ___   ___ ___ | | |  _ \ _ __ _____  ___   _ 
 | |_ / _ \ / _ \  | |_) | '__/ _ \| __/ _ \ / __/ _ \| | | |_) | '__/ _ \ \/ / | | |
 |  _| (_) | (_) | |  __/| | | (_) | || (_) | (_| (_) | | |  __/| | | (_) >  <| |_| |
 |_|  \___/ \___/  |_|   |_|  \___/ \__\___/ \___\___/|_| |_|   |_|  \___/_/\_\\__, |
                                                                               |___/ 
```

is a Golang implementation for proxy receiver over Foo protocol.

Table of Contents
-----------------

* [Overview](#overview)
* [Getting Started](#getting-started)
    + [Prerequisites](#prerequisites)
    + [Installation](#installation)
    + [Test Driver](#test-driver)
* [Tests](#tests)
* [Coding - __Structure & Design__](#coding---structure--design)
* [Todo](#todo)

Overview
--------

Foo is a simple, powerful, extensible, cloud-native request/response protocol.
Foo protocol messages are UTF-8 encoded strings with the following BNF-ish definition:

```go
Msg  := Type <whitespace> Seq [<whitespace> Data] '\n'
Type := "REQ" | "ACK" | "NAK"
Seq  := <integer>
Data := <string without newline>
```

> When the Foo server receives a TCP connection initiated by a Foo client,
A session is started and a set of messages are exchanged between the two parties,
till the client terminates the connection eventually.

The clients sign their messages to the server with "REQ" along with the required action,
and the server might respond with acknowledgement "ACK" or denial "NAK" signature.

##### Example

```html
A->B: <connect>
A<-B: REQ 1 Hey\n
A<-B: ACK 1 Hello\n
A->B: REQ 2 Hey there\n
A<-B: ACK 2\n
A->B: REQ 3 Hey\n
A->B: REQ 4 Hey\n
A->B: REQ 5 Hey\n
A<-B: ACK 3 What\n
A<-B: ACK 4 What do you want\n
A->B: REQ 6 Hey\n
A<-B: NAK 5 Stop it\n
A<-B: NAK 6 Stop doing that\n
A->B: <disconnect>
```

The proxy has a reporting features like:

- Reporting stats in JSON format to stdout when sending `SIGUSR2` signal to the process.
- Reporting stats in JSON format over HTTP [`/stats`][16] or [`/metrics`][17].
- Health check over HTTP [`/health`][18] or [`/status`][19].
- Data recovery after failure using a sliding window of `10s` time frame.

##### JSON Sample Response

```json
{"msg_total":10,"msg_req":10,"msg_ack":8,"msg_nak":2,"request_rate_1s":0.005,"request_rate_10s":0.004,"response_rate_1s":0.004,"response_rate_10s":0.003}

```

Getting Started
---------------

### Prerequisites

* [Golang][1] installation, having [$GOPATH][2] properly set.
* Optional, required only if using Docker approach.
    + [Docker][3]
    + [Docker Compose][4]


### Installation

To install [**Foo Protocol Proxy**][5]

```bash
$ go get github.com/ahmedkamals/foo-protocol-proxy
```

### Test Driver

You can use the following steps as a testing procedure

  * **Server**
    ```bash
    $ bin/server-$(uname -s | tr '[:upper:]' '[:lower:]') -listen ":8001"
    ```

  * **Proxy**
    - Normal Approach
    
        Running proxy on the host directly.
        
        ```bash
        $ make help setup get-deps
          export ARCH=amd64
          export FORWARDING_PORT=":8001"
          export LISTENING_PORT=":8002"
          export HTTP_ADDRESS="0.0.0.0:8088"
          export RECOVERY_PATH="data/recovery.json"
          export ARGS="-forward $FORWARDING_PORT -listen $LISTENING_PORT -http $HTTP_ADDRESS -recovery-path $RECOVERY_PATH"
          make run args=$ARGS
        ```
        
        **`Environment`**:
        + `OS` - the current operating system, e.g. (linux, darwin, ...etc.)
        + `ARCH` - the current system architecture, e.g. (amd64, 386)
            
        **`Params`**:           
        + `FORWARDING_PORT` - e.g. `":8001"`
        + `LISTENING_PORT` - e.g. `":8002"`
        + `HTTP_ADDRESS` - e.g. `"0.0.0.0:8088"`
        + `RECOVERY_PATH` - e.g. `"data/recovery.json"`
        
        **Sending `SIGUSR2` Signal**
                  
        ```bash
        $ make kill args=-SIGUSR2
        ```
                   
    - Docker Approach
       
       Running proxy through docker container. currently Linux only works fine.
       
       ```bash
       $ export IMAGE_PREFIX="ahmedkamal"
         export IMAGE_TAG=$(git describe --abbrev=0 | cut -d "v" -f 2 2> /dev/null)
         export FORWARDING_PORT=8001
         export LISTENING_PORT=8002
         export HTTP_ADDRESS=8088
         export RECOVERY_PATH="data/recovery.json"
         export ARGS="-f $FORWARDING_PORT -l $LISTENING_PORT -h $HTTP_ADDRESS -r $RECOVERY_PATH -p $IMAGE_PREFIX -t $IMAGE_TAG"
         make deploy args=$ARGS
       ```
        
       **`Params`**:
       + `f` - e.g. `8001`
       + `l` - e.g. `8002`
       + `h` - e.g. `8088`
       + `r` - e.g. `"data/recovery.json"`
       + `p` - e.g. `"ahmedkamal"`
       + `t` - e.g. `"0.0.1"`
       
       **Sending `SIGUSR2` Signal**
         
       ```bash
       $ make docker-kill args=-SIGUSR2
         make docker-logs
       ```
       
  * **Multiple Client Connections**
    ```bash
    $ export OS=`uname -s | tr '[:upper:]' '[:lower:]'`
      for i in {0..1000}
      do 
         bin/client-${OS} -connect "localhost:8002";
      done
    ```

## Tests
    
Not all items covered, just made few examples.
    
```bash
$ make test
```

## Coding - __Structure & Design__

| Item                    | Description                                                                                                                                              |
| :---:                   | :---                                                                                                                                                     |
| [`Dispatcher`][6]       | parses configuration, builds configuration object, and passes it to the proxy.                                                                           |
| [`Proxy`][7]            | orchestrates the interactions between the components.                                                                                                    |
| [`Listner`][8]          | awaits for client connections, and on every new connection, a `BridgeConnection` instance is created.                                                    |
| [`BridgeConnection`][9] | acts as Bi-directional communication object, that passes data forward and backward to the server.                                                        |
| [`Analyzer`][10]        | performs `analysis` by sniffing all the data read and written from the server.                                                                           |
| [`Stats`][11]           | wraps stats data the would be flushed to stdout upon request.                                                                                            |
| [`TimeTable`][12]       | contains snapshot of aggregated number of requests/responses at specific timestamp.                                                                      |
| [`Saver`][13]           | handles reading/writing data.                                                                                                                            |
| [`Recovery`][14]        | handles storing/retrieval of backup data.                                                                                                                |
| [`HTTPServer`][15]      | reports metrics/stats over HTTP using the path [`/stats`][16] or [`/mertrics`][17], also used for health check using [`/health`][18] or [`/status`][19]. |

## Todo
   - Resource pooling for connections, to enable reuse of a limited number of open connections with the server,
     and to requeue unused ones.
   - Performance and memory optimization.
   - More unit tests coverage.
   - Refactoring.

Enjoy!

[1]: https://golang.org/dl/ "Download Golang"
[2]: https://golang.org/doc/install "GOPATH Configuration"
[3]: https://docs.docker.com/engine/installation/ "Docker"
[4]: https://docs.docker.com/compose/install/ "Docker Compose"
[5]: https://github.com/ahmedkamals/foo-protocol-proxy "Source Code"
[6]: https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy/app#Dispatcher "Dispatcher"
[7]: https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy/app#Proxy "Proxy"
[8]: https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy/communication#Listener "Listener"
[9]: https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy/communication#BridgeConnection "BridgeConnection"
[10]: https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy/analysis#Analyzer "Analyzer"
[11]: https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy/analysis#Stats "Stats"
[12]: https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy/analysis#TimeTable "TimeTable"
[13]: https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy/persistence#Saver "Saver"
[14]: https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy/persistence#Recovery "Recovery"
[15]: https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy/app#HttpServer "HttpServer"
[16]: http://localhost:8088/stats "Stats" 
[17]: http://localhost:8088/metrics "Metrics" 
[18]: http://localhost:8088/health "Health" 
[19]: http://localhost:8088/status "Status" 