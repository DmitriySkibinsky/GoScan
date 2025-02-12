package disk

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// GetDiskTemperature возвращает температуру диска
func GetDiskTemperature(disk string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("smartctl", "-A", disk)
		out, err := cmd.Output()
		if err == nil {
			return parseSmartctlTemperature(string(out)), nil
		}
		// Если smartctl отсутствует, используем wmic
		cmd = exec.Command("wmic", "/namespace:\\root\\wmi", "PATH", "MSStorageDriver_ATAPISmartData", "get", "VendorSpecific")
	} else {
		cmd = exec.Command("smartctl", "-A", disk)
	}

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return parseWmicTemperature(string(out)), nil
}

// parseSmartctlTemperature парсит температуру из вывода smartctl
func parseSmartctlTemperature(output string) string {
	re := regexp.MustCompile(`(?m)^194\s+Temperature_Celsius.*?(\d+)\s*$`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1] + "°C"
	}
	return "Не удалось определить температуру"
}

// parseWmicTemperature парсит температуру из вывода wmic
func parseWmicTemperature(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 2 {
			if temp, err := strconv.Atoi(fields[2]); err == nil {
				return strconv.Itoa(temp) + "°C"
			}
		}
	}
	return "Не удалось определить температуру"
}

// TestDiskSpeed тестирует скорость чтения/записи диска
func testDiskSpeed(disk string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "chcp 65001 > nul && winsat disk -drive "+disk)
	} else {
		cmd = exec.Command("dd", "if=/dev/zero", "of="+disk+"/testfile", "bs=1M", "count=100", "oflag=direct")
	}
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func Metrics(disk string) {

	//if runtime.GOOS == "windows" {
	//	disk = "C"
	//} else {
	//	disk = "/dev/sda"
	//}

	temp, err := GetDiskTemperature(disk)
	if err != nil {
		fmt.Println("Ошибка получения температуры:", err)
	} else {
		fmt.Println("Температура диска:", temp)
	}

	speed, err := testDiskSpeed(disk)
	if err != nil {
		fmt.Println("Ошибка тестирования скорости:", err)
	} else {
		fmt.Println("Результаты тестирования скорости:", speed)
	}
}
