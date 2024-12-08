package multitouch

import (
	"fmt"

	afero "github.com/spf13/afero"
)

type options struct {
	fs afero.Fs
}

type Option func(options *options) error

// FileTree represents a file or directory
type FileTree struct {
	// the name of the file or directory
	Name string
	// optional content of the file, ignored if this is a directory
	Content string
	// if nil, this is a file
	// if not nil, this is a directory
	// and the children are the files and directories
	Children []FileTree
}

// Touch creates a directory structure based on the given tree inside a temporary directory
func Touch(tree []FileTree, opts ...Option) (afero.Fs, error) {

	options, err := calculateOptions(opts)
	if err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	if len(tree) == 0 {
		return nil, fmt.Errorf("tree must have at least one element")
	}

	err = createTree(options.fs, tree)
	if err != nil {
		return nil, fmt.Errorf("failed to create tree: %w", err)
	}

	return options.fs, nil
}

// WithFileSystem sets a custom filesystem to use
func WithFileSystem(fs afero.Fs) Option {
	return func(options *options) error {
		options.fs = fs
		return nil
	}
}

// WithBasePath sets the base path for the filesystem
func WithBasePath(path string) Option {
	return func(options *options) error {
		options.fs = afero.NewBasePathFs(options.fs, path)
		return nil
	}
}

func calculateOptions(opts []Option) (options, error) {
	options := options{
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

func createTree(dest afero.Fs, tree []FileTree) error {
	for _, child := range tree {
		err := createSingle(dest, child)
		if err != nil {
			return err
		}
	}
	return nil
}

func createSingle(dest afero.Fs, tree FileTree) error {
	if tree.Children == nil {
		// create a file
		file, err := dest.Create(tree.Name)
		if err != nil {
			return err
		}
		file.WriteString(tree.Content)
		file.Close()

	} else {

		// create a directory
		newDest := afero.NewBasePathFs(dest, tree.Name)
		err := dest.Mkdir(tree.Name, 0755)
		if err != nil {
			return err
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
