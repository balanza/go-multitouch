package supertouch

import (
	"testing"

	afero "github.com/spf13/afero"
)

func TestErrorOnEmpty(t *testing.T) {
	_, err := Touch([]TreeNode{})
	if err == nil {
		t.Fatalf("Create([]TreeNode{}) should return an error")
	}
}

func TestCreateSingleFileInTempDir(t *testing.T) {
	tree := Tree(EmptyFile("file.txt"))

	fs, err := Touch(tree)
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	if _, err := fs.Stat("file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
}

func TestDeepStructure(t *testing.T) {
	deepTree := Tree(
		Dir("dir1",
			EmptyFile("file1.txt"),
			Dir("dir1_1",
				EmptyFile("file1_1.txt"))),
		Dir("dir2",
			File("file2.txt", "Hello, World!"),
		),
	)

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
	const baseDir string = "/foo"

	rootFs := afero.NewMemMapFs()
	err := rootFs.Mkdir(baseDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %s", err)
	}

	fs, err := Touch(
		Tree(EmptyFile("file.txt")),
		WithFileSystem(rootFs),
		WithBasePath(baseDir),
	)
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	if _, err := fs.Stat("/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
	if _, err := rootFs.Stat(baseDir + "/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
}

func TestWithRealFileSystem(t *testing.T) {
	const baseDir string = "/tmp/foo"

	rootFs := afero.NewOsFs()
	err := rootFs.MkdirAll(baseDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %s", err)
	}

	fs, err := Touch(
		Tree(EmptyFile("file.txt")),
		WithFileSystem(rootFs),
		WithBasePath(baseDir),
	)
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	if _, err := fs.Stat("/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
	if _, err := rootFs.Stat(baseDir + "/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}

	err = rootFs.RemoveAll(baseDir)
	if err != nil {
		t.Fatalf("Failed to cleanup: %s", err)
	}
}

func TestErrorIfBasePathDoesNotExist(t *testing.T) {
	testFs := []afero.Fs{afero.NewMemMapFs(), afero.NewOsFs()}

	for _, rootFs := range testFs {
		t.Run(rootFs.Name(), func(t *testing.T) {
			const baseDir string = "/foo"

			rootFs := afero.NewMemMapFs()

			_, err := Touch(
				Tree(EmptyFile("file.txt")),
				WithFileSystem(rootFs),
				WithBasePath(baseDir),
			)
			if err == nil {
				t.Fatalf("Should return an error")
			}
		})
	}
}

func TestNestedDirectories(t *testing.T) {
	deepTree := Tree(
		Dir("dir1/dir1_1",
			EmptyFile("file1_1.txt"),
		),
		Dir("/dir2/di2_1",
			EmptyFile("file2_1.txt"),
		),
	)

	rootDir, err := Touch(deepTree)
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	files := []string{
		"dir1/dir1_1/file1_1.txt",
	}

	for _, file := range files {
		if _, err := rootDir.Stat(file); err != nil {
			t.Fatalf("File %s should exist", file)
		}
	}
}
