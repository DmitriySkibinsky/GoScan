package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
)

// getTempDirs возвращает список стандартных временных директорий в Windows
func getTempDirs() []string {
	username := os.Getenv("USERNAME")
	return []string{
		"C:\\Windows\\Temp",
		fmt.Sprintf("C:\\Users\\%s\\AppData\\Local\\Temp", username),
	}
}

// findTempFiles ищет файлы в указанных временных директориях
func findTempFiles() []map[string]interface{} {
	var tempFiles []map[string]interface{}
	for _, dir := range getTempDirs() {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Пропускаем ошибки доступа
			}
			if !info.IsDir() && info.Size() >= 8 {
				tempFiles = append(tempFiles, map[string]interface{}{
					"time": info.ModTime().Format("2006-01-02 15:04:05"),
					"name": info.Name(),
					"size": info.Size(),
				})
			}
			return nil
		})
	}
	return tempFiles
}

func Temp() {
	tempFiles := findTempFiles()
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "Время\tИмя\tРазмер (байт)")
	fmt.Fprintln(w, "-------------------------------------")
	for _, file := range tempFiles {
		fmt.Fprintf(w, "%s\t%s\t%d\n", file["time"], file["name"], file["size"])
	}
	w.Flush()
}
