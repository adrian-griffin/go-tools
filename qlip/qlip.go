package main

import (
	"fmt"
	"os"
	"strings"
	"flag"
	"path/filepath"
	"os/user"
)

func main() {
	// Determine user and generate homedir path for storing .qlips.txt
	usr, _ := user.Current()
	qlipsTxtDirectory := usr.HomeDir

	// Qlip storage file location
	qlipsPath := filepath.Join(qlipsTxtDirectory, ".qlips.txt")

	// Available flags
	add := flag.String("a", "", "Add new list item")
	list := flag.Bool("l", false, "List items")
	del := flag.Int("d", -1, "Delete item by index")

	flag.Parse()

	// Cases for each flag
	switch{
		case *add != "":
			addListItem(*add, qlipsPath)
		case *list:
			listItems(qlipsPath)
		case *del >= 0:
			delListItem(*del, qlipsPath)
		default:
			fmt.Println("No flag passed. Please use -a, -l, or -d.")
	}
}

func addListItem(item, qlipsPath string) {
	// Opens file 'qlipsPath' in APPEND mode, CREATES file if it does not exist, in WRITEONLY mode, and with the file permission 0644
	//// which allows the owner to read & write, and for others to read only
	f, err := os.OpenFile(qlipsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	// Defers function close to always end with the file access close, to ensure proper resource cleanup
	defer f.Close()
	if _, err := f.WriteString(item + "\n"); err != nil{
		fmt.Println("Error writing to file:", err)
	}
	fmt.Println("Added:", item)
}

func listItems(qlipsPath string) {
	// Opens file for reading
	data, err := os.ReadFile(qlipsPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Splits qlips file into array of string items and iterates over them to print the line to console
	items := strings.Split(string(data), "\n")

	// Looping over range "items" for each item with index "i"
	for i, item := range items {
		if item != "" {
			fmt.Printf("%d: %s\n", i, item)
		}
	}
}

func delListItem(index int, qlipsPath string) {
	// Opens file for storage as 'data', writes err to console
	data, err := os.ReadFile(qlipsPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	// Splits data into array of strings at newline delim
	items := strings.Split(string(data), "\n")

	// Checks for valid index value, reports error if invalid
	if index < 0 || index >= len(items) {
		fmt.Println("Invalid index, no such item")
		return
	}

	// Deleting item from list, creates slice of all items before index, and another slice with all items after index, then appends them together
	items = append(items[:index], items[index+1:]...)

	// Overwrites 'data' with joined string of remaining 'items' array entries
	data = []byte(strings.Join(items, "\n"))

	// Writes data back to file for persistence
	if err := os.WriteFile(qlipsPath, data, 0644); err != nil {
		fmt.Println("Error writing to file:", err)
	}

	// Print to console deletion confirmation
	fmt.Println("Deleted item at index", index)

}
