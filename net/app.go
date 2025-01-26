package net

import (
	"fmt"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"log"
	"os"
	"text/tabwriter"
)

func ScanProcessNetworkAppUsage() {
	// Получаем список всех процессов
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Ошибка при получении списка процессов: %v", err)
	}

	// Получаем сетевые соединения
	connections, err := net.Connections("all")
	if err != nil {
		log.Fatalf("Ошибка при получении сетевых соединений: %v", err)
	}

	// Создаем карту для хранения информации о соединениях по PID
	connectionMap := make(map[int32][]net.ConnectionStat)
	for _, conn := range connections {
		connectionMap[conn.Pid] = append(connectionMap[conn.Pid], conn)
	}

	// Выводим таблицу с результатами
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PID\tProcess Name\tLocal Address\tRemote Address\tStatus")

	for _, proc := range processes {
		pid := proc.Pid
		if conns, ok := connectionMap[pid]; ok {
			name, err := proc.Name()
			if err != nil {
				name = "Unknown"
			}
			for _, conn := range conns {
				fmt.Fprintf(w, "%d\t%s\t%s:%d\t%s:%d\t%s\n",
					pid,
					name,
					conn.Laddr.IP, conn.Laddr.Port,
					conn.Raddr.IP, conn.Raddr.Port,
					conn.Status,
				)
			}
		}
	}

	w.Flush()
}

func ScanProcessNetworkUsageAppFull() {
	// Получаем список всех процессов
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Ошибка при получении списка процессов: %v", err)
	}

	// Получаем сетевые соединения
	connections, err := net.Connections("all")
	if err != nil {
		log.Fatalf("Ошибка при получении сетевых соединений: %v", err)
	}

	// Создаем карту для хранения информации о соединениях по PID
	connectionMap := make(map[int32][]net.ConnectionStat)
	for _, conn := range connections {
		connectionMap[conn.Pid] = append(connectionMap[conn.Pid], conn)
	}

	// Выводим таблицу с результатами
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PID\tProcess Name\tUser\tCPU Usage\tMemory Usage\tLocal Address\tRemote Address\tStatus\tProtocol")

	for _, proc := range processes {
		pid := proc.Pid
		if conns, ok := connectionMap[pid]; ok {
			name, err := proc.Name()
			if err != nil {
				name = "Unknown"
			}

			// Получаем информацию о пользователе
			username, err := proc.Username()
			if err != nil {
				username = "Unknown"
			}

			// Получаем использование CPU
			cpuPercent, err := proc.CPUPercent()
			if err != nil {
				cpuPercent = 0
			}

			// Получаем использование памяти
			memInfo, err := proc.MemoryInfo()
			memoryUsage := uint64(0)
			if err == nil && memInfo != nil {
				memoryUsage = memInfo.RSS // Resident Set Size (память в RAM)
			}

			for _, conn := range conns {
				// Определяем протокол (TCP/UDP)
				protocol := "Unknown"
				switch conn.Type {
				case 1: // TCP
					protocol = "TCP"
				case 2: // UDP
					protocol = "UDP"
				}

				fmt.Fprintf(w, "%d\t%s\t%s\t%.2f%%\t%d KB\t%s:%d\t%s:%d\t%s\t%s\n",
					pid,
					name,
					username,
					cpuPercent,
					memoryUsage/1024, // Переводим байты в килобайты
					conn.Laddr.IP, conn.Laddr.Port,
					conn.Raddr.IP, conn.Raddr.Port,
					conn.Status,
					protocol,
				)
			}
		}
	}

	w.Flush()
}
