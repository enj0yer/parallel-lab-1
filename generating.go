package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func createFile(filename string, nums int) error {
	file, err := os.Create(fmt.Sprintf("./data/%s", filename))
	if err != nil {
		return fmt.Errorf("unable to create file with name %s: %w", filename, err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	for i := 1; i <= nums; i++ {
		_, err := writer.WriteString(fmt.Sprintf("%d\n", i))
		if err != nil {
			return fmt.Errorf("unable to write in file %s: %w", filename, err)
		}
	}
	return writer.Flush()
}

func readFile(filename string) ([]int, error) {
	file, err := os.Open(fmt.Sprintf("./data/%s", filename))
	if err != nil {
		return nil, fmt.Errorf("unable to read from file %s: %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var numbers []int

	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("unable to read number from file %s: %w", filename, err)
		}
		numbers = append(numbers, num)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error %s: %w", filename, err)
	}

	return numbers, nil
}

func generate(filename string, size int) {
	fmt.Printf("Started creating file with %d elements inside\n", size)
	if err := createFile(filename, size); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("File generated successfully")
}
