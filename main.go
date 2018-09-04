package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"log"
	"os/exec"
		"os"
		"path"
	"io/ioutil"
		"strings"
	"github.com/pelletier/go-toml"
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

	var hook Hook

	log.Println("time is: ", time.Now().String())
	if strings.Contains(c.Request.UserAgent() ,"Github"){
		hook = &GithubWebHook{}
		log.Println("pingtai","Github")
	} else {
		hook = &GiteeWebhook{}
		log.Println("pingtai","Gitee")
	}
	c.Bind(&hook)

	log.Println("password: ", hook.getPwd())
	log.Println("url: ", hook.getUrl())
	log.Println("name: ", hook.getName())

	p, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	var cmd *exec.Cmd
	for _, project := range conf.Projects {

		if project.Name == hook.getName() {
			log.Println("is ", project.Name)
			//全局密码与独立密码的判断
			if   project.Password != hook.getPwd() && conf.Password != hook.getPwd() {
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "error with password not match",
				})
				log.Println("stop because password error")
				return
			}

			if  Exists(path.Join(p, project.Cmd)){
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

	if	cmd ==nil {
		cmd = exec.Command("echo", "unknown repo: "+hook.getName())
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
		log.Println( err)
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
	Name string `toml:"name"`
	Path string `toml:"path"`
	Cmd  string `toml:"cmd"`
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

type Hook interface {
	getUrl() string
	getName() string
	getPwd() string
}