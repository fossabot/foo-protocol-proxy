Foo Protocol Proxy [![CircleCI](https://circleci.com/gh/ahmedkamals/foo-protocol-proxy.svg?style=svg)](https://circleci.com/gh/ahmedkamals/foo-protocol-proxy "Build Status")
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fahmedkamals%2Ffoo-protocol-proxy.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fahmedkamals%2Ffoo-protocol-proxy?ref=badge_shield)
==================

[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE  "License")
[![GitHub tag](https://img.shields.io/github/tag/ahmedkamals/foo-protocol-proxy.svg?style=flat)](https://github.com/ahmedkamals/foo-protocol-proxy/releases  "Version Tag")
[![Travis CI](https://travis-ci.org/ahmedkamals/foo-protocol-proxy.svg)](https://travis-ci.org/ahmedkamals/foo-protocol-proxy "Cross Build Status [Linux, OSx]") 
[![codecov](https://codecov.io/gh/ahmedkamals/foo-protocol-proxy/branch/master/graph/badge.svg)](https://codecov.io/gh/ahmedkamals/foo-protocol-proxy)
[![Coverage Status](https://coveralls.io/repos/github/ahmedkamals/foo-protocol-proxy/badge.svg?branch=master)](https://coveralls.io/github/ahmedkamals/foo-protocol-proxy?branch=master  "Code Coverage")
[![Go Report Card](https://goreportcard.com/badge/github.com/ahmedkamals/foo-protocol-proxy)](https://goreportcard.com/report/github.com/ahmedkamals/foo-protocol-proxy  "Go Report Card")
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/c282df1ff33c43ddb5da1d7fe4e85674)](https://www.codacy.com/app/ahmedkamals/foo-protocol-proxy?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ahmedkamals/foo-protocol-proxy&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy?status.svg)](https://godoc.org/github.com/ahmedkamals/foo-protocol-proxy "API Documentation")
[![Docker Pulls](https://img.shields.io/docker/pulls/ahmedkamal/foo-protocol-proxy.svg?maxAge=604800)](https://hub.docker.com/r/ahmedkamal/foo-protocol-proxy/ "Docker Pulls")
[![DepShield Badge](https://depshield.sonatype.org/badges/ahmedkamals/foo-protocol-proxy/depshield.svg)](https://depshield.github.io "DepShield")
[![Join the chat at https://gitter.im/ahmedkamals/foo-protocol-proxy](https://badges.gitter.im/ahmedkamals/foo-protocol-proxy.svg)](https://gitter.im/ahmedkamals/foo-protocol-proxy?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge "Let's discuss")

```bash
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

### Example

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

* Reporting stats in JSON format to stdout when sending `SIGUSR2` signal to the process.
* Reporting stats in JSON format over HTTP [`/stats`][16] or [`/metrics`][17].
* Health check over HTTP [`/health`][18] or [`/status`][19].
* Data recovery after failure using a sliding window of `10s` time frame.

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

  * __Server__
    ```bash
    $ build/assets/bin/server-$(uname -s | tr '[:upper:]' '[:lower:]') -listen ":8001"
    ```

  * __Proxy__
    - Normal Approach
    
        Running proxy on the host directly.
        
        ```bash
        $ make help get-deps go-install
          export OS=$(uname -s | tr '[:upper:]' '[:lower:]')
          export ARCH=amd64
          export FORWARDING_PORT=":8001"
          export LISTENING_PORT=":8002"
          export HTTP_ADDRESS="0.0.0.0:8088"
          export RECOVERY_PATH=".data/recovery.json"
          export ARGS="-forward $FORWARDING_PORT -listen $LISTENING_PORT -http $HTTP_ADDRESS -recovery-path $RECOVERY_PATH"
          make run OS=$OS ARCH=$ARCH args=$ARGS
        ```
        
        **`Environment`**:
        + `OS` - the current operating system, e.g. (linux, darwin, ...etc.) - default is the current operating system.
        + `ARCH` - the current system architecture, e.g. (amd64, 386)= default is amd64
            
        **`Params`**:           
        + `FORWARDING_PORT` - e.g. `":8001"`
        + `LISTENING_PORT` - e.g. `":8002"`
        + `HTTP_ADDRESS` - e.g. `"0.0.0.0:8088"`
        + `RECOVERY_PATH` - e.g. `".data/recovery.json"`
        
        **Sending `SIGUSR2` Signal**
                  
        ```bash
        $ make kill args=-SIGUSR2
        ```
                   
    - Docker Approach
       
       Running proxy through docker container. currently Linux only works fine.
       
       ```bash
       $ export ARCH="amd64"
         export FORWARDING_PORT=8001
         export LISTENING_PORT=8002
         export HTTP_ADDRESS=8088
         export RECOVERY_PATH=".data/recovery.json"
         export ARGS="-f $FORWARDING_PORT -l $LISTENING_PORT -h $HTTP_ADDRESS -r $RECOVERY_PATH"
         make deploy ARCH=$ARCH args=$ARGS
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
       
  * __Multiple Client Connections__
    ```bash
    $ for i in {0..100}
      do 
         docker start foo.client;
      done
    ```

## Tests
    
Not all items covered, just made few examples.
    
```bash
$ make test
```

## Coding - __Structure & Design__

| Item                    | Description                                                                                                                                                                 |
| :---:                   | :---                                                                                                                                                                        |
| [`Dispatcher`][6]       | parses configuration, builds configuration object, and passes it to the proxy.                                                                                              |
| [`Proxy`][7]            | orchestrates the interactions between the components.                                                                                                                       |
| [`Listner`][8]          | awaits for client connections, and on every new connection, a `BridgeConnection` instance is created.                                                                       |
| [`BridgeConnection`][9] | acts as Bi-directional communication object, that passes data forward and backward to the server.                                                                           |
| [`Analyzer`][10]        | performs `analysis` by sniffing all the data read and written from the server.                                                                                              |
| [`Stats`][11]           | wraps stats data the would be flushed to stdout upon request.                                                                                                               |
| [`TimeTable`][12]       | contains snapshot of aggregated number of requests/responses at specific timestamp.                                                                                         |
| [`Saver`][13]           | handles reading/writing data.                                                                                                                                               |
| [`Recovery`][14]        | handles storing/retrieval of backup data.                                                                                                                                   |
| [`HTTPServer`][15]      | reports metrics/stats over HTTP using the path [`/stats`][16] or [`/mertrics`][17], also used for health check using [`/health`][18] [`/heatrbeat`][19] or [`/status`][20]. |


<details>
<summary>Todo:</summary>
   - Resource pooling for connections, to enable reuse of a limited number of open connections with the server,
     and to requeue unused ones.
   - Performance and memory optimization.
   - More unit tests coverage.
   - Refactoring.
</details>

Enjoy!
[![Analytics](http://www.google-analytics.com/__utm.gif?utmwv=4&utmn=869876874&utmac=UA-136526477-1&utmcs=ISO-8859-1&utmhn=github.com&utmdt=Foo%20Protocol%20Proxy&utmcn=1&utmr=0&utmp=/ahmedkamals/foo-protocol-proxy?utm_source=www.github.com&utm_campaign=Foo+Protocol+Proxy&utm_term=Foo+Protocol+Proxy&utm_content=Foo+Protocol+Proxy&utm_medium=repository&utmac=UA-136526477-1)]()

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
[19]: http://localhost:8088/heartbeat "Heartbeat"
[20]: http://localhost:8088/status "Status"


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fahmedkamals%2Ffoo-protocol-proxy.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fahmedkamals%2Ffoo-protocol-proxy?ref=badge_large)