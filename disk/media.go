package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/tabwriter"
)

// workerPool задает количество горутин для обработки файлов
const workerPool = 8
const maxPathLength = 55 // Максимальная длина пути в таблице

// findFilesByExtension ищет файлы с заданным расширением асинхронно
func findFilesByExtension(root, ext string) ([]string, error) {
	filesChan := make(chan string, 100) // Канал для передачи найденных файлов
	var wg sync.WaitGroup

	// Запускаем обработку директорий в отдельной горутине
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Ошибка доступа к %s: %v\n", path, err)
				return nil
			}
			if !info.IsDir() && filepath.Ext(path) == "."+ext {
				relPath, err := filepath.Rel(root, path) // Относительный путь
				if err == nil {
					filesChan <- relPath
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Ошибка при сканировании директорий: %v\n", err)
		}
		close(filesChan) // Закрываем канал после завершения обхода
	}()

	// Обрабатываем файлы в нескольких горутинах (worker pool)
	var files []string
	var mu sync.Mutex

	// Запускаем worker-горутины
	for i := 0; i < workerPool; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range filesChan {
				mu.Lock()
				files = append(files, file) // Потокобезопасное добавление
				mu.Unlock()
			}
		}()
	}

	wg.Wait() // Ждем завершения всех горутин
	return files, nil
}

// shortenPath обрезает путь до maxPathLength символов
func shortenPath(path string) string {
	runes := []rune(path) // Переводим в срез рун, чтобы избежать проблем с Unicode
	if len(runes) > maxPathLength {
		return string(runes[:maxPathLength-3]) + "..." // Добавляем троеточие
	}
	return path
}

// printFilesTable выводит таблицу с информацией о файлах
func printFilesTable(root string, files []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Размер (MB)\tПуть\tНазвание файла")

	var wg sync.WaitGroup
	fileChan := make(chan string, len(files))
	results := make([]string, 0, len(files))
	var mu sync.Mutex

	// Заполняем канал путями файлов
	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)

	// Создаем worker-горутины для получения размеров файлов
	for i := 0; i < workerPool; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				fullPath := filepath.Join(root, file)
				info, err := os.Stat(fullPath)
				if err != nil {
					fmt.Printf("Ошибка при получении информации о файле %s: %v\n", fullPath, err)
					continue
				}
				sizeMB := float64(info.Size()) / (1024 * 1024) // Конвертируем в MB
				shortPath := shortenPath(filepath.Dir(file))   // Обрезаем путь до 35 символов

				// Добавляем в общий список результатов (синхронизация через мьютекс)
				mu.Lock()
				results = append(results, fmt.Sprintf("%.2f\t%s\t%s\n", sizeMB, shortPath, filepath.Base(file)))
				mu.Unlock()
			}
		}()
	}

	wg.Wait() // Ждем завершения всех воркеров

	// Выводим таблицу одним потоком
	for _, line := range results {
		fmt.Fprint(w, line)
	}

	w.Flush() // Выводим таблицу
}

// Media запускает поиск файлов и их вывод в таблицу
func Media(root, ext string) {
	files, err := findFilesByExtension(root, ext)
	if err != nil {
		fmt.Printf("Ошибка при поиске файлов: %v\n", err)
		return
	}

	printFilesTable(root, files)
}
