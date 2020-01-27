package shared

import "fmt"

type Range struct {
	Start, // 0 or -1 == start of []
	End int64 // -1 == end of []
}

type ByStart []*Range

func (a ByStart) Len() int { return len(a) }

func (a ByStart) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByStart) Less(i, j int) bool { return a[i].Start < a[j].Start }

func (r *Range) String() string {
	return fmt.Sprintf("%d-%d", r.Start, r.End)
}

func (r *Range) Size() int64 {
	return r.End - r.Start
}
