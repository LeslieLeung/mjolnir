package srt

import (
	"bufio"
	"bytes"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"os"
	"strconv"
	"strings"
)

type Block struct {
	Sequence int
	Start    string
	End      string
	Text     string
}

type File struct {
	Blocks []Block
}

func Parse(filename string) (*File, error) {
	file, err := os.Open(filename)
	println("Parsing file: ", filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := transform.NewReader(file, unicode.BOMOverride(unicode.UTF8.NewDecoder()))

	scanner := bufio.NewScanner(reader)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := bytes.IndexAny(data, "\n\r"); i >= 0 {
			if data[i] == '\r' && i+1 < len(data) && data[i+1] == '\n' {
				return i + 2, data[0:i], nil
			}
			return i + 1, data[0:i], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	})
	blocks := make([]Block, 0)

	var currentBlock *Block
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" {
			if currentBlock != nil {
				blocks = append(blocks, *currentBlock)
				currentBlock = nil
			}
			continue
		}

		if currentBlock == nil {
			seq, err := strconv.Atoi(line)
			if err != nil {
				return nil, err
			}
			currentBlock = &Block{Sequence: seq}
		} else if currentBlock.Start == "" {
			times := strings.Split(line, " --> ")
			currentBlock.Start = times[0]
			currentBlock.End = times[1]
		} else {
			currentBlock.Text += line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &File{Blocks: blocks}, nil
}
