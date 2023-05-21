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

type Nfile struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Extension    string `json:"extension"`
	Path         string `json:"path"`
	DateModified string `json:"date_modified"`
	Size         string `json:"size"`
	Hash         string `json:"hash"`
}

type NDirectory struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	DateModified   string `json:"date_modified"`
	Size           string `json:"size"`
	FileCount      int64  `json:"file_count"`
	DirectoryCount int64  `json:"directory_count"`
	ParentID       int64  `json:"parent_id"`
	NScanID        int64  `json:"nscan_id"`
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
	var nt mysql.NullTime
	var size int64
	query := c.DefaultQuery("q", "mario")
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

	fmt.Printf("query: '%v', type: '%v', after: '%v', before: '%v'\n", query, t, after, before)

	// files
	var files []Nfile

	db := getDB()

	var sq = `SELECT 	n.id as ID,
						n.name,
						n.extension,
						n.path,
						n.size,
						n.date_modified as DateModified,
						n.hash
				FROM nfile as n
				WHERE lower(n.name) like ?`

	switch t {
	case "image":
		sq += " and n.extension in ('jpeg', 'png', 'jpg', 'bmp')"
	case "doc":
		sq += " and n.extension in ('doc', 'docx', 'odt', 'pdf')"
	case "sheet":
		sq += " and n.extension in ('xls', 'xlsx', 'ods')"
	case "audio":
		sq += " and n.extension in ('mp3', 'ogg', 'wma', 'arm', 'wav')"
	case "video":
		sq += " and n.extension in ('mp4', 'mkv', 'avi', 'wmv')"
	case "zip":
		sq += " and n.extension in ('zip', 'rar', '7z', 'gz')"
	}

	sq += " and n.date_modified between ? and ? limit ? offset ?"

	//limit = per_page
	//offset = (page-1)*per_page

	rows, err := db.Query(sq,
		"%"+strings.ToLower(query)+"%", after, before, perPage, (page-1)*perPage)
	defer db.Close()

	if err != nil {
		files = nil
	} else {
		for rows.Next() {
			var r Nfile
			rows.Scan(&r.ID, &r.Name, &r.Extension, &r.Path, &size, &nt, &r.Hash)
			r.Size = sizeString(size)
			if nt.Valid {
				r.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
			}

			files = append(files, r)
		}
	}

	// directories
	var directories []NDirectory
	db2 := getDB()
	rows, err = db2.Query(`SELECT d.id as ID, d.name, d.path, d.date_modified as DateModified, d.size, d.file_count, d.directory_count, d.parent_id, d.nscan_id FROM ndirectory as d WHERE lower(d.name) like ? and ? limit ? offset ?`,
		"%"+strings.ToLower(query)+"%", perPage, (page-1)*perPage)
	defer db2.Close()

	if err != nil {
		directories = nil
	} else {
		for rows.Next() {
			var d NDirectory
			rows.Scan(&d.ID, &d.Name, &d.Path, &nt, &size, &d.FileCount, &d.DirectoryCount, &d.ParentID, &d.NScanID)
			d.Size = sizeString(size)
			if nt.Valid {
				d.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
			}

			directories = append(directories, d)
		}
	}

	// result
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	if files == nil && directories == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"query":             query,
			"total_directories": 0,
			"total_files":       0,
			"page":              0,
			"per_page":          0,
			"directories":       nil,
			"files":             nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"query":             query,
			"total_directories": len(directories),
			"total_files":       len(files),
			"page":              page,
			"per_page":          perPage,
			"directories":       directories,
			"files":             files,
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
		err := db.QueryRow("SELECT d.id as ID, d.name, d.path, d.date_modified as DateModified, d.size, d.file_count, d.directory_count, d.parent_id, d.nscan_id FROM ndirectory as d WHERE d.id = ?", id).Scan(&d.ID, &d.Name, &d.Path, &nt, &size, &d.FileCount, &d.DirectoryCount, &d.ParentID, &d.NScanID)
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
        d.parent_id, 
        d.nscan_id 
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
				rows.Scan(&r.ID, &r.Name, &r.Path, &nt, &size, &r.FileCount, &r.DirectoryCount, &r.ParentID, &r.NScanID)
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

	// /search?q=a&type=[image|doc|sheet|audio|video|zip|any]&after=1000-01-01&before=9999-12-31&page=1&per_page=50
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
