package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
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

type nscan struct {
	directories []ndirectory
	files       map[int64][]nfile
	fileCount   int64
	fileSize    int64
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s: %v\n", what, time.Since(start))
	}
}

func scan(paths []string, result *nscan) {
	p := ""
	var ds int64
	var id int64 = 1
	for i := 0; i < len(paths); i++ {
		p = paths[i]

		name := path.Base(p)

		parent := ndirectory{id: id, name: name, path: p}
		id++
		result.directories = append(result.directories, parent)

		files, err := ioutil.ReadDir(p)
		if err != nil {
			fmt.Printf("io error: %s\n", err)
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}

		ds = 0
		for _, fileinfo := range files {
			fp := path.Join(p, fileinfo.Name())
			if fileinfo.IsDir() {
				paths = append(paths, fp)
			} else {
				nfile := nfile{
					name:         fileinfo.Name(),
					extension:    strings.Trim(filepath.Ext(fileinfo.Name()), "."),
					path:         p,
					modified:     fileinfo.ModTime(),
					size:         fileinfo.Size(),
					ndirectoryID: parent.id,
				}
				result.files[parent.id] = append(result.files[parent.id], nfile)
				result.fileCount++
				ds += nfile.size
			}
		}

		parent.size = ds
		result.fileSize += parent.size
	}
}

// func handle(queue chan *Request) {
//     for r := range queue {
//         process(r)
//     }
// }

// func Serve(clientRequests chan *Request, quit chan bool) {
//     // Start handlers
//     for i := 0; i < MaxOutstanding; i++ {
//         go handle(clientRequests)
//     }
//     <-quit  // Wait to be told to exit.
// }

// var sem = make(chan int, MaxOutstanding)

// func Serve(queue chan *Request) {
//     for req := range queue {
//         sem <- 1
//         go func(req *Request) {
//             process(req)
//             <-sem
//         }(req)
//     }
// }

func process(result *nscan) {
	var f int64 = 1
	sem := make(chan int, 100)
	var wg sync.WaitGroup
	var index int = 1
	for i, dir := range result.directories {
		fmt.Printf("[%06d/%06v] %v\n", i+1, len(result.directories), dir.path)
		wg.Add(1)
		sem <- 1
		go func(files []nfile, index int) {
			defer wg.Done()
			for _, file := range files {
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
				//fmt.Printf("%-5v %v : [%06d/%06v] \t%v\n", "", index, f, result.fileCount, file.name)
				fmt.Printf("%-5v %v f:%v - %v\n", "", index, f, file.name)
				f++
			}
			<-sem
		}(result.files[dir.id], index)
		index++
	}
	wg.Wait()
}

func main() {
	fmt.Printf("process started\n")
	defer elapsed("process finished")()

	var result nscan
	result.files = make(map[int64][]nfile)
	fmt.Printf("processing directories\t")
	scan([]string{"/media/mario/etc/ordenar-ultimo-scan/music"}, &result)
	fmt.Printf("[ok]\n")

	fmt.Printf("total %v  items: %v files in %v directories, total size: %vGB\n", int64(len(result.directories))+result.fileCount, len(result.directories), result.fileCount, result.fileSize/1000/1000/1000)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	process(&result)
}
