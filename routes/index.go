package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddIndexRoutes(rg *gin.RouterGroup) {
	rg.GET("", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(
			`
            <!DOCTYPE html>
            <html>
                <head>
                    <title>Login Demo</title>
                    <button id="login">Login</button>
                </head>
            </html>

            <script>

                document.getElementById('login').addEventListener('click', async () => {
                    window.location.href = "/api/login";
                });
            </script>
            `))
	})
}
