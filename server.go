package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func writeARP(handle *pcap.Handle, iface *net.Interface, sourceAddr, targetAddr netip.Addr, vlanID uint16) error {
	vlan := layers.Dot1Q{
		VLANIdentifier: vlanID,
	}

	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: sourceAddr.AsSlice(),
		DstHwAddress:      net.HardwareAddr{0, 0, 0, 0, 0, 0},
		DstProtAddress:    targetAddr.AsSlice(),
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	if err := gopacket.SerializeLayers(buf, opts, &eth, &vlan, &arp); err != nil {
		return err
	}
	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func main() {
	// Ağ arayüzü adı
	ifaceName := "en0"

	// Ağ arayüzünü al
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding interface: %v\n", err)
		return
	}

	// Kaynak ve hedef IP adresleri
	sourceAddr := netip.MustParseAddr("192.168.51.93")
	targetAddr := netip.MustParseAddr("192.168.51.255")

	// pcap handle oluşturma
	handle, err := pcap.OpenLive(ifaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening pcap handle: %v\n", err)
		return
	}
	defer handle.Close()

	// ARP paketini gönder
	vlanID := uint16(0) // VLAN ID kullanılmıyorsa 0 olabilir
	if err := writeARP(handle, iface, sourceAddr, targetAddr, vlanID); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing ARP packet: %v\n", err)
		return
	}

	fmt.Println("ARP broadcast message sent successfully.")
}
