package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		//Encode the data
		postBody, _ := json.Marshal(map[string]string{
			"user":     "root",
			"password": "123456",
			"path":     "./server.exe",
		})
		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post("http://localhost:8000/", "application/json", responseBody)

		//Handle Error
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": resp.Status,
				"code":    resp.StatusCode,
				"data":    "",
			})
			return
		}

		//Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		bd := string(body)
		c.JSON(http.StatusOK, gin.H{
			"message": "发布成功",
			"code":    0,
			"data":    bd,
		})
	})

	//监听端口默认为8080
	r.Run()

}
