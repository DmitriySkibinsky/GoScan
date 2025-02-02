package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"
	"time"
)

// getAppSizesInDir возвращает информацию о приложениях в указанной директории
func getAppSizesInDirSort(root string) (map[string]AppInfo, error) {
	apps := make(map[string]AppInfo)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Ошибка доступа к %s: %v\n", path, err)
			return filepath.SkipDir
		}
		if info.IsDir() && path != root {
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
			return filepath.SkipDir
		}
		return nil
	})

	return apps, err
}

// FindUnusedApps находит приложения, которые не использовались более года
func FindUnusedApps() {
	paths := []string{
		`C:\\Program Files`,
		`C:\\Program Files (x86)`,
		`C:\\Users\\` + os.Getenv("USERNAME") + `\\AppData`,
	}
	yearAgo := time.Now().AddDate(0, -6, 0) // Дата год назад

	var unusedApps []AppInfo

	for _, path := range paths {
		fmt.Printf("Анализ директории: %s\n", path)
		apps, err := getAppSizesInDirSort(path)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			continue
		}

		for _, app := range apps {
			if app.LastModified.Before(yearAgo) {
				unusedApps = append(unusedApps, app)
			}
		}
	}

	// Сортируем от более не использовавшегося к менее
	sort.Slice(unusedApps, func(i, j int) bool {
		return unusedApps[i].LastModified.Before(unusedApps[j].LastModified)
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Приложение\tРазмер (MB)\tФайлов\tПоследнее изменение")

	for _, app := range unusedApps {
		fmt.Fprintf(w, "%s\t%.2f\t%d\t%s\n",
			app.Name,
			app.SizeMB,
			app.FileCount,
			app.LastModified.Format("2006-01-02 15:04:05"),
		)
	}

	w.Flush()
}
