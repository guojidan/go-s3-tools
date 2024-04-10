package file

import (
	"bufio"
	"log"
	"os"
	"sync"
)

func Read_line(ch chan<- string, wg *sync.WaitGroup, file_path *string) {
	defer wg.Done()
	defer close(ch)

	file, err := os.Open(*file_path)
	if err != nil {
		log.Panicln("Can not open file {}, err: {}", file_path, err)
	}
	defer file.Close()

	// read line one by one
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ch <- line
		log.Println(line)
	}

	// check if scanner have err
	if err := scanner.Err(); err != nil {
		log.Println("scan file err: {}", err)
		return
	}
}
