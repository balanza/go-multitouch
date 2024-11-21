# go-multitouch

go-multitouch is a Go package that allows you to scaffold a file structure from a descriptor. It is similar to executing `mkdir` and `touch` commands recursively. This package is primarily designed for testing file access packages.

## Installation

```
go get -u github.com/balanza/go-multitouch@latest
```

## Usage

```go
import (
    mtouch "github.com/balanza/go-multitouch"
)

func main() {
    tempDir, err := mtouch.Create([]mtouch.FileTree{
		{
			Name: "file.txt",
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
			Children: []mtouch.FileTree{},
		},
	})
}
```
