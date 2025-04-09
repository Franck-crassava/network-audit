package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func pingHost(ip string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	// Utilise la commande ping -n 1 (Windows) ou ping -c 1 (Linux/Mac)
	cmd := exec.Command("ping", "-n", "1", "-w", "500", ip) // Windows syntax
	output, err := cmd.CombinedOutput()
	if err == nil && strings.Contains(string(output), "TTL=") {
		results <- ip
	}
}

func main() {
	baseIP := "192.168.1."
	var wg sync.WaitGroup
	results := make(chan string, 256)

	start := time.Now()
	fmt.Println("🔍 Scanning local network...")

	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s%d", baseIP, i)
		wg.Add(1)
		go pingHost(ip, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		names, err := net.LookupAddr(res)
		if err != nil || len(names) == 0 {
			fmt.Printf("[+] Host Up: %s\n", res)
		} else {
			fmt.Printf("[+] Host Up: %s (%s)\n", res, names[0])
		}
	}

	fmt.Printf("✅ Scan completed in %s\n", time.Since(start))
}

