package stack

import "errors"

type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(val T) {
	s.data = append(s.data, val)
}

func (s *Stack[T]) Pop() (val T, err error) {
	if len(s.data) == 0 {
		return val, errors.New("can not pop from empty stack")
	}

	val = s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]

	return val, nil
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}
