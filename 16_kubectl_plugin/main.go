package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

const version = "0.1.0"

// TODO: sort-by ? How to handle the below scenarios?
func main() {
	p := flag.String("p", "", "Filter by the pod name (default:empty means all pods)")
	d := flag.String("d", "", "Filter by the deployment name (default:empty means all deployments)")
	n := flag.String("n", "", "Filter by namespace name (default:empty means all namespaces)")
	v := flag.Bool("v", false, "Show the plugin version")
	show := flag.String("print", "all", "Define what will be printed. Valid values all|pods|hpas|nodes ")
	csv := flag.Bool("csv-output", false, "Save the result to files with format 'kubectl-snapshot-<date>-<pods|hpas|nohpa|nodes|all>.csv'")
	debug := flag.Bool("debug", false, "Show debug info")
	flag.Parse()
	printFlags(*p, *d, *n, *v, *show, *csv, *debug)

	if *v || *debug {
		fmt.Println("Plugin Version: ", version)
		if *v {
			os.Exit(0)
		}
	}
	csvFilePrefix := ""
	if *csv {
		now := time.Now()
		csvFilePrefix = now.Format("kubectl-snapshot-2006-01-02-1504")
	}

	// Pods with resource usage (top) ..
	podList := RetrievePods(*n)
	if *p != "" {
		podList = filterPod2(podList, func(pod Pod) bool { return pod.Metadata.Name == *p })
	} else if *d != "" {
		podList = filterPod2(podList, func(pod Pod) bool { return pod.GetDeploymentName() == *d })
	}

	// Hpas, use podList to confirm resource usgage ..
	hpaList := RetrieveHpas(*n, podList)
	if *p != "" {
		hpaList = filterHpa(hpaList, func(h Hpa) bool { return h.ContainsPod(*p) })
	} else if *d != "" {
		hpaList = filterHpa(hpaList, func(h Hpa) bool { return h.RefToDeployment(*d) })
	}

	// Deployments for non-hpas, use podList to confirm resource usgage ..
	deploymentList := RetrieveDeployments(*n, podList)
	if *p != "" {
		deploymentList = filterDeployment(deploymentList, func(deploy Deployment) bool { return deploy.ContainsPod(*p) })
	} else if *d != "" {
		deploymentList = filterDeployment(deploymentList, func(deploy Deployment) bool { return deploy.Name == *d })
	}
	hpaMap := make(map[string]Hpa)
	for _, hpa := range hpaList {
		hpaMap[hpa.Namespace+"|"+hpa.ReferenceName] = hpa
	}
	deploymentWithoutHpa := []Deployment{}
	for _, deploy := range deploymentList {
		if _, hasHpa := hpaMap[deploy.GetDeploymentKey()]; !hasHpa {
			deploymentWithoutHpa = append(deploymentWithoutHpa, deploy)
		}
	}

	// Print standard io or send to csv files ..
	switch *show {
	case "pod":
	case "pods":
		printPodsTab(podList, csvFilePrefix, *debug)
	case "hpa":
	case "hpas":
		printHpaTab(hpaList, csvFilePrefix, *debug)
		printNoHpaTab(deploymentWithoutHpa, csvFilePrefix, *debug)
	case "node":
	case "nodes":
		printNodesTab(podList, csvFilePrefix, *debug)
	default:
		printPodsTab(podList, csvFilePrefix, *debug)
		printHpaTab(hpaList, csvFilePrefix, *debug)
		printNoHpaTab(deploymentWithoutHpa, csvFilePrefix, *debug)
		printNodesTab(podList, csvFilePrefix, *debug)
	}

}

func printFlags(p string, d string, n string, v bool, show string, csv bool, debug bool) {
	if debug {
		fmt.Println("---------------------------------------------")
		fmt.Println("[debug] FLAGS: ")
		fmt.Println("   -p [POD] is: ", p)
		fmt.Println("   -d [DEPLOYMENT] is: ", d)
		fmt.Println("   -o [NAMESPACE] is: ", n)
		fmt.Println("   -v [VERSION] is: ", v)
		fmt.Println("   -print [PRINT IN STANDARD OUTPUT] is: ", show)
		fmt.Println("   -csv-output [SAVE TO FILES] is: ", csv)
		fmt.Println("---------------------------------------------")
		fmt.Println()
	}
}

func printPodsTab(podList []Pod, csvFilePrefix string, debug bool) {
	result := Wrapper{Pods: podList}

	if csvFilePrefix == "" || debug {
		formatHeader := "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n"
		formatValues := "%v\t%v\t%vm\t%vm\t%0.2f%%\t%vMi\t%vMi\t%0.2f%%\t%vm\t%vMi\n"
		fmt.Println("\nPODs SNAPSHOT:")
		w := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, formatHeader, "Namespace", "Pod Name", "Requests CPU (m)", "TOP CPU (m)", "Usage CPU (%)", "Requests Memory (Mi)", "TOP Memory (Mi)", "Usage Memory (%)", "Limits CPU (m)", "Limitis Memory (Mi)")
		fmt.Fprintf(w, formatHeader, "---------", "--------", "----------------", "-----------", "-------------", "--------------------", "---------------", "----------------", "--------------", "-------------------")
		for _, pod := range result.Pods {
			fmt.Fprintf(w, formatValues, pod.Metadata.Namespace, pod.Metadata.Name, pod.GetRequestsMilliCPU(), pod.GetTopMilliCPU(), pod.GetUsageCPU(), pod.GetRequestsMiMemory(), pod.GetTopMiMemory(), pod.GetUsageMemory(), pod.GetLimitsMilliCPU(), pod.GetLimitsMiMemory())
		}
		fmt.Fprintf(w, formatHeader, " ", " ", "----------------", "-----------", "-------------", "--------------------", "---------------", "----------------", "--------------", "-------------------")
		fmt.Fprintf(w, formatValues, " ", " ", result.GetRequestsMilliCPU(), result.GetTopMilliCPU(), result.GetUsageCPU(), result.GetRequestsMiMemory(), result.GetTopMiMemory(), result.GetUsageMemory(), result.GetLimitsMilliCPU(), result.GetLimitsMiMemory())
		w.Flush()
	}

	if csvFilePrefix != "" {
		file, err := os.Create(csvFilePrefix + "-pods.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		header := []string{"Namespace", "Pod Name", "Requests CPU (m)", "TOP CPU (m)", "Usage CPU (%)", "Requests Memory (Mi)", "TOP Memory (Mi)", "Usage Memory (%)", "Limits CPU (m)", "Limitis Memory (Mi)"}
		err = writer.Write(header)
		if err != nil {
			log.Fatal(err)
		}
		for _, pod := range result.Pods {
			line := []string{pod.Metadata.Namespace, pod.Metadata.Name, strconv.Itoa(pod.GetRequestsMilliCPU()), strconv.Itoa(pod.GetTopMilliCPU()), fmt.Sprintf("%.2f", pod.GetUsageCPU()), strconv.Itoa(pod.GetRequestsMiMemory()), strconv.Itoa(pod.GetTopMiMemory()), fmt.Sprintf("%.2f", pod.GetUsageMemory()), strconv.Itoa(pod.GetLimitsMilliCPU()), strconv.Itoa(pod.GetLimitsMiMemory())}
			err := writer.Write(line)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func printHpaTab(hpaList []Hpa, csvFilePrefix string, debug bool) {
	if csvFilePrefix == "" || debug {
		formatHeader := "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n"
		formatValues := "%v\t%v\t%v\t%v\t%v\t%v\t%vm\t%vm\t%0.2f%%\t%vMi\t%vMi\t%0.2f%%\t%vm\t%vMi\n"
		fmt.Println("\nHPAs SNAPSHOT:")
		w := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, formatHeader, "Namespace", "Hpa Name", "Reference", "Target", "Replicas (Min/Max/Actual)", "# Pods ->", "Requests CPU (m)", "TOP CPU (m)", "Usage CPU (%)", "Requests Memory (Mi)", "TOP Memory (Mi)", "Usage Memory (%)", "Limits CPU (m)", "Limitis Memory (Mi)")
		fmt.Fprintf(w, formatHeader, "---------", "--------", "---------", "------", "-------------------------", "---------", "----------------", "-----------", "-------------", "--------------------", "---------------", "----------------", "--------------", "-------------------")
		for _, hpa := range hpaList {
			wp := Wrapper{Pods: hpa.Pods}
			replicas := fmt.Sprintf("%d/%d/%d", hpa.MinPods, hpa.MaxPods, hpa.Replicas)
			fmt.Fprintf(w, formatValues, hpa.Namespace, hpa.Name, hpa.GetReference(), hpa.GetUsageAndTarget(), replicas, len(hpa.Pods), wp.GetRequestsMilliCPU(), wp.GetTopMilliCPU(), wp.GetUsageCPU(), wp.GetRequestsMiMemory(), wp.GetTopMiMemory(), wp.GetUsageMemory(), wp.GetLimitsMilliCPU(), wp.GetLimitsMiMemory())
		}
		fmt.Fprintf(w, formatHeader, " ", " ", " ", "------", "-------------------------", "---------", "----------------", "-----------", "-------------", "--------------------", "---------------", "----------------", "--------------", "-------------------")
		w.Flush()
	}

	if csvFilePrefix != "" {
		file, err := os.Create(csvFilePrefix + "-hpas.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		header := []string{"Namespace", "Hpa Name", "Reference", "Hpa Use", "Hpa Target", "Min Replicas", "Max Replicas", "Actual Replicas", "# Pods ->", "Requests CPU (m)", "TOP CPU (m)", "Usage CPU (%)", "Requests Memory (Mi)", "TOP Memory (Mi)", "Usage Memory (%)", "Limits CPU (m)", "Limitis Memory (Mi)"}
		err = writer.Write(header)
		if err != nil {
			log.Fatal(err)
		}
		for _, hpa := range hpaList {
			wp := Wrapper{Pods: hpa.Pods}
			hpaUse := "<unknown>"
			if hpa.UsageCPU != -1 {
				hpaUse = strconv.Itoa(hpa.UsageCPU)
			}
			line := []string{hpa.Namespace, hpa.Name, hpa.GetReference(), hpaUse, strconv.Itoa(hpa.Target), strconv.Itoa(hpa.MinPods), strconv.Itoa(hpa.MaxPods), strconv.Itoa(hpa.Replicas), strconv.Itoa(len(hpa.Pods)), strconv.Itoa(wp.GetRequestsMilliCPU()), strconv.Itoa(wp.GetTopMilliCPU()), fmt.Sprintf("%.2f", wp.GetUsageCPU()), strconv.Itoa(wp.GetRequestsMiMemory()), strconv.Itoa(wp.GetTopMiMemory()), fmt.Sprintf("%.2f", wp.GetUsageMemory()), strconv.Itoa(wp.GetLimitsMilliCPU()), strconv.Itoa(wp.GetLimitsMiMemory())}
			err := writer.Write(line)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func printNoHpaTab(deploymentWithoutHpa []Deployment, csvFilePrefix string, debug bool) {
	if csvFilePrefix == "" || debug {
		formatHeader := "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n"
		formatValues := "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%vm\t%vm\t%0.2f%%\t%vMi\t%vMi\t%0.2f%%\t%vm\t%vMi\n"
		fmt.Println("\nNO HPA SNAPSHOT:")
		w := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, formatHeader, "Namespace", "Deployment Name", "Ready", "Up To Date", "Avaliable", "Age", "#Pods ->", "Requests CPU (m)", "TOP CPU (m)", "Usage CPU (%)", "Requests Memory (Mi)", "TOP Memory (Mi)", "Usage Memory (%)", "Limits CPU (m)", "Limitis Memory (Mi)")
		fmt.Fprintf(w, formatHeader, "---------", "---------------", "-----", "----------", "---------", "---", "--------", "----------------", "-----------", "-------------", "--------------------", "---------------", "----------------", "--------------", "-------------------")
		for _, deploy := range deploymentWithoutHpa {
			wp := Wrapper{Pods: deploy.Pods}
			ready := fmt.Sprintf("%d/%d", deploy.Replicas, deploy.ReplicasExpected)
			fmt.Fprintf(w, formatValues, deploy.Namespace, deploy.Name, ready, deploy.UpToDate, deploy.Avaliable, deploy.Age, len(deploy.Pods), wp.GetRequestsMilliCPU(), wp.GetTopMilliCPU(), wp.GetUsageCPU(), wp.GetRequestsMiMemory(), wp.GetTopMiMemory(), wp.GetUsageMemory(), wp.GetLimitsMilliCPU(), wp.GetLimitsMiMemory())
		}
		fmt.Fprintf(w, formatHeader, " ", " ", "-----", "----------", "---------", "---", "--------", "----------------", "-----------", "-------------", "--------------------", "---------------", "----------------", "--------------", "-------------------")
		w.Flush()
	}

	if csvFilePrefix != "" {
		file, err := os.Create(csvFilePrefix + "-nohpa.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		header := []string{"Namespace", "Deployment Name", "Replicas", "Expected Replicas", "Up To Date", "Avaliable", "Age", "#Pods ->", "Requests CPU (m)", "TOP CPU (m)", "Usage CPU (%)", "Requests Memory (Mi)", "TOP Memory (Mi)", "Usage Memory (%)", "Limits CPU (m)", "Limitis Memory (Mi)"}
		err = writer.Write(header)
		if err != nil {
			log.Fatal(err)
		}
		for _, deploy := range deploymentWithoutHpa {
			wp := Wrapper{Pods: deploy.Pods}
			line := []string{deploy.Namespace, deploy.Name, strconv.Itoa(deploy.Replicas), strconv.Itoa(deploy.ReplicasExpected), strconv.Itoa(deploy.UpToDate), strconv.Itoa(deploy.Avaliable), deploy.Age, strconv.Itoa(len(deploy.Pods)), strconv.Itoa(wp.GetRequestsMilliCPU()), strconv.Itoa(wp.GetTopMilliCPU()), fmt.Sprintf("%.2f", wp.GetUsageCPU()), strconv.Itoa(wp.GetRequestsMiMemory()), strconv.Itoa(wp.GetTopMiMemory()), fmt.Sprintf("%.2f", wp.GetUsageMemory()), strconv.Itoa(wp.GetLimitsMilliCPU()), strconv.Itoa(wp.GetLimitsMiMemory())}
			err := writer.Write(line)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func printNodesTab(podList []Pod, csvFilePrefix string, debug bool) {
	result := Wrapper{Pods: podList}
	podsInNodes := make(map[string][]Pod)
	for _, pod := range result.Pods {
		nodeName := pod.Spec.NodeName
		if pods, ok := podsInNodes[nodeName]; ok {
			podsInNodes[nodeName] = append(pods, pod)
		} else {
			podsInNodes[nodeName] = []Pod{pod}
		}
	}

	if csvFilePrefix == "" || debug {
		fmt.Println("\n\nNODEs SNAPSHOT:")
		formatHeader := "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n"
		formatValues := "%v\t%v\t%vm\t%vm\t%0.2f%%\t%vMi\t%vMi\t%0.2f%%\t%vm\t%vMi\n"
		tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintf(tw, formatHeader, "Node", "Num Pods In Node", "Requests CPU (m)", "TOP CPU (m)", "Usage CPU (%)", "Requests Memory (Mi)", "TOP Memory (Mi)", "Usage Memory (%)", "Limits CPU (m)", "Limitis Memory (Mi)")
		fmt.Fprintf(tw, formatHeader, "----", "----------------", "----------------", "-----------", "-------------", "--------------------", "---------------", "----------------", "--------------", "-------------------")
		min := 999
		max := 0
		total := 0
		for nodeName, pods := range podsInNodes {
			nPods := len(pods)
			total += nPods
			if nPods > max {
				max = nPods
			}
			if min > nPods {
				min = nPods
			}
			w := Wrapper{Pods: pods}
			fmt.Fprintf(tw, formatValues, nodeName, nPods, w.GetRequestsMilliCPU(), w.GetTopMilliCPU(), w.GetUsageCPU(), w.GetRequestsMiMemory(), w.GetTopMiMemory(), w.GetUsageMemory(), w.GetLimitsMilliCPU(), w.GetLimitsMiMemory())
		}
		avg := 0
		if len(podsInNodes) > 0 {
			avg = total / len(podsInNodes)
		} else {
			min = 0
		}
		fmt.Fprintf(tw, formatHeader, " ", "----------------", "----------------", "-----------", "-------------", "--------------------", "---------------", "----------------", "--------------", "-------------------")
		summary := fmt.Sprintf("Min:%d/Max:%d/Avg:%d", min, max, avg)
		fmt.Fprintf(tw, formatValues, " ", summary, result.GetRequestsMilliCPU(), result.GetTopMilliCPU(), result.GetUsageCPU(), result.GetRequestsMiMemory(), result.GetTopMiMemory(), result.GetUsageMemory(), result.GetLimitsMilliCPU(), result.GetLimitsMiMemory())
		tw.Flush()

		if debug {
			fmt.Println()
			fmt.Println("---------------------------------------------")
			fmt.Println("[debug] PODS IN EACH NODE: ")
			for nodeName, pods := range podsInNodes {
				fmt.Printf(" - %s\n   [ ", nodeName)
				for _, pod := range pods {
					fmt.Printf("%s   ", pod.GetPodKey())
				}
				fmt.Println("]")
			}
			fmt.Println("---------------------------------------------")
		}
	}

	if csvFilePrefix != "" {
		file, err := os.Create(csvFilePrefix + "-nodes.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		header := []string{"Node", "Num Pods In Node", "Requests CPU (m)", "TOP CPU (m)", "Usage CPU (%)", "Requests Memory (Mi)", "TOP Memory (Mi)", "Usage Memory (%)", "Limits CPU (m)", "Limitis Memory (Mi)"}
		err = writer.Write(header)
		if err != nil {
			log.Fatal(err)
		}
		for nodeName, pods := range podsInNodes {
			nPods := len(pods)
			w := Wrapper{Pods: pods}
			line := []string{nodeName, strconv.Itoa(nPods), strconv.Itoa(w.GetRequestsMilliCPU()), strconv.Itoa(w.GetTopMilliCPU()), fmt.Sprintf("%.2f", w.GetUsageCPU()), strconv.Itoa(w.GetRequestsMiMemory()), strconv.Itoa(w.GetTopMiMemory()), fmt.Sprintf("%.2f", w.GetUsageMemory()), strconv.Itoa(w.GetLimitsMilliCPU()), strconv.Itoa(w.GetLimitsMiMemory())}
			err := writer.Write(line)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// Wrapper contains a list of pods
type Wrapper struct {
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
	if top == 0 && requests != 0 {
		return float32(0)
	} else if requests == 0 {
		return float32(100)
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
	if top == 0 && requests != 0 {
		return float32(0)
	} else if requests == 0 {
		return float32(100)
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

func filterPod2(podList []Pod, test func(Pod) bool) (ret []Pod) {
	for _, pod := range podList {
		if test(pod) {
			ret = append(ret, pod)
		}
	}
	return
}

func filterHpa(hpaList []Hpa, test func(Hpa) bool) (ret []Hpa) {
	for _, hpa := range hpaList {
		if test(hpa) {
			ret = append(ret, hpa)
		}
	}
	return
}

func filterDeployment(deploymentList []Deployment, test func(Deployment) bool) (ret []Deployment) {
	for _, deploy := range deploymentList {
		if test(deploy) {
			ret = append(ret, deploy)
		}
	}
	return
}
