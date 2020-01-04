package controllers

import (
	"blog/system"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os/exec"
)

func Go(c *gin.Context) {
	//获取可执行文件路径hook.sh
	command := system.GetConfiguration().ShellPath
	serect:="135696cc92c1d9c0a74e956a4594652b"
	buff:=make([]byte,1024)
	n,_:=c.Request.Body.Read(buff)
	mac:=hmac.New(sha1.New,[]byte(serect))
    for;n>0;{
		mac.Write(buff[0:n])
		n,_=c.Request.Body.Read(buff)
	}
	sha1:="sha1="+hex.EncodeToString(mac.Sum(nil))

    signature :=c.GetHeader("X-Hub-Signature")
    if sha1 == signature {
    	doShell(command)
		c.JSON(200,"auto deploy success")
	}else {
		c.JSON(500,"signature error")
	}


}

func doShell(command string) {
	cmd := exec.Command("/bin/bash", "-c", command)

	err := cmd.Start()
	if err != nil {
		logrus.Errorf("Execute Shell:%s failed with error:%s", command, err.Error())
		return
	}

	logrus.Infof("Execute Shell %s finished", command)
}
