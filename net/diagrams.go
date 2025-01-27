package net

import (
	"fmt"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"log"
	"time"
)

// Функция для сканирования сетевой активности процессов
func ScanProcessNetworkPer() map[int32]uint64 {
	// Получаем список всех процессов
	_, err := process.Processes()
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
	time.Sleep(5 * time.Second)

	// Получаем конечные значения сетевой активности
	endCounters, err := net.IOCounters(true)
	if err != nil {
		log.Fatalf("Ошибка при получении конечной сетевой статистики: %v", err)
	}

	// Вычисляем разницу в байтах для каждого интерфейса
	bytesSentMap := make(map[int32]uint64)
	for i := range startCounters {
		bytesSent := endCounters[i].BytesSent - startCounters[i].BytesSent

		// Распределяем данные по процессам
		for pid, conns := range connectionMap {
			for _, conn := range conns {
				// Если соединение использует текущий интерфейс, добавляем данные
				if conn.Laddr.IP == startCounters[i].Name {
					bytesSentMap[pid] += bytesSent
				}
			}
		}
	}

	return bytesSentMap
}

// Функция для отображения процентной диаграммы
func DisplayTrafficDiagram() {
	// Получаем данные о сетевой активности
	bytesSentMap := ScanProcessNetworkPer()

	// Находим общее количество отправленных байт
	var totalBytesSent uint64
	for _, bytesSent := range bytesSentMap {
		totalBytesSent += bytesSent
	}

	// Если трафика нет, выводим сообщение
	if totalBytesSent == 0 {
		fmt.Println("Сетевой трафик не обнаружен.")
		return
	}

	// Находим процесс с максимальным трафиком
	var maxPid int32
	var maxBytesSent uint64
	for pid, bytesSent := range bytesSentMap {
		if bytesSent > maxBytesSent {
			maxPid = pid
			maxBytesSent = bytesSent
		}
	}

	// Получаем имя процесса с максимальным трафиком
	proc, err := process.NewProcess(maxPid)
	var maxProcName string
	if err == nil {
		maxProcName, _ = proc.Name()
	} else {
		maxProcName = "Unknown"
	}

	// Выводим диаграмму
	fmt.Println("\nДиаграмма сетевого трафика:")
	fmt.Printf("[%s] %d байт (%.2f%%)\n", maxProcName, maxBytesSent, float64(maxBytesSent)/float64(totalBytesSent)*100)

	// Выводим остальные процессы
	for pid, bytesSent := range bytesSentMap {
		if pid == maxPid {
			continue // Пропускаем процесс с максимальным трафиком
		}

		proc, err := process.NewProcess(pid)
		var procName string
		if err == nil {
			procName, _ = proc.Name()
		} else {
			procName = "Unknown"
		}

		// Вычисляем процент от общего трафика
		percent := float64(bytesSent) / float64(totalBytesSent) * 100
		if percent > 0 {
			fmt.Printf("[%s] %d байт (%.2f%%)\n", procName, bytesSent, percent)
		}
	}
}
