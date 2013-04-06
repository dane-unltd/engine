package core

import "sync/atomic"
import "runtime"

type Mutation struct {
	Id    EntId
	Value StateVar
}

const bufSize uint64 = 1 << 4
const indexMaskMut uint64 = bufSize - 1

type MutBuf struct {
	padding1 [8]uint64
	wAck     uint64
	padding2 [8]uint64
	wp       uint64
	padding3 [8]uint64
	rp       uint64
	padding4 [8]uint64
	data     [bufSize]Mutation
	padding5 [8]uint64
}

func NewMutBuf() *MutBuf {
	return &MutBuf{wAck: 0, wp: 1, rp: 1}
}

func (b *MutBuf) Write(value Mutation) {
	var ix = atomic.AddUint64(&b.wp, 1) - 1
	for ix > (b.rp + bufSize - 2) {
		runtime.Gosched()
	}
	b.data[ix&indexMaskMut] = value
	for !atomic.CompareAndSwapUint64(&b.wAck, ix-1, ix) {
		runtime.Gosched()
	}
}

func (b *MutBuf) Read() (Mutation, bool) {

	if b.rp > b.wAck {
		return Mutation{}, false
	}
	b.rp++
	return b.data[(b.rp-1)&indexMaskMut], true
}
