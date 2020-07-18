package main

func newStack() *stack {
	return &stack{
		make([]int, 0),
	}
}

func (s *stack) push(a addrType) {
	s.data = append(s.data, a)
}

func (s *stack) pop() addrType {
	result := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]

	return result
}

func (s *stack) peek(offset int) addrType {
	return s.data[len(s.data)-1-offset]
}
