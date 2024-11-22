package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Ullaakut/nmap/v2"
	"github.com/fatih/color"
)

// Config represents the structure of the servers.json file
type Config struct {
	Targets []string `json:"targets"`
	Ports   []string `json:"ports"`
}

// expandCIDR expands a CIDR notation into a slice of IP addresses
func expandCIDR(cidr string) ([]string, error) {
	var ips []string
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		ips = append(ips, ip.String())
	}
	// Remove network and broadcast addresses
	if len(ips) > 2 {
		return ips[1 : len(ips)-1], nil
	}
	return ips, nil
}

// incrementIP increments an IP address
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func main() {
	// Read the configuration file
	configFile, err := os.ReadFile("servers.json")
	if err != nil {
		log.Fatalf("Failed to read servers.json: %v", err)
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Failed to parse servers.json: %v", err)
	}

	// Create or open the log file
	logFile, err := os.OpenFile("scan_results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Set up logging to both the file and console
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Loading and processing targets...")
	// Prepare the list of targets
	var targets []string
	for _, target := range config.Targets {
		if _, _, err := net.ParseCIDR(target); err == nil {
			// Target is a CIDR notation
			fmt.Printf("Expanding CIDR range: %s\n", target)
			expandedIPs, err := expandCIDR(target)
			if err != nil {
				log.Printf("Failed to expand CIDR %s: %v", target, err)
				continue
			}
			targets = append(targets, expandedIPs...)
		} else {
			// Target is a single IP or hostname
			targets = append(targets, target)
		}
	}

	fmt.Printf("\nScan configuration:\n")
	fmt.Printf("Total targets to scan: %d\n", len(targets))
	fmt.Printf("Ports to scan: %s\n", config.Ports)
	fmt.Println("\nInitiating scan...")

	// Perform the scan
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(targets...),
		nmap.WithPorts(config.Ports...),
		nmap.WithContext(ctx),
	)
	if err != nil {
		log.Fatalf("Failed to create scanner: %v", err)
	}

	fmt.Printf("Starting Nmap scan (timeout: 30 seconds)...\n\n")

	result, warnings, err := scanner.Run()
	if err != nil {
		log.Fatalf("Failed to run scan: %v", err)
	}
	if warnings != nil {
		log.Printf("Warnings: %v", warnings)
	}

	// Process and log the results
	for _, host := range result.Hosts {
		if len(host.Addresses) == 0 {
			continue
		}
		address := host.Addresses[0].Addr
		fmt.Printf("Host: %s\n", address)
		log.Printf("Host: %s\n", address)

		for _, port := range host.Ports {
			state := port.State.String()
			service := port.Service.Name
			output := fmt.Sprintf("  Port %d/%s: %s (%s)\n", port.ID, port.Protocol, state, service)

			// Colorize output based on port state
			switch state {
			case "open":
				color.Green(output)
			case "closed":
				color.Red(output)
			default:
				color.Yellow(output)
			}

			log.Print(output)
		}
		fmt.Println()
		log.Println()
	}
}
