package service

import (
	"contrplatform/configs"
	"contrplatform/pkg/upload"
	"errors"
	"mime/multipart"
	"os"
)

//上传文件
func (svc *Service) UploadFile(fileType upload.FileType,fileHeader *multipart.FileHeader,
	pathPrefixes string) error {
	fileName := fileHeader.Filename
	//判断拓展名
	if !upload.CheckContainExt(fileType,fileName){
		return errors.New("文件扩展名错误")
	}
	//文件保存路径: /storage/uploads/sessionID/
	uploadSavePath := configs.UploadSavePath+"/"+pathPrefixes
	if upload.CheckSavePathNotExist(uploadSavePath) {
		if err := upload.CreatSavePath(uploadSavePath,os.ModePerm);err!=nil{
			return errors.New("创建文件夹失败")
		}
	}
	file,err := fileHeader.Open()
	if err != nil {
		return err
	}
	if upload.CheckOutMaxSize(fileType,file){
		return errors.New("文件大小过大")
	}
	//检测文件是否权限不足
	if upload.CheckNotPermission(uploadSavePath) {
		return errors.New("insufficient file permissions")
	}
	dstPath:=uploadSavePath+"/"+fileName
	if err:=upload.SaveFile(fileHeader,dstPath);err!=nil{
		return err
	}
	return nil
}
