package brainfuck

// This file contains the Instructions for vanilla Brainfuck language.

import (
	"errors"
	"fmt"
	"io"
)

var (
	symbolToInstruction = map[rune]Instruction{
		'>': incrPtr,
		'<': decrPtr,
		'+': incrData,
		'-': decrData,
		'.': writeOutput,
		',': readInput,
		'[': beginLoop,
		']': endLoop,
	}
)

// Increment (increase by one) the byte at the data pointer.
func incrData(s *State) error {
	s.Data[s.DataPtr]++
	return nil
}

// Decrement (decrease by one) the byte at the data pointer.
func decrData(s *State) error {
	s.Data[s.DataPtr]--
	return nil
}

// Expand data if necessary, so s.DataPtr is always in range of s.Data.
func expandData(s *State) {
	if len(s.Data) <= s.DataPtr {
		s.Data = append(s.Data, make([]rune, 1+s.DataPtr-len(s.Data))...)
	}
}

// Increment the data pointer (to point to the next cell to the right).
func incrPtr(s *State) error {
	s.DataPtr++
	expandData(s)

	return nil
}

// Decrement the data pointer (to point to the next cell to the left).
func decrPtr(s *State) error {
	if s.DataPtr == 0 {
		return errors.New("can not decrement data pointer: data pointer is zero")
	}

	s.DataPtr--

	return nil
}

// Output the byte at the data pointer.
func writeOutput(s *State) error {
	_, err := s.OutputWriter.WriteRune(s.Data[s.DataPtr])
	if err != nil {
		err = fmt.Errorf("failed to write output: %w", err)
		return err
	}

	return nil
}

// Accept one byte of input, storing its value in the byte at the data pointer.
func readInput(s *State) error {
	val, _, err := s.InputReader.ReadRune()
	if err != nil {
		if err != io.EOF {
			err = fmt.Errorf("failed to read output: %w", err)
			return err
		}
		// do not change value on EOF

		return nil
	}

	s.Data[s.DataPtr] = val

	return nil
}

// If the byte at the data pointer is zero, then instead of moving the instruction pointer forward to the next command, jump it forward to the command after the matching ] command.
func beginLoop(s *State) error {
	s.Stack.Push(s.ProgramCounter)

	if s.Data[s.DataPtr] != 0 {
		return nil
	}

	for i := s.ProgramCounter + 1; i < len(s.Code); i++ {
		switch s.Code[i] {
		case '[':
			s.Stack.Push(i)
		case ']':
			pc, err := s.Stack.Pop()
			if err != nil {
				return fmt.Errorf("can not end loop at %d: no matching '[' found", i)
			}

			if pc == s.ProgramCounter {
				// found the matching ] for the original [
				s.ProgramCounter = i + 1

				return nil
			}
		}
	}

	return fmt.Errorf("can not skip loop: no matching ']' found")
}

// If the byte at the data pointer is nonzero, then instead of moving the instruction pointer forward to the next command, jump it back to the command after the matching [ command.
func endLoop(s *State) error {
	pc, err := s.Stack.Pop()
	if err != nil {
		return errors.New("can not end loop: no matching '[' found")
	}

	if s.Data[s.DataPtr] == 0 {
		return nil
	}

	s.ProgramCounter = pc
	return nil
}
