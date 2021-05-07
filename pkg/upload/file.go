package upload

import (
	"contrplatform/configs"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

type FileType int

const TypeContract = iota + 1


func CheckSavePathNotExist(dst string) bool {
	_,err := os.Stat(dst)
	return os.IsNotExist(err)
}

func CreatSavePath(dst string, perm os.FileMode) error {
	if err:=os.MkdirAll(dst,perm); err != nil {
		return err
	}
	return nil
}


func CheckContainExt(t FileType,filename string) bool {
	ext := path.Ext(filename)
	switch t {
	case TypeContract:
		return strings.ToLower(ext)==strings.ToLower(configs.UploadContrExt)
	}
	return false
}

//检测文件权限是否不足够
func CheckNotPermission(dst string) bool {
	_, err:=os.Stat(dst)
	return os.IsPermission(err)
}

func CheckOutMaxSize(t FileType, f multipart.File) bool {
	content,_ := ioutil.ReadAll(f)
	size := len(content)
	switch t {
	case TypeContract:
		return size >= configs.UploadContrMaxSize*1024*1024
	}
	return false
}

func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out,err:=os.Create(dst)
	defer out.Close()
	_,err = io.Copy(out,src)
	return err
}

func GetDirFiles(dirPath string) []string {
	//files,_:=filepath.Glob(dirPath+"/*")
	files ,_ := ioutil.ReadDir(dirPath)
	filesName := make([]string,0)
	for _, file := range files {
		filesName = append(filesName, file.Name())
	}
	return filesName
}

func FileUpload(fileHeader *multipart.FileHeader, fileType FileType,
	pathSuffix string) error {
	fileName := fileHeader.Filename
	//判断拓展名
	if !CheckContainExt(fileType,fileName){
		return errors.New("文件扩展名错误")
	}
	//判断文件保存路径是否存在并创建
	//文件保存路径: /storage/uploads/sessionID/
	uploadSavePath := configs.UploadSavePath+"/"+ pathSuffix
	if CheckSavePathNotExist(uploadSavePath) {
		if err := CreatSavePath(uploadSavePath,os.ModePerm);err!=nil{
			return errors.New("创建文件夹失败")
		}
	}
	//检测文件是否权限不足
	if CheckNotPermission(uploadSavePath) {
		return errors.New("文件权限不足")
	}
	file,err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	//判断文件大小
	if CheckOutMaxSize(fileType,file){
		return errors.New("文件大小过大")
	}
	dstPath:=uploadSavePath+"/"+fileName
	//上传文件
	if err:=SaveFile(fileHeader,dstPath);err!=nil{
		return err
	}
	return nil
}

