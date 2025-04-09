package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
	"runtime"
)

// Structure pour le résultat JSON
type ScanResult struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	PortOpen bool   `json:"port_open"`
}

func pingHost(ip string, wg *sync.WaitGroup, results chan<- ScanResult, port int, subnet string) {
	defer wg.Done()
	var cmd *exec.Cmd

	// Utilise la commande ping selon le système d'exploitation
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "500", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "1", ip)
	}

	output, err := cmd.CombinedOutput()
	if err == nil && strings.Contains(string(output), "TTL=") {
		// Si le ping réussit, on teste le port si un port est spécifié
		portOpen := false
		if port != 0 {
			portOpen = checkPort(ip, port)
		}
		names, err := net.LookupAddr(ip)
		hostname := ip
		if err == nil && len(names) > 0 {
			hostname = names[0]
		}
		results <- ScanResult{IP: ip, Hostname: hostname, PortOpen: portOpen}
	}
}

func checkPort(ip string, port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 1*time.Second)
	if err == nil {
		conn.Close()
		return true
	}
	return false
}

func saveResults(results []ScanResult, outputFile string) {
	// Enregistrer les résultats dans un fichier JSON ou CSV
	ext := strings.ToLower(outputFile[len(outputFile)-4:])
	if ext == ".json" {
		file, _ := os.Create(outputFile)
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		encoder.Encode(results)
	} else if ext == ".csv" {
		file, _ := os.Create(outputFile)
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()
		writer.Write([]string{"IP", "Hostname", "Port Open"})
		for _, result := range results {
			writer.Write([]string{result.IP, result.Hostname, fmt.Sprintf("%t", result.PortOpen)})
		}
	}
}

func main() {
	// Définir les flags CLI
	subnet := flag.String("subnet", "192.168.1.", "Base subnet to scan (e.g., 192.168.1.)")
	port := flag.Int("port", 0, "Check if a specific port is open on active hosts (optional)")
	output := flag.String("output", "", "Save scan results to a file (JSON or CSV based on extension)")
	flag.Parse()

	var wg sync.WaitGroup
	results := make(chan ScanResult, 256)
	var scanResults []ScanResult

	start := time.Now()
	fmt.Println("🔍 Scanning local network...")

	// Scanner les hôtes dans la plage 1-254
	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s%d", *subnet, i)
		wg.Add(1)
		go pingHost(ip, &wg, results, *port, *subnet)
	}

	// Attendre la fin de tous les goroutines
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collecte des résultats
	for res := range results {
		scanResults = append(scanResults, res)
		fmt.Printf("[+] Host Up: %s (%s) Port Open: %t\n", res.IP, res.Hostname, res.PortOpen)
	}

	// Affichage du temps de scan
	fmt.Printf("✅ Scan completed in %s\n", time.Since(start))

	// Sauvegarder les résultats si demandé
	if *output != "" {
		saveResults(scanResults, *output)
		fmt.Printf("✅ Results saved to %s\n", *output)
	}
}
