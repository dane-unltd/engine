package core

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/binary"
	"log"
	"net"
	"net/http"
	"time"
)

type wsClientConn struct {
	lastRx  Tick
	lastAck Tick
	cmdBuf  *CmdBuf
	id      EntId
	ws      *websocket.Conn
	inBuf   [1500]byte
}

type WsServer struct {
	sim     *Sim
	clients []*wsClientConn
	newConn chan *wsClientConn
	path    string
	fps     float64
}

func NewWsServer(path string, fps float64) *WsServer {
	server := new(WsServer)
	server.path = path
	server.clients = make([]*wsClientConn, 0, 2)
	server.fps = fps

	server.sim = NewSim(NewIdGen())
	server.sim.AddState("input", NewIdMap())
	server.sim.AddState("logins", NewIdMap())
	server.sim.AddState("disconnects", NewIdList())

	server.newConn = make(chan *wsClientConn)
	return server
}

func (srv *WsServer) AddTransFunc(pr int, f TransFunc) {
	srv.sim.AddTransFunc(pr, f)
}

func (srv *WsServer) AddState(name string, state State) {
	srv.sim.AddState(name, state)
}

func (srv *WsServer) SetSerInfo(si SerInfo) {
	srv.sim.SetSerInfo(si)
}

func (srv *WsServer) Run() {

	wsHandler := func(ws *websocket.Conn) {
		cmd := UserCmd{}

		cl := &wsClientConn{}
		cl.cmdBuf = NewCmdBuf()
		cl.ws = ws
		cl.cmdBuf.Write(UserCmd{})

		log.Println("incoming connection")
		log.Println("clients:", srv.clients)
		srv.newConn <- cl
		for {
			pkt := cl.inBuf[0:]
			n, err := ws.Read(pkt)
			pkt = pkt[0:n]
			if err != nil {
				log.Println(err)
				break
			}
			buf := bytes.NewBuffer(pkt)
			err = binary.Read(buf, binary.LittleEndian, &cl.lastAck)
			if err != nil {
				log.Println(err)
				break
			}

			err = binary.Read(buf, binary.LittleEndian, &cmd)
			if err != nil {
				log.Println(err)
				break
			}

			if cmd.Time > cl.lastRx {
				cl.lastRx = cmd.Time
			} else {
				cmd.Time = cl.lastRx
			}
			cl.cmdBuf.Write(cmd)

		}
	}

	addr, err := net.ResolveUnixAddr("unix", srv.path)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	http.Handle("/ws/", websocket.Handler(wsHandler))

	go func() {
		log.Fatal(http.Serve(ln, nil))
	}()

	frameNS := time.Duration((1e9 / srv.fps))
	clk := time.NewTicker(frameNS)
	//main loop
	for {
		select {
		case <-clk.C:
			srv.updateInputs()
			srv.clearLogins()
			srv.sim.Update()
			srv.sendUpdates()
		case cl := <-srv.newConn:
			cl.id = srv.sim.NewId()
			srv.sim.Mutate("logins", cl.id, Empty{})
			buf := &bytes.Buffer{}
			binary.Write(buf, binary.LittleEndian, srv.sim.Time())
			binary.Write(buf, binary.LittleEndian, cl.lastAck)
			binary.Write(buf, binary.LittleEndian, cl.lastRx)
			srv.sim.Diff(buf, cl.id, 0, srv.sim.Time())
			err := websocket.Message.Send(cl.ws, buf.Bytes())
			if err != nil {
				log.Println(err)
			} else {
				srv.clients = append(srv.clients, cl)
			}
		}
	}
}

func (srv *WsServer) clearLogins() {
	state := srv.sim.states.Read(srv.sim.time).(StateMap)
	logins := state["logins"].(IdMap)
	for id, _ := range logins {
		srv.sim.Mutate("logins", id, nil)
	}
	discs := state["disconnects"].(IdList)
	for id := range discs {
		srv.sim.Mutate("disconnects", id, nil)
	}
}

func (srv *WsServer) updateInputs() {
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
			log.Println("drop")
		}
		srv.sim.Mutate("input", cl.id, b.Peak(0))
	}
}

func (srv *WsServer) sendUpdates() {
	for i := 0; i < len(srv.clients); i++ {
		buf := &bytes.Buffer{}
		binary.Write(buf, binary.LittleEndian, srv.sim.Time())
		binary.Write(buf, binary.LittleEndian, srv.clients[i].lastAck)
		binary.Write(buf, binary.LittleEndian, srv.clients[i].lastRx)
		srv.sim.Diff(buf, srv.clients[i].id, srv.sim.Time()-1, srv.sim.Time())
		err := websocket.Message.Send(srv.clients[i].ws, buf.Bytes())
		if err != nil {
			log.Println(err)
			srv.sim.Mutate("disconnects", srv.clients[i].id, Empty{})
			srv.clients[i] = srv.clients[len(srv.clients)-1]
			srv.clients = srv.clients[:len(srv.clients)-1]
			i--
		}
	}
}
