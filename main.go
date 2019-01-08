package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml"
	"gopkg.in/go-playground/webhooks.v5/github"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

var conf Config

func main() {
	getConfig()
	app := gin.New()
	app.POST("/webhook", handleWebhook)
	app.Run(":" + conf.Port)
}

var curStatus = Free
var lock sync.Mutex

func handleWebhook(c *gin.Context) {
	getConfig()
	var giteeHook *GiteeWebhook
	var projectName string
	log.Println("time is: ", time.Now().String())
	log.Println(c.Request.UserAgent())
	log.Println(c.Request.Host)
	if strings.Contains(c.Request.UserAgent(), "GitHub") {
		log.Println("pingtai", "GitHub")
		hook, _ := github.New(github.Options.Secret(conf.Password))
		payload, err := hook.Parse(c.Request, github.PushEvent)
		if err != nil {
			if err == github.ErrHMACVerificationFailed || err == github.ErrMissingHubSignatureHeader {
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "error with password not match",
				})
				log.Println("stop because password error")
			} else {
				return
			}
		}

		switch payload.(type) {
		case github.PushPayload:
			release := payload.(github.PushPayload)
			projectName = release.Repository.Name
		}
	} else {
		giteeHook = &GiteeWebhook{}
		log.Println("pingtai", "Gitee")
		err := c.Bind(&giteeHook)
		if err != nil {
			log.Fatalln(err)
		}

		projectName = giteeHook.getName()
		log.Println("gitee password: ", giteeHook.getPwd())
		log.Println("gitee name: ", giteeHook.getName())
	}

	if len(projectName) == 0 {
		log.Println("no effective project name")
		return
	}

	p, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	var cmd *exec.Cmd
	for _, project := range conf.Projects {
		if project.Name == projectName {
			log.Println("is ", project.Name)

			if Exists(path.Join(p, project.Cmd)) {
				cmd = exec.Command("bash", path.Join(p, project.Cmd))
			} else if Exists(path.Join(project.Path, project.Cmd)) {
				cmd = exec.Command("bash", project.Cmd)
			} else {
				cmd = exec.Command("echo", "There is no file to exec")
			}
			cmd.Dir = project.Path
			break
		}
	}

	if cmd == nil {
		cmd = exec.Command("echo", "unknown repo, please check your database.toml ")
		return
	}

	if curStatus == Running {
		c.JSON(http.StatusOK, gin.H{
			"code": 940,
			"msg":  "running last cmd, later to do this",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1314520,
			"msg":  "success hook",
		})
	}
	go Exec(cmd)
}

func Exec(cmd *exec.Cmd) {
	curStatus = Running
	lock.Lock()
	defer lock.Unlock()
	log.Println("start exec --- ")
	var abc []byte
	var err error
	if abc, err = cmd.Output(); err != nil {
		log.Println(err)
		log.Println("err when exec --")
	}
	log.Println(string(abc))
	log.Println("end exec --- ")
	curStatus = Free
}

func getConfig() {
	file, err := ioutil.ReadFile("database.toml")
	if err != nil {
		log.Fatalln(err)
	}

	err = toml.Unmarshal(file, &conf)

	if err != nil {
		log.Fatalln(err)
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

type Status int

type Project struct {
	Name     string `toml:"name"`
	Path     string `toml:"path"`
	Cmd      string `toml:"cmd"`
	Provider string `toml:"provider"`
	Password string `toml:"pwd"`
}

type Config struct {
	Password string    `toml:"pwd"`
	Port     string    `toml:"port"`
	Projects []Project `toml:"data"`
}

const (
	Running Status = iota
	Free
	Error
)
