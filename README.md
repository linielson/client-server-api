# client-server-api

This challenge is part of the [GoExpert Course](https://goexpert.fullcycle.com.br/)

In this challenge let's apply what we learned about:
* http web server
* context
* database
* file operations in Go

Requirements:

* two files:
- client.go
- server.go

* client.go should make a HTTP request to server.go, requesting the Dollar quotation

* server.go should consume the rest API https://economia.awesomeapi.com.br/json/last/USD-BRL with the quotation Dollar to Real, and return the response to the client in json format.

* using the "context" package, the server must register each received quote in the SQLite database, and the maximum timeout to call the dollar quote API must be 200ms, and the maximum timeout to be able to persist the data in the database should be 10ms.

* client.go will only need to receive from server.go the current exchange rate (JSON "bid" field). Using the "context" package, client.go will have a maximum timeout of 300ms to receive the result from server.go.

* the client.go will have to save the current quote in a file "quotation.txt" in the format: Dollar: {value}

* the required endpoint generated by server.go for this challenge will be: /quotation and the port to be used by the HTTP server will be 8080.
