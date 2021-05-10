package upload

import (
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


func CheckContainExt(filename string, exts []string) bool {
	ext := path.Ext(filename)
	for _, e := range exts {
		if strings.ToLower(ext)!=strings.ToLower(e) {
			return false
		}
	}
	return true
}

//检测文件权限是否不足够
func CheckNotPermission(dst string) bool {
	_, err:=os.Stat(dst)
	return os.IsPermission(err)
}

func CheckOutMaxSize(f multipart.File, maxSize int) bool {
	content,_ := ioutil.ReadAll(f)
	return len(content)>maxSize*1024*1024
}

func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out,err:=os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_,err = io.Copy(out,src)
	return err
}

func GetDirFiles(dirPath string) []string {
	files ,_ := ioutil.ReadDir(dirPath)
	filesName := make([]string,0)
	for _, file := range files {
		filesName = append(filesName, file.Name())
	}
	return filesName
}

