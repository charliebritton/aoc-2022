package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// type File struct {
// 	Name string
// 	Size int
// }

// type Directory struct {
// 	Name  string
// 	Files []File
// 	Dirs  []Directory
// }

// var root = Directory{
// 	"root",
// 	[]File{},
// 	[]Directory{},
// }

// var current = []Directory{root}

// func changeDir(newDir string) {

// 	//fmt.Printf("Current directory: %v", current)

// 	if newDir == ".." {

// 		// Store current directory to the next level up in the tree
// 		current[len(current)-2].Dirs = append(current[len(current)-1].Dirs, current[len(current)-1])

// 		// Set the variable to one level higher
// 		current = current[:len(current)-1]

// 	} else {

// 		found := false

// 		// Loop through the current directory
// 		// To try and find the new dir, if exists, then append
// 		for _, d := range current[len(current)-1].Dirs {

// 			// Dir exists
// 			if d.Name == newDir {

// 				current = append(current, d)
// 				found = true
// 				break

// 			}

// 		}

// 		if !found {

// 			new := Directory{newDir, []File{}, []Directory{}}

// 			current[len(current)-1].Dirs = append(current[len(current)-1].Dirs, new)

// 			current = append(current)

// 		}

// 	}

// }

type FileType int

const (
	Directory = iota + 1
	File
)

const (
	Bound      = 100000
	TotalSpace = 70000000
	FreeSpace  = 30000000
)

type Node struct {
	parent   *Node
	name     string
	fileType FileType
	fileSize int
	children []*Node
}

func main() {

	f, err := os.OpenFile("../data/day7a", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	root := Node{
		nil, "root", Directory, 0, []*Node{},
	}

	current = &root

	ignoreFirstFlag := false

	for fileScanner.Scan() {

		line := fileScanner.Text()

		if !ignoreFirstFlag {
			ignoreFirstFlag = true
			continue
		}

		tokens := strings.Split(line, " ")

		// Command
		if line[0] == '$' {

			switch tokens[1] {
			case "cd":
				changeDir(tokens[2])
				// case "ls":
				// 	fmt.Println("Listing Directory")
			}

			// Buffer output of ls
		} else {

			if tokens[0] == "dir" {
				// //fmt.Printf("Making directory by changing into and leaving...\n")
				changeDir(tokens[1])
				changeDir("..")
			} else {
				updateFileSize(tokens[0], tokens[1])
			}

		}

		//fmt.Printf("Current is %s\n", printDirList(*current))
		// printTree(root, 0)

	}

	printTree(root, 0)

	fmt.Printf("\nTOTAL BELOW BOUND: %d\n", walkTree(root, 100000))

	fmt.Printf("USED:  %d\nAVAIL: %d\nNEED:  %d\n", root.fileSize, TotalSpace-root.fileSize, FreeSpace-(TotalSpace-root.fileSize))

	need := FreeSpace - (TotalSpace - root.fileSize)
	// Get all directory sizes into a list and find the smallest that is > than NEED

	fmt.Printf("Delete dir with size %d\n", smallestDirToDelete(root, need))

}

var current *Node

// Change Directory, updates the current
func changeDir(name string) {

	//fmt.Printf("CHANGEDIR %s -> %s\n", printDirList(*current), printDirList(*current)+name)

	// Up the tree
	if name == ".." {
		current = current.parent
		return
	}

	for _, c := range current.children {

		if c.name == name {
			current = c
			return
		}

	}

	newDir := Node{
		current,
		name,
		Directory,
		0,
		[]*Node{},
	}

	current.children = append(current.children, &newDir)
	current = &newDir

}

// Prints a listing based on the current directory
func printDirList(dir Node) string {

	if dir.parent == nil {

		return "/"

	} else {

		if dir.fileType == Directory {
			return printDirList(*dir.parent) + dir.name + "/"
		} else {
			return printDirList(*dir.parent) + dir.name
		}

	}

}

// Updates a filesize for a given file, which may be new
func updateFileSize(size string, fileName string) {

	sizeInt, _ := strconv.Atoi(size)

	fileExists := false

	for _, f := range current.children {

		if f.fileType == File && f.name == fileName {

			fileExists = true

			//fmt.Printf("CHSIZE %s: %d -> %d", f.name, f.fileSize, sizeInt)
			f.fileSize = sizeInt

		}
	}

	// Create new file
	if !fileExists {

		//fmt.Printf("CREATE %s with size %d\n", printDirList(*current)+fileName, sizeInt)
		current.children = append(current.children, &Node{current, fileName, File, sizeInt, []*Node{}})

	}

	// Current is the pointer to the directory
	updateSizes(current)

}

func updateSizes(node *Node) {

	// Propogate file size change up tree, just visit each parent and sum all directories
	totalSize := 0

	for _, n := range node.children {

		//fmt.Printf("+ %d (%s) ", n.fileSize, n.name)
		totalSize += n.fileSize

	}
	//fmt.Printf("= %d\n", totalSize)

	node.fileSize = totalSize
	//fmt.Printf("PRSIZE %s: %d\n", printDirList(*node), node.fileSize)

	if node.parent != nil {
		//fmt.Printf("Found a parent, continuing up tree\n")
		updateSizes(node.parent)
	}

}

func printFile(n Node) string {

	if n.fileType == File {
		return fmt.Sprintf("- %s (file, size=%d)\n", n.name, n.fileSize)
	} else {
		return fmt.Sprintf("- %s (dir, size=%d)\n", n.name, n.fileSize)
	}
}

func printTree(n Node, depth int) {

	fmt.Printf(strings.Repeat("  ", depth) + printFile(n))
	for _, child := range n.children {
		printTree(*child, depth+1)
	}

}

func walkTree(n Node, bound int) int {

	total := 0

	for _, c := range n.children {

		// Skip files
		if c.fileType != Directory {
			continue
		}

		if c.fileSize < bound {
			total += c.fileSize
		}

		total += walkTree(*c, bound)

	}
	return total

}

func walkTree2(n Node) []int {

	rtnArr := []int{}

	for _, c := range n.children {

		// Skip files
		if c.fileType != Directory {
			continue
		}

		rtnArr = append(rtnArr, c.fileSize)

		rtnArr = append(rtnArr, walkTree2(*c)...)

	}
	return rtnArr

}

func smallestDirToDelete(n Node, need int) int {

	allDirSizes := walkTree2(n)
	sort.Ints(allDirSizes)

	for _, v := range allDirSizes {

		if v > need {
			return v
		}

	}
	return 0
}
