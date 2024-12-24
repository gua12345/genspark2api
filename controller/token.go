package controller

import (
    "io/ioutil"
    "net/http"
    "os"
    "genspark2api/common/config"
    "github.com/gin-gonic/gin"
    "sync"
    "bufio"
)

var (
    tokenFileMutex sync.Mutex  // 添加文件写入锁
)

const tokenFilePath = "/app/genspark2api/data/token.txt"

type TokenController struct{}

// 验证路径密码
func validatePathPassword(c *gin.Context) bool {
    password := c.Param("password")
    if password != config.TokenOperationPassword {
        c.JSON(http.StatusUnauthorized, gin.H{
            "code": http.StatusUnauthorized,
            "message": "无效的访问密码",
            "data": nil,
        })
        return false
    }
    return true
}

// GetTokens 获取所有token
func (t *TokenController) GetTokens(c *gin.Context) {
    if !validatePathPassword(c) {
        return
    }

    content, err := ioutil.ReadFile(tokenFilePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": http.StatusInternalServerError,
            "message": "读取文件失败",
            "data": nil,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code": http.StatusOK,
        "message": "获取成功",
        "data": string(content),
    })
}

// AppendToken 追加token
func (t *TokenController) AppendToken(c *gin.Context) {
    if !validatePathPassword(c) {
        return
    }

    var req struct {
        Token string `json:"token" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "message": "无效的请求参数",
            "data": nil,
        })
        return
    }

    // 检查token是否为空
    if req.Token == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "message": "token不能为空",
            "data": nil,
        })
        return
    }

    // 使用互斥锁保护文件写入
    tokenFileMutex.Lock()
    defer tokenFileMutex.Unlock()

    // 以追加模式打开文件
    f, err := os.OpenFile(tokenFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": http.StatusInternalServerError,
            "message": "无法打开文件",
            "data": nil,
        })
        return
    }
    defer f.Close()

    // 使用缓冲写入提高性能
    writer := bufio.NewWriter(f)
    if _, err := writer.WriteString(req.Token + "\n"); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": http.StatusInternalServerError,
            "message": "写入文件失败",
            "data": nil,
        })
        return
    }

    // 确保数据写入磁盘
    if err := writer.Flush(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": http.StatusInternalServerError,
            "message": "写入文件失败",
            "data": nil,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code": http.StatusOK,
        "message": "Token添加成功",
        "data": nil,
    })
}

// ClearTokens 清空token文件
func (t *TokenController) ClearTokens(c *gin.Context) {
    if !validatePathPassword(c) {
        return
    }

    // 以只写模式打开文件，并清空内容
    if err := os.WriteFile(tokenFilePath, []byte(""), 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": http.StatusInternalServerError,
            "message": "清空文件失败",
            "data": nil,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code": http.StatusOK,
        "message": "Token文件已清空",
        "data": nil,
    })
}

// TokenPage 返回token管理页面
func (t *TokenController) TokenPage(c *gin.Context) {
    if !validatePathPassword(c) {
        return
    }
    
    html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Token 管理系统</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: "Microsoft YaHei", "微软雅黑", sans-serif;
            background-color: #f7f1e6;
            margin: 0;
            padding: 20px;
            display: flex;
            justify-content: center;
            min-height: 100vh;
        }
        .container {
            background-color: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 4px 20px rgba(139, 69, 19, 0.1);
            max-width: 800px;
            width: 100%;
            margin-top: 20px;
            position: relative;
            border: 1px solid rgba(139, 69, 19, 0.2);
        }
        .title {
            text-align: center;
            color: #8b4513;
            font-size: 28px;
            margin-bottom: 30px;
            padding-bottom: 15px;
            border-bottom: 2px solid #d4a682;
        }
        .stats {
            text-align: center;
            margin-bottom: 25px;
            color: #5c3317;
            font-size: 18px;
            padding: 10px;
            background-color: rgba(212, 166, 130, 0.1);
            border-radius: 4px;
        }
        .input-group { margin-bottom: 20px; width: 100%; }
        textarea {
            width: 100%;
            height: 150px;
            padding: 15px;
            border: 1px solid #d4a682;
            border-radius: 4px;
            font-size: 14px;
            margin-bottom: 10px;
            resize: vertical;
            font-family: monospace;
            background-color: #fff9f0;
        }
        .button-group {
            display: flex;
            justify-content: center;
            gap: 15px;
            margin-bottom: 20px;
        }
        button {
            background-color: #8b4513;
            color: white;
            border: none;
            padding: 10px 25px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            min-width: 120px;
        }
        button:disabled {
            background-color: #cccccc;
            cursor: not-allowed;
        }
        .clear-btn { background-color: #c41e3a; }
        .message {
            position: fixed;
            top: 20px;
            left: 50%;
            transform: translateX(-50%);
            padding: 12px 25px;
            border-radius: 4px;
            display: none;
            z-index: 1000;
        }
        .success { background-color: #dff0d8; color: #3c763d; }
        .error { background-color: #f2dede; color: #a94442; }
        .progress {
            text-align: center;
            margin-top: 10px;
            color: #8b4513;
            display: none;
        }
        .progress-detail {
            font-size: 14px;
            color: #666;
            margin-top: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="title">Token 管理系统</div>
        <div class="stats">当前 Token 数量: <span id="tokenCount">0</span></div>
        <div class="input-group">
            <textarea id="newTokens" placeholder="请输入Token，每行一个"></textarea>
            <div class="tips">
                使用说明：
                <br>• 每行输入一个Token
                <br>• 支持一次性添加多个Token
                <br>• 使用 Ctrl + Enter 快捷键添加
            </div>
        </div>
        <div class="progress" id="progress">
            处理进度: <span id="progressText">0/0</span>
            <div class="progress-detail">
                成功: <span id="successCount">0</span>
                失败: <span id="failCount">0</span>
            </div>
        </div>
        <div class="button-group">
            <button id="addButton" onclick="addTokens()">批量添加</button>
            <button id="clearButton" class="clear-btn" onclick="clearTokens()">清空所有</button>
        </div>
        <div id="message" class="message"></div>
    </div>

    <script>
        const password = window.location.pathname.split("/")[1];
        const addButton = document.getElementById("addButton");
        const clearButton = document.getElementById("clearButton");
        const progress = document.getElementById("progress");
        const progressText = document.getElementById("progressText");
        const successCountElement = document.getElementById("successCount");
        const failCountElement = document.getElementById("failCount");

        function showMessage(text, isError = false) {
            const msg = document.getElementById("message");
            msg.textContent = text;
            msg.style.display = "block";
            msg.className = "message " + (isError ? "error" : "success");
            setTimeout(() => {
                msg.style.display = "none";
            }, 3000);
        }

        function updateProgress(completed, total, success, fail) {
            progressText.textContent = completed + "/" + total;
            successCountElement.textContent = success;
            failCountElement.textContent = fail;
        }

        async function loadTokens() {
            try {
                const response = await fetch("/" + password + "/token/list");
                const text = await response.text();
                const data = JSON.parse(text);
                if (data.code === 200 && data.data) {
                    const tokens = data.data.trim().split("\n").filter(t => t);
                    document.getElementById("tokenCount").textContent = tokens.length;
                }
            } catch (error) {
                showMessage("加载失败: " + error.message, true);
            }
        }

        function delay(ms) {
            return new Promise(resolve => setTimeout(resolve, ms));
        }

        async function processToken(token, retries = 3) {
            for (let attempt = 0; attempt < retries; attempt++) {
                try {
                    const response = await fetch("/" + password + "/token/append", {
                        method: "POST",
                        headers: {"Content-Type": "application/json"},
                        body: JSON.stringify({token: token.trim()})
                    });

                    if (!response.ok && response.status === 429) {
                        await delay(1000 * (attempt + 1));
                        continue;
                    }

                    const result = await response.json();
                    return {
                        success: result.code === 200,
                        message: result.message
                    };
                } catch (error) {
                    if (attempt === retries - 1) {
                        return {
                            success: false,
                            message: error.message
                        };
                    }
                    await delay(1000 * (attempt + 1));
                }
            }
            return {
                success: false,
                message: "重试次数已用完"
            };
        }

        async function addTokens() {
            const textarea = document.getElementById("newTokens");
            const tokens = textarea.value.trim().split("\n").filter(t => t.trim());
            
            if (tokens.length === 0) {
                showMessage("请输入至少一个Token", true);
                return;
            }

            addButton.disabled = true;
            clearButton.disabled = true;
            progress.style.display = "block";
            let completed = 0;
            let successCount = 0;
            let failCount = 0;

            const batchSize = 20;

            try {
                for (let i = 0; i < tokens.length; i += batchSize) {
                    const batch = tokens.slice(i, i + batchSize);
                    const batchPromises = batch.map(token => 
                        processToken(token)
                            .then(result => {
                                completed++;
                                if (result.success) {
                                    successCount++;
                                } else {
                                    failCount++;
                                }
                                updateProgress(completed, tokens.length, successCount, failCount);
                                return result;
                            })
                    );

                    await Promise.all(batchPromises);
                    if (i + batchSize < tokens.length) {
                        await delay(1000);
                    }
                }

                if (successCount > 0) {
                    showMessage("成功添加 " + successCount + " 个Token，失败 " + failCount + " 个");
                    textarea.value = "";
                    loadTokens();
                } else {
                    showMessage("添加失败", true);
                }
            } catch (error) {
                showMessage("添加失败: " + error.message, true);
            } finally {
                addButton.disabled = false;
                clearButton.disabled = false;
                progress.style.display = "none";
            }
        }

        async function clearTokens() {
            if (!confirm("确定要清空所有 Token 吗？此操作不可恢复！")) return;
            
            clearButton.disabled = true;
            try {
                const response = await fetch("/" + password + "/token/clear", {
                    method: "POST"
                });
                const data = await response.json();
                if (data.code === 200) {
                    showMessage("所有 Token 已清空");
                    loadTokens();
                } else {
                    showMessage(data.message || "清空失败", true);
                }
            } catch (error) {
                showMessage("清空失败: " + error.message, true);
            } finally {
                clearButton.disabled = false;
            }
        }

        document.addEventListener("keydown", function(e) {
            if (e.ctrlKey && e.key === "Enter") {
                addTokens();
            }
        });

        loadTokens();
    </script>
</body>
</html>`
    c.Header("Content-Type", "text/html; charset=utf-8")
    c.String(http.StatusOK, html)
}
