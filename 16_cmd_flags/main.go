package main

import (
	"flag"
	"fmt"
)

func main() {
	pod := flag.String("p", "", "The pod name (default:empty - means all pods)")
	ns := flag.String("n", "default", "Return all pod in the given namespace (default: default)")
	allNamespaces := flag.Bool("all-namespaces", false, "Returns all pods in all Namespaces (default: false)")
	flag.Parse()

	fmt.Println("POD is: ", *pod)
	fmt.Println("NS is: ", *ns)
	fmt.Println("All Namespaces is: ", *allNamespaces)
	fmt.Println("tail:", flag.Args())

	res := GetPodResources(*pod, *ns)
	fmt.Printf("go struct: %+v\n", res)
}
