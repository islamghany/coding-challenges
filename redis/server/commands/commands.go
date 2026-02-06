package commands

import "redis/foundation/store"

var (
	InvalidArguments = "Invalid arguments"
	WrongType        = "WRONGTYPE Operation against a key holding the wrong kind of value"
)

type Commander struct {
	Store *store.Store
}

func NewCommander(store *store.Store) *Commander {
	return &Commander{
		Store: store,
	}
}

func (cmdr *Commander) Flush() {
	cmdr.Store.FlushAll()
}
