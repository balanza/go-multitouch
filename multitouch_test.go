package multitouch

import (
	"testing"

	afero "github.com/spf13/afero"
)

func TestErrorOnEmpty(t *testing.T) {
	_, err := Touch([]FileTree{})
	if err == nil {
		t.Fatalf("Create([]FileTree{}) should return an error")
	}
}

func TestCreateSingleFileInTempDir(t *testing.T) {
	tree := FileTree{
		Name: "file.txt",
	}

	fs, err := Touch([]FileTree{tree})
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	if _, err := fs.Stat("file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
}

func TestDeepStructure(t *testing.T) {
	deepTree := []FileTree{
		{
			Name: "dir1",
			Children: []FileTree{
				{
					Name: "file1.txt",
				},
				{
					Name: "dir1_1",
					Children: []FileTree{
						{
							Name: "file1_1.txt",
						},
					},
				},
			},
		},
		{
			Name: "dir2",
			Children: []FileTree{
				{
					Name:    "file2.txt",
					Content: "Hello, World!",
				},
			},
		},
	}

	rootDir, err := Touch(deepTree)
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	files := []string{
		"dir1/file1.txt",
		"dir1/dir1_1/file1_1.txt",
		"dir2/file2.txt",
	}

	for _, file := range files {
		if _, err := rootDir.Stat(file); err != nil {
			t.Fatalf("File %s should exist", file)
		}
	}

	content, err := afero.ReadFile(rootDir, "dir2/file2.txt")
	if err != nil {
		t.Fatalf("Failed to read file")
	}
	if contentStr := string(content); contentStr != "Hello, World!" {
		t.Fatalf("File content is wrong: %s", contentStr)
	}
}

func TestWithBasePath(t *testing.T) {
	tree := FileTree{
		Name: "file.txt",
	}

	rootFs := afero.NewMemMapFs()

	fs, err := Touch([]FileTree{tree},
		WithFileSystem(rootFs),
		WithBasePath("/tmp"),
	)
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	if _, err := fs.Stat("/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
	if _, err := rootFs.Stat("/tmp/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
}

func TestWithRealFileSystem(t *testing.T) {
	tree := FileTree{
		Name: "file.txt",
	}

	rootFs := afero.NewOsFs()

	fs, err := Touch([]FileTree{tree},
		WithFileSystem(rootFs),
		WithBasePath("/tmp"),
	)
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	if _, err := fs.Stat("/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
	if _, err := rootFs.Stat("/tmp/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}

	err = rootFs.RemoveAll("/tmp")
	if err != nil {
		t.Fatalf("Failed to cleanup")
	}
}
