package main

import (
	"bufio"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/bradfitz/slice"
)

// Hpa struct
type Hpa struct {
	Namespace     string
	Name          string
	ReferenceKind string
	ReferenceName string
	UsageCPU      int
	Target        int
	MinPods       int
	MaxPods       int
	Replicas      int
	Age           string
	Pods          []Pod
}

// GetDeploymentKey should work for most of the cases
func (h Hpa) GetDeploymentKey() string {
	return h.Namespace + "|" + h.GetReference()
}

// GetReference ..
func (h Hpa) GetReference() string {
	return h.ReferenceKind + "/" + h.ReferenceName
}

// GetUsageAndTarget ..
func (h Hpa) GetUsageAndTarget() string {
	if h.UsageCPU == -1 {
		return "<unknown>/" + strconv.Itoa(h.Target) + "%"
	}
	return strconv.Itoa(h.UsageCPU) + "%/" + strconv.Itoa(h.Target) + "%"
}

// IsDeployment ..
func (h Hpa) IsDeployment() bool {
	return h.ReferenceKind == "Deployment"
}

// RefToDeployment ..
func (h Hpa) RefToDeployment(deployment string) bool {
	return h.IsDeployment() && h.ReferenceName == deployment
}

// ContainsPod ..
func (h Hpa) ContainsPod(pod string) bool {
	set := make(map[string]bool, len(h.Pods))
	for _, pod := range h.Pods {
		set[pod.Metadata.Name] = true
	}
	_, ok := set[pod]
	return ok
}

// RetrieveHpas executes kubectl get hpas command
// if ns is empty, then all namespaces are used
func RetrieveHpas(nsFilter string, podList []Pod) []Hpa {
	cmd := "kubectl get hpa --all-namespaces --no-headers"
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	data := string(out)

	deploymentMap := make(map[string][]Pod)
	replicaSetMap := make(map[string][]Pod)
	for _, pod := range podList {
		deployment := pod.GetDeploymentdKey()
		if deploymentMap[deployment] == nil {
			deploymentMap[deployment] = []Pod{pod}
		} else {
			deploymentMap[deployment] = append(deploymentMap[deployment], pod)
		}

		replicaset := pod.GetReplicaSetKey()
		if replicaSetMap[replicaset] == nil {
			replicaSetMap[replicaset] = []Pod{pod}
		} else {
			replicaSetMap[replicaset] = append(replicaSetMap[replicaset], pod)
		}
	}

	hpaList := buildHpaList(data, nsFilter)
	resultList := []Hpa{}
	for _, hpa := range hpaList {
		key := hpa.Namespace + "|" + hpa.ReferenceName
		if hpa.ReferenceKind == "Deployment" {
			hpa.Pods = deploymentMap[key]
		} else if hpa.ReferenceKind == "ReplicaSet" {
			hpa.Pods = replicaSetMap[key]
		} else {
			// not implemented - return empty pods
		}
		resultList = append(resultList, hpa)
	}
	slice.Sort(resultList, func(i, j int) bool {
		return (resultList[i].Namespace + resultList[i].Name) < (resultList[j].Namespace + resultList[j].Name)
	})
	return resultList
}

func buildHpaList(data string, nsFilter string) []Hpa {
	hpaMap := buildHpaMap(data, nsFilter)
	var hpas []Hpa
	for _, v := range hpaMap {
		hpas = append(hpas, v)
	}
	return hpas
}

func buildHpaMap(data string, nsFilter string) map[string]Hpa {
	scanner := bufio.NewScanner(strings.NewReader(data))
	hpaMap := make(map[string]Hpa)
	for scanner.Scan() {
		//reg, _ := regexp.Compile(`(\S*)\s*(\S*)\s*(\S*)\/(\S*)\s*(\S*)\%?\/(\S*)%\s*(\S*)\s*(\S*)\s*(\S*)\s*(\S*)\s*`)
		reg, _ := regexp.Compile(`(\S*)\s*(\S*)\s*(\S*)\/(\S*)\s*((\S*)%|(<unknown>))\/(\S*)%\s*(\S*)\s*(\S*)\s*(\S*)\s*(\S*)\s*`)
		txt := scanner.Text()
		groups := reg.FindStringSubmatch(txt)
		mamespace := groups[1]
		if nsFilter == "" || nsFilter == mamespace {
			usageCPU := -1
			if groups[6] != "" {
				usageCPU, _ = strconv.Atoi(groups[6])
			}

			target, _ := strconv.Atoi(groups[8])
			minPods, _ := strconv.Atoi(groups[9])
			maxPods, _ := strconv.Atoi(groups[10])
			replicas, _ := strconv.Atoi(groups[11])
			hpa := Hpa{
				Namespace:     mamespace,
				Name:          groups[2],
				ReferenceKind: groups[3],
				ReferenceName: groups[4],
				UsageCPU:      usageCPU,
				Target:        target,
				MinPods:       minPods,
				MaxPods:       maxPods,
				Replicas:      replicas,
				Age:           groups[12],
			}
			hpaMap[hpa.GetDeploymentKey()] = hpa
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return hpaMap
}
