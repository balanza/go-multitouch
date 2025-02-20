# Supertouch

Supertouch is a Go package that allows you to scaffold a file structure from a descriptor. It is similar to executing `mkdir` and `touch` commands recursively. This package is primarily designed for testing file access packages.

## Installation

```
go get -u github.com/balanza/supertouch@latest
```

## Usage

```go
import (
  s "github.com/balanza/supertouch"
)

func main() {
  // define the file tree declaratively
  tree := s.Tree(
    s.File("file.txt", "Lorem ipsum"),
    s.Dir("sub_dir",
      s.EmptyFile("other_file.txt"),
    ),
    s.Dir("empty_dir"),
  )

  // creates the file tree in a memory file system
  fs, err := s.Touch(tree)

  // creates the file tree in the provided directory of the memory file system
  fs, err := s.Touch(tree, s.WithBasePath("/my/path"))

  // creates the file tree in the provided directory of the current file system
  fs, err := s.Touch(tree,
    s.WithFileSystem(afero.NewOsFs()),
    s.WithBasePath("/my/path"),
  )
}
```

## Acknowledgement

It turns out there is already [a package](https://github.com/afdezl/supertouch) called Supertouch as well that serves a similar purpose, but in Python. 
It's an involuntary reference but you know, great minds think alike.
