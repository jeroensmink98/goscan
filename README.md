# GoScan

GoScan is a Go-based application designed to scan specified hosts and IP ranges for open ports, providing results both in the console with colored output and in a log file.

## Features

- Scans individual hosts and IP ranges (CIDR notation).
- Allows user-defined ports for scanning.
- Outputs results to the console with color-coded statuses.
- Logs detailed results to `scan_results.txt`.

## Prerequisites

- [Go](https://golang.org/dl/) installed on your system.
- [Nmap](https://nmap.org/download.html) installed, as the application utilizes Nmap for scanning.

## Installation

1. Clone the repository:

2. Navigate to the project directory:

   ```bash
   cd goscan
   ```

3. Install the required Go packages:

   ```bash
   go get github.com/Ullaakut/nmap/v2
   go get github.com/fatih/color
   ```

## Configuration

Create a `servers.json` file in the project directory to specify the targets and ports:

```json
{
  "targets": [
    "192.168.1.1",
    "example.com",
    "10.0.0.0/24"
  ],
  "ports": ["22", "80", "443"]
}
```

- `targets`: List of IP addresses, hostnames, or CIDR notations to scan.
- `ports`: List of ports to scan on each target.

## Usage

1. Build the application:

   ```bash
   go build -o goscan
   ```

2. Run the application:

   ```bash
   ./goscan
   ```

The application will read the `servers.json` file, perform the scan, and display the results in the console with colored output. Detailed results will be appended to `scan_results.txt`.

In Nmap scan results, each port on a host is assigned a status that indicates its accessibility and the presence of services.

Common Port Statuses Recognized by Nmap:

- Open: A service is listening on the port, and it is accessible.

- Closed: No service is listening on the port, but it is accessible. This indicates that the port could be opened in the future.

- Filtered: Nmap cannot determine whether the port is open due to packet filtering, typically by a firewall, which prevents probes from reaching the port.

- Unfiltered: The port is accessible, but Nmap cannot determine whether it is open or closed. This state is reported when a port is accessible but Nmap's probes do not provide enough information to determine its status.

- Open|Filtered: Nmap cannot determine whether the port is open or filtered. This occurs when no response is received, leaving the port's status ambiguous.

- Closed|Filtered: Nmap cannot determine whether the port is closed or filtered. This state is less common and indicates ambiguity in the port's status.

## Notes

- Ensure that Nmap is installed and accessible in your system's PATH.
- The application requires appropriate permissions to perform network scans.
- Use this tool responsibly and only scan networks and hosts you have permission to assess.
