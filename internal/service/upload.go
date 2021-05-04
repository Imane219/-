package service

import (
	"contrplatform/pkg/upload"
	"errors"
	"mime/multipart"
	"os"
	"strconv"
)

func (svc *Service) UploadFile(fileType upload.FileType, file multipart.File,
	fileHeader *multipart.FileHeader) error {
	fileName := fileHeader.Filename
	if !upload.CheckContainExt(fileType,fileName){
		return errors.New("文件扩展名错误")
	}
	uploadSavePath := upload.GetSavePath()+"/"+
		strconv.Itoa(int(svc.GetSessionID().(uint16)))
	if upload.CheckSavePathNotExist(uploadSavePath) {
		if err := upload.CreatSavePath(uploadSavePath,os.ModePerm);err!=nil{
			return errors.New("创建文件夹失败")
		}
	}
	if upload.CheckOutMaxSize(fileType,file){
		return errors.New("文件大小过大")
	}
	dstPath:=uploadSavePath+"/"+fileName
	if err:=upload.SaveFile(fileHeader,dstPath);err!=nil{
		return err
	}
	return nil
}
