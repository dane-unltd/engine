// +build ignore

package core

import "net"
import "encoding/gob"
import "time"

import "fmt"

//Receiver on client side. Receives inputs from the server
type CltRx struct {
	nextCmd, currCmd map[uint32]UsrCmd
	cmd              UsrCmd
	dec              *gob.Decoder
	localId          uint32
	offset, latency  time.Duration
	frameNS          time.Duration
	startTime        time.Time
	cmdMap           map[string]uint
}

func NewCltRx(conn net.Conn) *CltRx {
	rx := CltRx{
		make(map[uint32]UsrCmd),
		make(map[uint32]UsrCmd),
		UsrCmd{},
		gob.NewDecoder(conn),
		0,
		0, 0,
		0,
		time.Now(),
		make(map[string]uint),
	}

	return &rx
}

func (rx *CltRx) RxLocalId() {
	rx.dec.Decode(&rx.localId)
	rx.dec.Decode(&rx.frameNS)
	fmt.Println("frameNS:", rx.frameNS)
}

//blocking receive function
func (rx *CltRx) Receive() error {
	var srvTime, n uint32
	var lastRx int

	err := rx.dec.Decode(&srvTime)
	if err != nil {
		return err
	}
	err = rx.dec.Decode(&lastRx)
	if err != nil {
		return err
	}

	cltTimeNS := time.Now().Sub(rx.startTime) + rx.offset
	lastRxNS := time.Duration(lastRx) * rx.frameNS
	srvTimeNS := time.Duration(srvTime) * rx.frameNS

	rx.latency = 99*rx.latency/100 + (cltTimeNS-lastRxNS)/100
	rx.offset += (srvTimeNS - lastRxNS) / 100

	err = rx.dec.Decode(&n)
	if err != nil {
		return err
	}

	for i := uint32(0); i < n; i++ {
		rx.cmd = UsrCmd{}
		err = rx.dec.Decode(&rx.cmd)
		if err != nil {
			return err
		}

		rx.nextCmd[rx.cmd.Id] = rx.cmd
	}
	return nil
}

func (rx *CltRx) Id() string {
	return "CltRx"
}

func (rx *CltRx) Init(sim Sim, res *ResMgr) {
	cmds := make([]string, 32)
	res.ReadConfig("cmds", &(cmds))
	for i, cmd := range cmds {
		rx.cmdMap[cmd] = uint(i)
	}
}

func (rx *CltRx) Swap() {
	temp := rx.currCmd
	rx.currCmd = rx.nextCmd
	rx.nextCmd = temp
}

func (rx *CltRx) Update() {
}

func (rx *CltRx) Cmd(id uint32) UsrCmd {
	return rx.currCmd[id]
}

func (rx *CltRx) Time() uint32 {
	return uint32((time.Now().Sub(rx.startTime) + rx.offset) / rx.frameNS)
}

func (rx *CltRx) LocalId() uint32 {
	return rx.localId
}

func (rx *CltRx) FrameNS() time.Duration {
	return rx.frameNS
}

func (rx *CltRx) Active(id uint32, cmd string) bool {
	bit, ok := rx.cmdMap[cmd]
	if ok {
		return rx.currCmd[id].Actions&(1<<bit) > 0
	}
	return false
}

func (rx *CltRx) Point(id uint32) (x, y int) {
	return rx.currCmd[id].X, rx.currCmd[id].Y
}
