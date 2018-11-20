// +build main

package main

import (
	"bufio"
	"flag"
	"fmt"
	mrand "math/rand"
	"os"

	. "github.com/strickyak/lithium18"
)

var GOAL = flag.String("goal", "count", "Which goal to strive for")

func EmitNum(w *bufio.Writer, x Num) {
	fmt.Fprintf(w, "%d ", x)
	w.Flush()
}

func main() {
	flag.Parse()
	if *SEED != 0 {
		mrand.Seed(int64(*SEED))
	}
	w := bufio.NewWriter(os.Stdout)
	defer fmt.Fprintf(w, "\n")
	defer w.Flush()

	want := make([]Num, *OutputMAX)
	switch *GOAL {
	case "count":
		for i := 0; i < *OutputMAX; i++ {
			want[i] = Num(i + 1)
		}
	case "fib":
		want[0] = 1
		want[1] = 1
		for i := 2; i < *OutputMAX; i++ {
			want[i] = want[i-1] + want[i-2]
		}
	default:
		panic("bad goal")
	}

	fmt.Printf("GOAL: ")
	for _, e := range want {
		fmt.Printf("%d ", e)
	}
	fmt.Printf("\n")

	pop := InitialPopulation()
	for gen := 0; gen < 1000000; gen++ {
		ScorePopulation(pop, want)
		fmt.Printf("\n%d: [", gen)
		for i := 0; i < 10; i++ {
			fmt.Printf("%.4g ", pop.Creatures[i].Score)
		}
		fmt.Printf("]\n%d: (", gen)
		for _, x := range pop.Creatures[0].Output {
			fmt.Printf("%d ", x)
		}
		fmt.Printf(")\n")

		for i := 0; i < 10; i++ {
			for j := 0; j < 20 && j < len(pop.Creatures[i].Code); j++ {
				op := pop.Creatures[i].Code[j]
				fmt.Printf("%x.%s ", op, FormatOpcode(op))
			}
			fmt.Printf("\n")
		}

		MutatePopulation(pop, len(pop.Creatures)/20)
	}
}
