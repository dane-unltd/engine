package entstate

import (
	"bytes"
	"encoding/binary"
	"github.com/dane-unltd/engine/bitmask"
	"io"
)

func Serialize(buf io.Writer, serAll bool, active []bool) {
	nEnts := 0
	for _, act := range active {
		if act {
			nEnts++
		}
	}
	binary.Write(buf, binary.LittleEndian, uint32(nEnts))

	bufTemp := &bytes.Buffer{}
	for id, act := range active {
		if !act {
			continue
		}
		binary.Write(buf, binary.LittleEndian, EntId(id))
		bufTemp.Reset()

		bitMask := bitmask.New(len(networkedComps))
		for i, compId := range networkedComps {
			v := stateComps[compId].Val(EntId(id))
			if serAll || !oldStateComps[compId].Equal(v, EntId(id)) {
				bitMask.Set(i)
				binary.Write(bufTemp, binary.LittleEndian, v)
			}

		}
		buf.Write(bitMask)
		buf.Write(bufTemp.Bytes())
	}

}

func Deserialize(buf io.Reader) {
	var nEnts uint32
	binary.Read(buf, binary.LittleEndian, &nEnts)

	for i := uint32(0); i < nEnts; i++ {
		var id uint32
		err := binary.Read(buf, binary.LittleEndian, &id)
		if err != nil {
			panic(err)
		}

		bitMask := bitmask.New(len(networkedComps))
		_, err = buf.Read(bitMask)
		if err != nil {
			panic(err)
		}
	}

	//TODO
}
