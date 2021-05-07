package routers

import (
	"contrplatform/configs"
	"contrplatform/internal/service"
	"contrplatform/pkg/app"
	"contrplatform/pkg/errcode"
	"contrplatform/pkg/upload"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

func (d Detection) Upload(c *gin.Context) {
	response := app.NewResponse(c)
	form, err := c.MultipartForm()
	if err != nil {
		errRsp := errcode.InvalidParams.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}

	svc := service.New(c)
	var id string
	if ids, ok := form.Value["id"]; !ok {
		id = svc.GeneID()
	} else if len(ids) > 1 {
		errRsp := errcode.InvalidParams.WithDetails("id不唯一")
		response.ToErrorResponse(errRsp)
		return
	} else {
		id = ids[0]
	}

	if err :=svc.UploadContracts(form.File["file"],id); err != nil {
		errRsp:=errcode.ErrorUploadContractsFail.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}

	response.ToResponse(gin.H{
		"id":            id,
		"uploaded_file": getIdFilesName(id),
	})
}

func getIdFilesName(id string) []string {
	path, _ := filepath.Abs(configs.UploadSavePath + "/" + id)
	return upload.GetDirFiles(path)
}
