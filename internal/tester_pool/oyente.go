package tester_pool

import (
	"contrplatform/configs"
	"encoding/json"
	"io/ioutil"
	"strings"
)

type OyenteJson struct {
	Vul *VulJson	`json:"vulnerabilities"`
	Coverage string	`json:"evm_code_coverage"`
}

type VulJson struct {
	Callstack []string	`json:"callstack"`
	TimeDependency []string	`json:"time_dependency"`
	Reentrancy []string	`json:"reentrancy"`
	IntegerOverflow []string	`json:"integer_overflow"`
	ParityMultsig []string	`json:"parity_multsig_bug_2"`
	IntegerUnderflow []string	`json:"integer_underflow"`
	MoneyConcurrency [][]string `json:"money_concurrency"`
}

type OyenteOutput struct {
	Coverage string `json:"coverage"`
}

func (t *TesterPool) GetOyenteOutputs(id string, outputs map[string]*Output) error {
	dirPath := configs.OyenteOutputPath+"/"+id+"/"
	filesInfo,err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _,fileInfo := range filesInfo {
		if fileInfo.IsDir() {
			continue
		}
		if !strings.HasSuffix(fileInfo.Name(), ".json") {
			continue
		}
		content,err := ioutil.ReadFile(dirPath+fileInfo.Name())
		if err != nil {
			return err
		}
		var oyenteJson OyenteJson
		err = json.Unmarshal(content,&oyenteJson)
		if err != nil {
			return err
		}
		label:=strings.TrimSuffix(fileInfo.Name(),".json")
		label=strings.Replace(label,"^%",":",1)
		if _,ok := outputs[label]; !ok {
			labels := strings.Split(label,":")
			outputs[label] = NewOutput(labels[0],labels[1])
		}
		outputs[label].Oyente.Coverage = oyenteJson.Coverage+"%"
		outputsVul := outputs[label].Vul
		oyenteVul := oyenteJson.Vul
		outputsVul.Callstack = len(oyenteVul.Callstack)>0
		outputsVul.ParityMultsig = len(oyenteVul.ParityMultsig)>0
		outputsVul.MoneyConcurrency = len(oyenteVul.MoneyConcurrency)>0
		outputsVul.TimeDependency = outputsVul.TimeDependency || len(oyenteVul.TimeDependency)>0
		outputsVul.Reentrancy = outputsVul.Reentrancy || len(oyenteVul.Reentrancy)>0
		outputsVul.IntegerOverflow = outputsVul.IntegerOverflow || len(oyenteVul.IntegerOverflow)>0
		outputsVul.IntegerUnderflow = outputsVul.IntegerUnderflow || len(oyenteVul.IntegerUnderflow)>0
	}
	return nil
}
