package service

import (
	"math/rand"
	"strings"
)

var sessionIDs = make(map[uint16]bool)

func (svc *Service) geneSessionID() {
	sessionId := uint16(rand.Uint32() >> 16)
	for _, ok:=sessionIDs[sessionId];ok;{
		sessionId = uint16(rand.Uint32() >> 16)
	}
	sessionIDs[sessionId]=true
	svc.Session.Set("id",sessionId)
}

func (svc *Service) InitSession() error {
	svc.geneSessionID()
	svc.makeSessionFiles()
	return svc.Session.Save()
}

func (svc *Service) GetSessionID() interface{} {
	return svc.Session.Get("id")
}


func (svc *Service) makeSessionFiles()  {
	svc.Session.Set("files","")
}

func (svc *Service) SetSessionFile(fileName string) error {
	filesStr := svc.Session.Get("files").(string)
	files := strings.Split(filesStr,"$")
	for _,file := range files {
		if file == fileName {
			return nil
		}
	}
	if filesStr!=""{
		filesStr+="$"
	}
	filesStr += fileName
	svc.Session.Set("files",filesStr)
	return svc.Session.Save()
}

func (svc *Service) GetSessionFiles() []string {
	filesStr := svc.Session.Get("files").(string)
	return strings.Split(filesStr,"$")
}
