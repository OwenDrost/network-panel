package controller

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
)

// InstallScript serves the node install shell script from the container filesystem.
// Path: /install.sh
func InstallScript(c *gin.Context) {
    // Ensure file exists in working directory (copied by Dockerfile)
    const path = "install.sh"
    if _, err := os.Stat(path); err != nil {
        c.String(http.StatusNotFound, "install.sh not found")
        return
    }
    c.Header("Content-Type", "text/x-shellscript; charset=utf-8")
    c.File(path)
}

