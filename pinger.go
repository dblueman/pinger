package pinger

import (
   "crypto/rand"
   "encoding/binary"
   "fmt"
   "net"
   "os"
   "time"

   "golang.org/x/net/icmp"
   "golang.org/x/net/ipv4"
)

type Pinger struct {
   buf     [1500]byte
   msg     icmp.Message
   timeout time.Duration
   target  string
   udpaddr net.UDPAddr
   id      int
   seq     int
}

const (
   protocolICMP = 1
)

var (
   conn *icmp.PacketConn
)

func randShort() (int, error) {
   var b [2]byte
   _, err := rand.Read(b[:])
   if err != nil {
      return -1, err
   }

   return int(binary.BigEndian.Uint16(b[:])), nil
}

func NewPinger(target string, timeout time.Duration) (*Pinger, error) {
   var err error

   if conn == nil {
   	conn, err = icmp.ListenPacket("udp4", "0.0.0.0")
   	if err != nil {
   		return nil, err
   	}
   }

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
      Code: 0,
	}

   ips, err := net.LookupIP(target)
   if err != nil {
      return nil, err
   }

   if len(ips) == 0 {
      return nil, fmt.Errorf("No IP addreses for %s!", target)
   }

   ip := ips[0]

   udpaddr := net.UDPAddr{IP: ip}
   id, err := randShort()
   if err != nil {
      return nil, err
   }

   return &Pinger{target: ip.String(), udpaddr: udpaddr, timeout: timeout, msg: msg, id: id}, nil
}

func (p *Pinger) Ping() (bool, error) {
   p.seq++

   p.msg.Body = &icmp.Echo{
      ID: p.id,
      Seq: p.seq,
      // no Body
   }

	wb, err := p.msg.Marshal(nil)
	if err != nil {
		return false, err
	}

	_, err = conn.WriteTo(wb, &p.udpaddr)
   if err != nil {
		return false, err
	}

   err = conn.SetReadDeadline(time.Now().Add(p.timeout))
   if err != nil {
      return false, err
   }

   for {
   	n, peer, err := conn.ReadFrom(p.buf[:])
   	if err != nil {
         if os.IsTimeout(err) {
            return false, nil
         }

   		return false, err
   	}

      if peer.String() != p.target+":0" {
         continue
      }

   	msg, err := icmp.ParseMessage(protocolICMP, p.buf[:n])
   	if err != nil {
   		return false, err
   	}

      if msg.Type == ipv4.ICMPTypeEchoReply {
         return true, nil
      }
   }
}
