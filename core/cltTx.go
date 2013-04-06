// +build ignore

package core

import "net"
import "encoding/gob"

//import "fmt"

//Receiver on client side. Receives inputs from the server
type CltTx struct {
	enc   *gob.Encoder
	input CmdSrc
}

func NewCltTx(conn net.Conn) *CltTx {
	tx := CltTx{
		gob.NewEncoder(conn),
		nil,
	}
	return &tx
}

func (tx *CltTx) Id() string {
	return "CltTx"
}

func (tx *CltTx) Init(sim Sim, res *ResMgr) {
	tx.input = sim.Comp("Input").(CmdSrc)
}

func (tx *CltTx) Swap() {
}

func (tx *CltTx) Update() {
	cmd := tx.input.Cmd(0)
	//fmt.Println("CltTx:", cmd)
	err := tx.enc.Encode(cmd)
	if err != nil {
		panic(err)
	}
}
