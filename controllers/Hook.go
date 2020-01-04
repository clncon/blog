package controllers

import (
	"blog/system"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os/exec"
)

func Go(c *gin.Context) {
	//获取可执行文件路径hook.sh
	command := system.GetConfiguration().ShellPath
	doShell(command)
	c.JSON(200, "success")
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
