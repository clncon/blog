package controllers

import (
	"blog/system"
	"crypto/hmac"
	"crypto/sha1"
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

	mac.Write(buff[0:n])
    result:=string(mac.Sum(nil))

    signature :=c.GetHeader("X-Hub-Signature")
	doShell(command)
	c.JSON(200, signature+":"+result)
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
