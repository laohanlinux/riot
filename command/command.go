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

func (cm Comand) doGet(key []byte) []byte {
	return nil
}

func (cm Comand) doSet(key, value []byte) error {
	return nil
}

func (cm Comand) doDel(key []byte) error {
	return nil
}

func (cm Comand) doInfo() error {
	return nil
}
