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
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s: %v\n", what, time.Since(start))
	}
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

func insertItem(db *sql.DB, name string, size int64, modifiedDate time.Time, createdDate time.Time, paths string) (int64, error) {
	// stmt, err := db.Prepare("INSERT item SET name=?, size=?, modified_date=?, created_date=?, type=?, paths=?")
	// if err != nil {
	// 	return -1, err
	// }
	// defer stmt.Close()

	// res, err := stmt.Exec(name, size, modifiedDate, createdDate, t, paths)
	// if err != nil {
	// 	return -1, err
	// }

	// return res.LastInsertId()

	return 1, nil
}

func scan(paths []string, db *sql.DB) (int, int) {
	p := ""
	fc := 0
	dc := 0

	stmtFile, err := db.Prepare("INSERT INTO `nfile` (`name`, `extension`, `path`, `modified`, `size`, `hash`, `ndirectory_id`) VALUES (?, ?, ?, ?, ?, ?, ?)")
	defer stmtFile.Close()
	if err != nil {
		fmt.Printf("error preparing file insert: %v\n", err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	stmtDirectory, err := db.Prepare("INSERT INTO `ndirectory` (`name`, `path`, `parent_id`, `size`, `count`) VALUES (?, ?, ?, ?, ?)")
	defer stmtFile.Close()
	if err != nil {
		fmt.Printf("error preparing directory insert: %v\n", err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	for i := 0; i < len(paths); i++ {
		p = paths[i]
		files, err := ioutil.ReadDir(p)
		if err != nil {
			log.Fatal(err)
		}

		for _, fileinfo := range files {
			if fileinfo.IsDir() {
				fmt.Printf("DIRECTORY:\n\tfull path:%v\n\tmod time:%v\n\tname:%v\n\tsize:%v\n\tmode:%v\n\n", p, fileinfo.ModTime(), fileinfo.Name(), fileinfo.Size(), fileinfo.Mode())
				paths = append(paths, path.Join(p, fileinfo.Name()))
				insertItem(db, fileinfo.Name(), fileinfo.Size(), fileinfo.ModTime(), time.Now(), p)
				_, err = stmtDirectory.Exec(fileinfo.Name(), p, 0, fileinfo.Size(), fileinfo.ModTime())
				if err != nil {
					fmt.Printf("error inserting order: %v, %v\n%v\n", order.ID, err, order.toString())
					bufio.NewReader(os.Stdin).ReadBytes('\n')
				}
				dc++
			} else {
				fmt.Printf("FILE:\n\tfull path:%v\n\tmod time:%v\n\tname:%v\n\tsize:%v\n\tmode:%v\n\n", path.Join(p, fileinfo.Name()), fileinfo.ModTime(), fileinfo.Name(), fileinfo.Size(), fileinfo.Mode())
				insertItem(db, fileinfo.Name(), fileinfo.Size(), fileinfo.ModTime(), time.Now(), p)
				fc++
			}
		}
	}

	return fc, dc
}

func main() {
	fmt.Printf("process started\n")
	defer elapsed("process finished")()

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/nostalgia")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	f, d := scan([]string{"/Users/marmoreno/Downloads"}, db)

	fmt.Printf("total %v files in %v directories\n", f, d)
}
