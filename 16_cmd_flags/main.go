package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

func main() {
	p := flag.String("p", "", "The pod name (default:empty - means all pods)")
	d := flag.String("d", "", "The deployment name (default:empty - means all deployments)")
	n := flag.String("n", "default", "Return all pod in the given namespace (default: default)")
	allNamespaces := flag.Bool("all-namespaces", false, "Returns all pods in all Namespaces (default: false)")
	flag.Parse()

	fmt.Println("POD is: ", *p)
	fmt.Println("DEPLOYMENT is: ", *d)
	fmt.Println("NS is: ", *n)
	fmt.Println("All Namespaces is: ", *allNamespaces)
	fmt.Println("tail:", flag.Args())

	if *allNamespaces {
		*n = ""
	}
	podList := RetrievePods(*n)

	var result Wrapper
	if *p != "" {
		result = buildPod(podList, *p)
	} else if *d != "" {
		result = buildDeployment(podList, *d)
	} else {
		result = Wrapper{Type: "All Pods", Pods: podList}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", "Namespace", "Pod Name", "Requests CPU (m)", "Requests Memory (Mi)", "Limits CPU (m)", "Limitis Memory (Mi)")
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", "---------", "--------", "----------------", "--------------------", "--------------", "-------------------")
	for _, pod := range result.Pods {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", pod.Metadata.Namespace, pod.Metadata.Name, pod.GetRequestsMilliCPU(), pod.GetRequestsMiMemory(), pod.GetLimitsMilliCPU(), pod.GetLimitsMiMemory())
	}
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", "---------", "--------", "----------------", "--------------------", "--------------", "-------------------")
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", nil, nil, result.GetRequestsMilliCPU(), result.GetRequestsMiMemory(), result.GetLimitsMilliCPU(), result.GetLimitsMiMemory())
	w.Flush()
}

// Wrapper contains a list of pods
type Wrapper struct {
	Type string
	Pods []Pod
}

// GetRequestsMilliCPU total
func (d Wrapper) GetRequestsMilliCPU() int {
	total := 0
	for _, p := range d.Pods {
		total += p.GetRequestsMilliCPU()
	}
	return total
}

// GetRequestsMiMemory total
func (d Wrapper) GetRequestsMiMemory() int {
	total := 0
	for _, p := range d.Pods {
		total += p.GetRequestsMiMemory()
	}
	return total
}

// GetLimitsMilliCPU total
func (d Wrapper) GetLimitsMilliCPU() int {
	total := 0
	for _, p := range d.Pods {
		total += p.GetLimitsMilliCPU()
	}
	return total
}

// GetLimitsMiMemory total
func (d Wrapper) GetLimitsMiMemory() int {
	total := 0
	for _, p := range d.Pods {
		total += p.GetLimitsMiMemory()
	}
	return total
}

func buildPod(podList []Pod, p string) Wrapper {
	pods := make(map[string]Pod)
	for _, pod := range podList {
		pods[pod.Metadata.Name] = pod
	}
	return Wrapper{Type: "Pod", Pods: []Pod{pods[p]}}
}

func buildDeployment(podList []Pod, d string) Wrapper {
	pods := make(map[string][]Pod)
	for _, pod := range podList {
		deploymentName := pod.GetDeploymentName()
		if pods[deploymentName] == nil {
			pods[deploymentName] = []Pod{pod}
		} else {
			pods[deploymentName] = append(pods[deploymentName], pod)
		}
	}
	return Wrapper{Type: "Deployment", Pods: pods[d]}
}
