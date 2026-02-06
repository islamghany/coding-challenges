package calculator

type Stack[T any] struct {
	arr      []T
	capacity int
	length   int
}

func NewStack[T any](capacity int) *Stack[T] {
	return &Stack[T]{
		arr:      make([]T, capacity),
		capacity: capacity,
		length:   0,
	}
}

func (s *Stack[T]) Push(val T) {
	if s.length >= s.capacity {
		s.grow()
	}

	s.arr[s.length] = val
	s.length++

}
func (s *Stack[T]) grow() {
	newArr := make([]T, s.capacity*2)
	copy(newArr, s.arr)
	s.arr = newArr
	s.capacity *= 2
}
func (s *Stack[T]) Pop() T {
	if s.length == 0 {
		var zero T
		return zero

	}
	s.length--
	return s.arr[s.length]
}

func (s *Stack[T]) Peek() T {
	if s.length == 0 {
		var zero T
		return zero
	}
	return s.arr[s.length-1]
}

func (s *Stack[T]) Clear() {
	s.length = 0
}

func (s *Stack[T]) IsEmpty() bool {
	return s.length == 0
}
func (s *Stack[T]) ToArray() []T {
	a := make([]T, s.length)
	copy(a, s.arr[:s.length])
	return a
}
