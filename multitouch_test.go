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

	const BASE_DIR string = "/foo"

	tree := FileTree{
		Name: "file.txt",
	}

	rootFs := afero.NewMemMapFs()
	rootFs.Mkdir(BASE_DIR, 0755)

	fs, err := Touch([]FileTree{tree},
		WithFileSystem(rootFs),
		WithBasePath(BASE_DIR),
	)
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	if _, err := fs.Stat("/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
	if _, err := rootFs.Stat(BASE_DIR + "/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
}

func TestWithRealFileSystem(t *testing.T) {
	const BASE_DIR string = "/tmp/foo"

	tree := FileTree{
		Name: "file.txt",
	}

	rootFs := afero.NewOsFs()
	rootFs.MkdirAll(BASE_DIR, 0755)

	fs, err := Touch([]FileTree{tree},
		WithFileSystem(rootFs),
		WithBasePath(BASE_DIR),
	)
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	if _, err := fs.Stat("/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
	if _, err := rootFs.Stat(BASE_DIR + "/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}

	err = rootFs.RemoveAll(BASE_DIR)
	if err != nil {
		t.Fatalf("Failed to cleanup: %s", err)
	}
}

func TestErrorIfBasePathDoesNotExist(t *testing.T) {

	testFs := []afero.Fs{afero.NewMemMapFs(), afero.NewOsFs()}

	for _, rootFs := range testFs {
		t.Run(rootFs.Name(), func(t *testing.T) {
			const BASE_DIR string = "/foo"

			tree := FileTree{
				Name: "file.txt",
			}

			rootFs := afero.NewMemMapFs()

			_, err := Touch([]FileTree{tree},
				WithFileSystem(rootFs),
				WithBasePath(BASE_DIR),
			)
			if err == nil {
				t.Fatalf("Should return an error")
			}
		})
	}

}

func TestNestedDirectories(t *testing.T) {
	const BASE_DIR string = "/foo"

	tree := FileTree{
		Name: "my/nested/dir",
		Children: []FileTree{
			{Name: "file.txt"},
		},
	}

	rootFs := afero.NewMemMapFs()
	rootFs.Mkdir(BASE_DIR, 0755)

	fs, err := Touch([]FileTree{tree},
		WithFileSystem(rootFs),
		WithBasePath(BASE_DIR),
	)
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	if _, err := fs.Stat("my"); err != nil {
		t.Fatalf("Directory my should exist")
	}
	if _, err := fs.Stat("my/nested"); err != nil {
		t.Fatalf("Directory nested should exist")
	}
	if _, err := fs.Stat("my/nested/dir"); err != nil {
		t.Fatalf("Directory dir should exist")
	}
	if _, err := fs.Stat("my/nested/dir/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
}
 