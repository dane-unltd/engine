package idgen

type IdGen struct {
	maxId    uint32
	freeIds  []uint32
	idOut    chan uint32
	idIn     chan uint32
	maxIdInc func()
}

func New(idInc func()) *IdGen {
	g := IdGen{}
	g.idOut = make(chan uint32)
	g.idIn = make(chan uint32, 8)
	g.maxIdInc = idInc
	go g.run()
	return &g
}

func (g *IdGen) Next() uint32 {
	return <-g.idOut
}

func (g *IdGen) Free(id uint32) {
	g.idIn <- id
}

func (g *IdGen) run() {
	g.maxIdInc()
	g.maxIdInc()
	currId := g.maxId + 1
	for {
		select {
		case id := <-g.idIn:
			g.freeIds = append(g.freeIds, id)
			currId = id
		case g.idOut <- currId:
			if currId > g.maxId {
				g.maxId = currId
				g.maxIdInc()
				currId++
			} else {
				g.freeIds = g.freeIds[:len(g.freeIds)-1]
				if len(g.freeIds) > 0 {
					currId = g.freeIds[len(g.freeIds)-1]
				} else {
					currId = g.maxId + 1
				}
			}
		}
	}
}
