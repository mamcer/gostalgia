package main

import (
	"bufio"
	"crypto/sha1"
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

func scan(paths []string) (int, int) {
	p := ""
	fc := 0
	dc := 0
	var parentID int64 = 1
	for i := 0; i < len(paths); i++ {
		p = paths[i]
		dfc := 0
		var ds int64

		fmt.Printf("%v\n", p)
		name := path.Base(p)

		fileStat, err := os.Stat(p)
		if err != nil {
			log.Fatal(err)
		}

		parent := ndirectory{name: name, path: p, modified: fileStat.ModTime(), parentID: parentID, size: 10, count: 10}
		if err != nil {
			fmt.Printf("error inserting parent directory '%v': %v\n", parent.path, err)
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}

		parent.id = 0

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
				fmt.Printf("\t%-120v", fileinfo.Name())
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
				fmt.Printf("%v\n", "[ok]")
				fc++
				dfc++
				ds += nfile.size
			}
		}
	}

	return fc, dc
}

func main() {
	fmt.Printf("process started\n")
	defer elapsed("process finished")()

	// db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/nostalgia")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()

	f, d := scan([]string{"/media/mario/etc/ordenar-ultimo-scan/varios-scan"}) //, db)

	fmt.Printf("total %v files in %v directories\n", f, d)
}
