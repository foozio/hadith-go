// This is a stub to guide users to enable gRPC build.
package main

import "fmt"

func main() {
    fmt.Println("gRPC server build is disabled by default.")
    fmt.Println("To enable: generate proto code and build with -tags grpc.")
    fmt.Println("See: Makefile targets 'proto' and 'grpc'.")
}

