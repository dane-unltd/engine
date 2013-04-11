package helpers

import (
	"bytes"
	"encoding/binary"
	"github.com/dane-unltd/engine/core"
	"github.com/dane-unltd/linalg/matrix"
	"io"
)

type VecMap map[core.EntId]matrix.VecD

func NewVecMap() VecMap {
	v := make(VecMap)
	return v
}

func (v VecMap) Mutate(id core.EntId, value interface{}) {
	if value == nil {
		delete(v, id)
		return
	}
	v[id] = value.(matrix.VecD)
}

func (v VecMap) Clone() core.State {
	vNew := NewVecMap()
	for id := range v {
		vNew[id] = v[id].Copy().(matrix.VecD)
	}
	return vNew
}

func (v VecMap) SerDiff(buf io.Writer, newEnts []core.EntId, newSt core.State) {
	vNew := newSt.(VecMap)
	n := len(newEnts)
	nBytes := n / 8
	if n%8 != 0 {
		nBytes += 1
	}
	bitMask := make([]byte, nBytes)
	bufTemp := &bytes.Buffer{}
	for i, id := range newEnts {
		byIx := i / 8
		bitIx := uint(i % 8)
		_, ok := v[id]
		if !(ok && v[id].Equals(vNew[id])) {
			bitMask[byIx] |= 1 << bitIx
			binary.Write(bufTemp, binary.LittleEndian, vNew[id])
		}
	}
	buf.Write(bitMask)
	buf.Write(bufTemp.Bytes())
}

func (v VecMap) DeserDiff(buf io.Reader, newEnts []core.EntId) {
	//TODO
}
