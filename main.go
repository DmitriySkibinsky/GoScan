package main

import (
	"fmt"
	"os"
	"scan/net"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Используем switch для обработки различных ключей
	switch os.Args[1] {
	case "-m":
		if len(os.Args) < 3 {
			printUsage()
			return
		}
		switch os.Args[2] {
		case "scan":
			net.ScanNetworkTraffic()
		case "stats":
			net.ScanProcessNetworkAppUsage()
		case "fulstats":
			net.ScanProcessNetworkUsageAppFull()
		case "datastats":
			net.ScanProcessNetworkUsageAppData()
		case "piddata":
			net.GetNetworkUsage(6656)
		case "per":
			net.DisplayTrafficDiagram()

		default:
			fmt.Println("Неизвестная команда для ключа -m")
			printUsage()
		}
	case "-h", "--help":
		printUsage()
	default:
		fmt.Println("Неизвестный ключ")
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Использование:")
	fmt.Println("  go run main.go -m scan   - Сканировать сетевой трафик")
	fmt.Println("  go run main.go -m stats  - Показать статистику сетевого трафика")
	fmt.Println("  go run main.go -h        - Показать справку")
}
