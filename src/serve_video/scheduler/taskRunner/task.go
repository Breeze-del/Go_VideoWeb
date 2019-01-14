package taskRunner

import (
	"errors"
	"log"
	"os"
	"serve_video/scheduler/dbops"
	"sync"
)

// 延时删除视频

func VideoClearDispatcher(dc dataChan) error {
	res, err := dbops.ReadVideoDeletionRecord(3)
	if err != nil {
		log.Printf("Viedo clear dispathcer error :%s", err)
	}

	if len(res) == 0 {
		return errors.New("All tasks finished")
	}
	for _, id := range res {
		dc <- id
	}
	return nil
}

func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var wait sync.WaitGroup
	var err error
forloop:
	for {
		select {
		case vid := <-dc:
			wait.Add(1)
			go func(id interface{}) {
				defer wait.Done()
				if err := deleteVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
				if err := dbops.DelVideoDeletionRecoed(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)
		default:
			break forloop
		}
	}
	wait.Wait()
	errMap.Range(func(key, value interface{}) bool {
		err = value.(error)
		if err != nil {
			return false
		}
		return true
	})
	return err
}

// 删除viedo文件
func deleteVideo(vid string) error {
	err := os.Remove(VIDEO_PATH + vid)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Delete video err %s", err)
		return err
	} else {
		return nil
	}
}
