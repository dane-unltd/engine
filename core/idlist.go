package core

import (
	"bytes"
	"encoding/binary"
	"io"
)

type IdList map[EntId]Empty

func NewIdList() IdList {
	return make(IdList)
}

func (il IdList) Val(id EntId) StateVar {
	return il[id]
}

func (il IdList) Mutate(id EntId, value interface{}) {
	if value == nil {
		delete(il, id)
		return
	}

	il[id] = Empty{}
}

func (il IdList) Clone() State {
	iNew := make(IdList)
	for i := range il {
		iNew[i] = Empty{}
	}
	return iNew
}

func (il IdList) SerDiff(buf io.Writer, newEnts []EntId, newSt State) {
	newIl := newSt.(IdList)

	bufTemp := &bytes.Buffer{}

	n := 0
	for _, id := range newEnts {
		if _, ok := newIl[id]; ok {
			binary.Write(bufTemp, binary.LittleEndian, id)
			n++
		}
	}
	binary.Write(buf, binary.LittleEndian, uint32(n))
	buf.Write(bufTemp.Bytes())
}

func (il IdList) DeserDiff(buf io.Reader, newEnts []EntId) {
	for id := range il {
		delete(il, id)
	}

	var n uint32
	var id EntId
	binary.Read(buf, binary.LittleEndian, &n)
	for i := 0; i < int(n); i++ {
		binary.Read(buf, binary.LittleEndian, &id)
		il[id] = Empty{}
	}
}
