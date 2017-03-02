package renders

import (
	"container/list"
)

func Range(l int) []struct{} {
	return make([]struct{}, l)
}

func RangeN(start, last int) map[int]struct{} {
	var res = make(map[int]struct{})
	for i := start; i < last; i++ {
		res[i] = struct{}{}
	}
	return res
}

func Add(left, right int) int {
	return left + right
}

func Sub(left, right int) int {
	return left - right
}

func List(l *list.List) chan interface{} {
	e := l.Front()
	c := make(chan interface{})
	go func() {
		for e != nil {
			c <- e.Value
			e = e.Next()
		}
		close(c)
	}()
	return c
}
