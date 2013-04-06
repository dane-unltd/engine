package core

import (
	"bytes"
	"encoding/binary"
	"io"
)

type IdMap map[EntId]StateVar

func NewIdMap() IdMap {
	v := make(IdMap)
	return v
}

func (v IdMap) Mutate(id EntId, value interface{}) {
	if value == nil {
		delete(v, id)
		return
	}
	v[id] = value.(StateVar).Copy()
}

func (v IdMap) Val(id EntId) StateVar {
	val, ok := v[id]
	if ok {
		return val
	}
	return nil
}

func (v IdMap) Clone() State {
	vNew := NewIdMap()
	for id := range v {
		vNew[id] = v[id].Copy()
	}
	return vNew
}

func (v IdMap) SerDiff(buf io.Writer, newEnts []EntId, newSt State) {
	vNew := newSt.(IdMap)
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

func (v IdMap) DeserDiff(buf io.Reader, newEnts []EntId) {
	//TODO
}
