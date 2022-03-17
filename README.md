# memCacheLabs

A Golang-based implementation of an in-memory caching library as-a-service; built to be language-agnostic and accessible via http/s requests. <br>

###### (a) In Development

The service, memCacheLabs, acts as an integrative, lightweight, plug-and-playable caching layer, living between the server's application and database. <br>

###### (b.1) Directory

`├── client` <br>
`│ └── controllers` <br>
`│ └── models` <br>
`│ └── router` <br>
`│ └── main.go` <br>
`├── config` <br>
`│ └── redis.conf` <br>
`├── dockerfile` <br>
`├── LICENSE` <br>
`├── README.md` <br>

###### (b.2) run

Navigate to ./client,

`├── client` <br>

1. If no prior runs, proceed to pull in application dependecies by running `go get`. If dependencies already installed, proceed to step 2.
2. To build, run `go build main.go`.
3. To start the service, run `go run main.go`.<br>

###### (c) API Documentation

Read the Docs: https://documenter.getpostman.com/view/19335839/UVeAwUu2
<br>

###### < > with ♥ using Go

###### Author: P William Hozier
