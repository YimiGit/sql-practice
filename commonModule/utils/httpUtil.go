package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// TestConcurrentRequest 并发请求测试
func TestConcurrentRequest() {
	body := []byte(`{ "username":"yimi","password":"yimipass"}`)

	// 设置请求的 Content-Type
	contentType := "application/json"

	c := make(chan struct{}, 0)

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				// 创建新的请求对象，并设置请求 URL、Body 和 Content-Type
				req, err := http.NewRequest("POST", "http://127.0.0.1:19600/login", bytes.NewBuffer(body))
				if err != nil {
					fmt.Println("创建请求失败:", err)
					return
				}
				req.Header.Set("Content-Type", contentType)

				// 发送请求并获取响应
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("发送请求失败:", err)
					return
				}
				//defer resp.Body.Close()

				// 读取响应内容
				respBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("读取响应失败:", err)
					return
				}

				// 输出响应内容
				fmt.Println(string(respBody))
			}
		}()
	}
	c <- struct{}{}
}
