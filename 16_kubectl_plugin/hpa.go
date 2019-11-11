package main

import (
	"bufio"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
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
}

// // GetDeploymentName should work for most of the cases
// func (t Top) GetDeploymentName() string {
// 	reg, _ := regexp.Compile(`(.*)-([^-]*)-([^-]*)`)
// 	result := reg.FindStringSubmatch(t.Pod)
// 	return result[1]
// }

// // GetMilliCPU total pod cpu
// func (t Top) GetMilliCPU() int {
// 	total := 0
// 	for _, c := range t.Containers {
// 		str := strings.ReplaceAll(c.CPU, "m", "")
// 		milli, _ := strconv.Atoi(str)
// 		total += milli
// 	}
// 	return total
// }

// // GetMiMemory returns the memory in Mi
// func (t Top) GetMiMemory() int {
// 	total := 0
// 	for _, c := range t.Containers {
// 		reg, _ := regexp.Compile(`(\d*)(.*)`)
// 		groups := reg.FindStringSubmatch(c.Memory)
// 		memory, _ := strconv.Atoi(groups[1])
// 		total += memory
// 	}
// 	return total
// }

// RetrieveHpaMap executes kubectl get hpas command
// if ns is empty, then all namespaces are used
// returns key =  hpa.Namespace + "|" + hpa.ReferenceKind + "/" + hpa.ReferenceName
func RetrieveHpaMap(ns string) map[string]Hpa {
	cmd := "kubectl get hpa --all-namespaces --no-headers"
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	data := string(out)
	return buildHpaMap(data, ns)
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
		reg, _ := regexp.Compile(`(\S*)\s*(\S*)\s*(\S*)\/(\S*)\s*(\S*)%?\/(\S*)%\s*(\S*)\s*(\S*)\s*(\S*)\s*(\S*)\s*`)
		txt := scanner.Text()
		groups := reg.FindStringSubmatch(txt)
		mamespace := groups[1]
		if nsFilter == "" || nsFilter == mamespace {
			usageCPU, err := strconv.Atoi(groups[5])
			if err != nil {
				usageCPU = 0
			}
			target, _ := strconv.Atoi(groups[6])
			minPods, _ := strconv.Atoi(groups[7])
			maxPods, _ := strconv.Atoi(groups[8])
			replicas, _ := strconv.Atoi(groups[9])
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
				Age:           groups[10],
			}
			key := hpa.Namespace + "|" + hpa.ReferenceKind + "/" + hpa.ReferenceName
			hpaMap[key] = hpa
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return hpaMap
}
