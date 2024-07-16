package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.ListenPacket("udp4", ":5050")
	if err != nil {
		fmt.Println("Dinleme hatası: ", err)
	}

	defer conn.Close()
	addr, err := net.ResolveUDPAddr("udp4", "192.168.51.255:5050")
	if err != nil {
		fmt.Println("UDP adres hatası: ", err)
	}

	_, err = conn.WriteTo([]byte("broadcast"), addr)
	if err != nil {
		fmt.Println("Yazma hatası: ", err)
	}

	for {
		buffer := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, remoteAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Okuma zaman aşımında hata verir, sadece zaman aşımı hatalarını göz ardı edebiliriz
				continue
			}
			fmt.Println("Okuma hatası: ", err)
			break
		}

		// remoteAddr IP adresini yazdırmak için *net.UDPAddr türüne dönüştürülür
		udpAddr, ok := remoteAddr.(*net.UDPAddr)
		if !ok {
			fmt.Println("Adres türü hatası")
			continue
		}

		fmt.Printf("Cevap alındı from %s: %s\n", udpAddr.IP.String(), string(buffer[:n]))
	}
}
