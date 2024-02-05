package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "scans a specific directory and commits info to nostalgia database",
		Long:  ``,
		Run:   scan,
	}
)

type Scan struct {
	files       []FileItem
	directories []DirItem
	root        *DirNode
}

type FileItem struct {
	name string
}

type DirItem struct {
	name string
}

type DirNode struct {
	info  *DirItem
	files []FileItem
	leafs []*DirNode
}

func (dn *DirNode) AddLeaf(l *DirNode) {
	dn.leafs = append(dn.leafs, l)
}

func (dn *DirNode) AddFile(f FileItem) {
	dn.files = append(dn.files, f)
}

func (s *Scan) AddFile(f FileItem) {
	s.files = append(s.files, f)
}

func (s *Scan) AddDirectory(d DirItem) {
	s.directories = append(s.directories, d)
}

func read(p string) *Scan {
	d := DirItem{name: p}
	r := DirNode{info: &d}
	s := &Scan{root: &r}
	s.directories = append(s.directories, d)

	var dirs []*DirNode = []*DirNode{&r}
	for i := 0; i < len(dirs); i++ {
		files, err := os.ReadDir(dirs[i].info.name)
		if err != nil {
			fmt.Printf("error reading directory path: [%v] - %v", dirs[i], err)
		}

		for _, dirEntry := range files {
			// exclude hidden files
			if dirEntry.Name()[0] == '.' {
				continue
			}

			fp := path.Join(p, dirEntry.Name())
			if dirEntry.IsDir() {
				fmt.Printf("dir: '%v'\n", fp)
				d := DirItem{name: fp}
				l := DirNode{info: &d}
				dirs[i].AddLeaf(&l)
				dirs = append(dirs, &l)
				s.AddDirectory(d)
			} else {
				f := FileItem{name: fp}
				dirs[i].AddFile(f)
				s.AddFile(f)
				fmt.Printf("file: '%v'\n", fp)
			}
		}
	}

	return s
}

func printDirNode(c *DirNode, p *DirNode) {
	if p != nil {
		fmt.Printf("parent dir: '%v'\n", p.info.name)
	}

	fmt.Printf("current dir: '%v'\n", c.info.name)
	fmt.Printf("files:\n")
	for _, f := range c.files {
		fmt.Printf("	'%v'\n", f.name)
	}
	for _, l := range c.leafs {
		printDirNode(l, c)
	}
}

func printScan(s *Scan) {
	fmt.Printf("total files: %v\n", len(s.files))
	fmt.Printf("total directories: %v\n", len(s.directories))

	fmt.Printf("root dir: '%v'\n", s.root.info.name)
	printDirNode(s.root, nil)
}

func scan(ccmd *cobra.Command, args []string) {
	sp := viper.GetString("scan_path")
	fmt.Printf("hello there nostalgia config: %v, tags: %v\n", sp, tags)

	s := read(sp)

	printScan(s)
}
