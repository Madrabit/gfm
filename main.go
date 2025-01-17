package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	if err := execCommand(); err != nil {
		log.Fatal(err)
	}
}

func execCommand() error {
	command, err := getArgs(0)
	if err != nil {
		return err
	}
	switch command {
	case "pwd":
		return pwd()
	case "ls":
		return ls()
	case "touch":
		fileName, err := getArgs(1)
		if err != nil {
			return errors.New("no file name provided")
		}
		return touch(fileName)
	case "rm":
		delTarget, err := getArgs(1)
		if err != nil {
			return errors.New("no file name provided")
		}
		return rm(delTarget)
	case "mkdir":
		dirName, err := getArgs(1)
		if err != nil {
			return errors.New("no dir name provided")
		}
		return mkdir(dirName)
	case "echo":
		text, err := getArgs(1)
		if err != nil {
			return err
		}
		arrow, err := getArgs(2)
		if err != nil {
			return err
		}
		if arrow == ">" {
			file, err := getArgs(2)
			if err != nil {
				return err
			}
			return writeToFile(text, file)

		}
	case "cat":
		fileName, err := getArgs(1)
		if err != nil {
			return err
		}
		return readFromFile(fileName)
	case "cp":
		src, err := getArgs(1)
		if err != nil {
			return err
		}
		dst, err := getArgs(2)
		if err != nil {
			return err
		}
		return copyFile(src, dst)
	case "mv":
		src, err := getArgs(1)
		if err != nil {
			return err
		}
		dst, err := getArgs(2)
		if err != nil {
			return err
		}
		return move(src, dst)
	case "wc":
		src, err := getArgs(1)
		if err != nil {
			return err
		}
		return wordsCounter(src)
	default:
		return errors.New("unknown command")
	}
	return nil
}

func pwd() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	} else {
		fmt.Println(dir)
	}
	return nil
}

func ls() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return errors.New("cannot get current directory")
	}
	folder, err := os.Open(currentDir)
	if err != nil {
		return errors.New("cannot open directory")
	}
	defer folder.Close()
	files, err := folder.Readdir(-1)
	if err != nil {
		return errors.New("cannot read directory")
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}
	return nil
}

func touch(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return errors.New("error creat file")
	}
	defer file.Close()
	fmt.Println(fileName)
	return nil
}

func rm(fileName string) error {
	if err := os.RemoveAll(fileName); err != nil {
		return errors.New("error remove file")
	}
	return nil
}

func mkdir(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return errors.New("error create directory file")
	}
	return nil
}

func writeToFile(str string, fileName string) error {
	err := os.WriteFile(fileName, []byte(str), os.ModePerm)
	if err != nil {
		return errors.New("error write to file")
	}
	return nil
}

func readFromFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return errors.New("cant open file")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		return errors.New("error stdout data")
	}
	return nil
}

func copyFile(src string, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return errors.New("can't read file")
	}
	err = os.WriteFile(dst, data, os.ModePerm)
	if err != nil {
		return errors.New("can't write to file")
	}
	return nil
}

func move(str string, path string) error {
	dst := path + "\\" + str
	err := os.Rename(str, dst)
	if err != nil {
		return errors.New("cant move file")
	}
	return nil
}

func getArgs(argIdx int) (string, error) {
	arg := flag.Arg(argIdx)
	if arg == "" {
		return "", errors.New("no arg entered")
	}
	return arg, nil
}

func wordsCounter(fileName string) error {
	bytes, err := os.Stat(fileName)
	if err != nil {
		return errors.New("cant read byte size from file")
	}
	size := bytes.Size()
	file, err := os.Open(fileName)
	if err != nil {
		return errors.New("cant open file")
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var linesCounter int
	var wordsCounter int
	for {
		line, err := reader.ReadString('\n')
		wordsCounter += len(strings.Fields(line))
		if err != nil {
			// Если ошибка EOF, завершаем цикл
			if err.Error() != "EOF" {
				return err
			}
			// Если есть последняя строка, увеличиваем счетчик
			if len(line) > 0 {
				linesCounter++
				wordsCounter += len(strings.Fields(line))
			}
			break
		}
		linesCounter++

	}
	fmt.Printf("%d %d %d", size, linesCounter, wordsCounter)
	return nil
}
