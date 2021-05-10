package detectorpool

import (
	"contrplatform/pkg/upload"
	"errors"
	"mime/multipart"
	"os"
)

func (dp *DetectorPool) UploadContracts(fileHeaders []*multipart.FileHeader, id string) error {
	tester, _ := dp.detector(id)
	for _, fileHeader := range fileHeaders {
		if err := tester.checkUploadState(); err != nil {
			return err
		}
		if err := dp.uploadContract(fileHeader, id); err != nil {
			return err
		}
	}
	return nil
}

func (dp *DetectorPool) uploadContract(fileHeader *multipart.FileHeader, id string) error {
	fileName := fileHeader.Filename
	//判断拓展名
	if !upload.CheckContainExt(fileName,poolSetting.UploadContrExts){
		return errors.New("文件扩展名错误")
	}
	//判断文件保存路径是否存在并创建
	//文件保存路径: /storage/uploads/sessionID/
	uploadSavePath := poolSetting.UploadSavePath+"/"+ id
	if upload.CheckSavePathNotExist(uploadSavePath) {
		if err := upload.CreatSavePath(uploadSavePath,os.ModePerm);err!=nil{
			return errors.New("创建文件夹失败")
		}
	}
	//检测文件是否权限不足
	if upload.CheckNotPermission(uploadSavePath) {
		return errors.New("文件权限不足")
	}
	file,err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	//判断文件大小
	if upload.CheckOutMaxSize(file,poolSetting.UploadContrMaxSize){
		return errors.New("文件大小过大")
	}
	dstPath:=uploadSavePath+"/"+fileName
	//上传文件
	if err:=upload.SaveFile(fileHeader,dstPath);err!=nil{
		return err
	}
	return nil
}