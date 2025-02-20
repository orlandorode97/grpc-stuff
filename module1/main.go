package main

import (
	"fmt"

	"github.com/orlandorode97/grpc-stuff/module1/proto"
)

func main() {
	person := &proto.Person{
		Name: "Orlando",
	}

	fmt.Printf("Hello: %v\n", person.GetName())
}
