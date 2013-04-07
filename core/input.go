package core

type CmdId byte

type UserCmd struct {
	Id      EntId
	Time    Tick
	Status  uint32
	Target  [3]float64
	Actions uint32
}

type CmdSrc interface {
	Cmd(id EntId) UserCmd
	Active(id EntId, cmd CmdId) bool
	Target(id EntId) [3]float64
}

type Input struct{}

func (u UserCmd) Copy() interface{} {
	return u
}

func (u UserCmd) Equals(v interface{}) bool {
	u2, ok := v.(UserCmd)
	if !ok {
		return ok
	}
	return u == u2
}

func (u UserCmd) Active(cmd CmdId) bool {
	return u.Actions&(1<<cmd) != 0
}

func NewInput() *Input {
	return &Input{}
}

func (i *Input) Cmd(id EntId) UserCmd {
	return UserCmd{}
}

func (i *Input) Active(id EntId, cmd CmdId) bool {
	return false
}

func (i *Input) Target(id EntId) [3]float64 {
	return [3]float64{0, 0, 0}
}
