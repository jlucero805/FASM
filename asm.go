package main

import (
	"fmt"
	"bufio"
	"log"
	"os"
	"strings"
	"strconv"
)

func main() {
	labels := make(map[string]int)
	commands := getCommands(os.Args[1], labels)
	registers := initializeRegisters()
	interpret(commands, registers, labels)
	displayRegisters(registers)
}

// Initialize registers according to RISC-V specification
func initializeRegisters() map[string]int {
	registers := make(map[string]int)

	// zero -> always 0
	registers["zero"] = 0

	// ra -> return address
	registers["ra"] = 0

	// sp -> stack pointer
	registers["sp"] = 0

	// t0 - t6 -> temp registers
	for i := 0; i < 7; i++ {
		name := "t" + strconv.Itoa(i)
		registers[name] = 0
	}

	// a0 - a6 -> function registers
	for i := 0; i < 7; i++ {
		name := "a" + strconv.Itoa(i)
		registers[name] = 0
	}

	// s0 - s11 -> permanent registers
	for i := 0; i < 12; i++ {
		name := "s" + strconv.Itoa(i)
		registers[name] = 0
	}

	return registers
}

func displayRegisters(registers map[string]int) {
	// zero -> always 0
	fmt.Println("zero:\t", registers["zero"])

	fmt.Println("ra:\t", registers["ra"], "\t\tsp:\t", registers["sp"])
	fmt.Println("-----------------------------------")

	for i := 0; i < 7; i++ {
		name1 := "t" + strconv.Itoa(i)
		name2 := "a" + strconv.Itoa(i)
		fmt.Println(name1 + ":\t", registers[name1], "\t\t" + name2 + ":\t", registers[name2])
	}

	fmt.Println("-----------------------------------")

	for i := 0; i < 6; i++ {
		name1 := "s" + strconv.Itoa(i)
		name2 := "s" + strconv.Itoa(i + 6)
		fmt.Println(name1 + ":\t", registers[name1], "\t\t" + name2 + ":\t", registers[name2])
	}
}

// Reads instructions from file into [][]
func getCommands(filename string, labels map[string]int) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	var text [][]string

	i := 0
	for scanner.Scan() {
		fullline := scanner.Text()
		line := strings.Fields(fullline)
		if len(line) > 0 {
			if line[0][0] != '/' {
				text = append(text, line)
				if line[0][0] == '$' {
					labels[line[0]] = i
				}
				i++
			}
		}
	}

	return text
}

func interpret(text [][]string, registers map[string]int, labels map[string]int) {
	pc := labels["$main"]
	var stack [8096]int
	for pc < len(text) {
		line := text[pc]
		switch line[0] {
		case "li":
			registers[line[1]], _ = strconv.Atoi(line[2])
		case "add":
			registers[line[1]] = registers[line[2]] + registers[line[3]]
		case "addi":
			num, _ := strconv.Atoi(line[3])
			registers[line[1]] = registers[line[2]] + num
		case "b", "beq", "bgt", "ble":
			pc = bFamilyHandler(line, registers, labels, pc)
		case "jal":
			registers["ra"] = pc
			pc = labels[line[1]]
		case "ret":
			pc = registers["ra"]
		case "sw":
			stack[registers[line[1]]] = registers[line[2]]
		case "lw":
			registers[line[1]] = stack[registers[line[2]]]
		case "ecall":
			eFamilyHandler(line, registers)
		case "debug":
			if line[1] == "stack" {
				fmt.Println(stack[0:10])
			} else {
				fmt.Println(registers[line[1]])
			}
		case "pc":
			fmt.Println(pc)
		default:
		}
		pc++
	}
}

// handles instructions of b family => comparison type instructions
func bFamilyHandler(line []string, registers map[string]int, labels map[string]int, pc int) int {
	switch line[0] {
	case "b":
		pc = labels[line[1]]
	case "bgt":
		left := registers[line[1]]
		right := registers[line[2]]
		if left > right {
			pc = labels[line[3]]
		}
	case "beq":
		if registers[line[1]] == registers[line[2]] {
			pc = labels[line[3]]
		}
	case "ble":
		if registers[line[1]] <= registers[line[2]] {
			pc = labels[line[3]]
		}
	default:
		break
	}
	return pc
}

// handles special instructions
func eFamilyHandler(line []string, registers map[string]int) {
	switch line[1] {
	case "1":
		// print int in register a0
		fmt.Println(registers["a0"])
	case "5":
		// read int into register a0
		var input int
		fmt.Scanln(&input)
		registers["a0"] = input
	default:
		break
	}
}