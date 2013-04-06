package core

import "fmt"
import "sync/atomic"
import "runtime"

const queueSize uint64 = 1 << 4
const indexMask uint64 = queueSize - 1

type CmdBuf struct {
	padding1 [8]uint64
	wAck     uint64
	padding2 [8]uint64
	wp       uint64
	padding3 [8]uint64
	rp       uint64
	padding4 [8]uint64
	data     [queueSize]UserCmd
	padding5 [8]uint64
}

func NewCmdBuf() *CmdBuf {
	return &CmdBuf{wAck: 0, wp: 1, rp: 1}
}

func (b *CmdBuf) Write(value UserCmd) {
	var ix = atomic.AddUint64(&b.wp, 1) - 1
	for ix > (b.rp + queueSize - 2) {
		runtime.Gosched()
	}
	b.data[ix&indexMask] = value
	for !atomic.CompareAndSwapUint64(&b.wAck, ix-1, ix) {
		runtime.Gosched()
	}
}

func (b *CmdBuf) Read() UserCmd {
	var ix = atomic.AddUint64(&b.rp, 1) - 1
	for ix > b.wAck {
		runtime.Gosched()
	}
	return b.data[ix&indexMask]
}

func (b *CmdBuf) Dump() {
	fmt.Printf("wAck: %3d, wp: %3d, rp: %3d, content:", b.wAck, b.wp, b.rp)
	for ix, value := range b.data {
		fmt.Printf("%5v : %5v", ix, value)
	}
	fmt.Print("\n")
}

func (b *CmdBuf) Peak(offset uint64) UserCmd {
	var ix = b.rp + offset
	if ix > b.wAck {
		panic("CmdBuf.Peak: out of bounds")
	}
	return b.data[ix&indexMask]
}

func (b *CmdBuf) HasNext() bool {
	return b.rp < b.wAck
}

func (b *CmdBuf) Step() {
	b.rp++
}
