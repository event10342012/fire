package main

import (
	"fmt"
	"os"
)

type content struct {
	filepath string
	data     string
}

func write(ct content) error {
	file, err := os.Create(ct.filepath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(file)

	file.WriteString(ct.data)
	return nil
}

func read(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Stdout.Write(data)
}

func main() {
	args := os.Args[1:]

	ct := content{
		args[0],
		args[1],
	}
	write(ct)
	read(ct.filepath)
}
