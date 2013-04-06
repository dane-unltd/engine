package core

import "fmt"
import "log"
import "sync/atomic"
import "runtime"

type EntId uint32

type IdGen struct {
	in      chan EntId
	out     chan EntId
	freeIds []EntId
	idCheck []bool
	maxId   EntId
}

func NewIdGen() *IdGen {
	ig := &IdGen{
		make(chan EntId),
		make(chan EntId),
		make([]EntId, 0, 100),
		make([]bool, 0, 100),
		1,
	}
	ig.idCheck = append(ig.idCheck, false)
	ig.idCheck = append(ig.idCheck, true)
	go ig.run()
	return ig
}

func (ig *IdGen) run() {
	currId := EntId(1)
	for {
		select {
		case ig.out <- currId:
			log.Println("free ids: ", ig.freeIds)
			if len(ig.freeIds) > 0 {
				currId = ig.freeIds[len(ig.freeIds)-1]
				ig.freeIds = ig.freeIds[:len(ig.freeIds)-1]
				ig.idCheck[currId] = true
			} else {
				ig.maxId++
				currId = ig.maxId
				ig.idCheck = append(ig.idCheck, true)
			}
		case freeId := <-ig.in:
			if freeId >= ig.maxId {
				continue
			}

			if ig.idCheck[freeId] == false {
				continue
			}

			ig.idCheck[freeId] = false
			ig.freeIds = append(ig.freeIds, freeId)
		}
	}
}

func (ig *IdGen) Next() EntId {
	return <-ig.out
}

func (ig *IdGen) Free(id EntId) {
	ig.in <- id
}

const queueSizeId uint64 = 1 << 10
const indexMaskId uint64 = queueSizeId - 1

type IdBuf struct {
	padding1 [8]uint64
	wAck     uint64
	padding2 [8]uint64
	wp       uint64
	padding3 [8]uint64
	rp       uint64
	padding4 [8]uint64
	data     [queueSizeId]EntId
	padding5 [8]uint64
}

func NewIdBuf() *IdBuf {
	return &IdBuf{wAck: 0, wp: 1, rp: 1}
}

func (b *IdBuf) Write(value EntId) {
	var ix = atomic.AddUint64(&b.wp, 1) - 1
	if ix > (b.rp + queueSizeId - 2) {
		panic("overflow")
	}
	b.data[ix&indexMaskId] = value
	for !atomic.CompareAndSwapUint64(&b.wAck, ix-1, ix) {
		runtime.Gosched()
	}
}

func (b *IdBuf) Read() EntId {
	var ix = atomic.AddUint64(&b.rp, 1) - 1
	if ix > b.wAck {
		panic("underflow")
	}
	return b.data[ix&indexMaskId]
}

func (b *IdBuf) Dump() {
	fmt.Printf("wAck: %3d, wp: %3d, rp: %3d, content:", b.wAck, b.wp, b.rp)
	for ix, value := range b.data {
		fmt.Printf("%5v : %5v", ix, value)
	}
	fmt.Print("\n")
}

func (b *IdBuf) HasNext() bool {
	return b.rp <= b.wAck
}

func (b *IdBuf) Full() bool {
	return b.wp > (b.rp + queueSizeId - 2)
}
