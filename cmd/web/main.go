package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type NResult struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Extension    string `json:"extension"`
	Path         string `json:"path"`
	DateModified string `json:"date_modified"`
	Size         string `json:"size"`
	ParentID     int64  `json:"parent_id"`
	ParentName   string `json:"parent_name"`
	Type         string `json:"type"`
}

// Configuration container
type Configuration struct {
	ApiPort          string
	WebPort          string
	DBDriverName     string
	DBDataSourceName string
}

var config Configuration

func getDB() *sql.DB {
	var err error
	db, err := sql.Open(config.DBDriverName, config.DBDataSourceName)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func ping(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func filesCount(c *gin.Context) {
	count := 0
	db := getDB()
	db.QueryRow("SELECT count(id) from nfile").Scan(&count)
	defer db.Close()

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
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

func search(c *gin.Context) {
	// GET contains=a&type=[image|doc|sheet|audio|video|zip|any]&only_directories=false&after=1000-01-01&before=9999-12-31&page=1&per_page=50
	var nt mysql.NullTime
	var size int64
	contains := c.DefaultQuery("contains", "mario")
	onlyDirs := c.DefaultQuery("only_directories", "false")
	t := c.DefaultQuery("type", "any")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		fmt.Printf("error converting page to int: %v", err)
		page = 1
	}
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	if err != nil {
		fmt.Printf("error converting per page to int: %v", err)
		perPage = 50
	}
	layout := "2006-01-02"

	a := c.Query("after")
	after, _ := time.Parse(layout, "1000-01-01")
	if a != "" {
		after, _ = time.Parse(layout, a)
	}

	b := c.Query("before")
	before, _ := time.Parse(layout, "9999-12-31")
	if b != "" {
		before, _ = time.Parse(layout, b)
	}

	fmt.Printf("contains: '%v', type: '%v', only_directories: '%v', after: '%v', before: '%v', page: '%v', per_page: '%v'\n", contains, t, onlyDirs, after, before, page, perPage)

	// files
	var results []NResult
	var total int64 = 0
	var rtype = ""

	db := getDB()

	var sq = `SELECT 	n.id as ID,
						n.name,
						n.extension,
						n.path,
						n.size,
						n.date_modified as DateModified,
				FROM nfile as n`

	var where = `WHERE lower(n.name) like ?`

	switch t {
	case "image":
		where += " and n.extension in ('jpeg', 'png', 'jpg', 'bmp')"
	case "doc":
		where += " and n.extension in ('doc', 'docx', 'odt', 'pdf')"
	case "sheet":
		where += " and n.extension in ('xls', 'xlsx', 'ods')"
	case "audio":
		where += " and n.extension in ('mp3', 'ogg', 'wma', 'arm', 'wav')"
	case "video":
		where += " and n.extension in ('mp4', 'mkv', 'avi', 'wmv')"
	case "zip":
		where += " and n.extension in ('zip', 'rar', '7z', 'gz')"
	}

	where += " and n.date_modified between ? and ?"
	sq += where + " limit ? offset ?"

	//limit = per_page
	//offset = (page-1)*per_page

	rtype = "file"
	rows, err := db.Query(sq,
		"%"+strings.ToLower(contains)+"%", after, before, perPage, (page-1)*perPage)
	db.QueryRow("SELECT count(id) from nfile " + fmt.Sprintf(where, "%"+strings.ToLower(contains)+"%", after, before)).Scan(&total)
	defer db.Close()

	if err != nil {
		results = nil
	} else {
		for rows.Next() {
			var r NResult
			rows.Scan(&r.ID, &r.Name, &r.Extension, &r.Path, &size, &nt, &r.ParentID, &r.ParentName, rtype)
			r.Size = sizeString(size)
			if nt.Valid {
				r.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
			}

			results = append(results, r)
		}
	}

	// // directories
	// var directories []NDirectory
	// db2 := getDB()
	// rows, err = db2.Query(`SELECT d.id as ID, d.name, d.path, d.date_modified as DateModified, d.size, d.file_count, d.directory_count, d.parent_id FROM ndirectory as d WHERE lower(d.name) like ? and d.date_modified between ? and ? limit ? offset ?`,
	// 	"%"+strings.ToLower(contains)+"%", after, before, perPage, (page-1)*perPage)
	// defer db2.Close()

	// if err != nil {
	// 	directories = nil
	// } else {
	// 	for rows.Next() {
	// 		var d NDirectory
	// 		rows.Scan(&d.ID, &d.Name, &d.Path, &nt, &size, &d.FileCount, &d.DirectoryCount, &d.ParentID)
	// 		d.Size = sizeString(size)
	// 		if nt.Valid {
	// 			d.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
	// 		}

	// 		directories = append(directories, d)
	// 	}
	// }

	// result
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	if results == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"results":  nil,
			"contains": contains,
			"page":     0,
			"per_page": 0,
			"total":    0,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"results":  results,
			"contains": contains,
			"page":     page,
			"per_page": perPage,
			"total":    len(total),
		})
	}
}

func fileController(c *gin.Context) {
	id := c.Param("id")

	var nt mysql.NullTime
	var size int64
	var r Nfile
	if id != "" {
		db := getDB()
		err := db.QueryRow("SELECT n.id as ID, n.name, n.extension, n.path, n.size, n.date_modified as DateModified, n.hash FROM nfile as n WHERE n.id = ?", id).Scan(&r.ID, &r.Name, &r.Extension, &r.Path, &size, &nt, &r.Hash)
		defer db.Close()

		if err != sql.ErrNoRows {
			r.Size = sizeString(size)
			if nt.Valid {
				r.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
			}

			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
			c.JSON(http.StatusOK, r)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
			c.JSON(http.StatusNotFound, struct{}{})
		}
	}

}

func directoriesController(c *gin.Context) {
	id := c.Param("id")

	var nt mysql.NullTime
	var size int64

	if id != "" {
		var d NDirectory
		db := getDB()
		err := db.QueryRow("SELECT d.id as ID, d.name, d.path, d.date_modified as DateModified, d.size, d.file_count, d.directory_count, d.parent_id FROM ndirectory as d WHERE d.id = ?", id).Scan(&d.ID, &d.Name, &d.Path, &nt, &size, &d.FileCount, &d.DirectoryCount, &d.ParentID)
		defer db.Close()
		if err != sql.ErrNoRows {
			d.Size = sizeString(size)
			if nt.Valid {
				d.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
			}

			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
			c.JSON(http.StatusOK, d)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
			c.JSON(http.StatusNotFound, struct{}{})
		}
	}
}

func directoryFilesController(c *gin.Context) {
	id := c.Param("id")

	var nt mysql.NullTime
	var size int64

	if id != "" {
		db := getDB()
		query := `SELECT  f.id as ID, 
        f.name, 
        f.extension, 
        f.path, 
        f.size, 
        f.date_modified as DateModified, 
        f.hash 
		FROM nfile AS f, ndirectory as d, nfile_ndirectory AS fd
		WHERE d.id = ? and fd.ndirectory_id = d.id and fd.nfile_id = f.id`

		rows, err := db.Query(query, id)
		defer db.Close()

		if err != nil {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
			c.JSON(http.StatusNotFound, struct{}{})
		} else {
			var files []Nfile
			for rows.Next() {
				var r Nfile
				rows.Scan(&r.ID, &r.Name, &r.Extension, &r.Path, &size, &nt, &r.Hash)
				r.Size = sizeString(size)
				if nt.Valid {
					r.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
				}

				files = append(files, r)
			}

			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
			c.JSON(http.StatusOK, files)
		}
	}
}

func directoryDirectoriesController(c *gin.Context) {
	id := c.Param("id")

	var nt mysql.NullTime
	var size int64

	if id != "" {
		db := getDB()
		query := `SELECT  d.id as ID, 
        d.name, 
        d.path, 
        d.date_modified as DateModified, 
        d.size, 
        d.file_count, 
        d.directory_count, 
        d.parent_id 
		FROM ndirectory as d
		WHERE d.parent_id = ?`

		rows, err := db.Query(query, id)
		if err != nil {
			fmt.Printf("error executing query: %v", err)
			return
		}
		defer db.Close()

		if err != nil {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
			c.JSON(http.StatusNotFound, struct{}{})
		} else {
			var directories []NDirectory
			for rows.Next() {
				var r NDirectory
				rows.Scan(&r.ID, &r.Name, &r.Path, &nt, &size, &r.FileCount, &r.DirectoryCount, &r.ParentID)
				r.Size = sizeString(size)
				if nt.Valid {
					r.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
				}

				directories = append(directories, r)
			}

			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
			c.JSON(http.StatusOK, directories)
		}
	}
}

func preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.JSON(http.StatusOK, struct{}{})
}

func main() {
	f, err := os.Open("config.json")
	if err != nil {
		fmt.Printf("error opening config.json: %v", err)
		return
	}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Printf("error decoding config.json: %v", err)
		return
	}

	router := gin.Default()

	v1 := router.Group("/v1")

	v1.GET("/ping", ping)
	v1.OPTIONS("/ping", preflight)

	v1.GET("/search", search)
	v1.OPTIONS("/search", preflight)

	v1.GET("/files/:id", fileController)
	v1.OPTIONS("/files/:id", preflight)

	v1.GET("/filescount", filesCount)
	v1.OPTIONS("/filescount", preflight)

	v1.GET("/directories/:id", directoriesController)
	v1.OPTIONS("/directories/:id", preflight)

	v1.GET("/directories/:id/files", directoryFilesController)
	v1.OPTIONS("/directories/:id/files", preflight)

	v1.GET("/directories/:id/directories", directoryDirectoriesController)
	v1.OPTIONS("/directories/:id/directories", preflight)

	go func() {
		http.Handle("/",
			http.StripPrefix("/",
				http.FileServer(http.Dir("./"))))
		log.Fatal(http.ListenAndServe(":"+config.WebPort, nil))
	}()

	_ = router.Run(":" + config.ApiPort)
}
