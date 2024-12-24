package config

import (
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
    "genspark2api/common/env"
)

var (
    ApiSecret = os.Getenv("API_SECRET")
    ApiSecrets = strings.Split(os.Getenv("API_SECRET"), ",")
    TokenOperationPassword = os.Getenv("TOKEN_OPERATION_PASSWORD")
    GSCookies []string
    AutoDelChat = env.Int("AUTO_DEL_CHAT", 0)
    AllDialogRecordEnable = os.Getenv("ALL_DIALOG_RECORD_ENABLE")
    RequestOutTime = os.Getenv("REQUEST_OUT_TIME")
    StreamRequestOutTime = os.Getenv("STREAM_REQUEST_OUT_TIME")
    SwaggerEnable = os.Getenv("SWAGGER_ENABLE")
    OnlyOpenaiApi = os.Getenv("ONLY_OPENAI_API")
    DebugEnabled = os.Getenv("DEBUG") == "true"
    RateLimitKeyExpirationDuration = 20 * time.Minute
    RequestOutTimeDuration = 5 * time.Minute
    RequestRateLimitNum = env.Int("REQUEST_RATE_LIMIT", 60)
    RequestRateLimitDuration int64 = 1 * 60
)

func init() {
    // 指定 token.txt 的路径（与 Dockerfile 中的路径保持一致）
    tokenPath := "/app/genspark2api/data/token.txt"

    content, err := os.ReadFile(tokenPath)
    if err != nil {
        log.Printf("无法读取 %s: %v", tokenPath, err)
        // 尝试在当前目录查找
        currentDir, _ := os.Getwd()
        alterPath := filepath.Join(currentDir, "token.txt")
        content, err = os.ReadFile(alterPath)
        if err != nil {
            log.Fatalf("无法读取token文件 %s 和 %s: %v", tokenPath, alterPath, err)
        }
    }

    for _, line := range strings.Split(string(content), "\n") {
        if trimmed := strings.TrimSpace(line); trimmed != "" {
            GSCookies = append(GSCookies, trimmed)
        }
    }
        // 添加检查
    if len(GSCookies) == 0 {
        log.Printf("警告: token.txt 是空的")
    } else {
        log.Printf("成功加载 %d 个token", len(GSCookies))
    }
        // 检查 TOKEN_OPERATION_PASSWORD 是否设置
    if TokenOperationPassword == "" {
        log.Printf("警告: TOKEN_OPERATION_PASSWORD 未设置，使用默认密码: admin")
        TokenOperationPassword = "admin"
    }
}
