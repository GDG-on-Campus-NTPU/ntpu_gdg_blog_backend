package routes

import (
	"blog/routerRegister"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	routerRegister.Register = append(routerRegister.Register, func(rg *gin.RouterGroup) {
		img := rg.Group("/img")

		img.GET("/:id", func(c *gin.Context) {
			resp, err := http.Get(fmt.Sprintf("https://drive.google.com/uc?export=view&id=%v", c.Param("id")))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch image"})
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				c.String(resp.StatusCode, fmt.Sprintf("Failed to fetch image: received status code %d", resp.StatusCode))
				return
			}
			contentType := resp.Header.Get("Content-Type")
			c.DataFromReader(http.StatusOK, resp.ContentLength, contentType, resp.Body, nil)
			c.Header("Cache-Control", "public, max-age=300")
			c.Header("Content-Type", contentType)
			c.Status(http.StatusOK)
			_, err = io.Copy(c.Writer, resp.Body)
			if err != nil {
				fmt.Printf("Error copying image data: %v\n", err)
			}
		})

	})
}
