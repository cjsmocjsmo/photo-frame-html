package main

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var GLOBAL_IDX int = 0
var GLOBAL_COUNT int = 0

func setEnvVariable() {
	err := os.Setenv("PFPICPATH", "/home/pipi/Pictures/MasterPicsResize_SPLIT/")
	if err != nil {
		log.Fatal("Failed to set environment variable:", err)
	}
	err2 := os.Setenv("PFDBPATH", "/home/pipi/photo-frame-html/picinfo.db")
	if err2 != nil {
		log.Fatal("Failed to set environment variable:", err2)
	}
}

func checkAndCreateDB() int {
	dbPath := os.Getenv("PFDBPATH")

	// Check if PFDBPATH exists
	_, err := os.Stat(dbPath)
	if err == nil {
		//connect to the db
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal("Failed to open database:", err)
		}
		defer db.Close()

		var lastEntry int
		err = db.QueryRow("SELECT MAX(pfid) FROM picinfo").Scan(&lastEntry)
		if err != nil {
			log.Println("Failed to get last entry:", err)
		}

		var lastPfidx int
		err = db.QueryRow("SELECT pfidx FROM picinfo WHERE pfid = ?", lastEntry).Scan(&lastPfidx)
		if err != nil {
			log.Println("Failed to get last pfidx:", err)
		}

		return lastPfidx

	}

	// Create new PFDBPATH
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to create database:", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS picinfo (
		pfid INTEGER PRIMARY KEY AUTOINCREMENT,
		pfidx INTEGER NOT NULL,
		pfpath TEXT NOT NULL UNIQUE,
		pfhttp TEXT NOT NULL UNIQUE
	)`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	return 0
}

func ScanForImages() []string {
	rootFolder := os.Getenv("PFPICPATH")
	file_list := []string{}

	err := filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		println(path)
		ext := filepath.Ext(path)
		if ext == ".jpg" {
			file_list = append(file_list, path)
		}
		return nil
	})
	if err != nil {
		log.Println("Failed to walk directory:", err)
	}
	return file_list

}

func InsertAllImages(count int, images []string) {
	var wg sync.WaitGroup
	for _, image := range images {
		wg.Add(1)
		count += 1
		go InsertImage(count, image, &wg) // Pass the address of wg instead of dereferencing it
	}
	wg.Wait()
}

func InsertImage(count int, path string, wg *sync.WaitGroup) error {
	startTime := time.Now()
	pfidx := count
	pfpath := path

	_, file := filepath.Split(pfpath)
	pfhttp := "/media/" + file

	dbPath := os.Getenv("PFDBPATH")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM picinfo WHERE pfpath = ?", pfpath).Scan(&exists)
	if err != nil {
		return err
	}

	if exists == 0 {
		_, err = db.Exec(`INSERT INTO picinfo (pfidx, pfpath, pfhttp) VALUES (?, ?, ?)`, pfidx, pfpath, pfhttp)
		if err != nil {
			return err
		}
	}
	endTime := time.Now()
	endTime2 := endTime.Sub(startTime)
	EndTime := endTime2.Seconds()
	log.Printf("Inserted %s \nin %f seconds\n", pfhttp, EndTime)
	wg.Done()
	return nil
}

func SetGlobalCount() {
	dbPath := os.Getenv("PFDBPATH")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		GLOBAL_COUNT = 0
	}
	defer db.Close()

	var largestPfidx int
	err = db.QueryRow("SELECT MAX(pfidx) FROM picinfo").Scan(&largestPfidx)
	if err != nil {
		GLOBAL_COUNT = 0
	}

	GLOBAL_COUNT = largestPfidx
}

func init() {
	setEnvVariable()
	count := checkAndCreateDB()
	println(count)
	images := ScanForImages()
	println(images)
	// InsertAllImages(count, images)
	// SetGlobalCount()
}

type Template struct {
	templates *template.Template
}

type TemplateData struct {
	PFPath string
}

func main() {
	picpath := os.Getenv("PFPICPATH")
	println(picpath)
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = t
	e.GET("/", Index)
	// e.GET("/update", UpdateImage)
	e.GET("/fuck", Fuck)

	e.Static("/assets", "assets")
	e.Static("/media", picpath)
	e.Start(":8080")
}

func Fuck(c echo.Context) error {
	return c.Render(http.StatusOK, "fuck", "Worked")
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Index(c echo.Context) error {
	// picpath := os.Getenv("PFPICPATH")
// 	data := "/media/04728da2-0831-48ad-b75b-8b3cab8f9269.jpg"

// 	log.Println("hit index")
// 	return c.Render(http.StatusOK, "index", data)
// }

// func UpdateImage(c echo.Context) error {
	GLOBAL_IDX += 1
	if GLOBAL_IDX != GLOBAL_COUNT {

		dbPath := os.Getenv("PFDBPATH")
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal("Failed to open database:", err)
		}
		defer db.Close()

		var pfhttp string
		err = db.QueryRow("SELECT pfhttp FROM picinfo WHERE pfidx = ?", GLOBAL_IDX).Scan(&pfhttp)
		if err != nil {
			return err
		}
		// return c.JSON(http.StatusOK, pfhttp)
		return c.Render(http.StatusOK, "index", pfhttp)
	} else {
		GLOBAL_IDX = 1
		dbPath := os.Getenv("PFDBPATH")
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal("Failed to open database:", err)
		}
		defer db.Close()

		var pfhttp string
		err = db.QueryRow("SELECT pfhttp FROM picinfo WHERE pfidx = ?", GLOBAL_IDX).Scan(&pfhttp)
		if err != nil {
			return err
		}
		// return c.JSON(http.StatusOK, pfhttp)
		return c.Render(http.StatusOK, "index", pfhttp)
	}
}

