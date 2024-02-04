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
	files       *FileNode
	directories *DirNode
	root        *DirItem
}

type FileNode struct {
	name string
	next *FileNode
}

type DirNode struct {
	name string
	next *DirNode
}

type DirItem struct {
	name  string
	files *FileNode
	leafs *DirItem
}

func read(p string) Scan {
	var dirs []string = []string{p}
	for i := 0; i < len(dirs); i++ {
		files, err := os.ReadDir(dirs[i])
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
				dirs = append(dirs, fp)
			} else {
				fmt.Printf("file: '%v'\n", fp)
			}
		}
	}

	return Scan{}
}

func scan(ccmd *cobra.Command, args []string) {
	sp := viper.GetString("scan_path")
	fmt.Printf("hello there nostalgia config: %v, tags: %v\n", sp, tags)

	read(sp)
}
