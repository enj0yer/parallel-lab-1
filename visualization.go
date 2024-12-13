package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Result struct {
	Elements int
	Time     float64
}

type Parser struct {
	SequentialResults   map[int]float64
	SimultaneousResults map[int]map[int]float64
}

func newParser() *Parser {
	return &Parser{
		SequentialResults:   make(map[int]float64),
		SimultaneousResults: make(map[int]map[int]float64),
	}
}

func (p *Parser) parse(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	reSequential := regexp.MustCompile(`Sequentially processing (\d+) elements execution took ([\d.]+)(µs|ms)`)
	reSimultaneous := regexp.MustCompile(`Simultaneously processing (\d+) elements with (\d+) threads execution took ([\d.]+)(µs|ms)`)

	for scanner.Scan() {
		line := scanner.Text()
		if matches := reSequential.FindStringSubmatch(line); matches != nil {
			elements, _ := strconv.Atoi(matches[1])
			time, _ := strconv.ParseFloat(matches[2], 64)
			if matches[3] == "ms" {
				time *= 1000
			}
			p.SequentialResults[elements] = time
		} else if matches := reSimultaneous.FindStringSubmatch(line); matches != nil {
			elements, _ := strconv.Atoi(matches[1])
			threads, _ := strconv.Atoi(matches[2])
			time, _ := strconv.ParseFloat(matches[3], 64)
			if matches[4] == "ms" {
				time *= 1000
			}
			if p.SimultaneousResults[elements] == nil {
				p.SimultaneousResults[elements] = make(map[int]float64)
			}
			p.SimultaneousResults[elements][threads] = time
		}
	}

	return scanner.Err()
}

func (p *Parser) generateSequentialResults() error {
	pts := make(plotter.XYs, len(p.SequentialResults))
	i := 0
	for elements, time := range p.SequentialResults {
		pts[i].X = float64(elements)
		pts[i].Y = time
		i++
	}

	p1 := plot.New()

	p1.Title.Text = "Sequential Processing Time"
	p1.X.Label.Text = "Number of Elements"
	p1.Y.Label.Text = "Time (µs)"

	line, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}
	p1.Add(line)

	if err := p1.Save(4*vg.Inch, 4*vg.Inch, "./visualizations/sequential_results.png"); err != nil {
		return err
	}

	return nil
}

func (p *Parser) generateSimultaneousResults() error {
	for elements, threadsMap := range p.SimultaneousResults {
		pts := make(plotter.XYs, len(threadsMap))
		i := 0
		for threads, time := range threadsMap {
			pts[i].X = float64(threads)
			pts[i].Y = time
			i++
		}

		p1 := plot.New()

		p1.Title.Text = "Simultaneous Processing Time (Elements: " + strconv.Itoa(elements) + ")"
		p1.X.Label.Text = "Number of Threads"
		p1.Y.Label.Text = "Time (µs)"

		line, err := plotter.NewLine(pts)
		if err != nil {
			return err
		}
		p1.Add(line)

		if err := p1.Save(4*vg.Inch, 4*vg.Inch, "./visualizations/simultaneous_results_"+strconv.Itoa(elements)+".png"); err != nil {
			return err
		}
	}

	return nil
}

func Visualize(filename string) error {
	parser := newParser()
	if err := parser.parse(filename); err != nil {
		return fmt.Errorf("error while parsing file: %w", err)
	}
	if err := parser.generateSequentialResults(); err != nil {
		fmt.Printf("error while generating sequential file: %v\n", err)
	}
	if err := parser.generateSimultaneousResults(); err != nil {
		fmt.Printf("error while generating simultaneous results: %v", err)
	}
	return nil
}
