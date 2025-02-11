package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
)

type HashTable struct {
	buckets [][]string
	size    int
}

func NewHashTable(size int) *HashTable {
	return &HashTable{
		buckets: make([][]string, size),
		size:    size,
	}
}

func (ht *HashTable) Hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32()) % ht.size
}

func (ht *HashTable) Insert(s string) {
	index := ht.Hash(s)
	for _, str := range ht.buckets[index] {
		if str == s {
			return
		}
	}
	ht.buckets[index] = append(ht.buckets[index], s)
}

func (ht *HashTable) Contains(s string) bool {
	index := ht.Hash(s)
	for _, str := range ht.buckets[index] {
		if str == s {
			return true
		}
	}
	return false
}

func LoadFileToSet(filename string, ht *HashTable) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ht.Insert(line)
	}

	return scanner.Err()
}

func CompareSets(file1, file2 string) (bool, error) {
	ht1 := NewHashTable(1000)
	ht2 := NewHashTable(1000)

	if err := LoadFileToSet(file1, ht1); err != nil {
		return false, err
	}
	if err := LoadFileToSet(file2, ht2); err != nil {
		return false, err
	}

	for _, bucket := range ht1.buckets {
		for _, line := range bucket {
			if !ht2.Contains(line) {
				return false, nil
			}
		}
	}

	for _, bucket := range ht2.buckets {
		for _, line := range bucket {
			if !ht1.Contains(line) {
				return false, nil
			}
		}
	}

	return true, nil
}

func main() {
	file1 := "file1.txt"
	file2 := "file2.txt"

	equal, err := CompareSets(file1, file2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if equal {
		fmt.Println("Множини унікальних рядків співпадають.")
	} else {
		fmt.Println("Множини унікальних рядків не співпадають.")
	}
}
