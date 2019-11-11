package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// PodList struct
type PodList struct {
	Items []Pod
}

// Pod struct
type Pod struct {
	Metadata Metadata
	Spec     Spec
	Top      Top
}

// Metadata struct
type Metadata struct {
	Name      string
	Namespace string
}

// Spec struct
type Spec struct {
	NodeName   string
	Containers []struct {
		Name      string
		Resources struct {
			Requests Resource
			Limits   Resource
		}
	}
}

// Resource struct
type Resource struct {
	CPU    string
	Memory string
}

// GetMilliCPU returns the CPU in MilliCPU
func (r Resource) GetMilliCPU() int {
	if strings.Contains(r.CPU, "m") {
		str := strings.ReplaceAll(r.CPU, "m", "")
		milli, _ := strconv.Atoi(str)
		return milli
	}
	cpu, _ := strconv.ParseFloat(r.CPU, 64)
	milli := (int)(cpu * 1000)
	return milli
}

// GetMiMemory returns the memory in Mi
func (r Resource) GetMiMemory() int {
	reg, _ := regexp.Compile(`(\d*)(.*)`)
	groups := reg.FindStringSubmatch(r.Memory)
	memory, _ := strconv.Atoi(groups[1])
	suffix := groups[2]

	switch suffix {
	case "G":
		// http://extraconversion.com/data-storage-conversion-table/gigabytes-to-mebibytes.html
		return int(math.Round(float64(memory) * 953.67431640625))
	case "Gi":
		// http://extraconversion.com/data-storage-conversion-table/gibibytes-to-mebibytes.html
		return memory * 1024
	case "M":
		// http://extraconversion.com/data-storage-conversion-table/megabytes-to-mebibytes.html
		return int(math.Round(float64(memory) * 0.9537))
	case "Mi":
		return memory
	default:
		// http://extraconversion.com/data-storage-conversion-table/bytes-to-mebibytes.html
		return int(math.Round(float64(memory) * 9.53674E-7))
	}

	/*
		TODO:

		Limits and requests for memory are measured in bytes.
		You can express memory as a plain integer or as a fixed-point integer using one of these suffixes: E, P, T, G, M, K.
		You can also use the power-of-two equivalents: Ei, Pi, Ti, Gi, Mi, Ki. For example, the following represent roughly the same value:

		128974848, 129e6, 129M, 123Mi
	*/
}

// GetPodKey returns <namespace>-<pod name>
func (p Pod) GetPodKey() string {
	return p.Metadata.Namespace + "|" + p.Metadata.Name
}

// GetDeploymentdKey returns <namespace>-<pod name>
func (p Pod) GetDeploymentdKey() string {
	return p.Metadata.Namespace + "|" + p.GetDeploymentName()
}

// GetDeploymentName should work for most of the cases
func (p Pod) GetDeploymentName() string {
	reg, _ := regexp.Compile(`(.*)-([^-]*)-([^-]*)`)
	result := reg.FindStringSubmatch(p.Metadata.Name)
	return result[1]
}

// GetRequestsMilliCPU total
func (p Pod) GetRequestsMilliCPU() int {
	total := 0
	for _, c := range p.Spec.Containers {
		total += c.Resources.Requests.GetMilliCPU()
	}
	return total
}

// GetTopMilliCPU total
func (p Pod) GetTopMilliCPU() int {
	return p.Top.GetMilliCPU()
}

// GetUsageCPU %
func (p Pod) GetUsageCPU() float32 {
	return float32(p.GetTopMilliCPU()) / float32(p.GetRequestsMilliCPU()) * 100
}

// GetRequestsMiMemory total
func (p Pod) GetRequestsMiMemory() int {
	total := 0
	for _, c := range p.Spec.Containers {
		total += c.Resources.Requests.GetMiMemory()
	}
	return total
}

// GetTopMiMemory total
func (p Pod) GetTopMiMemory() int {
	return p.Top.GetMiMemory()
}

// GetUsageMemory %
func (p Pod) GetUsageMemory() float32 {
	return float32(p.GetTopMiMemory()) / float32(p.GetRequestsMiMemory()) * 100
}

// GetLimitsMilliCPU total
func (p Pod) GetLimitsMilliCPU() int {
	total := 0
	for _, c := range p.Spec.Containers {
		total += c.Resources.Limits.GetMilliCPU()
	}
	return total
}

// GetLimitsMiMemory total
func (p Pod) GetLimitsMiMemory() int {
	total := 0
	for _, c := range p.Spec.Containers {
		total += c.Resources.Limits.GetMiMemory()
	}
	return total
}

// RetrievePods executes kubectl get pods command
// if ns is empty, then all namespaces are used
func RetrievePods(ns string) []Pod {
	cmd := buildKubectlCmd(ns)
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	json := string(out)
	pods := buildPodList(json).Items
	return enrichPodsWithTopInfo(pods, ns)
}

func enrichPodsWithTopInfo(pods []Pod, ns string) []Pod {
	var podList []Pod
	topMap := RetrieveTopMap(ns)
	for _, pod := range pods {
		if top, ok := topMap[pod.GetPodKey()]; ok {
			pod.Top = top
		}
		podList = append(podList, pod)
	}
	return podList
}

func buildKubectlCmd(ns string) string {
	cmd := fmt.Sprintf("kubectl get pods --all-namespaces -o json")
	if ns != "" {
		cmd = fmt.Sprintf("kubectl get pods -n %s -o json", ns)
	}
	return cmd
}

func buildPodList(str string) PodList {
	pods := PodList{}
	err2 := json.Unmarshal([]byte(str), &pods)
	if err2 != nil {
		fmt.Println(err2.Error())
	}
	return pods
}
