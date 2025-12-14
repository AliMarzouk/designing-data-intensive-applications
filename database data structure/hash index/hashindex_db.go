package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var KEY_NOT_FOUND = errors.New("Key not found")

func readInt16FromFile(file *os.File) (int16, error) {
	intBytes := make([]byte, 2)
	_, err := file.Read(intBytes)
	if err != nil {
		if err == io.EOF {
			return 0, io.EOF
		}
		panic(err)
	}
	return int16FromBytes(&intBytes), nil
}

func readStringFromFile(file *os.File, length int64) (string, error) {
	strBytes := make([]byte, length)
	_, err := file.Read(strBytes)
	if err != nil {
		if err == io.EOF {
			return "", io.EOF
		}
		panic(err)
	}
	return string(strBytes), nil
}

func db_set(key, value string, index *map[string]int32, currentIndex *int32) error {
	file, err := os.OpenFile("database.data", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	binary.Write(file, binary.LittleEndian, int16(len(key)))
	binary.Write(file, binary.LittleEndian, int16(len(value)))
	file.WriteString(key)
	file.WriteString(value)

	(*index)[key] = *currentIndex
	*currentIndex += int32(4 + len(key) + len(value))

	return nil
}

func db_get(key string, index *map[string]int32) (string, error) {
	position, exists := (*index)[key]
	if !exists {
		return "", KEY_NOT_FOUND
	}

	file, err := os.OpenFile("database.data", os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Seek(int64(position), io.SeekStart)

	keyLen, _ := readInt16FromFile(file)
	valueLen, _ := readInt16FromFile(file)

	// read value and key

	keyStr, _ := readStringFromFile(file, int64(keyLen))
	valueStr, _ := readStringFromFile(file, int64(valueLen))

	if keyStr != key {
		panic("Corrupted memory: expected key :[" + key + "] got [" + keyStr + "]")
	}

	return valueStr, nil
}

func int16FromBytes(input *[]byte) int16 {
	var value int16
	value |= int16((*input)[0])
	value |= int16((*input)[1]) << 8
	return value
}

func restore() (*map[string]int32, int32) {
	index := map[string]int32{}
	var currentIndex int32 = 0

	if _, err := os.Stat("database.data"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("INFO: no initial file found")
		return &index, currentIndex
	}

	file, err := os.OpenFile("database.data", os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for {
		keyLen, err := readInt16FromFile(file)
		if err != nil {
			return &index, currentIndex
		}
		valueLen, err := readInt16FromFile(file)
		if err != nil {
			panic(err)
		}
		keyStr, err := readStringFromFile(file, int64(keyLen))
		if err != nil {
			panic(err)
		}
		_, err = file.Seek(int64(valueLen), io.SeekCurrent)
		if err != nil {
			panic(err)
		}
		index[keyStr] = currentIndex
		currentIndex += int32(4 + keyLen + valueLen)
	}
}

func main() {
	index, currentIndex := restore()

	stdinReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("db command line> ")
		input, readError := stdinReader.ReadString('\n')
		if readError == nil {
			inputParts := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(input), -1)
			switch inputParts[0] {
			case "set":
				if len(inputParts) == 3 {
					db_set(inputParts[1], inputParts[2], index, &currentIndex)
					fmt.Printf("SET %s: %s \n", inputParts[1], inputParts[2])
				} else {
					fmt.Printf("db_set: Wrong number of arguments\n")
				}
			case "get":
				if len(inputParts) == 2 {
					value, err := db_get(inputParts[1], index)
					if err == KEY_NOT_FOUND {
						fmt.Println("KEY NOT FOUND")
					} else {
						fmt.Println("VALUE:", value)
					}
				} else {
					fmt.Printf("db_get: Wrong number of arguments\n")
				}
			default:
				fmt.Printf("command unknown\n")
			}
		}
	}
}
