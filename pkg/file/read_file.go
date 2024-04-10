package file

import (
	"bufio"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

func Read_line(ch chan<- string, wg *sync.WaitGroup, file_path *string) {
	defer wg.Done()
	defer close(ch)

	file, err := os.Open(*file_path)
	if err != nil {
		logrus.Fatalln("Can not open file {}, err: {}", file_path, err)
	}
	defer file.Close()

	// read line one by one
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ch <- line
		logrus.Debugln(line)
	}

	// check if scanner have err
	if err := scanner.Err(); err != nil {
		logrus.Fatalln("scan file err: {}", err)
		return
	}
}
