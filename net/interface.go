package net

import (
	"fmt"
	"github.com/shirou/gopsutil/net"
	"os"
	"text/tabwriter"
	"time"
)

func ScanNetworkTraffic() {
	prevStats, err := net.IOCounters(true)
	if err != nil {
		fmt.Printf("Ошибка при получении статистики сети: %v\n", err)
		return
	}

	time.Sleep(1 * time.Second)

	currentStats, err := net.IOCounters(true)
	if err != nil {
		fmt.Printf("Ошибка при получении статистики сети: %v\n", err)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Interface\tBytes Sent\tBytes Received\tPackets Sent\tPackets Received")

	for i, current := range currentStats {
		prev := prevStats[i]
		bytesSent := current.BytesSent - prev.BytesSent
		bytesRecv := current.BytesRecv - prev.BytesRecv
		packetsSent := current.PacketsSent - prev.PacketsSent
		packetsRecv := current.PacketsRecv - prev.PacketsRecv

		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\n", current.Name, bytesSent, bytesRecv, packetsSent, packetsRecv)
	}

	w.Flush()
}
