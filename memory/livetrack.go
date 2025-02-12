package memory

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/mem"
)

func LiveTrackMem() {
	// Очищаем экран перед началом вывода
	clearScreen()

	// Создаем таблицу для вывода
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Параметр", "Значение"})

	// Печатаем данные о памяти в реальном времени
	for {
		printMemoryUsage(table)
		time.Sleep(1 * time.Second) // Пауза 1 секунда между измерениями
		clearScreen()               // Очищаем экран перед следующим выводом
	}
}

func printMemoryUsage(table *tablewriter.Table) {
	// Получаем информацию о памяти
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("Ошибка при получении информации о памяти: %v", err)
	}

	// Очищаем таблицу перед выводом новых данных
	table.ClearRows()

	// Рассчитываем процент использования свопа
	swapUsed := v.SwapTotal - v.SwapFree
	swapUsedPercent := 0.0
	if v.SwapTotal > 0 {
		swapUsedPercent = float64(swapUsed) / float64(v.SwapTotal) * 100
	}

	// Добавляем данные в таблицу
	table.Append([]string{"Общий объем памяти", fmt.Sprintf("%.2f GB", float64(v.Total)/1024/1024/1024)})
	table.Append([]string{"Свободная память", fmt.Sprintf("%.2f GB", float64(v.Available)/1024/1024/1024)})
	table.Append([]string{"Используемая память", fmt.Sprintf("%.2f GB (%.2f%%)", float64(v.Used)/1024/1024/1024, v.UsedPercent)})
	table.Append([]string{"Кэшированная память", fmt.Sprintf("%.2f GB", float64(v.Cached)/1024/1024/1024)})
	table.Append([]string{"Размер свопа", fmt.Sprintf("%.2f GB", float64(v.SwapTotal)/1024/1024/1024)})
	table.Append([]string{"Использование свопа", fmt.Sprintf("%.2f GB (%.2f%%)", float64(swapUsed)/1024/1024/1024, swapUsedPercent)})

	// Печатаем таблицу
	table.Render()
}

// clearScreen очищает экран, чтобы обновить вывод
func clearScreen() {
	// Используем ANSI escape последовательность для очистки экрана
	// Это работает в большинстве современных терминалов
	fmt.Print("\033[H\033[2J")
}

// Проверка на Windows
func isWindows() bool {
	return os.PathSeparator == '\\'
}
