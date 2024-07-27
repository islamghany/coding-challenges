package commands

import "redis/foundation/store"

var (
	InvalidArguments = "Invalid arguments"
)

type Commander struct {
	Store *store.Store
}

func NewCommander(store *store.Store) *Commander {
	return &Commander{
		Store: store,
	}
}
