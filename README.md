# go-multitouch

go-multitouch is a Go package that allows you to scaffold a file structure from a descriptor. It is similar to executing `mkdir` and `touch` commands recursively. This package is primarily designed for testing file access packages.

## Installation

```
go get -u github.com/balanza/go-multitouch@latest
```

## Usage

```go
import (
    "github.com/balanza/go-multitouch"
)

func main() {
	tree := []FileTree{
		{
			Name: "file.txt",
			Content: "Lorem ipsum"
		},
		{
			Name: "sub_dir",
			Children: []mtouch.FileTree{
				{
					Name: "other_file.txt",
				},
			},
		},
		{
			Name:     "empty_dir",
			Children: []FileTree{},
		},
	}

	// creates the file tree in a memory file system
    fs, err := Touch(tree)

	// created the file tree in the provided directory of the memory file system
	fs, err := Touch(tree, WithBasePath("/my/path"))

	// created the file tree in the provided directory of the current file system
	fs, err := Touch(tree,
			WithFileSystem(afero.NewOsFs()),
			WithBasePath("/my/path")
		)
}
```
