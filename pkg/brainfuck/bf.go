package brainfuck

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ferioss/brainfuck-go/pkg/stack"
)

// Interpreter is the main entity of this package. It contains an internal State, which is update by each instruction.
type Interpreter struct {
	symbolToInstruction map[rune]Instruction

	state *State

	debug bool
}

// State represents the state of the interpreter. Each Instruction that is executed will update the state, resulting in the next state.
// To extend the functionality of the language, new Instructions can be defined to modify the State in a custom way.
type State struct {
	Code []rune

	InputReader  *bufio.Reader
	OutputWriter *bufio.Writer

	ProgramCounter int
	Stack          stack.Stack[int]

	Data    []rune
	DataPtr int
}

// Instructions are functions that update the State.
type Instruction func(*State) error

// NewInterpreter creates a new Brainfuck interpreter which will read the code from codeReader. By default, it uses stdin/stdout for I/O. This behavior can be modified via options.
// Debug messages and custom instructions are also supported via options.
func NewInterpreter(codeReader io.Reader, options ...Option) (bf *Interpreter, err error) {
	// make a copy of the default symbol to instruction map
	s2i := make(map[rune]Instruction, len(symbolToInstruction))
	for k, v := range symbolToInstruction {
		s2i[k] = v
	}

	code, err := io.ReadAll(codeReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read code: %w", err)
	}

	state := &State{
		Code: []rune(string(code)),

		InputReader:  bufio.NewReader(os.Stdin),
		OutputWriter: bufio.NewWriter(os.Stdout),

		ProgramCounter: 0,
		Stack:          stack.Stack[int]{},

		Data:    make([]rune, 1),
		DataPtr: 0,
	}

	bf = &Interpreter{
		symbolToInstruction: s2i,

		state: state,
	}

	for _, o := range options {
		err := o.Apply(bf)
		if err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return bf, nil
}

// Runs the program.
func (bf *Interpreter) Run() error {
	defer func() {
		_, _ = bf.state.OutputWriter.WriteString("\n")
		_ = bf.state.OutputWriter.Flush()
	}()

	for bf.state.ProgramCounter < len(bf.state.Code) {
		instructionSymbol := bf.state.Code[bf.state.ProgramCounter]
		instruction, ok := bf.symbolToInstruction[instructionSymbol]
		if !ok {
			// instruction is invalid, ignore it
			bf.state.ProgramCounter++
			continue
		}

		if bf.debug {
			fmt.Fprintf(os.Stderr, "%v: %v \t cells: %v %v\n", bf.state.ProgramCounter, string(instructionSymbol), bf.state.DataPtr, bf.state.Data)
		}

		pc := bf.state.ProgramCounter
		err := bf.applyInstruction(instruction)
		if err != nil {
			err = fmt.Errorf("instruction %d (%v) failed to run: %w", pc, string(instructionSymbol), err)
			return err
		}
	}

	return nil
}

func (bf *Interpreter) applyInstruction(i Instruction) error {
	pc := bf.state.ProgramCounter

	err := i(bf.state)
	if err != nil {
		return err
	}

	if bf.state.ProgramCounter == pc {
		// do not increment program counter if the instruction changed it. This allows instructions to override normal execution flow.
		bf.state.ProgramCounter++
	}

	return nil
}
