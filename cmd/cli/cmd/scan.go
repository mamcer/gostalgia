package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mamcer/nostalgia/internal/pkg/hash"
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
	ID                int64     // scan id
	DateCreated       time.Time // scan creation date
	Duration          int64     // scan duration (in milliseconds)
	FileCount         int64     // file scan count
	DirectoryCount    int64     // directory scan count
	FileRepeatedCount int64     // file repeated scan count
	Status            int64     // scan status = done, inprogress, error

	files       []*FileItem
	directories []*DirItem
	root        *DirNode
}

type FileItem struct {
	ID           int64     // file id
	Name         string    // file name
	Extension    string    //file extension
	Path         string    // file path
	DateModified time.Time // file date modified
	Size         int64     // file size (in bytes)
	Hash         string    // file hash
}

type DirItem struct {
	ID             int64     // directory id
	Name           string    // directory name
	Path           string    // directory path
	DateModified   time.Time // date modified
	Size           int64     // directory size (in bytes)
	Exists         bool
	ExistingFileID int64
}

type DirNode struct {
	info  *DirItem
	files []*FileItem
	leafs []*DirNode
}

func (dn *DirNode) AddLeaf(l *DirNode) {
	dn.leafs = append(dn.leafs, l)
}

func (dn *DirNode) AddFile(f *FileItem) {
	dn.files = append(dn.files, f)
}

func (s *Scan) AddFile(f *FileItem) {
	s.files = append(s.files, f)
}

func (s *Scan) AddDirectory(d *DirItem) {
	s.directories = append(s.directories, d)
}

func read(p string) *Scan {
	d := &DirItem{
		Name: "root",
		Path: p,
	}
	r := &DirNode{info: d}
	s := &Scan{root: r}
	s.AddDirectory(d)

	var dirs []*DirNode = []*DirNode{r}
	for i := 0; i < len(dirs); i++ {
		files, err := os.ReadDir(dirs[i].info.Path)
		if err != nil {
			fmt.Printf("error reading directory path: [%v] - %v", dirs[i], err)
		}

		for _, dirEntry := range files {
			// exclude hidden files
			if dirEntry.Name()[0] == '.' {
				continue
			}

			fp := path.Join(dirs[i].info.Path, dirEntry.Name())
			fi, _ := dirEntry.Info()
			if dirEntry.IsDir() {
				d := &DirItem{
					Name:         dirEntry.Name(),
					Path:         fp,
					Size:         fi.Size(),
					DateModified: fi.ModTime(),
				}

				l := &DirNode{info: d}
				dirs[i].AddLeaf(l)
				s.AddDirectory(d)
				dirs = append(dirs, l)
			} else {
				f := &FileItem{
					Name:         dirEntry.Name(),
					Extension:    strings.Trim(filepath.Ext(dirEntry.Name()), "."),
					Path:         fp,
					DateModified: fi.ModTime(),
					Size:         fi.Size(),
				}
				dirs[i].AddFile(f)
				s.AddFile(f)
			}
		}
	}

	return s
}

func printDirNode(c *DirNode, p *DirNode) {
	fmt.Printf("current dir: '%v' size: %v\n", c.info.Name, sizeString(c.info.Size))
	fmt.Printf("files:\n")
	for _, f := range c.files {
		fmt.Printf("	'%v', hash: %v, size: %v\n", f.Name, f.Hash, sizeString(f.Size))
	}
	for _, l := range c.leafs {
		printDirNode(l, c)
	}
}

func printScan(s *Scan) {
	fmt.Printf("\ntotal files: %v\n", len(s.files))
	fmt.Printf("total directories: %v\n", len(s.directories))

	printDirNode(s.root, nil)
}

func hashFiles(s *Scan) *Scan {
	for _, f := range s.files {
		f.Hash, _ = hash.Calculate(f.Path)
	}

	return s
}

func sizeString(v int64) string {
	r := float64(v)
	u := 1000.0
	if v > int64(u) {
		r = r / u
		if r > u {
			r = r / u
			if r > u {
				r = r / u
				return fmt.Sprintf("%v GB", strconv.FormatFloat(r, 'f', 1, 64))
			} else {
				return fmt.Sprintf("%v MB", strconv.FormatFloat(r, 'f', 1, 64))
			}
		}
		return fmt.Sprintf("%v kB", strconv.FormatFloat(r, 'f', 1, 64))
	}

	return fmt.Sprintf("%v Bytes", strconv.FormatFloat(r, 'f', 1, 64))
}

func calculateSize(d *DirNode) int64 {
	var s int64 = 0
	for _, f := range d.files {
		s += f.Size
	}
	for _, l := range d.leafs {
		s += calculateSize(l)
	}

	d.info.Size = s
	return s
}

func size(s *Scan) *Scan {
	calculateSize(s.root)

	return s
}

func scan(ccmd *cobra.Command, args []string) {
	start := time.Now()

	sp := viper.GetString("scan_path")
	fmt.Printf("config: %v, tags: %v\n", sp, strings.Split(tags, ","))

	fmt.Printf("scan process started\n")

	// 	- read-directories
	// - file-structure (directories, files)
	// - hash
	// - update file size
	// - check-existing (exists, existing-id)
	// - persist
	// - copy-files

	// read
	fmt.Printf("\nreading directories...")
	s := read(sp)
	fmt.Printf("OK\n")
	fmt.Printf("directories: %v, files: %v\n", len(s.directories), len(s.files))

	// hash
	fmt.Printf("\nhashing files...")
	s = hashFiles(s)
	fmt.Printf("OK\n")

	// file size
	fmt.Printf("\nupdating file size...")
	s = size(s)
	fmt.Printf("OK\n")

	printScan(s)

	elapsed := time.Since(start)
	fmt.Printf("scan process finished: %v\n", elapsed)

}
