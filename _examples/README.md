phi examples
============

* [custom-handler](go.philip.id/phi/phi/blob/master/_examples/custom-handler/main.go) - Use a custom handler function signature
* [custom-method](go.philip.id/phi/phi/blob/master/_examples/custom-method/main.go) - Add a custom HTTP method
* [fileserver](go.philip.id/phi/phi/blob/master/_examples/fileserver/main.go) - Easily serve static files
* [graceful](go.philip.id/phi/phi/blob/master/_examples/graceful/main.go) - Graceful context signaling and server shutdown
* [hello-world](go.philip.id/phi/phi/blob/master/_examples/hello-world/main.go) - Hello World!
* [limits](go.philip.id/phi/phi/blob/master/_examples/limits/main.go) - Timeouts and Throttling
* [logging](go.philip.id/phi/phi/blob/master/_examples/logging/main.go) - Easy structured logging for any backend
* [rest](go.philip.id/phi/phi/blob/master/_examples/rest/main.go) - REST APIs made easy, productive and maintainable
* [router-walk](go.philip.id/phi/phi/blob/master/_examples/router-walk/main.go) - Print to stdout a router's routes
* [todos-resource](go.philip.id/phi/phi/blob/master/_examples/todos-resource/main.go) - Struct routers/handlers, an example of another code layout style
* [versions](go.philip.id/phi/phi/blob/master/_examples/versions/main.go) - Demo of `phi/render` subpkg


## Usage

1. `go get -v -d -u ./...` - fetch example deps
2. `cd <example>/` ie. `cd rest/`
3. `go run *.go` - note, example services run on port 3333
4. Open another terminal and use curl to send some requests to your example service,
   `curl -v http://localhost:3333/`
5. Read <example>/main.go source to learn how service works and read comments for usage
