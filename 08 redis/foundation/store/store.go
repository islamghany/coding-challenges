package store

type Store struct {
	Set *Set
}

func NewStore() *Store {
	return &Store{
		Set: NewSet(),
	}
}
