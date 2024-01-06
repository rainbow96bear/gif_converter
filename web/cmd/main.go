package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	// 정적 파일 서빙을 위한 미들웨어 추가
	r.Static("/converted", "./converted")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		// 업로드된 동영상을 저장
		videoPath := filepath.Join("uploads", file.Filename)
		if err := c.SaveUploadedFile(file, videoPath); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		// GIF로 변환
		gifPath := convertToGIF(videoPath, file.Filename)

		// 변환 결과를 웹 페이지에 표시
		c.HTML(http.StatusOK, "result.html", gin.H{
			"VideoPath": videoPath,
			"GIFPath":   "/converted/"+gifPath,
		})
	})

	r.Run(":8080")
}

func convertToGIF(videoPath, outputName string) string {
	currentTime := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("rainbowbear_%s", currentTime+".gif")
	outputPath := filepath.Join("converted", filename)
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=10,scale=320:-1", outputPath)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return filename
}
