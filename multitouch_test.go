package multitouch

import (
	"os"
	"testing"
)

func TestErrorOnEmpty(t *testing.T) {
	_, err := Create([]FileTree{})
	if err == nil {
		t.Fatalf("Create([]FileTree{}) should return an error")
	}
}

func TestCreateSingleFileInTempDir(t *testing.T) {
	tree := FileTree{
		Name: "file.txt",
	}

	tempDir, err := Create([]FileTree{tree})
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	if _, err := os.Stat(tempDir + "/file.txt"); err != nil {
		t.Fatalf("File should exist")
	}
}

func TestErrorOnDirNotExist(t *testing.T) {
	tree := FileTree{
		Name: "file.txt",
	}
	_, err := CreateInto("fake_dir", []FileTree{tree})

	if err == nil {
		t.Fatalf("Should fail on non-existing directory")
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
					Name: "file2.txt",
				},
			},
		},
	}

	rootDir, err := Create(deepTree)
	if err != nil {
		t.Fatalf("Should not return an error")
	}

	// Check if the file exists
	files := []string{
		rootDir + "/dir1/file1.txt",
		rootDir + "/dir1/dir1_1/file1_1.txt",
		rootDir + "/dir2/file2.txt",
	}

	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("File %s should exist", file)
		}
	}

}
