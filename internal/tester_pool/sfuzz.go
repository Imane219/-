package tester_pool

import (
	"contrplatform/configs"
	"contrplatform/pkg/line_reader"
	"io/ioutil"
	"os"
	"strings"
)

type SfuzzOutput struct {
	Duration     string `json:"duration"`	//运行时长
	Coverage     string `json:"coverage"`	//覆盖率
	Branches     int    `json:"branches"`	//分支数
	Predicates   int    `json:"predicates"`	//未覆盖分支数
	Tracebits    int    `json:"tracebits"` //已覆盖分支数
}

func (t *TesterPool) GetSfuzzOutputs(id string, outputs map[string]*Output) error {
	dirPath := configs.SfuzzOutputPath+"/"+id+"/"
	filesInfo,err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _,fileInfo := range filesInfo {
		if fileInfo.IsDir() {
			continue
		}
		if !strings.HasSuffix(fileInfo.Name(),".sol.txt"){
			continue
		}
		file,err:=os.Open(dirPath+fileInfo.Name())
		if err != nil {
			return err
		}
		r := reader.New(file)

		fileName := r.MustString()
		contractName := r.MustString()
		label := fileName + ":" + contractName
		if _,ok:= outputs[label]; !ok {
			outputs[label] = NewOutput(fileName,contractName)
		}
		outputSfuzz := outputs[label].Sfuzz
		outputSfuzz.Duration=r.MustString()
		outputSfuzz.Coverage=r.MustString()
		outputSfuzz.Branches=r.MustInt()
		outputSfuzz.Predicates=r.MustInt()
		outputSfuzz.Tracebits=r.MustInt()
		outputVul := outputs[label].Vul
		outputVul.GaslessSend = r.MustBool()
		outputVul.ExceptionDisorder = r.MustBool()
		outputVul.Reentrancy = outputVul.Reentrancy || r.MustBool()
		outputVul.TimeDependency = outputVul.TimeDependency ||r.MustBool()
		outputVul.NumberDependency = r.MustBool()
		outputVul.DelegateCall = r.MustBool()
		outputVul.FreezingEther = r.MustBool()
		outputVul.IntegerOverflow = outputVul.IntegerOverflow||r.MustBool()
		outputVul.IntegerUnderflow = outputVul.IntegerUnderflow || r.MustBool()
		_ = file.Close()
	}
	return nil
}
