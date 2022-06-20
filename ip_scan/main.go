package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func discover(addr string) bool {
	fmt.Println("Testing: ", addr)

	cmd := exec.Command("ping", addr)
	output, _ := cmd.Output()
	fmt.Println(string(output))

	if strings.Contains(string(output), "bytes=") == false {
		return false
	}
	return true
}

func main() {
	//var b4 uint8
	var discovered bool
	var dev_map map[string]bool

	// discover a device

	rand.Seed(time.Now().Unix())
	dev_map = make(map[string]bool)

	f, err := os.Create("addresses.json")
	if err != nil {
		log.Fatal("Failed to create file: ", err)
	}
	defer f.Close()

	for i := 0; i < 256; i++ {
		//mask := uint16(rand.Intn(0xFFFF))

		// local addresses
		addr := "192.168.1"

		// read in least sig byte
		//b4 = uint8(mask & 0xFF)

		addr = addr + "." + strconv.Itoa(i)

		discovered = discover(addr)

		if discovered {
			dev_map[addr] = true
		}
	}

	result, err := json.MarshalIndent(dev_map, "", "    ")
	if err != nil {
		log.Fatal("Failed to marshal: ", err)
	}

	fmt.Fprint(f, string(result))
}
