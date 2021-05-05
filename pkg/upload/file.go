package upload

import (
	"contrplatform/configs"
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