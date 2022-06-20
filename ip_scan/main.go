package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type dev_map struct {
	mu     sync.Mutex
	ip_map map[string]bool
}

func discover(addr string, device_map *dev_map) {
	fmt.Println("Testing: ", addr)

	cmd := exec.Command("ping", addr)
	output, _ := cmd.Output()

	if strings.Contains(string(output), "bytes=") == false {
		return
	}

	device_map.ip_map[addr] = true
}

func main() {
	device_map := &dev_map{
		ip_map: make(map[string]bool),
	}

	f, err := os.Create("addresses.json")
	if err != nil {
		log.Fatal("Failed to create file: ", err)
	}
	defer f.Close()

	var wg sync.WaitGroup

	// ping concurrently
	for i := 0; i < 256; i++ {
		wg.Add(1)
		go func(b4 int) {
			addr := "192.168.1" + "." + strconv.Itoa(b4)
			discover(addr, device_map)
			wg.Done()
		}(i)
	}
	wg.Wait()

	result, err := json.MarshalIndent(device_map.ip_map, "", "    ")
	if err != nil {
		log.Fatal("Failed to marshal: ", err)
	}

	fmt.Fprint(f, string(result))
}
