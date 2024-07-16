package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// Broadcast yapılacak bağlantı noktası
	broadcastPort := 5050

	// UDP bağlantısını oluştur
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: broadcastPort})
	if err != nil {
		fmt.Printf("Bağlantı hatası: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Broadcast yapılıyor...\n")

	// Broadcast mesajı
	broadcastMsg := []byte("Broadcast mesajı")

	// 192.168.51.255 adresine broadcast yap
	destAddr := &net.UDPAddr{IP: net.IPv4(192, 168, 51, 255), Port: broadcastPort}
	_, err = conn.WriteTo(broadcastMsg, destAddr)
	if err != nil {
		fmt.Printf("Broadcast yapma hatası: %v\n", err)
		return
	}

	// Broadcast mesajına cevap veren IP adreslerini topla
	ipAddresses := make(map[string]bool)
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	for {
		_, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break // Zaman aşımı hatası
			}
			fmt.Printf("Okuma hatası: %v\n", err)
			return
		}

		// Kendi IP adresimizi alalım
		localIP := getLocalIP()

		// Cevap veren cihazın IP adresini alalım
		ip := addr.IP.String()

		// Eğer cevap veren cihazın IP adresi kendi IP adresimiz değilse ekrana yazdıralım
		if ip != localIP {
			if !ipAddresses[ip] {
				ipAddresses[ip] = true
				fmt.Printf("Cevap alındı: %s\n", ip)
			}
		}

		// Daha fazla cevap beklemek için sürekli döngüye devam et
	}

	fmt.Println("Broadcast sonlandırıldı.")
}

func getLocalIP() string {
	// Tüm ağ arayüzlerini alalım
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Ağ arayüzleri alınamadı:", err)
		return ""
	}

	// İlk geçerli IPv4 adresini bulalım
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("IP adresleri alınamadı:", err)
			continue
		}

		for _, addr := range addrs {
			// IPv4 adresi mi kontrol edelim
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	fmt.Println("Geçerli bir IP adresi bulunamadı")
	return ""
}
