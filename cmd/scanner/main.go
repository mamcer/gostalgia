package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mamcer/nostalgia/internal/pkg/entities"
	"github.com/mamcer/nostalgia/internal/pkg/hash"
)

var physicalPath string = "/home/mario/stash"

func updateDirectorySize(db *sql.DB, rd int64, stmtUpdateDirectorySize *sql.Stmt) int64 {
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
		size += updateDirectorySize(db, did, stmtUpdateDirectorySize)
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
	db.QueryRow("SELECT nf.`size` FROM `nfile` as nf, `nfile_ndirectory` as nfs WHERE nf.id = nfs.nfile_id and nfs.ndirectory_id = ?", rd).Scan(&nfssizes)
	for i := 0; i < len(nfssizes); i++ {
		size += nfssizes[i]
	}

	rows3, err := db.Query("SELECT nf.`size` FROM `nfile` as nf, `nfile_ndirectory` as nfs WHERE nf.id = nfs.nfile_id and nfs.ndirectory_id = ?", rd)
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

	_, err = stmtUpdateDirectorySize.Exec(size, rd)
	if err != nil {
		fmt.Printf("error updating parent ndirectory: %v : %v\n", rd, err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	return size
}

func getFilePath(ns entities.Nscan, p string, s string) string {
	truncated := strings.Replace(p, ns.Name, "", 1)
	return path.Join(path.Join(s, fmt.Sprintf("%v", ns.ID)), truncated)
}

func copyFile(ns entities.Nscan, p string, fn string, stmtError *sql.Stmt) {
	pp := getFilePath(ns, p, physicalPath)
	err := os.MkdirAll(pp, 0755)
	if err != nil {
		_, _ = stmtError.Exec(fmt.Sprintf("failed to create directory: '%v' - %v", pp, err), ns.ID)
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
				_, _ = stmtError.Exec(fmt.Sprintf("failed to copy file: '%v' to '%v' - %v", ip, op, err), ns.ID)
			}
		} else {
			_, _ = stmtError.Exec(fmt.Sprintf("failed to create file to copy: %v - %v", op, err), ns.ID)
		}

	} else {
		_, _ = stmtError.Exec(fmt.Sprintf("failed to open file to copy: %v - %v", ip, err), ns.ID)
	}
}

func scan(root string, sname string, db *sql.DB) int {
	start := time.Now()
	fmt.Printf("scan process started\n")

	// nscan insert
	stmtScan, err := db.Prepare("INSERT INTO `nscan` (`date_created`, `status`, `name`, `retry_count`) VALUES (?, ?, ?, ?)")
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
	stmtDirectory, err := db.Prepare("INSERT INTO `ndirectory` (`name`, `path`, `date_modified`, `size`, `file_count`, `directory_count`, `parent_id`, `nscan_id`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmtDirectory.Close()
	if err != nil {
		fmt.Printf("error preparing ndirectory insert: %v\n", err)
		return 1
	}

	// ndirectory update
	stmtUpdateDirectory, err := db.Prepare("UPDATE `ndirectory` SET `size` = ?, `file_count` = ?, `directory_count` = ? WHERE id = ?")
	defer stmtUpdateDirectory.Close()
	if err != nil {
		fmt.Printf("error preparing ndirectory update: %v\n", err)
		return 1
	}

	// ndirectory size update
	stmtUpdateDirectorySize, err := db.Prepare("UPDATE `ndirectory` SET `size` = ? WHERE id = ?")
	defer stmtUpdateDirectorySize.Close()
	if err != nil {
		fmt.Printf("error preparing ndirectory update size: %v\n", err)
		return 1
	}

	// nfile_ndirectory insert
	stmtFileScan, err := db.Prepare("INSERT INTO `nfile_ndirectory` (`nfile_id`, `ndirectory_id`, `nscan_id`, `name`) VALUES (?, ?, ?, ?)")
	defer stmtFileScan.Close()
	if err != nil {
		fmt.Printf("error preparing nscan_nfile insert: %v\n", err)
		return 1
	}

	// nfile insert
	stmtFile, err := db.Prepare("INSERT INTO `nfile` (`name`, `extension`, `path`, `date_modified`, `size`, `hash`) VALUES (?, ?, ?, ?, ?, ?)")
	defer stmtFile.Close()
	if err != nil {
		fmt.Printf("error preparing nfile insert: %v\n", err)
		return 1
	}

	// nerror insert
	stmtError, err := db.Prepare("INSERT INTO `nerror` (`description`, `nscan_id`) VALUES (?, ?)")
	defer stmtError.Close()
	if err != nil {
		fmt.Printf("error preparing nerror insert: %v\n", err)
		return 1
	}

	// insert scan
	ns := entities.Nscan{DateCreated: time.Now(), Status: entities.InProgress, Name: root, RetryCount: 0}
	res, err := stmtScan.Exec(ns.DateCreated, ns.Status, ns.Name, ns.RetryCount)
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
		_, _ = stmtError.Exec(fmt.Sprintf("cannot stats root directory: '%v' - %v", root, err), ns.ID)
	}

	var ndirs []entities.Ndirectory = []entities.Ndirectory{{Name: sname, Path: getFilePath(ns, root, ""), Fpath: root, DateModified: fileStat.ModTime(), Size: 0, FileCount: 0, DirectoryCount: 0, ParentID: pid, NscanID: sid}}
	for i := 0; i < len(ndirs); i++ {
		p = ndirs[i].Fpath
		dfc := 0
		ddc := 0
		var ds int64

		parent := ndirs[i]
		res, err := stmtDirectory.Exec(parent.Name, parent.Path, parent.DateModified, parent.Size, parent.FileCount, parent.FileCount, parent.ParentID, parent.NscanID)
		if err != nil {
			fmt.Printf("error inserting parent directory '%v': %v\n", parent.Path, err)
			_, _ = stmtError.Exec(fmt.Sprintf("error inserting parent directory '%v' - %v", parent.Path, err), ns.ID)
			ec += 1
		}

		ppid := pid
		pid, _ = res.LastInsertId()
		parent.ID = pid

		if ppid == 1 {
			rid = parent.ID
			ns.RootDirectoryId = rid
			_, err = stmtUpdateScan.Exec(0, 0, 0, ns.Status, ns.RootDirectoryId, ns.RetryCount, ns.ID)
			if err != nil {
				fmt.Printf("error updating nscan: %v : %v\n", sid, err)
				_, _ = stmtError.Exec(fmt.Sprintf("error updating nscan: %v - %v", sid, err), ns.ID)
				ec += 1
			}
		}

		files, err := ioutil.ReadDir(p)
		if err != nil {
			fmt.Printf("error reading directory path: [%v] - %v", p, err)
			_, _ = stmtError.Exec(fmt.Sprintf("error reading directory path: [%v] - %v", p, err), ns.ID)
			ec += 1
		}

		for _, fileinfo := range files {
			fp := path.Join(p, fileinfo.Name())
			if fileinfo.IsDir() {
				if fileinfo.Name() != "." {
					d := entities.Ndirectory{Name: fileinfo.Name(), Path: getFilePath(ns, fp, ""), Fpath: fp, DateModified: fileinfo.ModTime(), Size: 0, FileCount: 0, ParentID: pid, NscanID: sid}
					ndirs = append(ndirs, d)
					dc++
					ddc++
				}
			} else {
				h, _ := hash.Calculate(fp)

				efi := 0
				efn := ""
				db.QueryRow("SELECT `id`, `name` FROM `nfile` WHERE `hash` = ?", h).Scan(&efi, &efn)
				if efi != 0 {
					// file exists
					efc += 1
					fmt.Printf("%.2f%% - [exists] %v\n", (float64(fc)+1)*100/float64(tfc), fp)
				} else {
					vp := getFilePath(ns, p, "")
					nfile := entities.Nfile{
						Name:         fileinfo.Name(),
						Extension:    strings.Trim(filepath.Ext(fileinfo.Name()), "."),
						Path:         vp,
						DateModified: fileinfo.ModTime(),
						Size:         fileinfo.Size(),
						Hash:         h,
					}
					res, err = stmtFile.Exec(nfile.Name, nfile.Extension, nfile.Path, nfile.DateModified, nfile.Size, nfile.Hash)
					if err != nil {
						fmt.Printf("[fail]\n%v\n", err)
						_, _ = stmtError.Exec(fmt.Sprintf("error inserting file [%v] - %v", nfile.Path, err), ns.ID)
						ec += 1
					}

					lastFileID, _ := res.LastInsertId()
					_, err := stmtFileScan.Exec(lastFileID, parent.ID, ns.ID, nfile.Name)
					if err != nil {
						fmt.Printf("error inserting file_scan '%v'- %v\n", efi, err)
						_, _ = stmtError.Exec(fmt.Sprintf("error inserting file_scan '%v' - %v", efi, err), ns.ID)
						ec += 1
					}

					// copy file
					//copyFile(ns, p, nfile.Name, stmtError)

					fmt.Printf("%.2f%% - [new] %v\n", (float64(fc)+1)*100/float64(tfc), fp)
				}

				fc++
				dfc++
				ds += fileinfo.Size()
			}
		}

		_, err = stmtUpdateDirectory.Exec(ds, dfc, ddc, parent.ID)
		if err != nil {
			fmt.Printf("error updating ndirectory: %v : %v\n", parent.ID, err)
			_, _ = stmtError.Exec(fmt.Sprintf("error updating ndirectory: %v - %v\n", parent.ID, err), ns.ID)
			ec += 1
		}
	}

	elapsed := time.Since(start)

	ns.FileCount = fc
	ns.DirectoryCount = dc
	ns.Status = entities.Done
	_, err = stmtUpdateScan.Exec(elapsed.Milliseconds(), ns.FileCount, ns.DirectoryCount, ns.Status, ns.RootDirectoryId, 0, ns.ID)
	if err != nil {
		fmt.Printf("error updating nscan: %v : %v\n", ns.ID, err)
		_, _ = stmtError.Exec(fmt.Sprintf("error updating nscan: %v : %v\n", ns.ID, err), ns.ID)
		ec += 1
	}

	var efp int64 = 0
	if ns.FileCount > 0 {
		efp = efc * 100 / ns.FileCount
	}
	fmt.Printf("process finished: %v\nscan_id: %v, files: %v, directories: %v, existing files: %v (%v%%), errors: %v\n", elapsed, ns.ID, ns.FileCount, ns.DirectoryCount, efc, efp, ec)

	// fmt.Printf("updating directory size...")
	//	ss := updateDirectorySize(db, ns.RootDirectoryId, stmtUpdateDirectorySize)
	// fmt.Printf("[ok]\ntotal size: %v bytes, %v\n", ss, files.SizeString(ss))

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
				db.QueryRow("SELECT `id` FROM `nscan` WHERE `name` = ?", p).Scan(&sid)
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
