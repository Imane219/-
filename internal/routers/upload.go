package routers

import (
	"contrplatform/configs"
	"contrplatform/internal/service"
	"contrplatform/internal/tester_pool"
	"contrplatform/pkg/app"
	"contrplatform/pkg/errcode"
	"contrplatform/pkg/upload"
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
)


func UploadContracts(c *gin.Context) {
	response := app.NewResponse(c)
	form, err := c.MultipartForm()
	if err != nil {
		errRsp := errcode.InvalidParams.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}

	svc := service.New(c)
	var id string
	if ids,ok := form.Value["id"]; !ok {
		id = svc.GeneID()
	}else if len(ids)>1 {
		errRsp := errcode.InvalidParams.WithDetails("id不唯一")
		response.ToErrorResponse(errRsp)
		return
	} else {
		id = ids[0]
		state := svc.State(id)
		if state != tester_pool.StateInit {
			details := fmt.Sprintf("id:%s, state:%s - not %s",id,state,
				tester_pool.StateInit)
			errRsp := errcode.ErrorTeseterState.WithDetails(details)
			response.ToErrorResponse(errRsp)
			return
		}
	}

	fileType := upload.TypeContract
		for _, fileHeader := range form.File["file"] {
		if err := svc.UploadFile(upload.FileType(fileType), fileHeader, id);
			err != nil {
			errRsp := errcode.ErrorUploadFileFail.WithDetails(err.Error())
			response.ToErrorResponse(errRsp)
			return
		}
	}

	response.ToResponse(gin.H{
		"id":id,
		"uploaded_file":getIdFilesName(id),
	})
}


func getIdFilesName(id string) []string {
	path,_:=filepath.Abs(configs.UploadSavePath+"/"+id)
	return upload.GetDirFiles(path)
}
