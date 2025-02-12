package memory

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/process"
)

func ProcTrack() {
	// Создаем таблицу для вывода
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"ID процесса", "Имя процесса", "Использование памяти (MB)", "CPU% использовано",
	})

	// Печатаем данные о процессах и их использовании памяти в реальном времени
	for {
		printProcessUsage(table)
		time.Sleep(1 * time.Second) // Пауза 1 секунда между измерениями
	}
}

func printProcessUsage(table *tablewriter.Table) {
	clearScreen()
	// Получаем список всех процессов
	procs, err := process.Processes()
	if err != nil {
		log.Fatalf("Ошибка при получении списка процессов: %v", err)
	}

	// Очищаем только данные, но не саму таблицу
	table.ClearRows()

	// Обрабатываем каждый процесс
	for _, proc := range procs {
		// Получаем информацию о процессе
		name, err := proc.Name()
		if err != nil {
			continue
		}

		memInfo, err := proc.MemoryInfo()
		if err != nil {
			continue
		}
		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			continue
		}

		// Добавляем данные в таблицу
		table.Append([]string{
			fmt.Sprintf("%d", proc.Pid),
			name,
			fmt.Sprintf("%.2f", float64(memInfo.RSS)/1024/1024),
			fmt.Sprintf("%.2f", cpuPercent),
		})
	}

	table.Render()
}
