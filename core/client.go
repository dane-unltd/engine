package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

type Client struct {
	inBuf       [1500]byte
	ipcBuf      [1500]byte
	sim         *Sim
	res         *ResMgr
	input       *Input
	lastRx      Tick
	offset, rtt time.Duration
	frameNS     time.Duration
	startTime   time.Time
	conn        net.Conn
	ipc         net.Conn
}

func NewClient(name, ipcPort string, host string) *Client {
	client := new(Client)

	//client.res = NewResMgr(name)

	var err error

	client.ipc, err = net.Dial("tcp", "localhost:"+ipcPort)
	if err != nil {
		panic(err)
	}

	client.conn, err = net.Dial("udp", host)
	if err != nil {
		panic(err)
	}

	client.input = NewInput()
	client.startTime = time.Now()
	client.frameNS = 1e9 / 10

	client.sim = NewSim(nil)

	return client
}

func (cl *Client) AddTransFunc(pr int, f TransFunc) {
	cl.sim.AddTransFunc(pr, f)
}

func (cl *Client) AddState(name string, state State) {
	cl.sim.AddState(name, state)
}

func (cl *Client) SetSerInfo(si SerInfo) {
	cl.sim.SetSerInfo(si)
}

func (cl *Client) Run() {
	go cl.viewerUpdate()

	for {
		refTime, buf := cl.receive()
		cl.sim.DeserDiff(buf, refTime, cl.lastRx)
	}
}

func (cl *Client) viewerUpdate() {
	var lastUpdate Tick
	var actions uint32
	for {
		pkt := cl.ipcBuf[0:]
		n, err := cl.ipc.Read(pkt)
		if err != nil {
			if err.Error() != "EOF" || n < 4 {
				panic(err)
			}
		}

		pkt = pkt[0:n]
		buf := bytes.NewBuffer(pkt)

		var temp uint32
		binary.Read(buf, binary.LittleEndian, &temp)

		actions = actions | temp

		//TODO: read input
		//TODO: send render data

		time.Sleep(1e9 / 100)

		if cl.Time() > lastUpdate {
			fmt.Println("sending actions:", actions)
			//TODO: update input
			cl.send()
			lastUpdate = cl.Time()
		}

		actions = 0
	}
}

func (cl *Client) receive() (Tick, io.Reader) {
	buf := &bytes.Buffer{}
	var srvTime Tick
	for {
		pkt := cl.inBuf[0:]
		fmt.Println("client incoming")
		n, err := cl.conn.Read(pkt)
		if err != nil {
			panic(err)
		}

		pkt = pkt[0:n]
		buf = bytes.NewBuffer(pkt)

		binary.Read(buf, binary.LittleEndian, &srvTime)
		if srvTime <= cl.lastRx {
			continue
		}
		break
	}

	cl.lastRx = srvTime

	var lastAck, lastRxSrv Tick
	binary.Read(buf, binary.LittleEndian, &lastAck)
	binary.Read(buf, binary.LittleEndian, &lastRxSrv)

	cltTimeNS := time.Now().Sub(cl.startTime) + cl.offset
	lastRxNS := time.Duration(lastRxSrv) * cl.frameNS
	srvTimeNS := time.Duration(srvTime) * cl.frameNS

	cl.offset += (srvTimeNS - lastRxNS) / 100
	cl.rtt = 99*cl.rtt/100 + (cltTimeNS-lastRxNS)/100

	fmt.Println("RTT:", cl.rtt)
	fmt.Println("offset:", cl.offset)
	fmt.Println(srvTimeNS, lastRxNS)

	return lastAck, buf
}

func (cl *Client) Time() Tick {
	return Tick((time.Now().Sub(cl.startTime) + cl.offset) / cl.frameNS)
}

func (cl *Client) send() {
	//TODO
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, cl.lastRx)
	binary.Write(buf, binary.LittleEndian, UserCmd{Time: cl.Time()})
	_, err := cl.conn.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}
}
