package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      string `json:"size"`
	FileCount int64  `json:"file_count"`
	ParentID  int64  `json:"parent_id"`
	NScanID   int64  `json:"nscan_id"`
}

var (
	db *sql.DB
)

func getDB() *sql.DB {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/nostalgia")
	if err != nil {
		panic(err.Error())
	}

	return db
}

func closeDB() {
	db.Close()
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
	getDB().QueryRow("SELECT count(id) from nfile").Scan(&count)
	defer closeDB()

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
			return fmt.Sprintf("%v MB", strconv.FormatFloat(r, 'f', 1, 64))
		}
		return fmt.Sprintf("%v kB", strconv.FormatFloat(r, 'f', 1, 64))
	}

	return fmt.Sprintf("%v", strconv.FormatFloat(r, 'f', 1, 64))
}

func search(c *gin.Context) {
	var nt mysql.NullTime
	var size int64
	query := c.DefaultQuery("q", "Default")

	// files
	var files []Nfile
	rows, err := getDB().Query(`SELECT n.id as ID, n.name, n.path, n.size, n.ndirectory_id as NDirectoryID, n.nscan_id as NScanID, n.date_modified as DateModified
								FROM nfile as n
								WHERE lower(n.name) like ?`,
		strings.ToLower(query)+"%")
	defer closeDB()

	if err != nil {
		files = nil
	} else {
		for rows.Next() {
			var r Nfile
			rows.Scan(&r.ID, &r.Name, &r.Path, &size, &r.NDirectoryID, &r.NScanID, &nt)
			r.Size = sizeString(size)
			if nt.Valid {
				r.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
			}

			files = append(files, r)
		}
	}

	// directories
	var directories []NDirectory
	rows, err = getDB().Query(`SELECT d.id as ID, d.name, d.size FROM ndirectory as d WHERE lower(d.name) like ?`,
		strings.ToLower(query)+"%")
	defer closeDB()

	if err != nil {
		directories = nil
	} else {
		for rows.Next() {
			var d NDirectory
			rows.Scan(&d.ID, &d.Name, &size)
			d.Size = sizeString(size)

			directories = append(directories, d)
		}
	}

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
		err := getDB().QueryRow("SELECT n.id as ID, n.name, n.path, n.size, n.ndirectory_id as NDirectoryID, n.nscan_id as NScanID, n.date_modified as DateModified FROM nfile as n WHERE n.id = ?", id).Scan(&r.ID, &r.Name, &r.Path, &size, &r.NDirectoryID, &r.NScanID, &nt)
		defer closeDB()
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
		err := getDB().QueryRow("SELECT d.name, d.parent_id FROM ndirectory as d WHERE d.id = ?", id).Scan(&dn, &dp)
		defer closeDB()
		if err != sql.ErrNoRows {
			rows, err := getDB().Query("SELECT n.id as ID, n.name, n.date_modified as DateModified, n.size FROM nfile as n WHERE n.ndirectory_id = ?", id)
			defer closeDB()

			if err != nil {
				files = nil
			} else {
				for rows.Next() {
					var r Nfile
					rows.Scan(&r.ID, &r.Name, &nt, &size)
					r.Size = sizeString(size)
					if nt.Valid {
						r.DateModified = fmt.Sprintf("%02d-%02d-%d", nt.Time.Day(), nt.Time.Month(), nt.Time.Year())
					}

					files = append(files, r)
				}
			}

			rows, err = getDB().Query("SELECT d.id as ID, d.name, d.size FROM ndirectory as d WHERE d.parent_id = ?", id)
			defer closeDB()

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
					rows.Scan(&d.ID, &d.Name, &size)
					d.Size = sizeString(size)

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
		log.Fatal(http.ListenAndServe(":80", nil))
	}()

	g.Run(":5000")
}
