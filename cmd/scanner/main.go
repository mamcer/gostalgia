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

const (
	Done int64 = iota
	InProgress
	Error
)

type nscan struct {
	ID                int64
	dateCreated       time.Time
	duration          int64
	fileCount         int64
	directoryCount    int64
	status            int64
	rootDirectoryPath string
	rootDirectoryId   int64
	retryCount        int64
}

type ndirectory struct {
	ID        int64
	name      string
	path      string
	size      int64
	fileCount int64
	parentID  int64
	nscanID   int64
}

type nfile struct {
	ID           int64
	name         string
	extension    string
	path         string
	dateModified time.Time
	size         int64
	hash         string
	ndirectoryID int64
	nscanID      int64
}

type nfilescan struct {
	ID           int64
	nfileID      int64
	ndirectoryID int64
	nscanID      int64
}

type nerror struct {
	ID          int64
	description string
	nscanID     int64
	retryCount  int64
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

func scan(paths []string, db *sql.DB) int {
	start := time.Now()
	fmt.Printf("scan process started\n")

	// nscan insert
	stmtScan, err := db.Prepare("INSERT INTO `nscan` (`date_created`, `status`, `root_directory_path`, `retry_count`) VALUES (?, ?, ?, ?)")
	defer stmtScan.Close()
	if err != nil {
		fmt.Printf("error preparing nscan insert: %v\n", err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(1)
	}

	// nscan update
	stmtUpdateScan, err := db.Prepare("UPDATE `nscan` SET `duration` = ?, `file_count` = ?, `directory_count` = ?, `status` = ?, `root_directory_id` = ?, `retry_count` = ? WHERE id = ?")
	defer stmtUpdateScan.Close()
	if err != nil {
		fmt.Printf("error preparing ndirectory update: %v\n", err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	// ndirectory insert
	stmtDirectory, err := db.Prepare("INSERT INTO `ndirectory` (`name`, `path`, `size`, `file_count`, `parent_id`, `nscan_id`) VALUES (?, ?, ?, ?, ?, ?)")
	defer stmtDirectory.Close()
	if err != nil {
		fmt.Printf("error preparing ndirectory insert: %v\n", err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	// ndirectory update
	stmtUpdateDirectory, err := db.Prepare("UPDATE `ndirectory` SET `size` = ?, `file_count` = ? WHERE id = ?")
	defer stmtUpdateDirectory.Close()
	if err != nil {
		fmt.Printf("error preparing ndirectory update: %v\n", err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	// nfile_nscan insert
	stmtFileScan, err := db.Prepare("INSERT INTO `nfile_nscan` (`nfile_id`, `ndirectory_id`, `nscan_id`) VALUES (?, ?, ?)")
	defer stmtFileScan.Close()
	if err != nil {
		fmt.Printf("error preparing nscan_nfile insert: %v\n", err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(1)
	}

	// nfile insert
	stmtFile, err := db.Prepare("INSERT INTO `nfile` (`name`, `extension`, `path`, `date_modified`, `size`, `hash`, `ndirectory_id`, `nscan_id`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmtFile.Close()
	if err != nil {
		fmt.Printf("error preparing nfile insert: %v\n", err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	// insert scan
	ns := nscan{dateCreated: time.Now(), status: InProgress, rootDirectoryPath: paths[0], retryCount: 0}
	res, err := stmtScan.Exec(ns.dateCreated, ns.status, ns.rootDirectoryPath, ns.retryCount)
	if err != nil {
		fmt.Printf("error inserting nscan: %v\n", err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	sid, _ := res.LastInsertId()
	ns.ID = sid

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
	dc := 1
	var pid int64 = 1
	var rid int64 = 1
	for i := 0; i < len(paths); i++ {
		p = paths[i]
		dfc := 0
		var ds int64

		name := path.Base(p)

		parent := ndirectory{name: name, path: p, size: 0, fileCount: 0, parentID: pid, nscanID: sid}
		res, err := stmtDirectory.Exec(parent.name, parent.path, parent.size, parent.fileCount, parent.parentID, parent.nscanID)
		if err != nil {
			fmt.Printf("error inserting parent directory '%v': %v\n", parent.path, err)
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}

		ppid := pid
		pid, _ = res.LastInsertId()
		parent.ID = pid

		if ppid == 1 {
			rid = parent.ID
			ns.rootDirectoryId = rid
			_, err = stmtUpdateScan.Exec(0, 0, 0, ns.status, ns.rootDirectoryId, ns.retryCount, ns.ID)
			if err != nil {
				fmt.Printf("error updating nscan: %v : %v\n", sid, err)
				bufio.NewReader(os.Stdin).ReadBytes('\n')
			}
		}

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

				efi := 0
				db.QueryRow("SELECT `id` FROM `nfile` WHERE `hash` = ?", h).Scan(&efi)
				if efi != 0 {
					// file exists
					_, err := stmtFileScan.Exec(efi, parent.ID, ns.ID)
					if err != nil {
						fmt.Printf("error inserting file_scan '%v': %v\n", efi, err)
						bufio.NewReader(os.Stdin).ReadBytes('\n')
					}
				} else {
					nfile := nfile{
						name:         fileinfo.Name(),
						extension:    strings.Trim(filepath.Ext(fileinfo.Name()), "."),
						path:         p,
						dateModified: fileinfo.ModTime(),
						size:         fileinfo.Size(),
						hash:         h,
						ndirectoryID: parent.ID,
						nscanID:      ns.ID,
					}
					_, err = stmtFile.Exec(nfile.name, nfile.extension, nfile.path, nfile.dateModified, nfile.size, nfile.hash, nfile.ndirectoryID, nfile.nscanID)
					if err != nil {
						fmt.Printf("[fail]\n%v\n", err)
					}
				}

				fc++
				dfc++
				ds += fileinfo.Size()
			}
		}

		_, err = stmtUpdateDirectory.Exec(ds, dfc, parent.ID)
		if err != nil {
			fmt.Printf("error updating parent ndirectory: %v : %v\n", parent.ID, err)
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}

	elapsed := time.Since(start)

	_, err = stmtUpdateScan.Exec(elapsed.Milliseconds(), fc, dc, Done, rid, 0, ns.ID)
	if err != nil {
		fmt.Printf("error updating nscan: %v : %v\n", ns.ID, err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	fmt.Printf("process finished: %v\n", elapsed)

	return 0
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
				_ = scan([]string{p}, db)
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
}
