package net

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func route() string {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "netsh wlan show networks mode=bssid")
	} else {
		cmd = exec.Command("iwlist", "wlan0", "scanning")
	}

	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Ошибка при выполнении: %v", err)
	}
	return string(output)
}

func parseAndPrintInfo(output string, fieldsToPrint []string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		for _, field := range fieldsToPrint {
			if strings.HasPrefix(line, field) {
				fmt.Printf("%s: %s\n", field, strings.Split(line, ": ")[1])
			}
		}
	}
	fmt.Println("-----------------------------")
}

func HotSpotScan() {
	output := route()
	fieldsToPrint := []string{"SSID", "BSSID", "Сигнал", "Тип радио", "Диапазон", "Канал"}
	parseAndPrintInfo(output, fieldsToPrint)
}

func HotSpotScanDatailed() {
	output := route()
	fieldsToPrint := []string{"SSID", "BSSID", "Сигнал", "Тип радио", "Диапазон", "Канал", "Тип шифрования", "Проверка подлинности"}
	parseAndPrintInfo(output, fieldsToPrint)
}

func HotSpotScanFull() {
	output := route()
	fieldsToPrint := []string{"SSID", "Тип сети", "Проверка подлинности", "Шифрование", "BSSID", "Сигнал", "Тип радио", "Диапазон", "Канал", "Подключенные станции", "Использование канала", "Средняя доступная емкость", "Поддерживается QoS MSCS", "Поддерживается сопоставление QoS", "Базовая скорость (мбит/с)", "Базовая скорость (мбит/с)"}
	parseAndPrintInfo(output, fieldsToPrint)
}
