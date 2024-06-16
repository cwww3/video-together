package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type FileInfo struct {
	Rid         string  `json:"rid"`
	Duration    float64 `json:"duration,string"`
	M3u8Content string  `json:"m3u8Content"`
	M3u8Url     string  `json:"m3u8Url"`
	CurrentTime float64 `json:"currentTime,string"`
	UnixMill    int64   `json:"unixMill,string"`
}

var m = make(map[string]FileInfo)
var rwLock sync.RWMutex

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/file", func(c *gin.Context) {
		var f FileInfo
		err := c.Bind(&f)
		if err != nil {
			fmt.Println("bind err", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err,
			})
			return
		}
		rwLock.RLock()
		_, ok := m[f.Rid]
		rwLock.RUnlock()
		if ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "exist",
			})
			return
		}

		rwLock.Lock()
		m[f.Rid] = f
		rwLock.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.GET("/progress/:rid", func(c *gin.Context) {
		rid := c.Param("rid")
		if rid == "test" {
			c.JSON(http.StatusOK, gin.H{
				"currentTime": fmt.Sprintf("%f", rand.Intn(15*60)),
			})
			return
		}

		d := c.Query("d")
		var diff float64
		if len(d) > 0 {
			diff, _ = strconv.ParseFloat(d, 64)
		}

		rwLock.RLock()
		f, ok := m[rid]
		rwLock.RUnlock()

		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not exist",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"currentTime": fmt.Sprintf("%f", f.CurrentTime+diff),
			"unixMill":    fmt.Sprintf("%d", f.UnixMill),
		})
	})
	r.PUT("/progress/:rid", func(c *gin.Context) {
		rid := c.Param("rid")

		rwLock.RLock()
		f, ok := m[rid]
		rwLock.RUnlock()

		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not exist",
			})
			return
		}

		var arg FileInfo
		err := c.Bind(&arg)
		if err != nil {
			fmt.Println("bind err", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		f.CurrentTime = arg.CurrentTime
		f.UnixMill = arg.UnixMill

		rwLock.Lock()
		m[rid] = f
		rwLock.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.LoadHTMLFiles("index.tmpl", "error.tmpl")
	r.GET("/file/:rid", func(c *gin.Context) {
		rid := c.Param("rid")

		if rid == "test" {
			c.HTML(200, "index.tmpl", gin.H{
				"url": "https://svip.high23-playback.com/20240602/13986_a397e1ae/index.m3u8",
				"rid": rid,
				"api": "http://localhost:8080",
			})
			return
		}

		rwLock.RLock()
		f, ok := m[rid]
		rwLock.RUnlock()

		if !ok {
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
				"message": "not exsit",
			})
			return
		}

		// 渲染模板并返回
		c.HTML(200, "index.tmpl", gin.H{
			"url":  f.M3u8Url,
			"rid":  rid,
			"api":  "http://localhost:8080",
			"diff": c.Query("diff"),
		})
	})
	r.Run(":8080")
}
