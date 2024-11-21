package multitouch

import (
	"fmt"
	"os"
)

// FileTree represents a file or directory
type FileTree struct {
	// the name of the file or directory
	Name string
	// if nil, this is a file
	// if not nil, this is a directory
	// and the children are the files and directories
	Children []FileTree
}

// Create creates a directory structure based on the given tree inside a temporary directory
func Create(tree []FileTree) (string, error) {
	if len(tree) == 0 {
		return "", fmt.Errorf("tree must have at least one element")
	}

	// create a temporary directory
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		return "", err
	}
	return CreateInto(dir, tree)
}

// CreateInto creates a directory structure based on the given tree inside the given directory
func CreateInto(dest string, tree []FileTree) (string, error) {
	if len(tree) == 0 {
		return "", fmt.Errorf("tree must have at least one element")
	}

	for _, child := range tree {
		err := createSub(dest, child)
		if err != nil {
			return "", err
		}
	}
	return dest, nil
}

func createSub(dest string, tree FileTree) error {
	if tree.Children == nil {
		// create a file
		file, err := os.Create(dest + "/" + tree.Name)
		if err != nil {
			return err
		}
		file.Close()

	} else {

		// create a directory
		newDest := dest + "/" + tree.Name
		err := os.Mkdir(newDest, 0755)
		if err != nil {
			return err
		}

		// recursively create the children
		CreateInto(newDest, tree.Children)
	}
	return nil
}
