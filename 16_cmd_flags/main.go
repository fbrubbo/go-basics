package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"
)

type Resource struct {
	Cpu    string
	Memory string
}

type Resources struct {
	Containers []struct {
		Resources struct {
			Requests Resource
			Limits   Resource
		}
	}
}

func main() {
	pod := flag.String("p", "", "The pod name (default:empty - means all pods)")
	ns := flag.String("n", "default", "Return all pod in the given namespace (default: default)")
	allNamespaces := flag.Bool("all-namespaces", false, "Returns all pods in all Namespaces (default: false)")
	flag.Parse()

	fmt.Println("POD is: ", *pod)
	fmt.Println("NS is: ", *ns)
	fmt.Println("All Namespaces is: ", *allNamespaces)
	fmt.Println("tail:", flag.Args())

	cmd := fmt.Sprintf("kubectl get pod %s -n %s -o json | jq -r '.spec'", *pod, *ns)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("Failed to execute command: %s", cmd)
	}
	s := string(out)
	fmt.Printf("combined out:\n%s\n", s)

	res := Resources{}
	err2 := json.Unmarshal([]byte(s), &res)
	if err2 != nil {
		fmt.Println(err2.Error())
	}
	fmt.Printf("Operation: %s", res)
}
