package brainfuck

import (
	"bufio"
	"fmt"
	"io"
)

type Option interface {
	Apply(bf *Interpreter) error
}

type OptionFunc func(bf *Interpreter) error

func (of OptionFunc) Apply(bf *Interpreter) error {
	return of(bf)
}

func WithDebugMessages() OptionFunc {
	return func(bf *Interpreter) error {
		bf.debug = true

		return nil
	}
}

func WithOutput(output io.Writer) OptionFunc {
	return func(bf *Interpreter) error {
		bf.state.OutputWriter = bufio.NewWriter(output)

		return nil
	}
}

func WithInput(input io.Reader) OptionFunc {
	return func(bf *Interpreter) error {
		bf.state.InputReader = bufio.NewReader(input)

		return nil
	}
}

func WithInstruction(symbol rune, instruction Instruction) OptionFunc {
	return func(bf *Interpreter) error {
		_, exists := bf.symbolToInstruction[symbol]
		if exists {
			// don't overwrite existing instructions
			return fmt.Errorf("can not redefine instruction %v", symbol)
		}

		bf.symbolToInstruction[symbol] = instruction

		return nil
	}
}
