package main

import (
	"fmt"
	"fpdns_client/methods"
)

func main() {
	_, err := methods.BehaviorTest("127.0.0.1", "127.0.0.1", "8888", "echodns.xyz")
	if err != nil {
		fmt.Println(err)
	}
}
