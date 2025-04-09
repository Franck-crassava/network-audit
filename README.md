# 🌐 Network Audit (Go)

A lightweight tool to scan your local network and detect active hosts using ICMP (ping).

## ⚙️ Features

- Pings all hosts on a local subnet (default: 192.168.1.0/24)
- Identifies active IPs and resolves hostnames
- Fast and concurrent execution (goroutines)
- Works on Windows (default `ping -n`) – easily adaptable for Unix

## 🚀 Usage

```bash
go run main.go
```
Or compile:

```bash
go build -o network-audit
./network-audit
```

## 👨‍💻 Author
Franck CRASSAVA – Cybersecurity & Network Architecture Student
