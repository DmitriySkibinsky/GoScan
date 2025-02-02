package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"
)

// AppInfo содержит информацию о приложении
type AppInfo struct {
	Name         string    // Название приложения
	SizeMB       float64   // Размер в мегабайтах
	FileCount    int       // Количество файлов
	LastModified time.Time // Дата последнего изменения
}

// getAppSizesInDir возвращает информацию о приложениях в указанной директории
func getAppSizesInDir(root string) (map[string]AppInfo, error) {
	apps := make(map[string]AppInfo)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Пропускаем директории, к которым нет доступа
			fmt.Printf("Ошибка доступа к %s: %v\n", path, err)
			return filepath.SkipDir
		}
		if info.IsDir() && path != root {
			// Предполагаем, что каждая директория — это приложение
			size, fileCount, lastModified, err := getDirInfo(path)
			if err != nil {
				fmt.Printf("Ошибка при расчете размера %s: %v\n", path, err)
				return filepath.SkipDir
			}
			sizeMB := bytesToMB(size)
			if sizeMB >= 20 { // Игнорируем приложения меньше 20 МБ
				apps[info.Name()] = AppInfo{
					Name:         info.Name(),
					SizeMB:       sizeMB,
					FileCount:    fileCount,
					LastModified: lastModified,
				}
			}
			return filepath.SkipDir // Пропускаем вложенные директории
		}
		return nil
	})
	return apps, err
}

// getDirInfo возвращает размер директории, количество файлов и дату последнего изменения
func getDirInfo(path string) (uint64, int, time.Time, error) {
	var size uint64
	var fileCount int
	var lastModified time.Time

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			// Пропускаем файлы, к которым нет доступа
			fmt.Printf("Ошибка доступа к файлу: %v\n", err)
			return nil
		}
		if !info.IsDir() {
			size += uint64(info.Size())
			fileCount++
			if info.ModTime().After(lastModified) {
				lastModified = info.ModTime()
			}
		}
		return nil
	})
	return size, fileCount, lastModified, err
}

// bytesToMB конвертирует байты в мегабайты
func bytesToMB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024)
}

func AppSize() {
	// Пути к директориям, где могут находиться приложения
	paths := []string{
		`C:\Program Files`,
		`C:\Program Files (x86)`,
		`C:\Users\` + os.Getenv("USERNAME") + `\AppData`,
	}

	// Создаем tabwriter для форматированного вывода таблицы
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Директория\tПриложение\tРазмер (MB)\tФайлов\tПоследнее изменение")

	for _, path := range paths {
		fmt.Printf("Анализ директории: %s\n", path)
		apps, err := getAppSizesInDir(path)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			continue
		}

		for _, app := range apps {
			fmt.Fprintf(w, "%s\t%s\t%.2f\t%d\t%s\n",
				path,
				app.Name,
				app.SizeMB,
				app.FileCount,
				app.LastModified.Format("2006-01-02 15:04:05"),
			)
		}
	}

	// Выводим таблицу
	w.Flush()
}
