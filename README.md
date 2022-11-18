# brainfuck-go

A [Brainfuck](https://en.wikipedia.org/wiki/Brainfuck) interpreter in Go.

## [Portability Issues](https://en.wikipedia.org/wiki/Brainfuck#Portability_issues)

### Cell size

Cells in this interpreter are considered as runes (32 bits), as it makes it simple to use cells as string characters. This means interacting with I/O is done via text rather than numbers.

### Array size

This implementation uses a dynamically-sized array, which is unlimited on the right. On the left size, the array begins at cell 0.

### End-of-line code

End of file code is considered to be 10.

### End-of-file behavior

Because it is easy to accommodate the "no change" behavior in Brainfuck programs, this interpreter does not change the value of the current cell when encountering EOF while reading input.
