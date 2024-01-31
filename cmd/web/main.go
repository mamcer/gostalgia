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

type Ndirectory struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	DateModified   string `json:"date_modified"`
	Size           string `json:"size"`
	FileCount      int64  // directory file count
	Fpath          string // current file path
	DirectoryCount int64  // directory directory count
	ParentID       int64  `json:"parent_id"`
}

type Nfile struct {
	ID           int64  // file id
	Name         string // file name
	Extension    string //file extension
	Path         string // file path
	DateModified string // file date modified
	Size         string // file size (in bytes)
	Hash         string // file hash
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
		page = 0
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

	//fmt.Printf("contains: '%v', type: '%v', only_directories: '%v', after: '%v', before: '%v', page: '%v', per_page: '%v'\n", contains, t, onlyDirs, after, before, page, perPage)

	// files
	var results []NResult
	var total int64 = 0
	var rtype = ""

	db := getDB()

	if onlyDirs == "true" {
		t = "any"
		rtype = "directory"
		db := getDB()
		rows, err := db.Query(`SELECT d.id as ID, d.name, d.path, d.date_modified as DateModified, d.size, d.parent_id, pd.name FROM ndirectory as d, ndirectory as pd WHERE lower(d.name) like ? and d.parent_id = pd.id and d.date_modified between ? and ? limit ? offset ?`,
			"%"+strings.ToLower(contains)+"%", after, before, perPage, (page-1)*perPage)
		defer db.Close()

		if err != nil {
			results = nil
		} else {
			for rows.Next() {
				var r NResult
				rows.Scan(&r.ID, &r.Name, &r.Path, &nt, &size, &r.ParentID, &r.ParentName)
				r.Size = sizeString(size)
				r.Type = rtype
				if nt.Valid {
					r.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
				}

				results = append(results, r)
			}
		}

		db.QueryRow("SELECT count(d.id) FROM ndirectory as d WHERE lower(d.name) like ? and d.date_modified between ? and ?", "%"+strings.ToLower(contains)+"%", after, before).Scan(&total)
	} else {
		var sq = `SELECT nf.id as ID, 
	nf.name,
	nf.extension,
	nf.path,
	nf.date_modified as DateModified,
	nf.size,
	nd.id,
	nd.name
	FROM nfile as nf, nfile_ndirectory as nfd, ndirectory as nd 
	WHERE lower(nf.name) like ? and nfd.nfile_id = nf.id and nfd.ndirectory_id = nd.id`

		var where = ""

		switch t {
		case "image":
			where += " and nf.extension in ('jpeg', 'png', 'jpg', 'bmp')"
		case "doc":
			where += " and nf.extension in ('doc', 'docx', 'odt', 'pdf')"
		case "sheet":
			where += " and nf.extension in ('xls', 'xlsx', 'ods')"
		case "audio":
			where += " and nf.extension in ('mp3', 'ogg', 'wma', 'arm', 'wav')"
		case "video":
			where += " and nf.extension in ('mp4', 'mkv', 'avi', 'wmv')"
		case "zip":
			where += " and nf.extension in ('zip', 'rar', '7z', 'gz')"
		case "any":
			where += " "
		default:
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
			if results == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "400",
					"title":  "invalid type",
					"detail": "valid types: image, doc, sheet, audio, video, zip",
				})
			}
			return
		}

		where += " and nf.date_modified between ? and ?"
		sq += where + " limit ? offset ?"

		//limit = per_page
		//offset = (page-1)*per_page

		rtype = "file"
		//fmt.Printf("query:'%v'\n", sq)
		rows, err := db.Query(sq,
			"%"+strings.ToLower(contains)+"%", after, before, perPage, (page-1)*perPage)
		defer db.Close()

		if err != nil {
			results = nil
		} else {
			for rows.Next() {
				var r NResult
				rows.Scan(&r.ID, &r.Name, &r.Extension, &r.Path, &nt, &size, &r.ParentID, &r.ParentName)
				r.Size = sizeString(size)
				r.Type = rtype
				if nt.Valid {
					r.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
				}

				results = append(results, r)
			}
		}

		// fmt.Printf("SELECT count(nf.id) FROM nfile as nf WHERE lower(nf.name) like ? " + where + "\n")

		db.QueryRow("SELECT count(nf.id) FROM nfile as nf, nfile_ndirectory as nfd, ndirectory as nd WHERE lower(nf.name) like ? and nfd.nfile_id = nf.id and nfd.ndirectory_id = nd.id "+where, "%"+strings.ToLower(contains)+"%", after, before).Scan(&total)
	}

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
			"total":    total,
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
		var d Ndirectory
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
			var directories []Ndirectory
			for rows.Next() {
				var r Ndirectory
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

	v1.GET("/files/count", filesCount)
	v1.OPTIONS("/files/count", preflight)

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
