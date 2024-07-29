package pgq

const (
	AddMsg = `INSERT INTO tgmsg (msg) VALUES (@msg);`
)

const (
	AddHistory = `INSERT INTO tghistory (msg) VALUES (@msg);`
)
