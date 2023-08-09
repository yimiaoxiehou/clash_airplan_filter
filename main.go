package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	r := gin.Default()
	r.GET("", func(c *gin.Context) {
		targetUrl := c.Query("url")
		request, _ := http.NewRequest("GET", targetUrl, nil)
		for k, v := range c.Request.Header {
			request.Header.Set(k, v[0])
		}

		client := &http.Client{
			Timeout: time.Second * 5, //超时时间
		}

		resp, err := client.Do(request)
		if err != nil {
			fmt.Println("出错了", err)
			return
		}

		defer resp.Body.Close()
		for k, v := range resp.Header {
			c.Header(k, v[0])
		}

		if info := resp.Header.Get("Subscription-Userinfo"); info != "" {
			m := make(map[string]string)
			for _, s := range strings.Split(info, ";") {
				s = strings.TrimSpace(s)
				l := strings.Split(s, "=")
				m[l[0]] = l[1]
			}
			upload, _ := strconv.Atoi(m["upload"])
			download, _ := strconv.Atoi(m["download"])
			total, _ := strconv.Atoi(m["total"])
			expire, _ := strconv.ParseInt(m["expire"], 10, 64)
			expireTime := time.Unix(expire, 0)
			now := time.Now()
			if now.After(expireTime) {
				c.AbortWithError(500, errors.New("this sub is expired."))
				c.Writer.Flush()
				return
			}
			rate := (float32(upload) + float32(download)) / float32(total)
			if rate > 0.92 {
				c.AbortWithError(500, errors.New("this sub use more than 92%"))
				c.Writer.Flush()
				return
			}
		}
		io.Copy(c.Writer, resp.Body)
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
