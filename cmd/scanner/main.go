package main

import (
	"bufio"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	ID                int64     // scan id
	dateCreated       time.Time // scan creation date
	duration          int64     // scan duration (in milliseconds)
	fileCount         int64     // file scan count
	directoryCount    int64     // directory scan count
	status            int64     // scan status = done, inprogress, error
	rootDirectoryPath string    // scan root directory path
	rootDirectoryId   int64     // scan directory id
	retryCount        int64     // scan retry count
}

type ndirectory struct {
	ID           int64     // directory id
	name         string    // directory name
	path         string    // directory path
	size         int64     // directory size (in bytes)
	fileCount    int64     // directory file count
	parentID     int64     // parent directory id
	nscanID      int64     // scan id
	fpath        string    // current file path
	dateModified time.Time // date modified
}

type nfile struct {
	ID           int64     // file id
	name         string    //file name
	extension    string    //file extension
	path         string    // file path
	dateModified time.Time // file date modified
	size         int64     // file size (in bytes)
	hash         string    // file hash
	ndirectoryID int64     // file directory id
	nscanID      int64     // scan id
}

type nfilescan struct {
	ID           int64 // file scan id
	nfileID      int64 // file id
	ndirectoryID int64 // directory id
	nscanID      int64 // scan id
}

type nerror struct {
	ID          int64  // error id
	description string // error description
	nscanID     int64  // scan id
	retryCount  int64  // retry count
}

var vstash string = "stash"
var pstash string = "/media/darkforce/stash"

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

func updateDirectorySize(db *sql.DB, rd int64, stmtUpdateDirectory *sql.Stmt) int64 {
	var size int64 = 0

	// directory ids
	var dids []int64
	rows, err := db.Query("SELECT `id` FROM `ndirectory` WHERE `parent_id` = ?", rd)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			panic(err)
		}

		dids = append(dids, id)
	}
	rows.Close()

	for _, did := range dids {
		size += updateDirectorySize(db, did, stmtUpdateDirectory)
	}

	var nfsizes []int64
	db.QueryRow("SELECT `size` FROM `nfile` WHERE `ndirectory_id` = ?", rd).Scan(&nfsizes)
	for i := 0; i < len(nfsizes); i++ {
		size += nfsizes[i]
	}

	var fc int64 = 0
	rows2, err := db.Query("SELECT `size` FROM `nfile` WHERE `ndirectory_id` = ?", rd)
	if err != nil {
		panic(err)
	}
	defer rows2.Close()
	for rows2.Next() {
		var s int64
		err = rows2.Scan(&s)
		if err != nil {
			panic(err)
		}

		fc += 1
		size += s
	}

	var nfssizes []int64
	db.QueryRow("SELECT nf.`size` FROM `nfile` as nf, `nfile_nscan` as nfs WHERE nf.id = nfs.nfile_id and nfs.ndirectory_id = ?", rd).Scan(&nfssizes)
	for i := 0; i < len(nfssizes); i++ {
		size += nfssizes[i]
	}

	rows3, err := db.Query("SELECT nf.`size` FROM `nfile` as nf, `nfile_nscan` as nfs WHERE nf.id = nfs.nfile_id and nfs.ndirectory_id = ?", rd)
	if err != nil {
		panic(err)
	}
	defer rows3.Close()
	for rows3.Next() {
		var s int64
		err = rows3.Scan(&s)
		if err != nil {
			panic(err)
		}

		fc += 1
		size += s
	}

	_, err = stmtUpdateDirectory.Exec(size, fc, rd)
	if err != nil {
		fmt.Printf("error updating parent ndirectory: %v : %v\n", rd, err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	return size
}

func getFilePath(ns nscan, p string, s string) string {
	return strings.Replace(p, ns.rootDirectoryPath, path.Join(s, fmt.Sprintf("%v", ns.ID)), 1)
}

func copyFile(ns nscan, p string, fn string, stmtError *sql.Stmt) {
	pp := getFilePath(ns, p, pstash)
	err := os.MkdirAll(pp, 0755)
	if err != nil {
		_, _ = stmtError.Exec(fmt.Sprintf("failed to create directory: '%v' - %v", pp, err), ns.ID, ns.retryCount)
	}
	ip := path.Join(p, fn)
	i, err := os.Open(ip)
	defer i.Close()
	if err == nil {
		op := path.Join(pp, fn)
		o, err := os.Create(op)
		defer o.Close()
		if err == nil {
			_, err = io.Copy(o, i)
			if err != nil {
				_, _ = stmtError.Exec(fmt.Sprintf("failed to copy file: '%v' to '%v' - %v", ip, op, err), ns.ID, ns.retryCount)
			}
		} else {
			_, _ = stmtError.Exec(fmt.Sprintf("failed to create file to copy: %v - %v", op, err), ns.ID, ns.retryCount)
		}

	} else {
		_, _ = stmtError.Exec(fmt.Sprintf("failed to open file to copy: %v - %v", ip, err), ns.ID, ns.retryCount)
	}
}

func scan(root string, sname string, db *sql.DB) int {
	start := time.Now()
	fmt.Printf("scan process started\n")

	// nscan insert
	stmtScan, err := db.Prepare("INSERT INTO `nscan` (`date_created`, `status`, `root_directory_path`, `retry_count`) VALUES (?, ?, ?, ?)")
	defer stmtScan.Close()
	if err != nil {
		fmt.Printf("error preparing nscan insert: %v\n", err)
		return 1
	}

	// nscan update
	stmtUpdateScan, err := db.Prepare("UPDATE `nscan` SET `duration` = ?, `file_count` = ?, `directory_count` = ?, `status` = ?, `root_directory_id` = ?, `retry_count` = ? WHERE id = ?")
	defer stmtUpdateScan.Close()
	if err != nil {
		fmt.Printf("error preparing ndirectory update: %v\n", err)
		return 1
	}

	// ndirectory insert
	stmtDirectory, err := db.Prepare("INSERT INTO `ndirectory` (`name`, `path`, `date_modified`, `size`, `file_count`, `parent_id`, `nscan_id`) VALUES (?, ?, ?, ?, ?, ?, ?)")
	defer stmtDirectory.Close()
	if err != nil {
		fmt.Printf("error preparing ndirectory insert: %v\n", err)
		return 1
	}

	// ndirectory update
	stmtUpdateDirectory, err := db.Prepare("UPDATE `ndirectory` SET `size` = ?, `file_count` = ? WHERE id = ?")
	defer stmtUpdateDirectory.Close()
	if err != nil {
		fmt.Printf("error preparing ndirectory update: %v\n", err)
		return 1
	}

	// nfile_nscan insert
	stmtFileScan, err := db.Prepare("INSERT INTO `nfile_nscan` (`nfile_id`, `ndirectory_id`, `nscan_id`) VALUES (?, ?, ?)")
	defer stmtFileScan.Close()
	if err != nil {
		fmt.Printf("error preparing nscan_nfile insert: %v\n", err)
		return 1
	}

	// nfile insert
	stmtFile, err := db.Prepare("INSERT INTO `nfile` (`name`, `extension`, `path`, `date_modified`, `size`, `hash`, `ndirectory_id`, `nscan_id`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmtFile.Close()
	if err != nil {
		fmt.Printf("error preparing nfile insert: %v\n", err)
		return 1
	}

	// nerror insert
	stmtError, err := db.Prepare("INSERT INTO `nerror` (`description`, `nscan_id`, `retry_count`) VALUES (?, ?, ?)")
	defer stmtError.Close()
	if err != nil {
		fmt.Printf("error preparing nerror insert: %v\n", err)
		return 1
	}

	// ntag insert
	stmtTag, err := db.Prepare("INSERT INTO `ntag` (`name`, `nfile_id`, `ndirectory_id`) VALUES (?, ?, ?)")
	defer stmtTag.Close()
	if err != nil {
		fmt.Printf("error preparing ntag insert: %v\n", err)
		return 1
	}

	// insert scan
	ns := nscan{dateCreated: time.Now(), status: InProgress, rootDirectoryPath: root, retryCount: 0}
	res, err := stmtScan.Exec(ns.dateCreated, ns.status, ns.rootDirectoryPath, ns.retryCount)
	if err != nil {
		fmt.Printf("error inserting nscan: %v\n", err)
		return 1
	}

	sid, _ := res.LastInsertId()
	ns.ID = sid

	p := ""
	var tfc int64 = 0
	ps := []string{root}
	for i := 0; i < len(ps); i++ {
		p = ps[i]

		files, _ := ioutil.ReadDir(p)

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
	var fc int64 = 0
	var dc int64 = 1
	var efc int64 = 0
	ec := 0
	var pid int64 = 1
	var rid int64 = 1

	fileStat, err := os.Stat(root)
	if err != nil {
		_, _ = stmtError.Exec(fmt.Sprintf("cannot stats root directory: '%v' - %v", root, err), ns.ID, ns.retryCount)
	}

	var ndirs []ndirectory = []ndirectory{{name: sname, path: getFilePath(ns, root, vstash), fpath: root, dateModified: fileStat.ModTime(), size: 0, fileCount: 0, parentID: pid, nscanID: sid}}
	for i := 0; i < len(ndirs); i++ {
		p = ndirs[i].fpath
		dfc := 0
		var ds int64

		parent := ndirs[i]
		res, err := stmtDirectory.Exec(parent.name, parent.path, parent.dateModified, parent.size, parent.fileCount, parent.parentID, parent.nscanID)
		if err != nil {
			fmt.Printf("error inserting parent directory '%v': %v\n", parent.path, err)
			_, _ = stmtError.Exec(fmt.Sprintf("error inserting parent directory '%v' - %v", parent.path, err), ns.ID, ns.retryCount)
			ec += 1
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
				_, _ = stmtError.Exec(fmt.Sprintf("error updating nscan: %v - %v", sid, err), ns.ID, ns.retryCount)
				ec += 1
			}
		}

		files, err := ioutil.ReadDir(p)
		if err != nil {
			fmt.Printf("error reading directory path: [%v] - %v", p, err)
			_, _ = stmtError.Exec(fmt.Sprintf("error reading directory path: [%v] - %v", p, err), ns.ID, ns.retryCount)
			ec += 1
		}

		for _, fileinfo := range files {
			fp := p + "/" + fileinfo.Name()
			if fileinfo.IsDir() {
				if fileinfo.Name() != "." {
					d := ndirectory{name: fileinfo.Name(), path: getFilePath(ns, fp, vstash), fpath: fp, dateModified: fileinfo.ModTime(), size: 0, fileCount: 0, parentID: pid, nscanID: sid}
					ndirs = append(ndirs, d)
					dc++
				}
			} else {
				h, _ := calculateHash(fp)

				efi := 0
				efn := ""
				db.QueryRow("SELECT `id`, `name` FROM `nfile` WHERE `hash` = ?", h).Scan(&efi, &efn)
				if efi != 0 {
					// file exists
					_, err := stmtFileScan.Exec(efi, parent.ID, ns.ID)
					if err != nil {
						fmt.Printf("error inserting file_scan '%v'- %v\n", efi, err)
						_, _ = stmtError.Exec(fmt.Sprintf("error inserting file_scan '%v' - %v", efi, err), ns.ID, ns.retryCount)
						ec += 1
					}

					tn := fileinfo.Name()
					if efn != tn {
						// add new file name as tag
						et := 0
						db.QueryRow("SELECT count(`id`) FROM `ntag` WHERE `name` = ? and nfile_id = ?", tn, efi).Scan(&et)

						if et == 0 {
							_, err = stmtTag.Exec(tn, efi, parent.ID)
							if err != nil {
								fmt.Printf("error inserting tag: %v '%v' -  %v\n", efi, tn, err)
								_, _ = stmtError.Exec(fmt.Sprintf("error inserting tag: %v '%v' -  %v\n", efi, tn, err), ns.ID, ns.retryCount)
								ec += 1
							}
						}
					}

					efc += 1
					fmt.Printf("%.2f%% - [exists] %v\n", (float64(fc)+1)*100/float64(tfc), fp)
				} else {
					vp := getFilePath(ns, p, vstash)
					nfile := nfile{
						name:         fileinfo.Name(),
						extension:    strings.Trim(filepath.Ext(fileinfo.Name()), "."),
						path:         vp,
						dateModified: fileinfo.ModTime(),
						size:         fileinfo.Size(),
						hash:         h,
						ndirectoryID: parent.ID,
						nscanID:      ns.ID,
					}
					_, err = stmtFile.Exec(nfile.name, nfile.extension, nfile.path, nfile.dateModified, nfile.size, nfile.hash, nfile.ndirectoryID, nfile.nscanID)
					if err != nil {
						fmt.Printf("[fail]\n%v\n", err)
						_, _ = stmtError.Exec(fmt.Sprintf("error inserting file [%v] - %v", nfile.path, err), ns.ID, ns.retryCount)
						ec += 1
					}

					// copy file
					copyFile(ns, p, nfile.name, stmtError)

					fmt.Printf("%.2f%% - [new] %v\n", (float64(fc)+1)*100/float64(tfc), fp)
				}

				fc++
				dfc++
				ds += fileinfo.Size()
			}
		}

		_, err = stmtUpdateDirectory.Exec(ds, dfc, parent.ID)
		if err != nil {
			fmt.Printf("error updating ndirectory: %v : %v\n", parent.ID, err)
			_, _ = stmtError.Exec(fmt.Sprintf("error updating ndirectory: %v - %v\n", parent.ID, err), ns.ID, ns.retryCount)
			ec += 1
		}
	}

	elapsed := time.Since(start)

	ns.fileCount = fc
	ns.directoryCount = dc
	ns.status = Done
	_, err = stmtUpdateScan.Exec(elapsed.Milliseconds(), ns.fileCount, ns.directoryCount, ns.status, ns.rootDirectoryId, 0, ns.ID)
	if err != nil {
		fmt.Printf("error updating nscan: %v : %v\n", ns.ID, err)
		_, _ = stmtError.Exec(fmt.Sprintf("error updating nscan: %v : %v\n", ns.ID, err), ns.ID, ns.retryCount)
		ec += 1
	}

	var efp int64 = 0
	if ns.fileCount > 0 {
		efp = efc * 100 / ns.fileCount
	}
	fmt.Printf("process finished: %v\nscan_id: %v, files: %v, directories: %v, existing files: %v (%v%%), errors: %v\n", elapsed, ns.ID, ns.fileCount, ns.directoryCount, efc, efp, ec)

	fmt.Printf("updating directory size...")
	ss := updateDirectorySize(db, ns.rootDirectoryId, stmtUpdateDirectory)
	fmt.Printf("[ok]\ntotal size: %v bytes, %v\n", ss, sizeString(ss))

	return ec
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
				sid := 0
				db.QueryRow("SELECT `id` FROM `nscan` WHERE `root_directory_path` = ?", p).Scan(&sid)
				if sid != 0 {
					fmt.Printf("There is an existing scan (%v) with root directory path: %v\n", sid, p)
					reader := bufio.NewReader(os.Stdin)
					for out := false; out == false; {
						fmt.Printf("Do you want to continue? [Y/n] ")
						k, _ := reader.ReadString('\n')
						if err == nil {
							if k == "\n" || k == "y\n" || k == "Y\n" {
								out = true
								sn := path.Base(p)
								if len(os.Args) > 3 && os.Args[3] != "" {
									sn = os.Args[3]
								}
								_ = scan(p, sn, db)
							} else if k == "n\n" || k == "N\n" {
								out = true
							}
						} else {
							out = true
						}
					}
				} else {
					sn := path.Base(p)
					if len(os.Args) > 3 && os.Args[3] != "" {
						sn = os.Args[3]
					}
					_ = scan(p, sn, db)
				}
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
