package tester_pool

import (
	"contrplatform/configs"
	"contrplatform/pkg/util"
	"math/rand"
	"os"
	"time"
)


type TesterPool struct {
	testerMap map[string]*Tester
}

func New() *TesterPool {
	return &TesterPool{
		testerMap: make(map[string]*Tester),
	}
}

func (t *TesterPool) GeneID() string {
	rand.Seed(time.Now().Unix())
	id := util.Encode32(time.Now().Unix())+
		util.Encode32(rand.Int63()>>32)
	t.testerMap[id]=&Tester{
		id: id,
		state: StateInit,
	}
	return id
}

func (t *TesterPool) State(id string) TesterState {
	tester,ok := t.testerMap[id]
	if !ok {
		return StateNull
	}
	return tester.State()
}



func (t *TesterPool) StartCmd(id, runTime string) error {
	return t.testerMap[id].startCmd(runTime)
}

func (t *TesterPool) StopCmd(id string) {
	t.testerMap[id].StopCmd()
}

func (t *TesterPool) Reset(id string) error {
	if err:=os.RemoveAll(configs.OyenteOutputPath+"/"+id);err!=nil{
		return err
	}
	if err:=os.RemoveAll(configs.SfuzzOutputPath+"/"+id);err!=nil{
		return err
	}
	if err:=os.RemoveAll(configs.UploadSavePath+"/"+id);err!=nil{
		return err
	}
	delete(t.testerMap,id)
	return nil
}
