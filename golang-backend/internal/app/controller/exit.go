package controller

import (
    "net/http"
    "time"
    "fmt"

    dbpkg "flux-panel/golang-backend/internal/db"
    "flux-panel/golang-backend/internal/app/response"
    "github.com/gin-gonic/gin"
)

// POST /api/v1/node/set-exit {nodeId, port, password, method?}
// Creates/updates an SS server service on the selected node with given port/password.
func NodeSetExit(c *gin.Context) {
    var p struct {
        NodeID   int64  `json:"nodeId" binding:"required"`
        Port     int    `json:"port" binding:"required"`
        Password string `json:"password" binding:"required"`
        Method   string `json:"method"`
        // optional extras
        Observer string                 `json:"observer"`
        Limiter  string                 `json:"limiter"`
        RLimiter string                 `json:"rlimiter"`
        Metadata map[string]interface{} `json:"metadata"`
    }
    if err := c.ShouldBindJSON(&p); err != nil {
        c.JSON(http.StatusOK, response.ErrMsg("参数错误")); return
    }
    if p.Port <= 0 || p.Port > 65535 || p.Password == "" {
        c.JSON(http.StatusOK, response.ErrMsg("无效的端口或密码")); return
    }
    if p.Method == "" { p.Method = "AEAD_CHACHA20_POLY1305" }

    // Ensure node exists
    var cnt int64
    dbpkg.DB.Table("node").Where("id = ?", p.NodeID).Count(&cnt)
    if cnt == 0 { c.JSON(http.StatusOK, response.ErrMsg("节点不存在")); return }

    // Build service config and push to node
    name := fmt.Sprintf("exit_ss_%d", p.Port)
    svc := buildSSService(name, p.Port, p.Password, p.Method, map[string]any{
        "observer": p.Observer,
        "limiter":  p.Limiter,
        "rlimiter": p.RLimiter,
        "metadata": p.Metadata,
    })
    if err := sendWSCommand(p.NodeID, "AddService", []map[string]any{svc}); err != nil {
        c.JSON(http.StatusOK, response.ErrMsg("发送到节点失败: "+err.Error())); return
    }

    _ = time.Now() // reserved for future persistence
    c.JSON(http.StatusOK, response.OkMsg("出口节点服务已创建/更新"))
}
