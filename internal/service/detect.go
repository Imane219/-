package service

import (
	"errors"
	"os/exec"
	"strconv"
)

type DetectRequest struct {
	RunTime int	`form:"run_time,default=120" binding:"max=6400"`
}

func (svc *Service) detect(param *DetectRequest) error {
	if svc.GetSessionID() == nil {
		return errors.New("缺少上传文件信息")
	}
	cmd := exec.Command("run_detect",strconv.Itoa(param.RunTime))
	if err:= cmd.Start(); err != nil {
		return err
	}
	return nil
}
