package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestBuildHpaList(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/hpa.txt")
	if err != nil {
		log.Fatal(err)
	}
	data := string(b)

	hpas := buildHpaList(data, "")
	ex := 18
	if l := len(hpas); l != ex {
		t.Fatalf("Test failed! found %d expected %d", l, ex)
	}
	for _, hpa := range hpas {
		if hpa.Namespace == "" || hpa.Name == "" || hpa.ReferenceKind == "" || hpa.ReferenceName == "" || hpa.Age == "" {
			t.Fatalf("Test failed! hpa must have all info")
		}
	}
}

// func TestBuildhpaTopDefaultNamespace(t *testing.T) {
// 	b, err := ioutil.ReadFile("test-data/top-many-hpas.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	data := string(b)

// 	hpas := buildTopList(data, "default")
// 	ex := 23
// 	if l := len(hpas); l != ex {
// 		t.Fatalf("Test failed! found %d expected %d", l, ex)
// 	}
// }

// func TestBuildOnehpaTop(t *testing.T) {
// 	b, err := ioutil.ReadFile("test-data/top-one-hpa.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	data := string(b)

// 	top := buildTopList(data, "")[0]
// 	expectedCPU := 32
// 	if cpu := top.GetMilliCPU(); cpu != expectedCPU {
// 		t.Fatalf("Test failed! %d but expected %d", cpu, expectedCPU)
// 	}
// 	expectedMemory := 25
// 	if mem := top.GetMiMemory(); mem != expectedMemory {
// 		t.Fatalf("Test failed! %d but expected %d", mem, expectedMemory)
// 	}
// }
