package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Names struct {
	names map[string][]string
}

type Summary struct {
	Key      string
	Quantity int
	Paths    []string
	Size     int64
}

func (n *Names) Do(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	n.names[info.Name()] = append(n.names[info.Name()], path)
	return nil
}

func (n *Names) Summary() {
	for k, v := range n.names {
		if len(v) > 1 {
			fmt.Printf("%s: %v %v\n\n", k, len(v), v)
		}
	}
}

type MDSum struct {
	names map[string][]string
}

func (n *MDSum) Do(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	sum := fmt.Sprintf("%x", md5.Sum(data))
	n.names[sum] = append(n.names[sum], path)
	if len(n.names[sum]) > 1 {
		fmt.Print("!")
	}
	fmt.Print(".")
	return nil
}

func (n *MDSum) Summary() {
	for k, v := range n.names {
		if len(v) > 1 {
			fmt.Printf("%s: %v %v\n\n", k, len(v), v)
		}
	}
}

func main() {
	flag.Parse()
	dir := flag.Arg(0)
	op := MDSum{map[string][]string{}}
	//op := Names{map[string][]string{}}
	err := filepath.Walk(dir, op.Do)
	if err != nil {
		panic(err)
	}
	op.Summary()
}
