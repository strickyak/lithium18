package Li

import (
	"fmt"
)

const N = 8

type Num int16

type Stack struct {
	X [N]Num
	P int
}

type Opcode byte

const (
	Alu Opcode = iota
	System
	Store
	Recall
	Call
	Jump
	BranchIfZero
	BranchUnlessZero
	Lit
	NUM_OPS
)

type AluOp byte

const (
	Add AluOp = iota
	Neg
	Mul
	Mod
	Dup
	Pop
	LE
	LT
)

type SystemOp byte

const (
	Emit1 SystemOp = iota
	Emit2
	Emit3
)

type Engine struct {
	Code []byte
	PC   int

	Global Stack
	Data   Stack
	Return Stack

	Emit func(Num)
}

func Inc(p int) int {
	return (p + 1) % N
}
func Dec(p int) int {
	return (p + N - 1) % N
}

func (s *Stack) Push(x Num) {
	s.X[s.P] = x
	s.P = Inc(s.P)
}
func (s *Stack) PushBool(x bool) {
	if x {
		s.Push(1)
	} else {
		s.Push(0)
	}
}
func (s *Stack) Pop() Num {
	s.P = Dec(s.P)
	return s.X[s.P]
}
func (s *Stack) Peek() Num {
	return s.X[Dec(s.P)]
}

func (e *Engine) Step() {
	opcode := e.Code[e.PC]
	e.SetPC(e.PC + 1)
	cmd, arg := (opcode/N)%byte(NUM_OPS), (opcode % N)

	switch Opcode(cmd) {
	case Alu:
		e.AluStep(AluOp(arg))
	case System:
		e.SystemStep(SystemOp(arg))
	case Store:
		e.Global.X[arg] = e.Data.Pop()
	case Recall:
		e.Data.Push(e.Global.X[arg])
	case Call:
		e.Return.Push(Num(e.PC))
		e.JumpToMark(arg)
	case Jump:
		e.JumpToMark(arg)
	case BranchIfZero:
		x := e.Data.Pop()
		if x == 0 {
			offset := (int(arg)+N)%N - (N / 2)
			if offset >= 0 {
				offset++
			}
			e.SetPC(e.PC + offset)
		}
	case BranchUnlessZero:
		x := e.Data.Pop()
		if x != 0 {
			offset := (int(arg)+N)%N - (N / 2)
			if offset >= 0 {
				offset++
			}
			e.SetPC(e.PC + offset)
		}
	case Lit:
		e.Data.Push(Num(arg))
	}
}
func (e *Engine) JumpToMark(arg byte) {
	n := len(e.Code)
	for i := 1; i < n/2; i++ {
		if e.Code[(e.PC+i)%n]%N == arg {
			e.SetPC(e.PC + i)
			break
		}
	}
}
func (e *Engine) AluStep(arg AluOp) {
	switch arg {
	case Add:
		x := e.Data.Pop()
		y := e.Data.Pop()
		e.Data.Push(x + y)
	case Neg:
		e.Data.Push(-e.Data.Pop())
	case Mul:
		x := e.Data.Pop()
		y := e.Data.Pop()
		e.Data.Push(x * y)
	case Mod:
		x := e.Data.Pop()
		y := e.Data.Pop()
		if y == 0 {
			e.Data.Push(0)
		} else if y < 0 { // Use abs(y)
			e.Data.Push(x % (-y))
		} else {
			e.Data.Push(x % y)
		}
	case Dup:
		e.Data.Push(e.Data.Peek())
	case Pop:
		e.Data.Pop()
	case LE:
		x := e.Data.Pop()
		y := e.Data.Pop()
		e.Data.PushBool(x <= y)
	case LT:
		x := e.Data.Pop()
		y := e.Data.Pop()
		e.Data.PushBool(x < y)
	}
}
func (e *Engine) SystemStep(arg SystemOp) {
	switch arg {
	case Emit1, Emit2, Emit3:
		x := e.Data.Pop()
		e.Emit(x)
	default:
		e.SetPC(int(e.Return.Pop()))
	}
}
func (e *Engine) SetPC(pc int) {
	lc := len(e.Code)
	e.PC = (pc%lc + lc) % lc
}

func FormatOpcode(opcode byte) string {
	cmd, arg := (opcode/N)%byte(NUM_OPS), (opcode % N)

	switch Opcode(cmd) {
	case Alu:
		return FormatAlu(AluOp(arg))
	case System:
		return FormatSystem(SystemOp(arg))
	case Store:
		return fmt.Sprintf("STO[%d]", arg)
	case Recall:
		return fmt.Sprintf("RCL[%d]", arg)
	case Call:
		return fmt.Sprintf("CAL[%d]", arg)
	case Jump:
		return fmt.Sprintf("JMP[%d]", arg)
	case BranchIfZero:
		offset := (int(arg)+N)%N - (N / 2)
		if offset >= 0 {
			offset++
		}
		return fmt.Sprintf("BZ[%d]", offset)
	case BranchUnlessZero:
		offset := (int(arg)+N)%N - (N / 2)
		if offset >= 0 {
			offset++
		}
		return fmt.Sprintf("BNZ[%d]", offset)
	case Lit:
		return fmt.Sprintf("LIT[%d]", arg)
	default:
		panic(opcode)
	}
}

func FormatAlu(arg AluOp) string {
	switch arg {
	case Add:
		return "ALU[Add]"
	case Neg:
		return "ALU[Neg]"
	case Mul:
		return "ALU[Mul]"
	case Mod:
		return "ALU[Mod]"
	case Dup:
		return "ALU[Dup]"
	case Pop:
		return "ALU[Pop]"
	case LE:
		return "ALU[LE]"
	case LT:
		return "ALU[LT]"
	default:
		panic(arg)
	}
}

func FormatSystem(arg SystemOp) string {
	switch arg {
	case Emit1, Emit2, Emit3:
		return "Emit"
	default:
		return "Return"
	}
}
