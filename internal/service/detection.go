package service

import (
	"contrplatform/internal/detectorpool"
	"strconv"
)

type StartDetectionRequest struct {
	ID      string `form:"id" binding:"required"`
	RunTime int    `form:"run_time,default=120" binding:"min=6,max=6400"`
}

func (svc *Service) StartDetection(param *StartDetectionRequest) error {
	return svc.pool.StartCmd(param.ID, strconv.Itoa(param.RunTime))
}

type GetResultRequest struct {
	ID string `form:"id" binding:"required"`
}

func (svc *Service) GetResult(param *GetResultRequest) ([]*detectorpool.Output, error) {
	outputs := make(map[string]*detectorpool.Output)
	if err := svc.pool.GetOyenteOutputs(param.ID, outputs); err != nil {
		return nil, err
	}
	if err := svc.pool.GetSfuzzOutputs(param.ID, outputs); err != nil {
		return nil, err
	}
	results := make([]*detectorpool.Output, 0, len(outputs))
	for _, output := range outputs {
		results = append(results, output)
	}
	return results, nil
}

type StopDetectionRequest struct {
	ID string `form:"id" binding:"required"`
}

func (svc *Service) StopDetection(param *StopDetectionRequest) error {
	return svc.pool.StopCmd(param.ID)
}

type ResetDetectionRequest struct {
	ID string `form:"id" binding:"required"`
}

func (svc *Service) ResetDetection(param *ResetDetectionRequest) error {
	return svc.pool.Reset(param.ID)
}
