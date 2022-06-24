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
	"time"
)

type dev_map struct {
	mu     sync.Mutex // most likely don't need because each goroutine is working with a different address
	ip_map map[string]bool
}

func discover(addr string, device_map *dev_map) {
	fmt.Println("Testing: ", addr)

	cmd := exec.Command("ping", addr)
	output, _ := cmd.Output()

	if strings.Contains(string(output), "bytes=") == false {
		return
	}

	fmt.Println("Connected: ", addr)
	device_map.ip_map[addr] = true
}

var num_workers = 255

func main() {
	device_map := &dev_map{
		ip_map: make(map[string]bool),
	}

	ip_chan := make(chan string)

	f, err := os.Create("addresses.json")
	if err != nil {
		log.Fatal("Failed to create file: ", err)
	}
	defer f.Close()

	// initialize workers
	for i := 0; i < num_workers; i++ {
		go func(worker_id int) {
			for ip := range ip_chan {
				discover(ip, device_map)
			}
		}(i + 1)
	}

	var wg2 sync.WaitGroup

	// ping concurrently
	for i := 0; i < 256; i++ {
		wg2.Add(1)
		go func(b4 int) {
			addr := "192.168.1." + strconv.Itoa(b4)
			ip_chan <- addr
			wg2.Done()
		}(i)
	}

	wg2.Wait()
	close(ip_chan)

	time.Sleep(time.Second * 1)
	result, err := json.MarshalIndent(device_map.ip_map, "", "    ")
	if err != nil {
		log.Fatal("Failed to marshal: ", err)
	}

	fmt.Fprint(f, string(result))
}
