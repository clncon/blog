package controllers

import (
	"blog/system"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os/exec"
)

func Go(c *gin.Context){
//获取可执行文件路径hook.sh
command:=system.GetConfiguration().ShellPath

cmd:=exec.Command("/bin/bash","-c",command)

output,err:=cmd.Output()
if err!=nil {
	logrus.Errorf("Execute Shell:%s failed with error:%s",command,err.Error())
	return
}


logrus.Infof("Execute Shell %s finished with output:\n%s",command,string(output))


}
