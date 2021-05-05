package tester_pool


type Output struct {
	FileName     string        `json:"file_name"`
	ContractName string        `json:"contract_name"`
	Sfuzz        *SfuzzOutput  `json:"sfuzz"`
	Oyente       *OyenteOutput `json:"oyente"`
	Vul	*VulOutput	`json:"vulnerabilities"`
}

type VulOutput struct {
	GaslessSend       bool `json:"gasless_send"`
	ExceptionDisorder bool `json:"exception_disorder"`
	Reentrancy        bool `json:"reentrancy"`
	TimeDependency    bool `json:"timestamp_dependency"`
	NumberDependency  bool `json:"block_number_dependency"`
	DelegateCall      bool `json:"delegate_call"`
	FreezingEther     bool `json:"freezing_ether"`
	IntegerOverflow   bool `json:"integer_overflow"`
	IntegerUnderflow  bool `json:"integer_underflow"`
	Callstack	bool	`json:"callstack"`
	ParityMultsig bool	`json:"parity_multsig_bug_2"`
	MoneyConcurrency bool `json:"money_concurrency"`
}

func NewOutput(fileName, contractName string) *Output {
	return &Output{
		FileName: fileName,
		ContractName: contractName,
		Sfuzz: &SfuzzOutput{},
		Oyente: &OyenteOutput{},
		Vul: &VulOutput{},
	}
}