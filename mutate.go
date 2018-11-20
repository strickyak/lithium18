package Li

func RandomInt(mod int) int {
	buf := RandomBytes(4)
	x := int(buf[0])
	x = 256*x + int(buf[1])
	x = 256*x + int(buf[2])
	x = 256*x + int(buf[3])
	return (x%mod + mod) % mod
}

type MutateHow int

const (
	INSERT MutateHow = iota
	DELETE
	SPLICE
	NUM_HOW
)

func MutatePopulation(pop *Population, numMutants int) {
	nc := len(pop.Creatures)

	for i := 0; i < numMutants; i++ {
		// Pick a victim from the second half.
		victim := RandomInt(nc/2) + (nc / 2)

		// Pick a donor from anywhere
		donor := RandomInt(nc)
		code := pop.Creatures[donor].Code

		// Pick a mutation
		how := RandomInt(int(NUM_HOW))
		switch MutateHow(how) {
		case INSERT:
			code = Insert1(code)
		case DELETE:
			code = Delete1(code)
		case SPLICE:
			donor2 := RandomInt(nc)
			code2 := pop.Creatures[donor2].Code
			code = Splice(code, code2)
		}
		pop.Creatures[victim] = &Creature{
			Code: code,
		}
	}
}

func Insert1(code []byte) []byte {
	lc := len(code)
	z := make([]byte, lc+1)
	i := RandomInt(lc + 1)
	if i > 0 {
		copy(z[:i], code[:i])
	}
	z[i] = byte(RandomInt(256))
	if i < lc {
		copy(z[i+1:], code[i:])
	}
	return z
}

func Delete1(code []byte) []byte {
	lc := len(code)
	if lc < 2 {
		return code
	}
	z := make([]byte, lc-1)
	i := RandomInt(lc)
	if i > 0 {
		copy(z[:i], code[:i])
	}
	if i+1 < lc {
		copy(z[i:], code[i+1:])
	}
	return z
}

func Splice(a, b []byte) []byte {
	la := len(a)
	lb := len(b)
	if la < 2 {
		return b
	}
	if lb < 2 {
		return a
	}
	ia := RandomInt(la-1) + 1
	ib := RandomInt(lb-1) + 1
	z := make([]byte, ia+ib)
	copy(z[:ia], a[:ia])
	copy(z[ia:], b[lb-ib:])
	return z
}
