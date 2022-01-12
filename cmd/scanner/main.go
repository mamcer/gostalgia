package main

import (
	"bufio"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type ndirectory struct {
	id       int64
	name     string
	path     string
	modified time.Time
	parentID int64
	size     int64
	count    int64
}

type nfile struct {
	id           int64
	name         string
	extension    string
	path         string
	modified     time.Time
	size         int64
	hash         string
	ndirectoryID int64
}

func calculateHash(filePath string) (string, error) {
	var sha string

	file, err := os.Open(filePath)
	if err != nil {
		return sha, err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return sha, err
	}

	bytes := hash.Sum(nil)[:20]
	sha = hex.EncodeToString(bytes)
	return sha, nil
}

func scan(paths []string, db *sql.DB) (int, int) {
	start := time.Now()
	fmt.Printf("scan process started\n")

	// stmtFile, err := db.Prepare("INSERT INTO `nfile` (`name`, `extension`, `path`, `modified`, `size`, `hash`, `ndirectory_id`) VALUES (?, ?, ?, ?, ?, ?, ?)")
	// defer stmtFile.Close()
	// if err != nil {
	// 	fmt.Printf("error preparing file insert: %v\n", err)
	// 	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// }

	// stmtDirectory, err := db.Prepare("INSERT INTO `ndirectory` (`name`, `path`, `modified`, `parent_id`, `size`, `count`) VALUES (?, ?, ?, ?, ?, ?)")
	// defer stmtDirectory.Close()
	// if err != nil {
	// 	fmt.Printf("error preparing directory insert: %v\n", err)
	// 	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// }

	// stmtUpdateDirectory, err := db.Prepare("UPDATE `ndirectory` SET `size` = ?, `count` = ? WHERE id = ?")
	// defer stmtUpdateDirectory.Close()
	// if err != nil {
	// 	fmt.Printf("error preparing directory update: %v\n", err)
	// 	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// }

	p := ""
	tfc := 0
	ps := []string{paths[0]}
	for i := 0; i < len(ps); i++ {
		p = ps[i]

		files, err := ioutil.ReadDir(p)
		if err != nil {
			log.Fatal(err)
		}

		for _, fileinfo := range files {
			fp := path.Join(p, fileinfo.Name())
			if fileinfo.IsDir() {
				ps = append(ps, fp)
			} else {
				tfc += 1
			}
		}
	}

	fmt.Printf("total file count: %v\n", tfc)

	p = ""
	fc := 0
	dc := 0
	var parentID int64 = 1
	for i := 0; i < len(paths); i++ {
		p = paths[i]
		dfc := 0
		var ds int64

		name := path.Base(p)

		fileStat, err := os.Stat(p)
		if err != nil {
			log.Fatal(err)
		}

		parent := ndirectory{name: name, path: p, modified: fileStat.ModTime(), parentID: parentID, size: 10, count: 10}
		//res, err := stmtDirectory.Exec(parent.name, parent.path, parent.modified, parent.parentID, parent.size, parent.count)
		if err != nil {
			fmt.Printf("error inserting parent directory '%v': %v\n", parent.path, err)
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}

		//parentID, _ := res.LastInsertId()
		//parent.id = parentID

		files, err := ioutil.ReadDir(p)
		if err != nil {
			log.Fatal(err)
		}

		for _, fileinfo := range files {
			fp := path.Join(p, fileinfo.Name())
			if fileinfo.IsDir() {
				paths = append(paths, fp)
				dc++
			} else {
				fmt.Printf("%+03v%% - %v\n", (fc+1)*100/tfc, fp)
				h, _ := calculateHash(fp)
				nfile := nfile{
					name:         fileinfo.Name(),
					extension:    strings.Trim(filepath.Ext(fileinfo.Name()), "."),
					path:         p,
					modified:     fileinfo.ModTime(),
					size:         fileinfo.Size(),
					hash:         h,
					ndirectoryID: parentID,
				}
				//_, err = stmtFile.Exec(nfile.name, nfile.extension, nfile.path, nfile.modified, nfile.size, nfile.hash, nfile.ndirectoryID)
				if err != nil {
					fmt.Printf("[fail]\n%v\n", err)
					bufio.NewReader(os.Stdin).ReadBytes('\n')
				}

				fc++
				dfc++
				ds += nfile.size
			}
		}

		//_, err = stmtUpdateDirectory.Exec(ds, dfc, parent.id)
		// if err != nil {
		// 	fmt.Printf("error updating directory: %v : %v\n", parent.id, err)
		// 	bufio.NewReader(os.Stdin).ReadBytes('\n')
		// }
	}

	elapsed := time.Since(start)
	fmt.Printf("process finished: %v\n", elapsed)

	return fc, dc
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/nostalgia")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if len(os.Args) > 1 {
		if os.Args[1] == "scan" {
			if len(os.Args) > 2 {
				p := os.Args[2]
				_, _ = scan([]string{p}, db)
			} else {
				fmt.Printf("you must provide a valid path to scan\n")
			}
		} else if os.Args[1] == "retry" {
			if len(os.Args) > 2 {
				fmt.Printf("%v %v\n", os.Args[1], os.Args[2])
			} else {
				fmt.Printf("you must provide a valid scan id to retry\n")
			}
		} else if os.Args[1] == "scans" {
			fmt.Printf("%v\n", os.Args[1])
		} else if os.Args[1] == "errors" {
			if len(os.Args) > 2 {
				fmt.Printf("%v %v\n", os.Args[1], os.Args[2])
			} else {
				fmt.Printf("you must provide a valid scan id to see errors\n")
			}
		} else {
			fmt.Printf("unknown command: %v\n", os.Args[1])
		}
	} else {
		fmt.Printf("usage:\nscan [path]\nretry [scan-id]\nscans\nerrors [scan-id]\n")
	}

	//f, d := scan([]string{"/media/darkforce"}, db)

	//fmt.Printf("total %v files in %v directories\n", f, d)
}
