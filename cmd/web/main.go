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
	NDirectoryID int64  `json:"ndirectory_id"`
	NScanID      int64  `json:"nscan_id"`
}

type NDirectory struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	DateModified string `json:"date_modified"`
	Size         string `json:"size"`
	FileCount    int64  `json:"file_count"`
	ParentID     int64  `json:"parent_id"`
	NScanID      int64  `json:"nscan_id"`
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
	if query == "" {
		query = "mario"
	}

	t := c.DefaultQuery("type", "any")

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
							n.Extension, 
							n.path, 
							n.size, 
							n.ndirectory_id as NDirectoryID, 
							n.nscan_id as NScanID, 
							n.date_modified as DateModified 
				FROM nfile as n 
				WHERE lower(n.name) like ?`
	switch t {
	case "image":
		sq += " and n.extension in ('jpeg', 'png', 'jpg', 'bmp')"
	case "doc":
		sq += " and n.extension in ('doc', 'docx', 'odt')"
	case "sheet":
		sq += " and n.extension in ('xls', 'xlsx', 'ods')"
	case "audio":
		sq += " and n.extension in ('mp3', 'ogg', 'wma', 'arm')"
	case "video":
		sq += " and n.extension in ('mp4', 'mkv', 'avi', 'wmv')"
	case "zip":
		sq += " and n.extension in ('zip', 'rar', '7z', 'gz')"
	}

	sq += " and n.date_modified between ? and ?"

	rows, err := db.Query(sq,
		"%"+strings.ToLower(query)+"%", after, before)
	defer db.Close()

	if err != nil {
		files = nil
	} else {
		for rows.Next() {
			var r Nfile
			rows.Scan(&r.ID, &r.Name, &r.Extension, &r.Path, &size, &r.NDirectoryID, &r.NScanID, &nt)
			r.Size = sizeString(size)
			if nt.Valid {
				r.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
			}

			files = append(files, r)
		}
	}

	// directories
	var directories []NDirectory
	// db2 := getDB()
	// rows, err = db2.Query(`SELECT d.id as ID, d.name, d.date_modified as DateModified, d.size FROM ndirectory as d WHERE lower(d.name) like ?`,
	// 	strings.ToLower(query)+"%")
	// defer db2.Close()

	// if err != nil {
	// 	directories = nil
	// } else {
	// 	for rows.Next() {
	// 		var d NDirectory
	// 		rows.Scan(&d.ID, &d.Name, &nt, &size)
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
	if files == nil && directories == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"query":       query,
			"directories": nil,
			"files":       nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"query":       query,
			"directories": directories,
			"files":       files,
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
		err := db.QueryRow("SELECT n.id as ID, n.name, n.path, n.size, n.ndirectory_id as NDirectoryID, n.nscan_id as NScanID, n.date_modified as DateModified FROM nfile as n WHERE n.id = ?", id).Scan(&r.ID, &r.Name, &r.Path, &size, &r.NDirectoryID, &r.NScanID, &nt)
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
	var files []Nfile
	var directories []NDirectory
	var dn string
	var dp int64
	if id != "" {
		db1 := getDB()
		err := db1.QueryRow("SELECT d.name, d.parent_id FROM ndirectory as d WHERE d.id = ?", id).Scan(&dn, &dp)
		defer db1.Close()
		if err != sql.ErrNoRows {
			db2 := getDB()
			rows, err := db2.Query("SELECT n.id as ID, n.name, n.path, n.date_modified as DateModified, n.size FROM nfile as n WHERE n.ndirectory_id = ?", id)
			defer db2.Close()

			if err != nil {
				files = nil
			} else {
				for rows.Next() {
					var r Nfile
					rows.Scan(&r.ID, &r.Name, &r.Path, &nt, &size)
					r.Size = sizeString(size)
					if nt.Valid {
						r.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
					}

					files = append(files, r)
				}
			}

			db3 := getDB()
			rows, err = db3.Query("SELECT d.id as ID, d.name, d.date_modified as DateModified, d.size FROM ndirectory as d WHERE d.parent_id = ?", id)
			defer db3.Close()

			if err != nil {
				c.Header("Access-Control-Allow-Origin", "*")
				c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
				c.JSON(http.StatusInternalServerError, gin.H{
					"name":        dn,
					"parent_id":   dp,
					"directories": nil,
					"files":       files,
				})
			} else {
				for rows.Next() {
					var d NDirectory
					rows.Scan(&d.ID, &d.Name, &nt, &size)
					d.Size = sizeString(size)
					if nt.Valid {
						d.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
					}

					directories = append(directories, d)
				}

				c.Header("Access-Control-Allow-Origin", "*")
				c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
				c.JSON(http.StatusOK, gin.H{
					"name":        dn,
					"parent_id":   dp,
					"directories": directories,
					"files":       files,
				})
			}

		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
			c.JSON(http.StatusNotFound, gin.H{
				"name":        "unknown",
				"parent_id":   1,
				"directories": nil,
				"files":       nil,
			})
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

	g := gin.Default()

	g.GET("/ping", ping)
	g.OPTIONS("/ping", preflight)

	g.GET("/search", search)
	g.OPTIONS("/search", preflight)

	g.GET("/files/:id", fileController)
	g.OPTIONS("/files/:id", preflight)

	g.GET("/filescount", filesCount)
	g.OPTIONS("/filescount", preflight)

	g.GET("/directories/:id", directoriesController)
	g.OPTIONS("/directories/:id", preflight)

	go func() {
		http.Handle("/",
			http.StripPrefix("/",
				http.FileServer(http.Dir("./"))))
		log.Fatal(http.ListenAndServe(":"+config.WebPort, nil))
	}()

	g.Run(":" + config.ApiPort)
}
