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
	tree := Tree(
		File("file.txt", "Lorem ipsum"),
		Dir("sub_dir", 
			EmptyFile("other_file.txt"),
		),
		Dir("empty_dir")
	)

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
