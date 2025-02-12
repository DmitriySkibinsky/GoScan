package memory

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"log"
)

func Info() {
	// Получаем информацию о памяти
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("Ошибка при получении информации о памяти: %v", err)
	}

	// Выводим характеристики оперативной памяти
	fmt.Printf("Общий объем памяти: %v MB\n", v.Total/1024/1024)
	fmt.Printf("Свободная память: %v MB\n", v.Available/1024/1024)
	fmt.Printf("Используемая память: %v MB\n", v.Used/1024/1024)
	fmt.Printf("Процент использования: %.2f%%\n", v.UsedPercent)
}

func InfoFull() {
	// Получаем информацию о памяти
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("Ошибка при получении информации о памяти: %v", err)
	}

	// Выводим характеристики оперативной памяти
	fmt.Printf("Общий объем памяти: %v MB\n", v.Total/1024/1024)
	fmt.Printf("Свободная память: %v MB\n", v.Available/1024/1024)
	fmt.Printf("Используемая память: %v MB\n", v.Used/1024/1024)
	fmt.Printf("Процент использования памяти: %.2f%%\n", v.UsedPercent)

	fmt.Printf("Свободная память (Free): %v MB\n", v.Free/1024/1024)
	fmt.Printf("Кэшированная память (Cached): %v KB\n", v.Cached/1024)
	fmt.Printf("Лимит коммитов (CommitLimit): %v KB\n", v.CommitLimit/1024)
	fmt.Printf("Объем зафиксированной памяти (Committed AS): %v KB\n", v.CommittedAS/1024)

	fmt.Printf("Свободная HighMemory (HighFree): %v KB\n", v.HighFree/1024)
	fmt.Printf("Общий объем HighMemory (HighTotal): %v KB\n", v.HighTotal/1024)

	fmt.Printf("Размер HugePage (HugePageSize): %v KB\n", v.HugePageSize)
	fmt.Printf("Свободные HugePages (HugePagesFree): %v\n", v.HugePagesFree)
	fmt.Printf("Общее количество HugePages (HugePagesTotal): %v\n", v.HugePagesTotal)

	fmt.Printf("Неактивная память (Inactive): %v KB\n", v.Inactive/1024)
	fmt.Printf("Очередь очистки (Laundry): %v KB\n", v.Laundry/1024)

	fmt.Printf("Свободная LowMemory (LowFree): %v KB\n", v.LowFree/1024)
	fmt.Printf("Общий объем LowMemory (LowTotal): %v KB\n", v.LowTotal/1024)

	fmt.Printf("Разделяемая память (Shared): %v KB\n", v.Shared/1024)
	fmt.Printf("Память, используемая ядром (Slab): %v KB\n", v.Slab/1024)
	fmt.Printf("Кэш подкачки (SwapCached): %v KB\n", v.SwapCached/1024)

	fmt.Printf("Свободная область VMalloc (VMallocChunk): %v KB\n", v.VMallocChunk/1024)
	fmt.Printf("Используемая VMalloc (VMallocUsed): %v KB\n", v.VMallocUsed/1024)
	fmt.Printf("Общий объем VMalloc (VMallocTotal): %v KB\n", v.VMallocTotal/1024)

	fmt.Printf("Зарезервированная проводная память (Wired): %v KB\n", v.Wired/1024)
	fmt.Printf("Записываемая память (Writeback): %v KB\n", v.Writeback/1024)
	fmt.Printf("Временный буфер записи (WritebackTmp): %v KB\n", v.WritebackTmp/1024)

}
