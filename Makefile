all: main.go processing.go visualization.go generating.go appliers.go
	go build -o main main.go processing.go visualization.go generating.go appliers.go

