package cmd

// ....
const (
	CmdGet          = "GET"
	CmdSet          = "SET"
	CmdDel          = "DEL"
	CmdShare        = "SHARE"
	CmdGetBucket    = "GET BUCKET"
	CmdCreateBucket = "CREATE BUCKET"
	CmdDelBucket    = "DEL BUCKET"
)

const (
	QsConsistent = iota
	QsRandom
)
