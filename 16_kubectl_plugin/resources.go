package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
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
	Name            string
	Namespace       string
	OwnerReferences []struct {
		Kind string
		Name string
	}
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
	return String2MilliCPU(r.CPU)
}

// GetMiMemory returns the memory in Mi
func (r Resource) GetMiMemory() int {
	return String2MiMemory(r.Memory)
}

// GetPodKey returns <namespace>-<pod name>
func (p Pod) GetPodKey() string {
	return p.Metadata.Namespace + "|" + p.Metadata.Name
}

// GetDeploymentdKey returns <namespace>-<pod name>
func (p Pod) GetDeploymentdKey() string {
	return p.Metadata.Namespace + "|" + p.GetDeploymentName()
}

// GetReplicaSetKey returns <namespace>-<pod name>
func (p Pod) GetReplicaSetKey() string {
	return p.Metadata.Namespace + "|" + p.GetReplicaSetName()
}

//senninha-quotation-redis-slave-0
// zoidberg-pentaho-report-1572104400-rklgx

const stafulsetPattern = `(.*)-(\d*)`
const deploymentPattern = `(.*)-([^-]*)-([^-]*)`
const jobPattern = `(.*)-([^-]*)`

// GetDeploymentName should work for most of the cases
func (p Pod) GetDeploymentName() string {
	name := p.Metadata.Name
	var reg *regexp.Regexp
	if match, _ := regexp.MatchString(deploymentPattern, name); match {
		reg, _ = regexp.Compile(deploymentPattern)
	} else if match, _ := regexp.MatchString(stafulsetPattern, name); match {
		reg, _ = regexp.Compile(stafulsetPattern)
	} else if p.Metadata.OwnerReferences != nil && p.Metadata.OwnerReferences[0].Kind == "Job" {
		reg, _ = regexp.Compile(jobPattern)
	}
	result := reg.FindStringSubmatch(name)
	return result[1]
}

// GetReplicaSetName should work for most of the cases
func (p Pod) GetReplicaSetName() string {
	if p.Metadata.OwnerReferences == nil {
		return "<no-references>"
	}
	return p.Metadata.OwnerReferences[0].Name
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
	top := float32(p.GetTopMilliCPU())
	requests := float32(p.GetRequestsMilliCPU())
	if top == 0 && requests != 0 {
		return 0
	} else if requests == 0 {
		return 100
	}
	return top / requests * 100
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
	top := float32(p.GetTopMiMemory())
	requests := float32(p.GetRequestsMiMemory())
	if top == 0 && requests != 0 {
		return 0
	} else if requests == 0 {
		return 100
	}
	return top / requests * 100
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
	err := json.Unmarshal([]byte(str), &pods)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	return pods
}
