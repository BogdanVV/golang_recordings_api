package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

func main() {
	var db *sql.DB

	// DB SETUP
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected")

	// API
	r := gin.Default()

	r.GET("/albums", func(c *gin.Context) {
		getAlbums(c, db)
	})
	r.GET("/albums/:id", func(c *gin.Context) {
		getAlbumById(c, db)
	})
	r.POST("/albums", func(c *gin.Context) {
		postAlbum(c, db)
	})
	r.DELETE("/albums/:id", func(c *gin.Context) {
		deleteAlbum(c, db)
	})
	r.PUT("/albums/:id", func(c *gin.Context) {
		updateAlbum(c, db)
	})

	r.Run("localhost:8080")
}

func updateAlbum(c *gin.Context, db *sql.DB) {
	var albumResult Album
	if err := c.BindJSON(&albumResult); err != nil {
		fmt.Println("err>>>", err)
	}
	fmt.Println("albumResult>>>", albumResult)

	if albumResult.Artist != "" {

	}

	// ====

	// id := c.Param("id")
	// var alb Album

	// row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	// if err := row.Scan(&alb.Artist, &alb.Price, &alb.Title, &alb.ID); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		fmt.Println("album not found")
	// 	}
	// 	fmt.Println("error>>>", err)
	// }
	// fmt.Println("alb>>>", alb.Artist, alb.ID, alb.Price, alb.Title)

	// v := reflect.ValueOf(alb)

	// fmt.Println("v>>>", v.Type())

	// values := make([]interface{}, v.NumField())

	// for i := 0; i < v.NumField(); i++ {
	// 	fmt.Println(v.Field(i))
	// 	values[i] = v.Field(i).Interface()
	// }

	// fmt.Println("values>>>", values)
}

func getAlbums(c *gin.Context, db *sql.DB) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			log.Fatal(err)
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumById(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)

	if err := row.Scan(&alb.ID, &alb.Artist, &alb.Title, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("getAlbumById %v: %v\n", id, err)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "The album was not found!"})
			return
		}
		fmt.Printf("getAlbumById %v: %v\n", id, err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Failed to read the album that was found"})
		return
	}
	c.IndentedJSON(http.StatusOK, alb)
}

func postAlbum(c *gin.Context, db *sql.DB) {
	var newAlbum Album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad fields"})
		return
	}

	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "Failed to insert album into db"})
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Failed to create ID for new identity"})
	}
	newAlbum.ID = id

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func deleteAlbum(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)

	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("no such album in db")
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no such album in db"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	_, err := db.Exec("DELETE FROM album WHERE ID = ?", id)
	if err != nil {
		fmt.Println("err", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "couldn't perform deletion of the album"})
		return
	}

	c.Status(200)
}

// ====================

// func getAlbumsByArtist(artist string) ([]Album, error) {
// 	var albums []Album

// 	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", artist)

// 	if err != nil {
// 		return nil, fmt.Errorf("albumsByArtist %q: %v", artist, err)
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		var alb Album

// 		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
// 			return nil, fmt.Errorf("albumsByArtist %q: %v", artist, err)
// 		}

// 		albums = append(albums, alb)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("albumsByArtist %q: %v\n", artist, err)
// 	}

// 	return albums, nil
// }

// TODO:
// $ export DBUSER=root
// $ export DBPASS=***

// TODO: move functions related to db to another file

// TODO: create Docker container with mysql

// TODO: add README

// TODO: add another CRUD operations
