package net

import (
	"fmt"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"log"
	"os"
	"text/tabwriter"
	"time"
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

func ScanProcessNetworkUsageAppData() {
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

	// Получаем начальные значения сетевой активности
	startCounters, err := net.IOCounters(true)
	if err != nil {
		log.Fatalf("Ошибка при получении начальной сетевой статистики: %v", err)
	}

	// Ждем 1 секунду для сбора данных
	time.Sleep(1 * time.Second)

	// Получаем конечные значения сетевой активности
	endCounters, err := net.IOCounters(true)
	if err != nil {
		log.Fatalf("Ошибка при получении конечной сетевой статистики: %v", err)
	}

	// Вычисляем разницу в байтах для каждого интерфейса
	bytesSentMap := make(map[int32]uint64)
	bytesRecvMap := make(map[int32]uint64)
	for i := range startCounters {
		bytesSent := endCounters[i].BytesSent - startCounters[i].BytesSent
		bytesRecv := endCounters[i].BytesRecv - startCounters[i].BytesRecv

		// Суммируем данные для каждого процесса
		for _, conn := range connections {
			if conn.Pid != 0 {
				bytesSentMap[conn.Pid] += bytesSent
				bytesRecvMap[conn.Pid] += bytesRecv
			}
		}
	}

	// Выводим таблицу с результатами
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PID\tProcess Name\tUser\tCPU Usage\tMemory Usage\tLocal Address\tRemote Address\tStatus\tProtocol\tBytes Sent\tBytes Recv")

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

			// Получаем количество отправленных и полученных байт
			bytesSent := bytesSentMap[pid]
			bytesRecv := bytesRecvMap[pid]

			for _, conn := range conns {
				// Определяем протокол (TCP/UDP)
				protocol := "Unknown"
				switch conn.Type {
				case 1: // TCP
					protocol = "TCP"
				case 2: // UDP
					protocol = "UDP"
				}

				fmt.Fprintf(w, "%d\t%s\t%s\t%.2f%%\t%d KB\t%s:%d\t%s:%d\t%s\t%s\t%d\t%d\n",
					pid,
					name,
					username,
					cpuPercent,
					memoryUsage/1024, // Переводим байты в килобайты
					conn.Laddr.IP, conn.Laddr.Port,
					conn.Raddr.IP, conn.Raddr.Port,
					conn.Status,
					protocol,
					bytesSent,
					bytesRecv,
				)
			}
		}
	}

	w.Flush()
}

func GetNetworkUsage(pid int32) {
	// Получаем процесс по PID
	_, err := process.NewProcess(pid)
	if err != nil {
		fmt.Printf("Процесс с PID %d не найден: %v\n", pid, err)
		return
	}

	// Получаем начальные значения сетевой активности
	startCounters, err := net.IOCounters(true)
	if err != nil {
		fmt.Printf("Ошибка при получении сетевой статистики: %v\n", err)
		return
	}

	// Ждем 1 секунду для сбора данных
	time.Sleep(1 * time.Second)

	// Получаем конечные значения сетевой активности
	endCounters, err := net.IOCounters(true)
	if err != nil {
		fmt.Printf("Ошибка при получении сетевой статистики: %v\n", err)
		return
	}

	// Вычисляем разницу в байтах
	var bytesSent, bytesRecv uint64
	for i := range startCounters {
		bytesSent += endCounters[i].BytesSent - startCounters[i].BytesSent
		bytesRecv += endCounters[i].BytesRecv - startCounters[i].BytesRecv
	}

	fmt.Printf("Процесс с PID %d:\n", pid)
	fmt.Printf("Отправлено: %d байт\n", bytesSent)
	fmt.Printf("Получено: %d байт\n", bytesRecv)
}
