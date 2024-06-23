package slicesext

type Slice[T any] []T

func (s *Slice[T]) Len() int {
	return len(*s)
}

func (s *Slice[T]) Push(v T) {
	*s = append(*s, v)
}

func (s *Slice[T]) Pop() (res T) {
	if size := s.Len(); size > 0 {
		res = (*s)[size-1]
		*s = (*s)[:size-1]
	}
	return
}

func (s *Slice[T]) Clone() *Slice[T] {
	res := make(Slice[T], s.Len())
	copy(res, *s)
	return &res
}

func (s *Slice[T]) Filter(fn func(T) bool) *Slice[T] {
	res := make(Slice[T], 0, s.Len())
	for _, v := range *s {
		if fn(v) {
			res.Push(v)
		}
	}
	return &res
}

func (s *Slice[T]) Find(fn func(T) bool) (res T, ok bool) {
	for _, v := range *s {
		if fn(v) {
			return v, true
		}
	}
	return
}

func (s *Slice[T]) Traverse(fn func(int, T)) {
	for i, v := range *s {
		fn(i, v)
	}
}
