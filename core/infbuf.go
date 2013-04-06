package core

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

const queueSizeRb uint64 = 1 << 4
const indexMaskRb uint64 = queueSizeRb - 1

type InfBuf struct {
	padding1 [8]uint64
	wAck     uint64
	padding2 [8]uint64
	wp       uint64
	padding3 [8]uint64
	rp       uint64
	padding4 [8]uint64
	data     [queueSizeRb]interface{}
	padding5 [8]uint64
}

func NewInfBuf() *InfBuf {
	return &InfBuf{wAck: 0, wp: 1, rp: 1}
}

func (b *InfBuf) Write(value interface{}) Tick {
	var ix = atomic.AddUint64(&b.wp, 1) - 1
	if ix > (b.rp + queueSizeRb - 2) {
		b.rp++
	}
	b.data[ix&indexMaskRb] = value
	for !atomic.CompareAndSwapUint64(&b.wAck, ix-1, ix) {
		runtime.Gosched()
	}
	return Tick(ix)
}

func (b *InfBuf) Read(t Tick) interface{} {
	ix := uint64(t)
	if ix > b.wAck || ix < b.rp {
		panic("index out of bounds")
	}
	return b.data[ix&indexMaskRb]
}

func (b *InfBuf) Reset(t Tick) {
	ix := uint64(t)
	if ix <= b.rp {
		panic("infbuf, reset")
	}
	b.wp = ix
	b.wAck = ix - 1
	if ix > (b.rp + queueSizeRb - 2) {
		b.rp = ix + 2 - queueSizeRb
	}
}

func (b *InfBuf) Dump() {
	fmt.Printf("wAck: %3d, wp: %3d, rp: %3d, content:", b.wAck, b.wp, b.rp)
	for ix, value := range b.data {
		fmt.Printf("%5v : %5v", ix, value)
	}
	fmt.Print("\n")
}
