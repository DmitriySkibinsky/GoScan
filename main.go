package main

import (
	"fmt"
	"os"
	"scan/disk"
	"scan/memory"
	"scan/net"
	"strconv"
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
		case "net":
			switch os.Args[3] {
			case "interface":
				net.ScanNetworkTraffic()
			case "stats":
				if len(os.Args) > 3 && os.Args[3] == "-l" {
					net.ScanProcessNetworkUsageAppFull()
				} else {
					net.ScanProcessNetworkAppUsage()
				}
			case "app":
				net.ScanProcessNetworkUsageAppData()
			case "piddata":
				pid, err := strconv.Atoi(os.Args[4])
				if err != nil {
					fmt.Printf("Ошибка преобразования: %v\n", err)
					return
				}
				net.GetNetworkUsage(int32(pid))
			case "per":
				net.DisplayTrafficDiagram()
			case "hotspot":
				if len(os.Args) < 4 {
					net.HotSpotScan()
				} else {
					switch os.Args[3] {
					case "-l":
						net.HotSpotScanDatailed()
					case "-ls", "-sl":
						net.HotSpotScanFull()
					default:
						fmt.Println("Неизвестный аргумент:", os.Args[3])
					}
				}
			case "device":
				net.GetDevicesInNetwork() //дополнить
			default:
				fmt.Println("Неизвестный ключ для режима Сеть")
				printUsage()
			}
		case "disk":
			switch os.Args[3] {
			case "app":
				disk.AppSize() //сделать расширенную версию
			case "unused":
				disk.FindUnusedApps()
			case "temp":
				disk.Temp()
			case "media":
				if len(os.Args) != 5 {
					fmt.Println("Нет ввода директории или расширения файла")
					fmt.Println("goscan -m media <Директория> <Тип файла>")
				} else {
					root := os.Args[4]
					ext := os.Args[5]
					disk.Media(root, ext)
				}
			case "metrics":
				disk.Metrics("C") // исправить
			case "fs":
				if len(os.Args) < 4 {
					disk.Fs()
				} else {
					if os.Args[3] == "-l" {
						disk.FsFull()
					} else {
						fmt.Println("Неверный ключ")
					}
				}

			default:
				fmt.Println("Неизвестный ключ для режима Диск")
				printUsage()
			}
		case "mem":
			switch os.Args[3] {
			case "info":
				if len(os.Args) < 4 {
					memory.Info()
				} else {
					if os.Args[3] == "-l" {
						memory.InfoFull()
					} else {
						fmt.Println("Неверный ключ")
					}
				}
			}
		case "track":
			memory.LiveTrackMem()
		case "proc":
			memory.ProcTrack() //сделать трекинг и сортировку
		}

		//case "meminfo":
		//	memory.InfoFull()
		//case "track":
		//	memory.ProcTrack()
		//default:
		//	fmt.Println("Неизвестная команда для ключа -m")
		//	printUsage()
		//}
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
