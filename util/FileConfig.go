package util

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type FileConfig struct {
	value map[string]interface{}
}

func (p *FileConfig) LoadFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	p.value = make(map[string]interface{})

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			index := strings.Index(line, "=")
			key := line[:index]
			value := line[index+1:]

			if value[0] == '"' {
				value = value[1:]
			}
			if value[len(value)-1] == '"' {
				value = value[:len(value)-1]
			}

			p.value[key] = value
		}
	}
}

func (p *FileConfig) Get(key string, defaultValue interface{}) interface{} {
	ret, exist := p.value[key]
	if exist == false {
		return defaultValue
	}
	return ret
}
