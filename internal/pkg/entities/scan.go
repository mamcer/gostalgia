package entities

import (
	"io/ioutil"
	"path"
)

type DirectoryNode struct {
	Name  string
	LNode *DirectoryNode
	RNode *DirectoryNode
}

type FileNode struct {
	Name string
	Next *FileNode
}

type Scan struct {
	Name           string
	Duration       int64
	FileCount      int64
	DirectoryCount int64
	RootNode       *DirectoryNode
	Files          *FileNode
}

func (s Scan) Scan(rootPath string) {
	dirPath := ""
	paths := []string{rootPath}
	for i := 0; i < len(paths); i++ {
		dirPath = paths[i]

		files, _ := ioutil.ReadDir(dirPath)

		for _, fi := range files {
			fp := path.Join(dirPath, fi.Name())
			if fi.IsDir() {
				paths = append(paths, fp)
			}
			// else {
			// 	tfc += 1
			// }
		}
	}
}
