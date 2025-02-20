package multitouch

import (
	"errors"
	"fmt"

	afero "github.com/spf13/afero"
)

type optionDict struct {
	fs afero.Fs
}

type Option func(options *optionDict) error

// TreeNode represents a file or directory
type TreeNode struct {
	// the name of the file or directory
	Name string
	// optional content of the file, ignored if this is a directory
	Content string
	// if nil, this is a file
	// if not nil, this is a directory
	// and the children are the files and directories
	Children []TreeNode
}

// Touch creates a directory structure based on the given tree inside a temporary directory
func Touch(tree []TreeNode, opts ...Option) (afero.Fs, error) {
	options, err := calculateOptions(opts)
	if err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	if len(tree) == 0 {
		return nil, errors.New("tree must have at least one element")
	}

	err = createTree(options.fs, tree)
	if err != nil {
		return nil, fmt.Errorf("failed to create tree: %w", err)
	}

	return options.fs, nil
}

// WithFileSystem sets a custom filesystem to use
func WithFileSystem(fs afero.Fs) Option {
	return func(options *optionDict) error {
		options.fs = fs
		return nil
	}
}

// WithBasePath sets the base path for the filesystem
func WithBasePath(path string) Option {
	return func(options *optionDict) error {
		dir, err := options.fs.Stat(path)
		if err != nil {
			return fmt.Errorf("base path must exists and be accessible: %w", err)
		}
		if !dir.IsDir() {
			return errors.New("base path must be a directory")
		}

		options.fs = afero.NewBasePathFs(options.fs, path)
		return nil
	}
}

func calculateOptions(opts []Option) (optionDict, error) {
	options := optionDict{
		fs: afero.NewMemMapFs(),
	}
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return options, err
		}
	}
	return options, nil
}

func createTree(dest afero.Fs, tree []TreeNode) error {
	for _, child := range tree {
		err := createSingle(dest, child)
		if err != nil {
			return err
		}
	}
	return nil
}

func createSingle(dest afero.Fs, tree TreeNode) error {
	if tree.Children == nil {
		err := createFile(dest, tree.Name, tree.Content)
		if err != nil {
			return err
		}
	} else {
		// create a directory
		newDest := afero.NewBasePathFs(dest, tree.Name)
		err := dest.Mkdir(tree.Name, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", tree.Name, err)
		}

		// recursively create the children
		for _, child := range tree.Children {
			err := createSingle(newDest, child)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createFile(dest afero.Fs, name string, content string) error {
	file, err := dest.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", name, err)
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write on file %s: %w", name, err)
	}
	return nil
}
