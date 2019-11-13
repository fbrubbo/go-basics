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
	csv := flag.Bool("csv-output", false, "Save the result to files with format 'kubectl-snapshot-<date>-<pods|hpas|nodes|all>.csv'")
	debug := flag.Bool("debug", false, "Show debug info")
	noHeaders := flag.Bool("no-headers", false, "When true, remove filters")
	flag.Parse()
	printFlags(*p, *d, *n, *v, *show, *csv, *debug, *noHeaders)

	if *v || *debug {
		fmt.Println("Plugin Version: ", version)
		os.Exit(0)
	}

	podList := RetrievePods(*n)
	var resultWrapper Wrapper
	if *p != "" {
		resultWrapper = filterPod(podList, *p)
	} else if *d != "" {
		resultWrapper = filterDeployment(podList, *d)
	} else {
		resultWrapper = Wrapper{Type: "All Pods", Pods: podList}
	}

	hpaList := RetrieveHpas(*n, resultWrapper.Pods)
	if *p != "" {
		hpaList = filterHpa(hpaList, func(h Hpa) bool { return h.ContainsPod(*p) })
	} else if *d != "" {
		hpaList = filterHpa(hpaList, func(h Hpa) bool { return h.RefToDeployment(*d) })
	}

	csvFilePrefix := ""
	if *csv {
		now := time.Now()
		csvFilePrefix = now.Format("kubectl-snapshot-2006-01-02-1504")
	}
	switch *show {
	case "pod":
	case "pods":
		printPodsTab(resultWrapper, csvFilePrefix)
	case "hpa":
	case "hpas":
		printHpaTab(hpaList, csvFilePrefix)
	case "node":
	case "nodes":
		printNodesTab(resultWrapper, *debug, csvFilePrefix)
	default:
		printPodsTab(resultWrapper, csvFilePrefix)
		printHpaTab(hpaList, csvFilePrefix)
		printNodesTab(resultWrapper, *debug, csvFilePrefix)
	}

}

func printFlags(p string, d string, n string, v bool, show string, csv bool, debug bool, noHeaders bool) {
	if debug {
		fmt.Println("---------------------------------------------")
		fmt.Println("[debug] FLAGS: ")
		fmt.Println("   -p [POD] is: ", p)
		fmt.Println("   -d [DEPLOYMENT] is: ", d)
		fmt.Println("   -o [NAMESPACE] is: ", n)
		fmt.Println("   -v [VERSION] is: ", v)
		fmt.Println("   -print [PRINT IN STANDARD OUTPUT] is: ", show)
		fmt.Println("   -csv-output [SAVE TO FILES] is: ", csv)
		fmt.Println("   -no-headers [NO HEADERS] is: ", noHeaders)
		fmt.Println("---------------------------------------------")
		fmt.Println()
	}
}

func printPodsTab(result Wrapper, csvFilePrefix string) {
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

func printHpaTab(hpaList []Hpa, csvFilePrefix string) {
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

func printNodesTab(result Wrapper, debug bool, csvFilePrefix string) {
	podsInNodes := make(map[string][]Pod)
	for _, pod := range result.Pods {
		nodeName := pod.Spec.NodeName
		if pods, ok := podsInNodes[nodeName]; ok {
			podsInNodes[nodeName] = append(pods, pod)
		} else {
			podsInNodes[nodeName] = []Pod{pod}
		}
	}
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

func filterPod(podList []Pod, p string) Wrapper {
	pods := make(map[string]Pod)
	for _, pod := range podList {
		pods[pod.Metadata.Name] = pod
	}
	return Wrapper{Type: "Pod", Pods: []Pod{pods[p]}}
}

func filterDeployment(podList []Pod, d string) Wrapper {
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

func filterHpa(hpaList []Hpa, test func(Hpa) bool) (ret []Hpa) {
	for _, hpa := range hpaList {
		if test(hpa) {
			ret = append(ret, hpa)
		}
	}
	return
}
