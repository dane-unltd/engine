package core

type Comp interface {
	Id() string

	Swap()
	Update()
	Init(sim Sim, res *ResMgr)
}
