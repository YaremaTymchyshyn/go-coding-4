package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var memory []int
var maxOutputWidth int

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please set memory size and max output width:")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	params := strings.Split(input, " ")

	memSize, _ := strconv.Atoi(params[0])
	maxOutputWidth, _ = strconv.Atoi(params[1])

	memory = make([]int, memSize)
	for i := range memory {
		memory[i] = -1
	}

	fmt.Println("Type 'help' for additional info.")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		parts := strings.Split(input, " ")
		command := parts[0]

		switch command {
		case "help":
			printHelp()
		case "exit":
			return
		case "print":
			printMemory()
		case "allocate":
			if len(parts) < 2 {
				fmt.Println("Please provide the number of cells to allocate.")
				continue
			}
			numCells, err := strconv.Atoi(parts[1])
			if err != nil || numCells <= 0 {
				fmt.Println("Invalid number of cells.")
				continue
			}
			allocateMemory(numCells)
		case "free":
			if len(parts) < 2 {
				fmt.Println("Please provide the block ID to free.")
				continue
			}
			blockID, err := strconv.Atoi(parts[1])
			if err != nil || blockID < 0 {
				fmt.Println("Invalid block ID.")
				continue
			}
			freeMemory(blockID)
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}

func printHelp() {
	fmt.Println(`Available commands:

 help  - show this help
 exit  - exit this program
 print - print memory blocks map
 allocate <num> - allocate <num> cells. Returns block first cell number
 free <num> - free block with first cell number <num>`)
}

func printMemory() {
	for i := 0; i < len(memory); i += maxOutputWidth {
		end := i + maxOutputWidth
		if end > len(memory) {
			end = len(memory)
		}
		fmt.Print("|")
		currentBlock := -1
		for j := i; j < end; j++ {
			if memory[j] == -1 {
				fmt.Print(" ")
			} else {
				if memory[j] != currentBlock {
					currentBlock = memory[j]
					fmt.Printf("%d", memory[j])
				} else {
					fmt.Print("x")
				}
			}
		}
		fmt.Println("|")
	}
}

func allocateMemory(numCells int) {
	start := -1
	for i := 0; i <= len(memory)-numCells; i++ {
		isFree := true
		for j := 0; j < numCells; j++ {
			if memory[i+j] != -1 {
				isFree = false
				break
			}
		}
		if isFree {
			start = i
			break
		}
	}

	if start == -1 {
		fmt.Println("Not enough memory to allocate.")
		return
	}

	for i := start; i < start+numCells; i++ {
		memory[i] = start
	}

	fmt.Println(start)
}

func freeMemory(blockID int) {
	for i := 0; i < len(memory); i++ {
		if memory[i] == blockID {
			memory[i] = -1
		}
	}
}
