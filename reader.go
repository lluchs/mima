// MIMA Assembler Reader

package mima

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Program struct {
	Marks        map[string]uint32
	Instructions map[uint32]Instruction
}

type Instruction struct {
	Op       string
	Argument string
}

// Parses the instruction's argument if it is a number.
func (i *Instruction) ParseArgument() (uint32, error) {
	return parseNumber(i.Argument)
}

func parseNumber(num string) (uint32, error) {
	hexRegex := regexp.MustCompile(`^\$[0-9a-fA-F]+$`)
	decRegex := regexp.MustCompile(`^\d+$`)
	switch {
	case hexRegex.MatchString(num):
		result, _ := strconv.ParseUint(num[1:], 16, 24)
		return (uint32)(result), nil
	case decRegex.MatchString(num):
		result, _ := strconv.ParseUint(num, 10, 24)
		return (uint32)(result), nil
	}
	return 0, errors.New("Not a number")
}

// Returns whether a string is empty, i.e. contains only whitespace.
func empty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// Parses a MIMA assembler program.
func Parse(in io.Reader) (*Program, error) {
	// Regex matching comments
	// ; comment
	commentRegex := regexp.MustCompile(`;.*$`)
	// Regex matching constants or load point directives
	// * = 123
	// MAX = $12F
	constantRegex := regexp.MustCompile(`^\s*(\*|\w+)\s*=\s*(\$[0-9a-fA-F]+|\d+)\s*$`)
	// Regex matching instructions with optional marks.
	// LABEL LDV FOO
	//       ADD EINS
	//       HALT
	instructionRegex := regexp.MustCompile(`^(\w*)\s+(\w+)(\s+(\$[0-9a-fA-F]+|\w+))?\s*$`)

	scanner := bufio.NewScanner(in)
	lineNum := 0
	program := &Program{Marks: make(map[string]uint32), Instructions: make(map[uint32]Instruction)}
	// Loading Point: where to insert instructions
	var lp uint32 = 0
	// Read input linewise.
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		// Remove comments.
		line = commentRegex.ReplaceAllString(line, "")
		// Skip empty lines.
		if empty(line) {
			continue
		}

		// Try to match the line.
		if match := constantRegex.FindStringSubmatch(line); match != nil {
			// Parse and validate the argument (required).
			arg, err := parseNumber(match[2])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Invalid number in line %d: %s", lineNum, err.Error()))
			}
			// Check whether the match is a loading point or a constant.
			if match[1] == "*" {
				// Adjust the loading point.
				lp = arg
			} else {
				// Save constant.
				program.Marks[match[1]] = arg
			}
			// Continue to avoid increasing the load point by 1.
			continue
		} else if match := instructionRegex.FindStringSubmatch(line); match != nil {
			// Handle instructions.
			// Check for a mark.
			if !empty(match[1]) {
				// Save the mark.
				program.Marks[match[1]] = lp
			}
			// Add an instruction.
			program.Instructions[lp] = Instruction{
				Op:       match[2],
				Argument: match[4],
			}
		} else {
			// Invalid line.
			return nil, errors.New(fmt.Sprint("Parse Error in line ", lineNum, ": ", line))
		}

		// Write the next instruction to the next loading point.
		lp++
	}
	return program, nil
}
