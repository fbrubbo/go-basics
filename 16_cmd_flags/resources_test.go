package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
)

type testResource struct {
	res      Resource
	expected int
}

func TestGetMilliCPU(t *testing.T) {
	tests := []testResource{
		testResource{res: Resource{CPU: "130m", Memory: "350"}, expected: 130},
		testResource{res: Resource{CPU: "1", Memory: "450"}, expected: 1000},
		testResource{res: Resource{CPU: "0.5", Memory: "500"}, expected: 500},
		testResource{res: Resource{CPU: "1.64", Memory: "640"}, expected: 1640},
	}

	log.Infof("%+v", tests)

	for i, test := range tests {
		log.Infof("test info %d -> %+v", i, test)
		if result := test.res.GetMilliCPU(); result != test.expected {
			t.Fatalf("Test failed! %d but expected %d", result, test.expected)
		}
	}
}

func TestGetMiMemory(t *testing.T) {
	tests := []testResource{
		testResource{res: Resource{CPU: "130m", Memory: "123Mi"}, expected: 123},
		testResource{res: Resource{CPU: "1", Memory: "129M"}, expected: 123},
		testResource{res: Resource{CPU: "0.5", Memory: "128974848"}, expected: 123},
	}

	log.Infof("%+v", tests)

	for i, test := range tests {
		log.Infof("test info %d -> %+v", i, test)
		if result := test.res.GetMiMemory(); result != test.expected {
			t.Fatalf("Test failed! %d but expected %d", result, test.expected)
		}
	}
}

func TestPodResources(t *testing.T) {
	b, err := ioutil.ReadFile("test.json")
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	pr := GetResources(str)

	expected := 200
	if result := pr.GetRequestsMilliCPU(); result != expected {
		t.Fatalf("Test failed! %d but expected %d", result, expected)
	}
	expected = 192
	if result := pr.GetRequestsMiMemory(); result != expected {
		t.Fatalf("Test failed! %d but expected %d", result, expected)
	}
	expected = 2200
	if result := pr.GetLimitsMilliCPU(); result != expected {
		t.Fatalf("Test failed! %d but expected %d", result, expected)
	}
	expected = 256
	if result := pr.GetLimitsMiMemory(); result != expected {
		t.Fatalf("Test failed! %d but expected %d", result, expected)
	}
}
