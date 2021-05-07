package configs

import (
	"github.com/gin-gonic/gin"
	"time"
)

const (
	RunMode = gin.DebugMode
	HttpPort = ":7777"

	DetectorExpiredTime = 5*time.Second

	UploadSavePath = "storage/uploads"
	UploadContrExt = ".sol"
	UploadContrMaxSize = 5

	TestScriptPath = "storage/tester/tmp.py"
	OyenteOutputPath = "storage/tester/oyenteoutput"
	SfuzzOutputPath = "storage/tester/sfuzzoutput"
	//SfuzzOutputPath = "E:\\Computer\\Go\\contrplatform\\storage\\tester\\sfuzzoutput"
)

