package net

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"strings"
	"time"
)

// Найти Wi-Fi интерфейс
func getWirelessInterface() (string, error) {
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		desc := strings.ToLower(iface.Description)
		if strings.Contains(desc, "wireless") || strings.Contains(desc, "wi-fi") || strings.Contains(desc, "беспроводная сеть") {
			return iface.Name, nil
		}
	}
	return "", fmt.Errorf("беспроводной интерфейс не найден")
}

// Отправка ARP-запросов
func sendARPRequests(handle *pcap.Handle, iface net.Interface) error {
	// Получаем IP-адрес и маску подсети интерфейса
	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	var ipnet *net.IPNet
	for _, addr := range addrs {
		if ipnet = addr.(*net.IPNet); ipnet.IP.To4() != nil {
			break
		}
	}
	if ipnet == nil {
		return fmt.Errorf("не удалось получить IPv4-адрес интерфейса")
	}

	// Создаем ARP-пакет
	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, // Broadcast
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: []byte(ipnet.IP.To4()),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}

	// Отправляем ARP-запросы на все IP-адреса в подсети
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	for ip := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		arp.DstProtAddress = []byte(ip)
		gopacket.SerializeLayers(buf, opts, &eth, &arp)
		if err := handle.WritePacketData(buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

// Увеличиваем IP-адрес на 1
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// Конвертация IP-адреса из байтов
func ipToString(ip []byte) string {
	return net.IP(ip).String()
}

// Сканирование сети (ARP-ответы)
func ScanDevices() {
	ifaceName, err := getWirelessInterface()
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
	fmt.Printf("Имя интерфейса: %s\n", ifaceName)

	handle, err := pcap.OpenLive(ifaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatalf("Ошибка при открытии интерфейса: %v", err)
	}
	defer handle.Close()

	// Получаем информацию о сетевом интерфейсе
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatalf("Ошибка при получении информации о интерфейсе: %v", err)
	}

	// Отправляем ARP-запросы
	if err := sendARPRequests(handle, *iface); err != nil {
		log.Fatalf("Ошибка при отправке ARP-запросов: %v", err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	fmt.Println("🔍 Сканирование сети (ожидание ARP-ответов)...\n")
	fmt.Println(" IP-адрес        | MAC-адрес")
	fmt.Println("-----------------|----------------------")

	// Слушаем ARP-ответы в течение 5 секунд
	timeout := time.After(5 * time.Second)
	for {
		select {
		case packet := <-packetSource.Packets():
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer != nil {
				arpPacket, _ := arpLayer.(*layers.ARP)
				fmt.Printf(" %-15s | %02X:%02X:%02X:%02X:%02X:%02X\n",
					ipToString(arpPacket.SourceProtAddress), // Преобразование IP
					arpPacket.SourceHwAddress[0], arpPacket.SourceHwAddress[1], arpPacket.SourceHwAddress[2],
					arpPacket.SourceHwAddress[3], arpPacket.SourceHwAddress[4], arpPacket.SourceHwAddress[5])
			}
		case <-timeout:
			fmt.Println("\nСканирование завершено.")
			return
		}
	}
}
