package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type Names struct {
	names map[string]*Summary
}

type Summary struct {
	Quantity int
	Paths    []string
	Size     int64
}

func (s *Summary) String() string {
	return fmt.Sprintf("%v %v %v", s.Quantity, s.Size, s.Paths)
}

func (s *Summary) AddFile(path string, info os.FileInfo) {
	s.Quantity += 1
	s.Paths = append(s.Paths, path)
	s.Size += info.Size()
}

type ByQuantity []Summary

func (a ByQuantity) Len() int           { return len(a) }
func (a ByQuantity) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByQuantity) Less(i, j int) bool { return a[i].Quantity < a[j].Quantity }

type BySize []Summary

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].Size < a[j].Size }

func (n *Names) Do(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	if n.names[info.Name()] == nil {
		n.names[info.Name()] = &Summary{}
	}
	n.names[info.Name()].AddFile(path, info)
	return nil
}

func (n *Names) Summary() {
	summary := make([]Summary, len(n.names))
	i := 0
	for _, v := range n.names {
		summary[i] = *v
		i++
	}
	sort.Sort(BySize(summary))
	for _, s := range summary {
		if s.Quantity > 1 {
			fmt.Println(s)
		}
	}
}

type MDSum struct {
	names map[string]*Summary
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
	if n.names[sum] == nil {
		n.names[sum] = &Summary{}
	}
	n.names[sum].AddFile(path, info)
	if n.names[sum].Quantity > 1 {
		fmt.Print("!")
	}
	fmt.Print(".")
	return nil
}

func (n *MDSum) Summary() {
	summary := make([]Summary, len(n.names))
	i := 0
	for _, v := range n.names {
		summary[i] = *v
		i++
	}
	sort.Sort(BySize(summary))
	for _, s := range summary {
		if s.Quantity > 1 {
			fmt.Println(s)
		}
	}
}

func main() {
	flag.Parse()
	dir := flag.Arg(0)
	op := MDSum{map[string]*Summary{}}
	//op := Names{map[string]*Summary{}}
	err := filepath.Walk(dir, op.Do)
	if err != nil {
		panic(err)
	}
	op.Summary()
}
