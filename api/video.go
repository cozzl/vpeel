package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"vpeel/api/cors"
	"vpeel/internal/log"
	"vpeel/internal/trans"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type videInfoS struct {
	Author    string      `json:"author"`
	Title     string      `json:"title"`
	Src       string      `json:"src"`
	Resources []ResourceS `json:"resource"`
}

type ResourceS struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

const (
	uploadFolder        = "/Users/markov/Documents/code/go_code/vpeel/data/videos"
	remoteVideoDataPath = "http://127.0.0.1:8080/video/data/"
)

func getVideoList(c *gin.Context) {
	var videoList []videInfoS
	resp := basicResponse{Code: http.StatusOK, Meeeage: "success"}
	defer func() {
		if resp.Code != http.StatusOK {
			c.JSON(resp.Code, resp)
		} else {
			c.JSON(resp.Code, videoList)
		}
	}()
	files, err := os.ReadDir(uploadFolder)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Meeeage = err.Error()
		return
	}
	videoList = make([]videInfoS, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			log.Logger.Warnf("%s is not dir", file.Name())
			continue
		}
		videoInfo := videInfoS{}
		name := file.Name()
		nameSplit := strings.Split(name, "_")
		if len(nameSplit) < 2 {
			log.Logger.Warnf("%s format error, miss _", name)
			continue
		}
		videoInfo.Title = nameSplit[0]
		videoInfo.Author = nameSplit[1]

		subFiles, err := os.ReadDir(filepath.Join(uploadFolder, name))
		if err != nil {
			log.Logger.Warnf("read dir %s error: %v", name, err)
			continue
		}
		for _, subFile := range subFiles {
			videoName := subFile.Name()
			videoNameSplit := strings.Split(strings.Split(videoName, ".")[0], "_")
			path := remoteVideoDataPath + name + "/" + videoName
			if videoInfo.Src == "" {
				videoInfo.Src = path
			}
			if len(videoNameSplit) < 3 {
				log.Logger.Warnf("%s format error, miss _", name)
				// videoInfo.Src = remoteVideoDataPath + name + "/" + videoName
				continue
			}

			if videoNameSplit[len(videoNameSplit)-1] == "480p" {
				videoInfo.Src = path
			}
			resource := ResourceS{
				Url:  path,
				Name: videoNameSplit[len(videoNameSplit)-1],
			}

			videoInfo.Resources = append(videoInfo.Resources, resource)
		}
		videoList = append(videoList, videoInfo)
	}

}

type basicResponse struct {
	Code    int    `json:"code"`
	Meeeage string `json:"meeeage"`
	Status  int    `json:"-"`
}

var Templates = map[string]trans.TransParam{
	"tp480": {
		Vcodec:     "libx264",
		Acodec:     "copy",
		Width:      640,
		Height:     480,
		Resolution: "480p",
	},
	"tp720": {
		Vcodec:     "libx264",
		Acodec:     "copy",
		Width:      1280,
		Height:     720,
		Resolution: "720p",
	},
	"tp1080": {
		Vcodec:     "libx264",
		Acodec:     "copy",
		Width:      1920,
		Height:     1080,
		Resolution: "1080p",
	},
}

func uploadVideo(c *gin.Context) {
	resp := basicResponse{Code: http.StatusOK, Meeeage: "success"}
	defer func() {
		c.JSON(resp.Code, resp)
	}()
	file, err := c.FormFile("file")
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Meeeage = err.Error()
		return
	}

	filename := filepath.Base(file.Filename)
	dir := strings.Split(filename, ".")[0]

	dstDir := filepath.Join(uploadFolder, dir)
	if err := os.Mkdir(dstDir, 0777); err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Meeeage = err.Error()
		return
	}

	dstFile := filepath.Join(dstDir, filename)
	if err := c.SaveUploadedFile(file, dstFile); err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Meeeage = err.Error()
		return
	}

	for name, tp := range Templates {
		task := &trans.TranscodeTask{
			ID:         uuid.New().String(),
			InputFile:  dstFile,
			OutputFile: filepath.Join(dstDir, dir+"_"+tp.Resolution+filepath.Ext(filename)),
			Param:      tp,
		}
		log.Logger.Debugf("template:%s, inputFile:%s, outputFile:%s, task:%s",
			name, task.InputFile, task.OutputFile, task.ID)
		trans.DefaultManager.Submit(task)
	}
}

func videoRouterInit(r *gin.Engine) {
	video := r.Group("/video").Use(cors.Cors())
	{
		video.GET("/list", getVideoList)
		video.POST("/upload", uploadVideo)
		video.Static("/data", uploadFolder)
	}

}

func init() {
	AddRouter(videoRouterInit)
}
