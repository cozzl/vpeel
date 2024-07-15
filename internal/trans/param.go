package trans

import (
	"reflect"
	"fmt"
) 

type TransParam struct {
	Vcodec     string `json:"vcodec,omitempty"`
	Acodec     string `json:"acodec,omitempty"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	Resolution string `json:"resolution,omitempty"`
	Fps        int    `json:"fps,omitempty"`
	Bitrate    int    `json:"bitrate,omitempty"`
	Gop        int    `json:"gop,omitempty"`
	Bframes    int    `json:"bframes,omitempty"`
	Filter     string `json:"filter,omitempty"`
	Thread     int    `json:"thread,omitempty"`
	CodecParam string `json:"codec_param,omitempty"`
	Profile    string `json:"profile,omitempty"`
	Preset     string `json:"preset,omitempty"`
}

func (p *TransParam) ToFFmpegArgs(inputFile, outputFile string) []string {
	args := []string{"-i", inputFile}

	appendArg := func(flag string, value interface{}) {
		v := reflect.ValueOf(value)
		// 检查 value 是否为空字符串或零值
		if (v.Kind() == reflect.String && v.String() == "") ||
			(v.Kind() >= reflect.Int && v.Kind() <= reflect.Float64 && v.IsZero()) {
			return
		}
		args = append(args, flag, fmt.Sprintf("%v", value))
	}

	appendArg("-vcodec", p.Vcodec)
	appendArg("-acodec", p.Acodec)
	if p.Width > 0 && p.Height > 0 {
		resolution := fmt.Sprintf("%d:%d", p.Width, p.Height)
		if p.Filter != "" {
			p.Filter = fmt.Sprintf("scale=%s,%s", resolution, p.Filter)
		} else {
			p.Filter = fmt.Sprintf("scale=%s", resolution)
		}
	}
	appendArg("-r", p.Fps)
	appendArg("-b:v", p.Bitrate)
	appendArg("-g", p.Gop)
	appendArg("-bf", p.Bframes)
	appendArg("-filter_complex", p.Filter)
	appendArg("-threads", p.Thread)
	if p.CodecParam != "" {
		switch p.Vcodec {
		case "libx264":
			args = append(args, "-x264-params", p.CodecParam)
		case "libx265":
			args = append(args, "-x265-params", p.CodecParam)
		case "libksc265":
			args = append(args, "-ksc265-params", p.CodecParam)
		}
	}
	appendArg("-profile:v", p.Profile)
	appendArg("-preset", p.Preset)

	args = append(args, outputFile)
	return args
}
