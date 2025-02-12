package net

import (
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"log"
	"os/exec"
	"strings"
)

func convertToUTF8(input []byte) (string, error) {
	// Преобразуем из Windows-1251 в UTF-8
	decoder := charmap.Windows1251.NewDecoder()
	utf8Bytes, _, err := transform.Bytes(decoder, input)
	if err != nil {
		return "", err
	}
	return string(utf8Bytes), nil
}

func GetDevicesInNetwork() {
	// Выполняем команду arp -a
	cmd := exec.Command("cmd", "/C", "arp -a")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("Ошибка при выполнении команды ARP:", err)
	}

	// Преобразуем вывод в строку с правильной кодировкой (Windows-1251 -> UTF-8)
	outputStr, err := convertToUTF8(output)
	if err != nil {
		log.Fatal("Ошибка при преобразовании кодировки:", err)
	}

	// Отладочная информация: выводим сырой вывод команды
	fmt.Println("Вывод команды arp -a:")
	fmt.Println(outputStr)

	// Обрабатываем вывод команды
	lines := strings.Split(outputStr, "\n")

	// Проходим по строкам и извлекаем нужные данные
	for _, line := range lines {
		// Фильтруем строки с информацией о MAC-адресах
		if strings.Contains(line, "dynamic") {
			parts := strings.Fields(line)
			if len(parts) > 3 {
				ip := parts[0]  // IP-адрес
				mac := parts[1] // MAC-адрес
				fmt.Printf("IP: %s, MAC: %s\n", ip, mac)
			}
		}
	}
}
