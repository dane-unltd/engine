package core

import (
	"encoding/binary"
	"io"
)

type SerInfo struct {
	EntSel func(EntId, StateMap) []EntId
	States []string
}

func Serialize(info SerInfo, buf io.Writer, playerId EntId,
	oldSt, newSt StateMap) {

	newEnts := info.EntSel(playerId, newSt)
	n := uint32(len(newEnts))

	binary.Write(buf, binary.LittleEndian, n)
	binary.Write(buf, binary.LittleEndian, newEnts)

	for _, name := range info.States {
		oldSt[name].SerDiff(buf, newEnts, newSt[name])
	}
}

func Deserialize(info SerInfo, buf io.Reader, refSt StateMap, mut MutFuncs) {
	var n uint32

	check := make(map[EntId]bool)
	oldEnts := info.EntSel(0, refSt)

	for _, id := range oldEnts {
		check[id] = false
	}

	err := binary.Read(buf, binary.LittleEndian, &n)
	if err != nil {
		panic(err)
	}

	newEnts := make([]EntId, n)
	err = binary.Read(buf, binary.LittleEndian, newEnts)
	if err != nil {
		panic(err)
	}

	for _, id := range newEnts {
		if _, ok := check[id]; ok {
			check[id] = true
		}
	}

	for id, flag := range check {
		if !flag {
			mut.Destroy(id)
		}
	}

	for _, name := range info.States {
		refSt[name].DeserDiff(buf, newEnts)
	}
}
