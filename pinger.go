package pinger

import (
   "math/rand"
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

   udpaddr := net.UDPAddr{IP: net.ParseIP(target)}
   return &Pinger{target: target, udpaddr: udpaddr, timeout: timeout, msg: msg, id: rand.Int()}, nil
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
