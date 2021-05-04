package configs

import "github.com/gin-gonic/gin"

const (
	RunMode = gin.DebugMode

	UploadSavePath = "storage/uploads"
	UploadContrExt = ".sol"
	UploadContrMaxSize = 5
)

