package setting

import "time"

type ServerSettingS struct {
	RunMode string
	HttpPort string
}

type AppSettingS struct {
	FrontendPath string
}

type LogSettingS struct {
	SavePath string
	FileName string
	FileExt string
	MaxSize int
	MaxAge int
}

type PoolSettingS struct {
	UploadSavePath string
	UploadContrExts []string
	UploadContrMaxSize int
	DetectorExpiredTime time.Duration
	TestScriptPath string
	OyenteOutputPath string
	SfuzzOutputPath string
}

//存放每一部分配置的映射
//键值是接口,实际上为配置结构体的地址
var sections = make(map[string]interface{})

//读取配置到结构体
func (s *Setting) ReadSection(k string, v interface{}) error {
	//对标题为k的配置部分解码到v中
	if err := s.vp.UnmarshalKey(k, v); err != nil {
		return err
	}
	//将配置记录到映射中
	if _, ok:=sections[k]; !ok {
		sections[k]=v
	}
	return nil
}

//读取所有部分的配置
func (s *Setting) ReloadAllSection() error {
	for k,v :=range sections{
		//由于sections中键值实际上是配置结构体的地址
		//因此此处读取配置后会存到全局的配置结构体中
		if err := s.ReadSection(k,v); err != nil {
			return err
		}
	}
	return nil
}