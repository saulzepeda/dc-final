package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"github.com/gin-gonic/gin"
	"encoding/base64"
	"strings"
	"time"
	"math/rand"

	"github.com/saulzepeda/dc-final/controller"
	"github.com/saulzepeda/dc-final/scheduler"
)

const characters = "123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var actualuser = ""
//var actualpassword = "pw123"
var defaultUser = "username"
var defaultPassword = "password"
var actualtoken = ""
var usernames []string
var passwords []string

var Jobs = make(chan scheduler.Job)

func GenerateToken(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = characters[rand.Intn(len(characters))]
    }
    return string(b)
}

func GenerateId(id_type int) string { //0 for workload; 1 for image
	//if id_type == 0{
		n := len(controller.Workloads)
	//}
	
	n++
	if n < 10 {
		return "000" + strconv.Itoa(n)
	} else if n < 100{
		return "00" + strconv.Itoa(n)
	} else {
		return "0" + strconv.Itoa(n)
	}
}

func FormatId(id int) string {
	if id < 10 {
		return "000" + strconv.Itoa(id)
	} else if id < 100{
		return "00" + strconv.Itoa(id)
	} else {
		return "0" + strconv.Itoa(id)
	}
}

func ValidateUsername(actualU string) bool {
	if actualU == defaultUser{
		return true
	} else {
		return false
	}
	/*
    for i := 0; i < len(usernames); i++ {
		if actualU == usernames[i] {
			return true
		}
	}
	return false*/
}

func ValidatePassword(actualU, actualP string) bool {
	if actualP == defaultPassword{
		return true
	} else {
		return false
	}
	/*pos := 0
    for i := 0; i < len(usernames); i++ {
		if actualU == usernames[i] {
			pos = i
		}
	}
	if actualP == passwords[pos]{
		return true
	} else {
		return false
	}*/
}

func Start() {
	r := gin.Default()
	/*
	r.GET("/signin", func(c *gin.Context) {
		authorization := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		todecod, _ := base64.StdEncoding.DecodeString(authorization[1])
		userdata := strings.SplitN(string(todecod), ":", 2)

		if !ValidateUsername(userdata[0]){
			usernames = append(usernames, userdata[0])	
			passwords = append(passwords, userdata[1])
			c.JSON(http.StatusOK, gin.H{
				"message": "Hi " + userdata[0] + " , your user has been created.",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Error, your username already exists.",
			})
		}

	})*/
	
	r.GET("/login", func(c *gin.Context) {
		authorization := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		todecod, _ := base64.StdEncoding.DecodeString(authorization[1])
		userdata := strings.SplitN(string(todecod), ":", 2)
		
		//Verify the username and password

		if ValidateUsername(userdata[0]) {
			if ValidatePassword(userdata[0], userdata[1]) {
				actualuser = userdata[0]
				tokenrand := GenerateToken(8)
				actualtoken = tokenrand

				c.JSON(http.StatusOK, gin.H{
					"message": "Hi " + actualuser + " , welcome to the  DPIP System",
					"token": tokenrand,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "The password is incorrect",
				})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "The username is not registered",
			})
		}
	})

	r.GET("/logout", func(c *gin.Context) {
		authorization := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		token := authorization[1]

		if token == actualtoken {
			
			actualtoken = ""
			c.JSON(http.StatusOK, gin.H{
				"message": "Bye " + actualuser + ", your token has been revoked",
			})
			actualuser = ""
		}
	})
	
	r.POST("/images", func(c *gin.Context) {
		authorization := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		token := authorization[1]

		if token == actualtoken {
			file, err := c.FormFile("data")
			wl_id := c.PostForm("workload-id")

			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
				return
			}

			fileName := filepath.Base(file.Filename)
			fmt.Println(fileName)
			f, err := os.Open(fileName)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"message": "Error opening the image",
					"filename": fileName,
				})
				return
			}

			_, ok := controller.Workloads[wl_id]
			if ok {
				img_id := FormatId(len(controller.Workloads[wl_id].Filtered_images) + 1)
				new_filename := img_id + "_original" + filepath.Ext(file.Filename)
				path := "images/" + controller.Workloads[wl_id].Name + "/" + new_filename

				err := c.SaveUploadedFile(file, path)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"message": "Error saving the image",
						"filename": new_filename,
					})
					return
				} else {
					c.JSON(http.StatusOK, gin.H{
						"message": "Image saved",
						"filename": fileName,
					})
				}
				new_job := scheduler.Job{
					Address: "localhost:50051", 
					RPCName: "image", 
					Filepath: path,
					Wl_id: wl_id,
					Filter_type: controller.Workloads[wl_id].Filter,
				}
				Jobs <- new_job
				time.Sleep(time.Second * 1)
				
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "The workload doesn't exist",
					"workload_id": wl_id,
				})
				return
			}

			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"message": "Error opening the image",
					"filename": fileName,
				})
				return
			}else{
				fi, _ := f.Stat()
				c.JSON(http.StatusOK, gin.H{
					"message": "An image has been successfully uploaded",
					"filename": fileName,
					"size": strconv.Itoa(int(fi.Size()/1000)) + " kb",
				})
			}

			f.Close()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "ERROR, you have to log in",
			})
			return
		}

	})
	
	r.GET("/status", func(c *gin.Context) {
		authorization := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		token := authorization[1]

		if token == actualtoken {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hi " + actualuser + " , the DPIP System is Up and Running",
				"time": time.Now().Format("2006-01-02T15:04:05+07:00"),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "ERROR, you have to log in",
			})
		}
		
	})

	r.POST("/workloads", func(c *gin.Context) {
		authorization := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		token := authorization[1]

		if token == actualtoken {
			
			wl_name := c.PostForm("workload-name")
			filter_type := c.PostForm("filter")

			wl_exists := false
			
			//Check if the workload name already exists
			for _, wl_actual := range controller.Workloads {
				if wl_actual.Name == wl_name {
					wl_exists = true
					break
				}
			}

			if(wl_exists){
				c.JSON(http.StatusOK, gin.H{
					"message": "The workload already exists.",
				})
			} else{
				_ = os.MkdirAll("images/" + wl_name + "/", 0755)

				wl_id := GenerateId(0)
				wl := controller.Workload{
					ID: wl_id,
					Filter: filter_type,
					Name: wl_name,
					Status: "scheduling",
					Running_jobs: 0,
					Filtered_images: []string{},
				}

				controller.Workloads[wl.ID] = wl

				c.JSON(http.StatusOK, gin.H{
					"workload_id": controller.Workloads[wl.ID].ID,
					"filter":   controller.Workloads[wl.ID].Filter,
					"workload_name": controller.Workloads[wl.ID].Name,
					"status": controller.Workloads[wl.ID].Status,
					"running_jobs": controller.Workloads[wl.ID].Running_jobs,
					"filtered_images": controller.Workloads[wl.ID].Filtered_images,
				})
			}

		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "ERROR, you have to log in",
			})
		}

	})

	r.GET("/workloads/:workload_id", func(c *gin.Context) {
		authorization := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		token := authorization[1]

		wl_id := c.Param("workload_id")

		if token == actualtoken {
			_, ok := controller.Workloads[wl_id]
			if ok {
				c.JSON(http.StatusOK, gin.H{
					"workload_id": controller.Workloads[wl_id].ID,
					"filter":   controller.Workloads[wl_id].Filter,
					"workload_name": controller.Workloads[wl_id].Name,
					"status": controller.Workloads[wl_id].Status,
					"running_jobs": controller.Workloads[wl_id].Running_jobs,
					"filtered_images": controller.Workloads[wl_id].Filtered_images,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "The workload id doesn't exist.",
				})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "ERROR, you have to log in",
			})
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
