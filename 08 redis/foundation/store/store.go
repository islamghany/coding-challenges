package store

import "redis/foundation/enconder/resp"

type Store struct {
	Set *Set
}

func NewStore() *Store {
	s := &Store{
		Set: NewSet(),
	}
	s.LoadAll()
	return s
}

// FlushAll save the data to the disk
func (s *Store) FlushAll() resp.RESPData {
	err := s.Set.Flush()
	if err != nil {
		return resp.NewError(err.Error())
	}
	return resp.NewSimpleString("OK")
}

func (s *Store) LoadAll() {
	s.Set.Load()
}
