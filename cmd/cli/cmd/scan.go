package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

type ScanStatus int64

const (
	InProgress ScanStatus = 0
	Done       ScanStatus = 1
	Error      ScanStatus = 2
)

type Scan struct {
	ID                        int64      // scan id
	DateCreated               time.Time  // scan creation date
	Duration                  int64      // scan duration (in milliseconds)
	DbFileRepeatedCount       int64      // db file repeated scan count
	SameScanFileRepeatedCount int64      // same scan file repeated scan count
	Status                    ScanStatus // scan status = inprogress, done, error

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
	FileExists   bool      // file exists
}

type DirItem struct {
	ID           int64     // directory id
	Name         string    // directory name
	Path         string    // directory path
	DateModified time.Time // date modified
	Size         int64     // directory size (in bytes)
}

type DirNode struct {
	info  *DirItem
	files []*FileItem
	leafs []*DirNode
}

type NTag struct {
	ID   int64  // ntag id
	Name string // ntag name
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
		Name:         "root",
		Path:         p,
		DateModified: time.Now(),
	}

	r := &DirNode{info: d}
	s := &Scan{root: r, DateCreated: time.Now(), Status: InProgress}
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
	fmt.Printf("current dir: '%v', date: %v, size: %v, parent: %v\n", c.info.Name, c.info.DateModified, sizeString(c.info.Size), p.info.Name)
	fmt.Printf("files:\n")
	for _, f := range c.files {
		fmt.Printf("	'%v', date: %v, hash: %v, size: %v\n", f.Name, f.DateModified, f.Hash, sizeString(f.Size))
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
	for i, f := range s.files {
		fmt.Printf("[%v/%v] hashing: '%v'\n", i+1, len(s.files), f.Name)
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

func checkExisting(s *Scan) *Scan {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/nostalgia")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var id int64
	for _, f := range s.files {
		id = 0
		db.QueryRow("SELECT `id` FROM `nfile` WHERE `hash` = ?", f.Hash).Scan(&id)
		if id != 0 {
			f.ID = id
			f.FileExists = true
			s.DbFileRepeatedCount += 1
		}
	}

	for _, fi := range s.files {
		for _, f := range s.files {
			if fi.Path != f.Path && fi.Hash == f.Hash {
				if s.SameScanFileRepeatedCount == 0 {
					fmt.Println("")
				}

				fmt.Printf("file repeated on same scan: '%v'\n", f.Path)
				s.SameScanFileRepeatedCount += 1
			}
		}
	}

	return s
}

func persistDirNode(dn *DirNode, pid int64, sid int64, db *sql.DB, rp string, tags []NTag) {
	// ndirectory insert
	stmtDirectory, err := db.Prepare("INSERT INTO `ndirectory` (`name`, `path`, `date_modified`, `size`, `file_count`, `directory_count`, `parent_id`) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("error preparing ndirectory insert: %v\n", err)
	}
	defer stmtDirectory.Close()

	// nfile_ndirectory insert
	stmtFileDirectory, err := db.Prepare("INSERT INTO `nfile_ndirectory` (`nfile_id`, `ndirectory_id`, `nscan_id`, `name`) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("error preparing nfile_ndirectory insert: %v\n", err)
	}
	defer stmtFileDirectory.Close()

	// nfile insert
	stmtFile, err := db.Prepare("INSERT INTO `nfile` (`name`, `extension`, `path`, `date_modified`, `size`, `hash`) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("error preparing nfile insert: %v\n", err)
	}
	defer stmtFile.Close()

	// ntag ndirectory insert
	stmtTagDirectory, err := db.Prepare("INSERT INTO `ntag_ndirectory` (`ntag_id`, `ndirectory_id`) VALUES (?, ?)")
	if err != nil {
		fmt.Printf("error preparing ntag_ndirectory insert: %v\n", err)
	}
	defer stmtTagDirectory.Close()

	// ntag nfile insert
	stmtTagFile, err := db.Prepare("INSERT INTO `ntag_nfile` (`ntag_id`, `nfile_id`) VALUES (?, ?)")
	if err != nil {
		fmt.Printf("error preparing ntag_nfile insert: %v\n", err)
	}
	defer stmtTagFile.Close()

	var id int64
	if dn.info.ID == 0 {
		nd := strings.Replace(dn.info.Path, rp, "", 1)
		id = 0
		db.QueryRow("SELECT `id` FROM `ndirectory` WHERE  path = ?", nd).Scan(&id)
		// check if the directory already have the tag
		if id != 0 {
			dn.info.ID = id
		} else {
			res, err := stmtDirectory.Exec(dn.info.Name, nd, dn.info.DateModified, dn.info.Size, len(dn.files), len(dn.leafs), pid)
			if err != nil {
				fmt.Printf("error inserting ndirectory: %v\n", dn)
			}
			dn.info.ID, err = res.LastInsertId()
			if err != nil {
				fmt.Printf("error defining ndirectory last insert id: %v\n", dn)
			}
		}

		for _, t := range tags {
			err = db.QueryRow("SELECT `id` FROM `ntag_ndirectory` WHERE  ntag_id = ? and ndirectory_id = ?", t.ID, dn.info.ID).Scan()
			// check if the directory already have the tag
			if err == sql.ErrNoRows {
				_, err = stmtTagDirectory.Exec(t.ID, dn.info.ID)
				if err != nil {
					fmt.Printf("error inserting ntag_ndirectory, directory_id: %v, tag_id:%v\n", dn.info.ID, t.ID)
				}
			}
		}
	}

	for _, f := range dn.files {
		if !f.FileExists {
			// check if its not repeated in the same scan
			var id int64 = 0
			db.QueryRow("SELECT `id` FROM `nfile` WHERE `hash` = ?", f.Hash).Scan(&id)
			if id != 0 {
				// its repeated but from previous scans, already in the database
				f.ID = id
				f.FileExists = true
			} else {
				res, err := stmtFile.Exec(f.Name, f.Extension, strings.Replace(f.Path, rp, "", 1), f.DateModified, f.Size, f.Hash)
				if err != nil {
					fmt.Printf("error inserting nfile: %v\n", f)
				}
				f.ID, err = res.LastInsertId()
				if err != nil {
					fmt.Printf("error defining nfile last insert id: %v\n", f)
				}
			}
		}

		err = db.QueryRow("SELECT `id` FROM `nfile_ndirectory` WHERE nfile_id = ? and ndirectory_id = ? and name = ?", f.ID, dn.info.ID, f.Name).Scan()
		if err == sql.ErrNoRows {
			_, err = stmtFileDirectory.Exec(f.ID, dn.info.ID, sid, f.Name)
			if err != nil {
				fmt.Printf("error inserting nfile_ndirectory: %v, %v\n", f, dn)
			}
		}

		for _, t := range tags {
			err = db.QueryRow("SELECT `id` FROM `ntag_nfile` WHERE  ntag_id = ? and nfile_id = ?", t.ID, f.ID).Scan()
			// check if the file already have the tag
			if err == sql.ErrNoRows {
				_, err = stmtTagFile.Exec(t.ID, f.ID)
				if err != nil {
					fmt.Printf("error inserting ntag_nfile, file_id: %v, tag_id:%v\n", f.ID, t.ID)
				}
			}
		}
	}

	for _, d := range dn.leafs {
		persistDirNode(d, dn.info.ID, sid, db, rp, tags)
	}

}

func getTags(ntags []string) []NTag {
	tags := []NTag{}

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/nostalgia")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtTag, err := db.Prepare("INSERT INTO `ntag` (`name`) VALUES (?)")
	if err != nil {
		fmt.Printf("error preparing ntag insert: %v\n", err)
	}
	defer stmtTag.Close()

	var id int64
	for _, nt := range ntags {
		id = 0
		tag := NTag{}
		tag.Name = nt
		db.QueryRow("SELECT `id` FROM `ntag` WHERE lower(`name`) = ?", strings.ToLower(nt)).Scan(&id)
		if id != 0 {
			// existing tag
			tag.ID = id
		} else {
			// non existing tag
			res, err := stmtTag.Exec(nt)
			if err != nil {
				fmt.Printf("error inserting ntag: %v\n", nt)
			} else {
				tag.ID, _ = res.LastInsertId()
			}
		}

		tags = append(tags, tag)
	}

	return tags
}

func persist(s *Scan, ntags []string) *Scan {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/nostalgia")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtScan, err := db.Prepare("INSERT INTO `nscan` (`date_created`, `duration`, `file_count`, `directory_count`, `db_file_repeated_count`, `same_scan_file_repeated_count`, `status`, `root_directory_id`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("error preparing nscan insert: %v\n", err)
	}
	defer stmtScan.Close()

	res, err := stmtScan.Exec(s.DateCreated, s.Duration, len(s.files), len(s.directories)-1, s.DbFileRepeatedCount, s.SameScanFileRepeatedCount, s.Status, s.root.info.ID)
	if err != nil {
		fmt.Printf("error inserting nscan: %v\n", err)
	}
	s.ID, _ = res.LastInsertId()

	tags := getTags(ntags)

	persistDirNode(s.root, s.root.info.ID, s.ID, db, s.root.info.Path+"/", tags)

	// nscan update
	stmtUpdateScan, err := db.Prepare("UPDATE `nscan` SET `status` = ? WHERE id = ?")
	if err != nil {
		fmt.Printf("error preparing ndirectory update: %v\n", err)
	}
	defer stmtUpdateScan.Close()

	_, err = stmtUpdateScan.Exec(Done, s.ID)
	if err != nil {
		fmt.Printf("error updating nscan: %v\n", s)
	}

	return s
}

func getSourceDirectories() []DirItem {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/nostalgia")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, name FROM ndirectory WHERE is_source = 1`)
	results := []DirItem{}
	if err == nil {
		for rows.Next() {
			di := DirItem{}
			rows.Scan(&di.ID, &di.Name)
			results = append(results, di)
		}
	}

	return results
}

func copyFiles(s *Scan, np string) *Scan {
	fmt.Println("")
	var fok, ef, sf, fe int64 // file_copy_ok, existing_file, skipped_file, file_with_error
	existing := []string{}
	for i, f := range s.files {
		if !f.FileExists {
			fmt.Printf("[%v/%v] copying file: '%v'", i+1, len(s.files), f.Name)
			rp := strings.Replace(f.Path, s.root.info.Path, "", 1)
			fp := path.Join(np, rp)
			if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
				err := os.MkdirAll(filepath.Dir(fp), 0755)
				if err != nil {
					fmt.Printf(" - failed to create directory: '%v' - %v\n", fp, err)
					fe += 1
				} else {
					i, err := os.Open(f.Path)
					if err == nil {
						defer i.Close()
						o, err := os.Create(fp)
						if err == nil {
							defer o.Close()
							_, err = io.Copy(o, i)
							if err != nil {
								fmt.Printf(" - failed to copy file: '%v' to '%v' - %v\n", f.Name, fp, err)
								fe += 1
							} else {
								fmt.Printf(" - OK\n")
								fok += 1
							}
						} else {
							fmt.Printf(" - failed to create file to copy: %v - %v\n", fp, err)
							fe += 1
						}

					} else {
						fmt.Printf(" - failed to open file to copy: %v - %v\n", f.Path, err)
						fe += 1
					}
				}
			} else {
				fmt.Printf(" - file already exists on disk: '%v'\n", fp)
				existing = append(existing, fp)
				ef += 1
			}
		} else {
			fmt.Printf("[%v/%v] skipping existing file: '%v'\n", i+1, len(s.files), f.Name)
			sf += 1
		}
	}

	fmt.Printf("files copied: %v, with error: %v, skipped: %v, existing on disk: %v\n", fok, fe, sf, ef)
	if len(existing) > 0 {
		fmt.Printf("existing files on disk:\n")
		for _, ex := range existing {
			fmt.Printf("%v\n", ex)
		}
	}

	return s
}

func scan(ccmd *cobra.Command, args []string) {
	start := time.Now()

	sd := getSourceDirectories()

	sp := viper.GetString("scan_path")
	np := viper.GetString("nostalgia_path")
	ntags := strings.Split(tags, ",")
	fmt.Printf("\nscan_path: %v\ntags: %v\nsource: %v\n", sp, strings.Join(ntags, "-"), source)
	fmt.Println("source directories:")
	var sourceID int64
	for _, di := range sd {
		fmt.Printf("id: %v, name: %s\n", di.ID, di.Name)
		if source == di.Name {
			sourceID = di.ID
		}
	}
	fmt.Println("")

	if sourceID == 0 {
		fmt.Printf("invalid source directory: '%v'\n", source)
		return
	}

	fmt.Printf("scan process started\n")

	// read
	partial := time.Now()
	fmt.Printf("\nreading directories...")
	s := read(sp)
	s.root.info.ID = sourceID
	elapsedpartial := time.Since(partial)
	fmt.Printf("OK (%v)\n", elapsedpartial)
	fmt.Printf("directories: %v, files: %v\n", len(s.directories), len(s.files))

	// hash
	partial = time.Now()
	fmt.Printf("\nhashing files...\n")
	s = hashFiles(s)
	elapsedpartial = time.Since(partial)
	fmt.Printf("OK (%v)\n", elapsedpartial)

	// file size
	partial = time.Now()
	fmt.Printf("\nupdating file size...")
	_ = size(s)
	elapsedpartial = time.Since(partial)
	fmt.Printf("OK (%v)\n", elapsedpartial)

	//printScan(s)

	// check existing
	partial = time.Now()
	fmt.Printf("\nchecking existing files...")
	_ = checkExisting(s)
	elapsedpartial = time.Since(partial)
	fmt.Printf("OK (%v)\n", elapsedpartial)
	fmt.Printf("db file repeated count: %v (%.0f%%)\n", s.DbFileRepeatedCount, float64(s.DbFileRepeatedCount*int64(100)/int64(len(s.files))))
	fmt.Printf("same scan file repeated count: %v (%.0f%%)\n", s.SameScanFileRepeatedCount, float64(s.SameScanFileRepeatedCount*int64(100)/int64(len(s.files))))

	// scan finished
	elapsed := time.Since(start)
	s.Duration = elapsed.Milliseconds()
	fmt.Printf("\nscan process finished: %v\n", elapsed)

	// fmt.Println("\npress enter key to continue")
	// fmt.Scanln()

	// persist changes
	partial = time.Now()
	fmt.Printf("\npersist changes...")
	_ = persist(s, ntags)
	elapsedpartial = time.Since(partial)
	fmt.Printf("OK (%v)\n", elapsedpartial)

	// copy files
	partial = time.Now()
	fmt.Printf("\ncopy files...")
	_ = copyFiles(s, path.Join(np, source))
	elapsedpartial = time.Since(partial)
	fmt.Printf("OK (%v)\n", elapsedpartial)
}
