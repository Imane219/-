package configs

import "github.com/gin-gonic/gin"

const (
	RunMode = gin.DebugMode

	UploadSavePath = "storage/uploads"
	UploadContrExt = ".sol"
	UploadContrMaxSize = 5

	TestScriptPath = "storage/tester/tmp.py"
	OyenteOutputPath = "storage/tester/oyenteoutput"
	SfuzzOutputPath = "storage/tester/sfuzzoutput"
	//SfuzzOutputPath = "E:\\Computer\\Go\\contrplatform\\storage\\tester\\sfuzzoutput"
)

