package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type clientConn struct {
	lastRx  Tick
	lastAck Tick
	cmdBuf  *CmdBuf
	id      EntId
}

type Server struct {
	sim     *Sim
	clients map[string]*clientConn
	cmds    map[EntId]UserCmd
	conn    net.PacketConn
	inBuf   [1500]byte
}

func NewServer(port string) *Server {
	server := new(Server)
	server.clients = make(map[string]*clientConn)
	server.cmds = make(map[EntId]UserCmd)

	server.sim = NewSim(NewIdGen())

	var err error
	server.conn, err = net.ListenPacket("udp", port)
	if err != nil {
		panic(err)
	}

	return server
}

func (srv *Server) AddTransFunc(pr int, f TransFunc) {
	srv.sim.AddTransFunc(pr, f)
}

func (srv *Server) AddState(name string, state State) {
	srv.sim.AddState(name, state)
}

func (srv *Server) SetSerInfo(si SerInfo) {
	srv.sim.SetSerInfo(si)
}

func (srv *Server) Run() {

	go srv.receive()

	frameNS := time.Duration((1e9 / 10))
	clk := time.NewTicker(frameNS)
	//main loop
	for {
		select {
		case <-clk.C:
			srv.updateInputs()
			srv.sim.Update()
			srv.sendUpdates()
		}
	}
}

func (srv *Server) receive() {
	cmd := UserCmd{}

	for {
		pkt := srv.inBuf[0:]
		n, addr, err := srv.conn.ReadFrom(pkt)
		fmt.Println("server incoming", n, addr, err)
		pkt = pkt[0:n]
		if err != nil {
			fmt.Println(pkt)
			panic(err)
		}
		buf := bytes.NewBuffer(pkt)
		cl, ok := srv.clients[addr.String()]
		if !ok {
			srv.clients[addr.String()] = &clientConn{
				0,
				0,
				NewCmdBuf(),
				srv.sim.NewId(),
			}
			cl = srv.clients[addr.String()]
			cl.cmdBuf.Write(UserCmd{})
			fmt.Println(addr.String())
			fmt.Println("clients:", srv.clients)
		}

		var ack Tick
		binary.Read(buf, binary.LittleEndian, &ack)
		if ack > cl.lastAck {
			cl.lastAck = ack
		}

		for {
			err := binary.Read(buf, binary.LittleEndian, &cmd)
			if err != nil {
				break
			}
			if cmd.Time > cl.lastRx {
				cl.cmdBuf.Write(cmd)
				cl.lastRx = cmd.Time
			} else {
				fmt.Println("drop")
			}
		}
	}
}

func (srv *Server) updateInputs() {
	for _, cl := range srv.clients {
		b := cl.cmdBuf
		steps := 0
		for b.HasNext() {
			if b.Peak(1).Time <= srv.sim.Time()+1 {
				b.Step()
				steps++
			} else {
				break
			}
		}
		if steps > 1 {
			fmt.Println("drop")
		}
		srv.cmds[cl.id] = b.Peak(0)
	}
}

func (srv *Server) sendUpdates() {
	for addr, cl := range srv.clients {
		buf := &bytes.Buffer{}
		binary.Write(buf, binary.LittleEndian, srv.sim.Time())
		binary.Write(buf, binary.LittleEndian, cl.lastAck)
		binary.Write(buf, binary.LittleEndian, cl.lastRx)
		srv.sim.Diff(buf, cl.id, cl.lastAck, srv.sim.Time())

		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			panic(err)
		}
		srv.conn.WriteTo(buf.Bytes(), udpAddr)
	}
}
