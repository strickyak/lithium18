package Li

import (
	crand "crypto/rand"
	"flag"
	"fmt"
	mrand "math/rand"
	"sort"
)

var NumSTEPS = flag.Int("s", 100000, "Max num steps per code execution")
var NumCREATURES = flag.Int("c", 500, "Num creatures per population")
var OutputMAX = flag.Int("o", 100, "Max num outputs per code execution")
var SEED = flag.Int("seed", 0, "Random generator seed")

type Population struct {
	Creatures []*Creature
}

type Creature struct {
	Code   []byte
	Output []Num
	Score  float64
}

func RunEngine(code []byte) []Num {
	var vec []Num
	enough := false
	e := &Engine{
		Code: code,
		Emit: func(x Num) {
			vec = append(vec, x)
			if len(vec) >= *OutputMAX {
				enough = true
			}
		},
	}
	for step := 0; step < *NumSTEPS && !enough; step++ {
		e.Step()
	}
	return vec
}
func RandomBytes(n int) []byte {
	var err error
	buf := make([]byte, n)
	if *SEED == 0 {
		_, err = crand.Read(buf)
	} else {
		_, err = mrand.Read(buf)
	}
	if err != nil {
		panic(err)
	}
	return buf
}
func RandomCode() []byte {
	buf := RandomBytes(1)
	codeLen := int(buf[0]) + 1
	return RandomBytes(codeLen)
}

func InitialPopulation() *Population {
	pop := &Population{
		Creatures: make([]*Creature, *NumCREATURES),
	}
	for i := 0; i < *NumCREATURES; i++ {
		pop.Creatures[i] = &Creature{
			Code: RandomCode(),
		}
	}
	return pop
}

func Score(got, want []Num, code []byte) float64 {
	n := len(want)
	var total float64
	for i := 0; i < n; i++ {
		var x Num = 0
		if i < len(got) {
			x = got[i]
		}
		var diff float64
		if float64(x) < float64(want[i]) {
			diff = float64(want[i]) - float64(x)
		} else {
			diff = float64(x) - float64(want[i])
		}
		// total += diff * float64(n-i) * float64(n-i)
		total += diff * diff * float64(n-i)
	}
	return total / float64(n)
	// return total/float64(n) + 0.01*float64(len(code))
}

func ScorePopulation(pop *Population, want []Num) {
	par := NewParallel(*NumCREATURES, func(in interface{}) interface{} {
		c := in.(*Creature)
		got := RunEngine(c.Code)
		c.Output = got
		c.Score = Score(got, want, c.Code)
		return c
	})

	tasks := 0
	for i, c := range pop.Creatures {
		if c.Score != 0 {
			continue
		}
		if i%10 == 9 {
			fmt.Printf(".")
		}
		/*
			got := RunEngine(c.Code)
			c.Output = got
			c.Score = Score(got, want, c.Code)
		*/
		par.Add1(c)
		tasks++
	}
	for i := 0; i < tasks; i++ {
		par.Wait1()
	}
	par.Finish()
	sort.Slice(pop.Creatures, func(i, j int) bool {
		return pop.Creatures[i].Score < pop.Creatures[j].Score
	})
}
