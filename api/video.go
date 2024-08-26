package api

import (
	"encoding/json"
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

type basicResponseS struct {
	Code    int    `json:"code"`
	Meeeage string `json:"meeeage"`
	Status  int    `json:"-"`
}

type deleteVideoReqS struct {
	Name string `json:"name"`
}

const (
	uploadFolder        = "/Users/markov/Documents/code/go_code/vpeel/data/videos"
	remoteVideoDataPath = "http://127.0.0.1:8080/video/data/"
)

var DefaultTemplates = map[string]trans.TransParam{
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

func getVideoList(c *gin.Context) {
	var videoList []videInfoS
	resp := basicResponseS{Code: http.StatusOK, Meeeage: "success"}
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
		videoInfo.Title = name
		// TODO: 如何将作者传进来
		videoInfo.Author = "unkown"

		subFiles, err := os.ReadDir(filepath.Join(uploadFolder, name))
		if err != nil {
			log.Logger.Warnf("read dir %s error: %v", name, err)
			continue
		}
		for _, subFile := range subFiles {
			videoName := subFile.Name()
			videoNamePre := strings.Split(videoName, ".")[0]
			path := remoteVideoDataPath + name + "/" + videoName
			if videoInfo.Src == "" {
				videoInfo.Src = path
			}
			resouceName := ""
			if videoNamePre == name {
				resouceName = "原画"
			} else {
				videoNameSplit := strings.Split(videoNamePre, "_")
				if len(videoNameSplit) < 2 {
					log.Logger.Warnf("video name split error: %s", videoNamePre)
					continue
				}
				resouceName = videoNameSplit[len(videoNameSplit)-1]
			}

			resource := ResourceS{
				Url:  path,
				Name: resouceName,
			}

			videoInfo.Resources = append(videoInfo.Resources, resource)
		}
		videoList = append(videoList, videoInfo)
	}

}

func deletetVideo(c *gin.Context) {
	var deleteVideoReq deleteVideoReqS
	resp := basicResponseS{Code: http.StatusOK, Meeeage: "success"}
	defer func() {
		c.JSON(resp.Code, resp)
	}()

	if err := c.ShouldBindBodyWithJSON(&deleteVideoReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Meeeage = err.Error()
		return
	}

	dir := filepath.Join(uploadFolder, deleteVideoReq.Name)
	if err := os.RemoveAll(dir); err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Meeeage = err.Error()
		return
	}
}

func uploadAndTransVideo(c *gin.Context) {
	resp := basicResponseS{Code: http.StatusOK, Meeeage: "success"}
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

	// 获取并解析 params 参数
	paramsJson := c.PostForm("params")
	var transParams map[string]trans.TransParam
	err = json.Unmarshal([]byte(paramsJson), &transParams)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Meeeage = err.Error()
		return
	}

	if len(transParams) == 0 {
		transParams = DefaultTemplates
	}

	for name, tp := range transParams {
		task := &trans.TranscodeTask{
			ID:         uuid.New().String(),
			InputFile:  dstFile,
			OutputFile: filepath.Join(dstDir, dir+"_"+name+filepath.Ext(filename)),
			Param:      tp,
		}
		log.Logger.Debugf("template:%s, inputFile:%s, outputFile:%s, task:%s, param:%+v",
			name, task.InputFile, task.OutputFile, task.ID, task.Param)
		trans.DefaultManager.Submit(task)
	}
}

func videoRouterInit(r *gin.Engine) {
	video := r.Group("/video").Use(cors.Cors())
	{
		video.OPTIONS("/*any", nil)
		video.GET("/list", getVideoList)
		video.Static("/data", uploadFolder)
		video.POST("/delete", deletetVideo)
		video.POST("/uploadAndTrans", uploadAndTransVideo)
	}

}

func init() {
	AddRouter(videoRouterInit)
}
