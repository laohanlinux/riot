package command

// ....
const (
	CmdGet = "get"
	CmdSet = "set"
	CmdDel = "del"
)

type Command struct {
	Op    string
	Key   string
	Value []byte
}

// DoGet returns value by specified key
func (cm Comand) DoGet(key []byte) []byte {
	return nil
}

func (cm Comand) DoSet(key, value []byte) error {
	return nil
}

func (cm Comand) DoDel(key []byte) error {
	return nil
}

func (cm Comand) doInfo() error {
	return nil
}

func (cm Command) coordinator() {

}
