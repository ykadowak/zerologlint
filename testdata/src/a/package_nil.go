package a

import "golang.org/x/exp/slices"

// Cases return nil on v.Package() that resulted in panic.
func package_nil() {
	s := []int{1, 2, 3}
	slices.Sort(s)
}
