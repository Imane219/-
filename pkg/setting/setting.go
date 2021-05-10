package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

//配置
type Setting struct {
	vp *viper.Viper
}

//新建配置
func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config") //配置文件名(无拓展名)
	vp.SetConfigType("yaml")   //配置文件拓展名
	vp.AddConfigPath("configs/")
	//根据配置参数来配置文件路径
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}
	//读取配置文件
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}
	s := &Setting{
		vp: vp,
	}
	s.WatchSettingChange() //监视文件
	return s, nil
}

//var Ch = make(chan bool)

//监视配置文件变更情况
func (s *Setting) WatchSettingChange() {
	//新建协程监视配置文件
	go func() {
		//设置热更新时触发的函数
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSection() //重新读取配置
			log.Print("Config file updated")
			//Ch<-true
		})
		//监视配置文件
		s.vp.WatchConfig()
	}()
}
