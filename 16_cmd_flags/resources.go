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
	result := reg.FindStringSubmatch(r.Memory)
	memory, _ := strconv.Atoi(result[1])
	suffix := result[2]

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

// PodResources struct
type PodResources struct {
	Metadata struct {
		Name string
	}
	Spec struct {
		NodeName   string
		Containers []struct {
			Resources struct {
				Requests Resource
				Limits   Resource
			}
		}
	}
}

// GetRequestsMilliCPU total
func (pr PodResources) GetRequestsMilliCPU() int {
	total := 0
	for _, c := range pr.Spec.Containers {
		total += c.Resources.Requests.GetMilliCPU()
	}
	return total
}

// GetRequestsMiMemory total
func (pr PodResources) GetRequestsMiMemory() int {
	total := 0
	for _, c := range pr.Spec.Containers {
		total += c.Resources.Requests.GetMiMemory()
	}
	return total
}

// GetLimitsMilliCPU total
func (pr PodResources) GetLimitsMilliCPU() int {
	total := 0
	for _, c := range pr.Spec.Containers {
		total += c.Resources.Limits.GetMilliCPU()
	}
	return total
}

// GetLimitsMiMemory total
func (pr PodResources) GetLimitsMiMemory() int {
	total := 0
	for _, c := range pr.Spec.Containers {
		total += c.Resources.Limits.GetMiMemory()
	}
	return total
}

// GetPodResources executes kubectl get pod command
func GetPodResources(name string, ns string) []PodResources {
	cmd := fmt.Sprintf("kubectl get pod %s -n %s -o json", name, ns)
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	json := string(out)
	fmt.Printf("JSON:\n%s\n", json)
	return GetResources(json)
}

// GetResources based on a json string
func GetResources(str string) PodResources {
	pod := PodResources{}
	err2 := json.Unmarshal([]byte(str), &pod)
	if err2 != nil {
		fmt.Println(err2.Error())
	}
	return pod
}
