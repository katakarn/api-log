package main

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/grafana/loki-client-go/loghttp"
    "github.com/sirupsen/logrus"
)

func main() {
    // Set up Logrus with JSON formatter
    logrus.SetFormatter(&logrus.JSONFormatter{})

    // Set up Loki as logrus hook
    hook, err := loghttp.New("http://loki:3100/loki/api/v1/push", "myapp")
    if err != nil {
        logrus.Fatal(err)
    }
    defer hook.Flush()
    logrus.AddHook(hook)

    // Set up Gin router
    router := gin.Default()

    // Endpoint for testing logging
    router.GET("/log", func(c *gin.Context) {
        logrus.Info("This is an info log from /log endpoint.")
        c.JSON(http.StatusOK, gin.H{"message": "Log sent to Loki successfully."})
    })

    // Start the server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    router.Run(":" + port)
}
