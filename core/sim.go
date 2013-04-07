package core

import "io"

type State interface {
	Mutate(id EntId, value interface{})
	Clone() State
	SerDiff(buf io.Writer, newEnts []EntId, newSt State)
	DeserDiff(buf io.Reader, newEnts []EntId)
}

type StateMap map[string]State
type TransFunc func(Tick, StateMap, MutFuncs)

type MutFuncs struct {
	Mutate  func(tp string, id EntId, value interface{})
	Destroy func(id EntId)
	NewId   func() EntId
}

type Sim struct {
	states     *InfBuf
	emptySt    StateMap
	mutBufs    map[string]*MutBuf
	transFuncs [10][]TransFunc
	idGen      *IdGen
	time       Tick
	mut        MutFuncs
	serInfo    SerInfo
}

func NewSim(idGen *IdGen) *Sim {
	sim := new(Sim)
	sim.states = NewInfBuf()
	sim.emptySt = make(StateMap)
	sim.mutBufs = make(map[string]*MutBuf)
	for i := 0; i < len(sim.transFuncs); i++ {
		sim.transFuncs[i] = make([]TransFunc, 0)
	}
	sim.idGen = idGen

	sim.mut.Destroy = func(id EntId) {
		for _, m := range sim.mutBufs {
			m.Write(Mutation{id, nil})
		}
		if idGen != nil {
			sim.idGen.Free(id)
		}
	}
	sim.mut.Mutate = func(tp string, id EntId, value interface{}) {
		sim.mutBufs[tp].Write(Mutation{id, value})
	}
	sim.mut.NewId = func() EntId {
		return idGen.Next()
	}

	st := make(StateMap)
	sim.time = sim.states.Write(st)
	return sim
}

func (sim *Sim) Time() Tick {
	return sim.time
}

func (sim *Sim) Mutate(tp string, id EntId, value interface{}) {
	sim.mut.Mutate(tp, id, value)
}

func (sim *Sim) Destroy(id EntId) {
	sim.mut.Destroy(id)
}

func (sim *Sim) NewId() EntId {
	return sim.mut.NewId()
}

func (sim *Sim) AddTransFunc(pr int, f TransFunc) {
	sim.transFuncs[pr] = append(sim.transFuncs[pr], f)
}

func (sim *Sim) SetSerInfo(si SerInfo) {
	sim.serInfo = si
}

func (sim *Sim) Diff(buf io.Writer, id EntId, t1, t2 Tick) {
	refState := sim.emptySt
	if t1 > 0 {
		refState = sim.states.Read(t1).(StateMap)
	}
	Serialize(sim.serInfo, buf, id, refState, sim.states.Read(t2).(StateMap))
}

func (sim *Sim) AddState(name string, state State) {
	st := sim.states.Read(sim.time).(StateMap)
	st[name] = state
	sim.mutBufs[name] = NewMutBuf()

	sim.emptySt[name] = state.Clone()
}

func (sim *Sim) DeserDiff(buf io.Reader, refTime, time Tick) {
	sim.states.Reset(time)
	refState := sim.emptySt
	if refTime > 0 {
		refState = sim.states.Read(refTime).(StateMap)
	}
	newState := make(StateMap)
	for s := range refState {
		newState[s] = refState[s].Clone()
	}
	Deserialize(sim.serInfo, buf, newState, sim.mut)
	sim.time = sim.states.Write(newState)
}

func (sim *Sim) Update() {
	done := make(chan struct{})
	currState := sim.states.Read(sim.time).(StateMap)
	newState := make(StateMap)
	for s := range currState {
		newState[s] = currState[s].Clone()
	}

	for pr, _ := range sim.transFuncs {
		for _, f := range sim.transFuncs[pr] {
			g := f
			go func() {
				g(sim.time, newState, sim.mut)
				done <- struct{}{}
			}()
		}

		for i := 0; i < len(sim.transFuncs[pr]); i++ {
			<-done
		}
		sim.applyMuts(newState)
	}
	sim.time = sim.states.Write(newState)
}

func (sim *Sim) applyMuts(st StateMap) {
	for s := range st {
		for m, ok := sim.mutBufs[s].Read(); ok; {
			st[s].Mutate(m.Id, m.Value)
			m, ok = sim.mutBufs[s].Read()
		}
	}
}
