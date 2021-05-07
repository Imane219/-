package detectorpool

import (
	"contrplatform/configs"
	"contrplatform/pkg/upload"
	"contrplatform/pkg/util"
	"log"
	"math/rand"
	"mime/multipart"
	"sync"
	"time"
)

type DetectorPool struct {
	detectorMap sync.Map
	timerMap sync.Map
}

func New() *DetectorPool {
	return &DetectorPool{}
}

func (dp *DetectorPool) UploadContracts(fileHeaders []*multipart.FileHeader,id string) error {
	tester,_ := dp.detector(id)
	for _, fileHeader := range fileHeaders {
		if err:=tester.checkUploadState();err!=nil {
			return err
		}
		if err:=upload.FileUpload(fileHeader,upload.TypeContract,id);err!=nil{
			return err
		}
	}
	return nil
}

func (dp *DetectorPool) GeneID() string {
	rand.Seed(time.Now().Unix())
	id := util.Encode32(time.Now().Unix()) +
		util.Encode32(rand.Int63()>>32)
	dt := newDetector(id)
	dp.detectorMap.Store(id,dt)

	go func() {
		timer:=time.NewTimer(configs.DetectorExpiredTime)
		dp.timerMap.Store(id,timer)
		<-timer.C
		_ = dt.delete()
		dp.detectorMap.Delete(id)
		dp.timerMap.Delete(id)
		log.Print("delete")
	}()
	return id
}

func (dp *DetectorPool) resetTimer(id string,duration time.Duration) {
	timer,_ := dp.timerMap.Load(id)
	timer.(*time.Timer).Reset(duration)
}


func (dp *DetectorPool) DetectorState(id string) DetectorState {
	dt, ok := dp.detectorMap.Load(id)
	if !ok {
		return StateNull
	}
	return dt.(*detector).State()
}

func (dp *DetectorPool) detector(id string) (*detector, error) {
	dt,ok:= dp.detectorMap.Load(id)
	if !ok {
		return nil,newStateError(StateNull,StateInit)
	}
	return dt.(*detector),nil
}

func (dp *DetectorPool) StartCmd(id, runTime string) error {
	dt,err := dp.detector(id)
	if err!=nil{
		return err
	}
	if err = dt.startCmd(id,runTime); err != nil {
		return err
	}
	dp.resetTimer(id,configs.DetectorExpiredTime)
	return nil
}

func (dp *DetectorPool) StopCmd(id string) error {
	dt,err := dp.detector(id)
	if err != nil {
		return err
	}
	if err = dt.stopCmd(); err != nil {
		return err
	}
	dp.resetTimer(id,configs.DetectorExpiredTime)
	return nil
}

func (dp *DetectorPool) Reset(id string) error {
	dt, err := dp.detector(id)
	if err != nil {
		return err
	}
	if err = dt.checkResetDetectionState(); err != nil {
		return err
	}
	if err = dt.delete(); err != nil {
		return err
	}
	dp.resetTimer(id,0)
	return nil
}

func (dp *DetectorPool) Delete() {
	dp.timerMap.Range(func(key, value interface{}) bool {
		return value.(*time.Timer).Reset(0)
	})
	time.Sleep(500*time.Millisecond)
}