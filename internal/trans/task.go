package trans

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"sync"
	"vpeel/internal/log"
)

var DefaultManager = NewTranscodeManager(3)

var ffmpegBin = "/Users/markov/Documents/code/go_code/vpeel/tool/ffmpeg/ffmpeg"

// 转码任务结构体
type TranscodeTask struct {
	ID         string
	InputFile  string
	OutputFile string
	Param      TransParam
	args       []string
}

// 转码管理器
type TranscodeManager struct {
	tasks       chan *TranscodeTask
	results     chan *TranscodeResult
	workerCount int
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// 转码结果结构体
type TranscodeResult struct {
	TaskID string
	Error  error
}

// 创建新的转码管理器
func NewTranscodeManager(workerCount int) *TranscodeManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &TranscodeManager{
		tasks:       make(chan *TranscodeTask, 1000),
		results:     make(chan *TranscodeResult, 1000),
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// 启动工作池
func (tm *TranscodeManager) Start() {
	log.Logger.Debugf("start trans")
	for i := 1; i <= tm.workerCount; i++ {
		tm.wg.Add(1)
		go tm.worker(i)
	}
}

// 停止工作池
func (tm *TranscodeManager) Stop() {
	log.Logger.Debugf("stop trans")
	tm.cancel()
	tm.wg.Wait()
	close(tm.results)
}

// 提交任务
func (tm *TranscodeManager) Submit(task *TranscodeTask) {
	log.Logger.Debugf("submit trans task: %s", task.ID)
	tm.tasks <- task
}

func (tm *TranscodeManager) Result() {
	// 处理结果
	for result := range tm.results {
		if result.Error != nil {
			log.Logger.Infof("task %s failed: %v\n", result.TaskID, result.Error)
		} else {
			log.Logger.Infof("task %s succeeded\n", result.TaskID)
		}
	}
}

// 处理任务的 worker
func (tm *TranscodeManager) worker(id int) {
	log.Logger.Debugf("trans worker %d runing", id)
	defer tm.wg.Done()
	for {
		select {
		case task := <-tm.tasks:
			log.Logger.Infof("worker %d starting task: %s\n", id, task.ID)
			err := tm.transcode(task)
			tm.results <- &TranscodeResult{TaskID: task.ID, Error: err}
			log.Logger.Infof("worker %d finished task: %s\n", id, task.ID)
		case <-tm.ctx.Done():
			return
		}
	}
}

// 执行转码任务
func (tm *TranscodeManager) transcode(task *TranscodeTask) error {
	task.args = task.Param.ToFFmpegArgs(task.InputFile, task.OutputFile)
	cmd := exec.CommandContext(tm.ctx, ffmpegBin, task.args...)
	log.Logger.Debugf("transcode trans task: %s, cmd: %s", task.ID, cmd.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("transcode failed: %w, output: %s", err, out.String())
	}
	return nil
}

