package main

import (
	"flag"
	"log"
	"os"

	"github.com/ferioss/brainfuck-go/pkg/brainfuck"
)

func main() {
	input := flag.String("i", "example.bf", "path to a brainfuck program")
	flag.Parse()

	file, err := os.Open(*input)
	defer func() { _ = file.Close() }()

	if err != nil {
		log.Fatal(err)
	}

	it, err := brainfuck.NewInterpreter(file,
		brainfuck.WithInput(os.Stdin),
		brainfuck.WithOutput(os.Stdout),

		// brainfuck.WithDebugMessages(), // print messages to stderr that are useful for debugging.

		brainfuck.WithInstruction('*', func(s *brainfuck.State) error { s.Data[s.DataPtr] *= s.Data[s.DataPtr]; return nil }), // custom instruction to square current value
		brainfuck.WithInstruction('~', func(s *brainfuck.State) error { s.DataPtr = 0; return nil }),                          // custom instruction to move data pointer to 0
	)
	if err != nil {
		log.Fatal(err)
	}

	err = it.Run()
	if err != nil {
		log.Fatal(err)
	}

}
