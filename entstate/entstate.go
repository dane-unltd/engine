package entstate

import (
	"github.com/dane-unltd/engine/idgen"
)

type CompId uint32
type EntId uint32

type StateComp interface {
	Copyer
	Clone() interface{}
	Zero(id EntId)
	Append(n uint32)
	Val(id EntId) interface{}
	Equal(v interface{}, id EntId) bool
}

type Copyer interface {
	Copy(src interface{})
}

//id generator
var idGen = idgen.New(incMaxEnts)

//state collection
var stateComps = make([]StateComp, 0)
var networkedComps = make([]CompId, 0)
var oldStateComps = make([]StateComp, 0)

func New() EntId {
	return EntId(idGen.Next())
}

func Delete(id EntId) {
	for i := range stateComps {
		stateComps[i].Zero(id)
	}
	idGen.Free(uint32(id))
}

func incMaxEnts(n uint32) {
	for i := range stateComps {
		stateComps[i].Append(n)
	}
	for i := range oldStateComps {
		if oldStateComps[i] != nil {
			oldStateComps[i].Append(n)
		}
	}
}

func RegisterComp(id CompId, networked bool, sc StateComp) {
	if len(stateComps) < int(id)+1 {
		temp := make([]StateComp, id+1)
		copy(temp, stateComps)
		stateComps = temp

		temp = make([]StateComp, id+1)
		copy(temp, oldStateComps)
		oldStateComps = temp
	}
	if stateComps[id] != nil {
		panic("two components with same id")
	}

	stateComps[id] = sc

	if networked {
		oldSc := sc.Clone()
		oldStateComps[id] = oldSc.(StateComp)
		networkedComps = append(networkedComps, id)
	}
}

func CopyState() {
	for i := range stateComps {
		if oldStateComps[i] != nil {
			oldStateComps[i].Copy(stateComps[i])
		}
	}
}
