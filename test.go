package main

import (
	"fmt"
)

type test struct {
	i int
	s string
}

func main() {
	f := []test{}

	for i := 0; i < 5; i++ {

		f = append(f, test{i: i, s: "hmm"})
	}

	comp := test{i: 0, s: "hmm"}

	for _, fun := range f {
		if fun != comp {
			fmt.Println(fun)
		}
	}
}
