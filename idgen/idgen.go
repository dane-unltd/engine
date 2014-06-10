package idgen

type useReq struct {
	id  uint32
	ret chan bool
}

type IdGen struct {
	maxId    uint32
	freeIds  []uint32
	idOut    chan uint32
	idIn     chan uint32
	idUse    chan useReq
	maxIdInc func(uint32)
}

func New(idInc func(uint32)) *IdGen {
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

func (g *IdGen) Use(id uint32) bool {
	req := useReq{
		id:  id,
		ret: make(chan bool, 1),
	}
	g.idUse <- req
	return <-req.ret
}

func (g *IdGen) run() {
	currId := g.maxId + 1
	g.maxIdInc(g.maxId)
	for {
		select {
		case id := <-g.idIn:
			g.freeIds = append(g.freeIds, id)
			currId = id
		case g.idOut <- currId:
			if currId > g.maxId {
				g.maxId = currId
				g.maxIdInc(g.maxId)
				currId++
			} else {
				g.freeIds = g.freeIds[:len(g.freeIds)-1]
				if len(g.freeIds) > 0 {
					currId = g.freeIds[len(g.freeIds)-1]
				} else {
					currId = g.maxId + 1
				}
			}
		case req := <-g.idUse:
			if req.id > g.maxId {
				for i := g.maxId + 1; i < req.id; i++ {
					g.freeIds = append(g.freeIds, i)
				}
				currId = g.freeIds[len(g.freeIds)-1]
				g.maxId = req.id
				g.maxIdInc(g.maxId)
				req.ret <- true
			} else {
				i := 0
				for ; i < len(g.freeIds); i++ {
					if g.freeIds[i] == req.id {
						break
					}
				}
				if i < len(g.freeIds) {
					g.freeIds[i] = g.freeIds[len(g.freeIds)-1]
					g.freeIds = g.freeIds[:len(g.freeIds)-1]
					currId = g.freeIds[len(g.freeIds)-1]
					req.ret <- true
				} else {
					req.ret <- false
				}
			}
		}
	}
}
