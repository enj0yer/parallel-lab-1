package main

import (
	"fmt"
	"os"
	"strconv"
)

func defaults() {
	fmt.Println("Available: ")
	fmt.Println("    gen <filename> <size>")
	fmt.Println("    seq <filename>")
	fmt.Println("    sim <filename> <threads>")
	fmt.Println("    vis <filename>")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No command provided")
		defaults()
	}
	command := os.Args[1]
	switch command {
	case "gen":
		if len(os.Args) < 4 {
			fmt.Println("usage: generate <filename> <size>")
			return
		}
		filename := os.Args[2]
		size, err := strconv.Atoi(os.Args[3])
		if err != nil && size <= 0 {
			fmt.Printf("Error: size must be a positive number: %v\n", err)
			return
		}
		generate(filename, size)

	case "seq":
		if len(os.Args) < 3 {
			fmt.Println("usage: seq <filename>")
			return
		}
		filename := os.Args[2]
		data, err := readFile(filename)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Println(CountExecutionTime(func() {
			ProcessSequentially(data, Pow2)
		}, fmt.Sprintf("Sequentially processing %d elements", len(data))))

	case "sim":
		if len(os.Args) < 4 {
			fmt.Println("usage: seq <filename> <threads>")
			return
		}
		filename := os.Args[2]
		threads, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Printf("Error: threads amount must be a pozitive number: %v", err)
			return
		}
		data, err := readFile(filename)

		if err != nil {
			fmt.Printf("Error: %v", err)
		}

		fmt.Println(CountExecutionTime(func() {
			ProcessSimultaneously(data, Pow2, threads)
		}, fmt.Sprintf("Simultaneously processing %d elements with %d threads", len(data), threads)))
	case "vis":
		if len(os.Args) < 3 {
			fmt.Println("usage: seq <filename>")
			return
		}
		filename := os.Args[2]
		fmt.Printf("Starting visualizing file %s\n", filename)
		if err := Visualize(filename); err != nil {
			fmt.Printf("Error during visualizing file %s: %v\n", filename, err)
			return
		}
		fmt.Printf("Successfully visualized file %s\n", filename)
	default:
		fmt.Println("Unknown command")
		defaults()
	}
}
