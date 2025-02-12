package disk

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"log"
)

// printDiskInfo выводит информацию о диске
func printDiskInfo(path, fsType string, full bool) {
	usage, err := disk.Usage(path)
	if err != nil {
		fmt.Printf("Ошибка при получении информации о диске %s: %v\n", path, err)
		return
	}

	fmt.Printf("Диск: %s\n", path)
	fmt.Printf("Файловая система: %s\n", fsType)
	fmt.Printf("Общий размер: %.2f GB\n", float64(usage.Total)/(1024*1024*1024))
	fmt.Printf("Свободное место: %.2f GB\n", float64(usage.Free)/(1024*1024*1024))
	fmt.Printf("Занятое место: %.2f GB (%.2f%%)\n", float64(usage.Used)/(1024*1024*1024), usage.UsedPercent)

	if full {
		fmt.Printf("Всего inodes: %d\n", usage.InodesTotal)
		fmt.Printf("Свободных inodes: %d\n", usage.InodesFree)
		fmt.Printf("Использованных inodes: %d (%.2f%%)\n", usage.InodesUsed, usage.InodesUsedPercent)
		fmt.Printf("Тип файловой системы: %s\n", usage.Fstype)
	}
}

// getAllDisksInfo получает информацию обо всех дисках в системе
func getAllDisksInfo(full bool) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Fatalf("Ошибка при получении списка дисков: %v", err)
	}

	for _, p := range partitions {
		printDiskInfo(p.Mountpoint, p.Fstype, full)
		fmt.Println("-----------------------------")
	}
}

// Fs выводит основную информацию о дисках
func Fs() {
	fmt.Println("=== Информация о дисках ===")
	getAllDisksInfo(false)
}

// FsFull выводит полную информацию о дисках
func FsFull() {
	fmt.Println("=== Полная информация о дисках ===")
	getAllDisksInfo(true)
}

// починить
func GetDiskIO() {
	ctx := context.Background()

	// Получаем данные обо всех доступных дисках (без аргумента)
	ioCounters, err := disk.IOCountersWithContext(ctx)
	if err != nil {
		log.Fatalf("Ошибка получения информации о диске: %v", err)
	}

	if len(ioCounters) == 0 {
		log.Println("Нет данных о дисках")
		return
	}

	for name, io := range ioCounters {
		fmt.Printf("Диск: %s\n", name)
		fmt.Printf("Число чтений: %d\n", io.ReadCount)
		fmt.Printf("Число записей: %d\n", io.WriteCount)
		fmt.Printf("Количество прочитанных байт: %d\n", io.ReadBytes)
		fmt.Printf("Количество записанных байт: %d\n", io.WriteBytes)
		fmt.Printf("Время выполнения операций чтения: %d мс\n", io.ReadTime)
		fmt.Printf("Время выполнения операций записи: %d мс\n", io.WriteTime)
		fmt.Println("---------------------------")
	}
}
