package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type PassWordLoginForm struct {
	User     string `form:"user" json:"user" binding:"required"`
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Path     string `form:"path" json:"path" binding:"required"`
}

func main() {
	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		passwordLoginForm := PassWordLoginForm{}
		if err := c.ShouldBind(&passwordLoginForm); err != nil {
			fmt.Printf("%s", err)
			return
		}
		dir, _ := os.Getwd()

		//var wkCount sync.WaitGroup
		//gitPull(&wkCount)
		gitCmd := exec.Command("git", "pull")
		var stderr bytes.Buffer
		gitCmd.Stderr = &stderr
		gitCmd.Dir = dir + "/demo.go-admin.cn/"
		err := gitCmd.Run()
		//_, err := gitCmd.CombinedOutput()
		//err := exec.Command("git", "clone", "https://github.com/GoAdminGroup/demo.go-admin.cn.git").Run()
		if err != nil {
			fmt.Println(stderr.String())
			return
		}
		//goGet(&wkCount)
		cmd := exec.Command("go", "get", "-u")
		cmd.Dir = dir + "/server/"
		fmt.Sprintf("%s", dir)
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		//goBuild(&wkCount)
		goCmd := exec.Command("go", "build", "-o", "../")
		goCmd.Dir = dir + "/server/"
		err = goCmd.Run()
		if err != nil {
			panic(err)
		}
		//wkCount.Wait()
		/*goCmd := exec.Command("scp", "-r", "../test/", "root@localhost:/home/")
		goCmd.Dir = "./server/"
		err := goCmd.Run()
		if err != nil {
			panic(err)
		}*/
		scp(passwordLoginForm)
		c.String(http.StatusOK, fmt.Sprintf("%s", passwordLoginForm))
	})
	r.Run(":8000")
}

func gitPull(wkCount *sync.WaitGroup) {
	wkCount.Add(1)
	go func() {
		defer wkCount.Done()
		err := exec.Command("git", "-C", "./demo.go-admin.cn/", "pull").Run()
		//err := exec.Command("git", "clone", "https://github.com/GoAdminGroup/demo.go-admin.cn.git").Run()
		if err != nil {
			panic(err)
		}
	}()
}

func goGet(wkCount *sync.WaitGroup) {
	wkCount.Add(1)
	go func() {
		defer wkCount.Done()
		cmd := exec.Command("go", "get", "-u", "github.com/gin-gonic/gin")
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}()
}

func goBuild(wkCount *sync.WaitGroup) {
	wkCount.Add(1)
	go func() {
		defer wkCount.Done()
		goCmd := exec.Command("go", "build", "-o", "../test/")
		goCmd.Dir = "./server/"
		err := goCmd.Run()
		if err != nil {
			panic(err)
		}
	}()
}

func scp(p PassWordLoginForm) {
	Conf := &ssh.ClientConfig{
		User: p.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(p.PassWord),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	Client, err := ssh.Dial("tcp", "10.18.13.112:22", Conf)
	if err != nil {
		fmt.Println(nil)
	}
	defer func(Client *ssh.Client) {
		err := Client.Close()
		if err != nil {

		}
	}(Client)

	if session, err := Client.NewSession(); err == nil {
		defer session.Close()

		go func() {
			Buf := make([]byte, 1024)
			w, _ := session.StdinPipe()
			defer w.Close()
			_, fileName := filepath.Split(p.Path)
			File, _ := os.Open(p.Path)
			info, _ := File.Stat()
			fmt.Fprintln(w, "C0644", info.Size(), fileName)
			for {
				n, err := File.Read(Buf)
				fmt.Fprint(w, string(Buf[:n]))
				if err != nil {
					if err == io.EOF {
						return
					} else {
						panic(err)
					}
				}
			}
			File.Close()
		}()
		if err := session.Run("/usr/bin/scp -qrt /home"); err != nil {
			if err != nil {
				if err.Error() != "Process exited with: 1. Reason was:  ()" {
					fmt.Println(err.Error())
				}
			}
		}
	}
}
