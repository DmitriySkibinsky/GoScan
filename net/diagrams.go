package net

import (
	"fmt"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"log"
	"strings"
	"time"
)

func ScanProcessNetworkPer() map[int32]uint64 {

	connections, err := net.Connections("all")
	if err != nil {
		log.Fatalf("Ошибка при получении соединений: %v", err)
	}

	connectionMap := make(map[int32][]net.ConnectionStat)
	for _, conn := range connections {
		if conn.Pid > 0 { // Игнорируем PID = 0
			connectionMap[conn.Pid] = append(connectionMap[conn.Pid], conn)
		}
	}

	startCounters, err := net.IOCounters(true)
	if err != nil {
		log.Fatalf("Ошибка при получении начальной статистики: %v", err)
	}

	time.Sleep(5 * time.Second) // Ждем 5 секунд для сбора данных

	endCounters, err := net.IOCounters(true)
	if err != nil {
		log.Fatalf("Ошибка при получении конечной статистики: %v", err)
	}

	// Карта для хранения трафика по процессам
	bytesSentMap := make(map[int32]uint64)

	for _, conn := range connections {
		if conn.Pid > 0 {
			for i := range startCounters {
				bytesSent := endCounters[i].BytesSent - startCounters[i].BytesSent
				bytesSentMap[conn.Pid] += bytesSent
			}
		}
	}

	return bytesSentMap
}

func formatProcessName(name string) string {
	if len(name) > 20 {
		return name[:17] + "..."
	}
	return name
}

// Функция для вывода таблицы с процентами и прогресс-барами
func DisplayTrafficDiagram() {
	// Получаем данные
	bytesSentMap := ScanProcessNetworkPer()

	// Подсчитываем общий трафик
	var totalBytesSent uint64
	for _, bytes := range bytesSentMap {
		totalBytesSent += bytes
	}

	// Проверяем, есть ли трафик
	if totalBytesSent == 0 {
		fmt.Println("Сетевой трафик не обнаружен.")
		return
	}

	fmt.Printf("\n%-6s | %-20s | %-15s | %-7s | %s\n", "PID", "Process Name", "Traffic (bytes)", "% Usage", "Progress")
	fmt.Println(strings.Repeat("-", 80))

	for pid, bytesSent := range bytesSentMap {
		proc, err := process.NewProcess(pid)
		var procName string
		if err == nil {
			procName, _ = proc.Name()
		} else {
			procName = "Unknown"
		}
		procName = formatProcessName(procName)
		percent := float64(bytesSent) / float64(totalBytesSent) * 100
		barLength := int(percent / 2)
		bar := strings.Repeat("█", barLength) + strings.Repeat("░", 50-barLength)

		// Выводим строку таблицы
		fmt.Printf("%-6d | %-20s | %-15d | %-7.2f%% | %s\n", pid, procName, bytesSent, percent, bar)
	}
}
