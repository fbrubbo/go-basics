package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

const version = "0.1.0"

func main() {
	p := flag.String("p", "", "Filter by the pod name (default:empty - means all pods)")
	d := flag.String("d", "", "Filter by the deployment name (default:empty - means all deployments)")
	n := flag.String("n", "default", "Filter by namespace name (default: default)")
	o := flag.String("o", "tab", "Show output as: tab | csv (default: tab)")
	v := flag.Bool("v", false, "Show the plugin version")
	allNamespaces := flag.Bool("all-namespaces", false, "No filter at all, returns all pods in all namespaces (default: false)")
	noHeaders := flag.Bool("no-headers", false, "When true, remove filters")
	flag.Parse()

	fmt.Println("---------------------------------------------")
	fmt.Println("FLAGS: ")
	fmt.Println("   -p [POD] is: ", *p)
	fmt.Println("   -d [DEPLOYMENT] is: ", *d)
	fmt.Println("   -n [OUTPUT] is: ", *o)
	fmt.Println("   -o [NAMESPACE] is: ", *n)
	fmt.Println("   -v [VERSION] is: ", *v)
	fmt.Println("   -no-headers [NO HEADERS] is: ", *noHeaders)
	fmt.Println("   -all-namespaces [All NAMESPACE] is: ", *allNamespaces)
	fmt.Println("---------------------------------------------")
	fmt.Println()

	if *v {
		fmt.Println("Plugin Version: ", version)
		os.Exit(0)
	}

	if *allNamespaces {
		*n = ""
	}
	var podList []Pod
	tmpPodList := RetrievePods(*n)
	topMap := RetrieveTopMap(*n)
	for _, pod := range tmpPodList {
		if top, ok := topMap[pod.GetPodKey()]; ok {
			pod.Top = top // enrich Pod with top info
		}
		podList = append(podList, pod)
	}

	var result Wrapper
	if *p != "" {
		result = buildPod(podList, *p)
	} else if *d != "" {
		result = buildDeployment(podList, *d)
	} else {
		result = Wrapper{Type: "All Pods", Pods: podList}
	}

	formatHeader := "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n"
	formatValues := "%v\t%v\t%vm\t%vm\t%0.2f%%\t%vMi\t%vMi\t%0.2f%%\t%vm\t%vMi\n"
	if *o == "tab" {
		fmt.Println("PODs SUMMARY:")
		w := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, formatHeader, "Namespace", "Pod Name", "Requests CPU (m)", "TOP CPU (m)", "Usage CPU (%)", "Requests Memory (Mi)", "TOP Memory (Mi)", "Usage Memory (%)", "Limits CPU (m)", "Limitis Memory (Mi)")
		fmt.Fprintf(w, formatHeader, "---------", "--------", "----------------", "-----------", "-------------", "--------------------", "---------------", "----------------", "--------------", "-------------------")
		for _, pod := range result.Pods {
			fmt.Fprintf(w, formatValues, pod.Metadata.Namespace, pod.Metadata.Name, pod.GetRequestsMilliCPU(), pod.GetTopMilliCPU(), pod.GetUsageCPU(), pod.GetRequestsMiMemory(), pod.GetTopMiMemory(), pod.GetUsageMemory(), pod.GetLimitsMilliCPU(), pod.GetLimitsMiMemory())
		}
		fmt.Fprintf(w, formatHeader, "---------", "--------", "----------------", "-----------", "-------------", "--------------------", "---------------", "----------------", "--------------", "-------------------")
		fmt.Fprintf(w, formatValues, nil, nil, result.GetRequestsMilliCPU(), result.GetTopMilliCPU(), result.GetUsageCPU(), result.GetRequestsMiMemory(), result.GetTopMiMemory(), result.GetUsageMemory(), result.GetLimitsMilliCPU(), result.GetLimitsMiMemory())
		w.Flush()
	} else if *o == "csv" {
		fmt.Println("Output csv not implemented!")
	} else {
		fmt.Println("Output (-o) parameter is invalid!")
	}

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

// GetTopMilliCPU total
func (d Wrapper) GetTopMilliCPU() int {
	total := 0
	for _, p := range d.Pods {
		total += p.Top.GetMilliCPU()
	}
	return total
}

// GetUsageCPU % usage
func (d Wrapper) GetUsageCPU() float32 {
	requests, top := 0, 0
	for _, p := range d.Pods {
		requests += p.GetRequestsMilliCPU()
		top += p.Top.GetMilliCPU()
	}
	return float32(top) / float32(requests) * 100
}

// GetRequestsMiMemory total
func (d Wrapper) GetRequestsMiMemory() int {
	total := 0
	for _, p := range d.Pods {
		total += p.GetRequestsMiMemory()
	}
	return total
}

// GetTopMiMemory total
func (d Wrapper) GetTopMiMemory() int {
	total := 0
	for _, p := range d.Pods {
		total += p.Top.GetMiMemory()
	}
	return total
}

// GetUsageMemory % usage
func (d Wrapper) GetUsageMemory() float32 {
	requests, top := 0, 0
	for _, p := range d.Pods {
		requests += p.GetRequestsMiMemory()
		top += p.Top.GetMiMemory()
	}
	return float32(top) / float32(requests) * 100
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
