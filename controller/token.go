package controller

import (
    "io/ioutil"
    "net/http"
    "os"
    "genspark2api/common/config"
    "github.com/gin-gonic/gin"
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

    // 写入新行
    if _, err := f.WriteString(req.Token + "\n"); err != nil {
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
    
    html := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Token 管理系统</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
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
        .container::before {
            content: '';
            position: absolute;
            top: 5px;
            left: 5px;
            right: 5px;
            bottom: 5px;
            border: 1px solid rgba(139, 69, 19, 0.1);
            border-radius: 6px;
            pointer-events: none;
        }
        .title {
            text-align: center;
            color: #8b4513;
            font-size: 28px;
            margin-bottom: 30px;
            padding-bottom: 15px;
            border-bottom: 2px solid #d4a682;
            position: relative;
        }
        .title::after {
            content: '';
            position: absolute;
            bottom: -2px;
            left: 50%;
            transform: translateX(-50%);
            width: 100px;
            height: 2px;
            background-color: #8b4513;
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
        .input-group {
            margin-bottom: 20px;
            width: 100%;
        }
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
            transition: border-color 0.3s, box-shadow 0.3s;
        }
        textarea:focus {
            outline: none;
            border-color: #8b4513;
            box-shadow: 0 0 5px rgba(139, 69, 19, 0.2);
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
            transition: all 0.3s ease;
            font-size: 16px;
            min-width: 120px;
        }
        button:hover {
            background-color: #5c3317;
            transform: translateY(-1px);
            box-shadow: 0 2px 5px rgba(0,0,0,0.2);
        }
        button:active {
            transform: translateY(0);
            box-shadow: none;
        }
        .clear-btn {
            background-color: #c41e3a;
        }
        .clear-btn:hover {
            background-color: #a01830;
        }
        .tips {
            color: #8b4513;
            font-size: 14px;
            margin: 5px 0 15px;
            padding: 10px;
            background-color: rgba(212, 166, 130, 0.1);
            border-radius: 4px;
            border-left: 3px solid #d4a682;
        }
        .token-list {
            max-height: 400px;
            overflow-y: auto;
            border: 1px solid #d4a682;
            border-radius: 4px;
            background-color: #fff9f0;
            padding: 10px;
        }
        .token-list::-webkit-scrollbar {
            width: 8px;
        }
        .token-list::-webkit-scrollbar-track {
            background: #fff9f0;
        }
        .token-list::-webkit-scrollbar-thumb {
            background-color: #d4a682;
            border-radius: 4px;
        }
        .token-item {
            padding: 10px;
            border-bottom: 1px solid #e8d5c5;
            display: flex;
            justify-content: space-between;
            align-items: center;
            transition: background-color 0.3s;
        }
        .token-item:last-child {
            border-bottom: none;
        }
        .token-item:hover {
            background-color: rgba(212, 166, 130, 0.1);
        }
        .token-text {
            word-break: break-all;
            font-family: monospace;
            color: #5c3317;
        }
        .message {
            position: fixed;
            top: 20px;
            left: 50%;
            transform: translateX(-50%);
            padding: 12px 25px;
            border-radius: 4px;
            display: none;
            z-index: 1000;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .success {
            background-color: #dff0d8;
            color: #3c763d;
            border: 1px solid #d6e9c6;
        }
        .error {
            background-color: #f2dede;
            color: #a94442;
            border: 1px solid #ebccd1;
        }
        @media (max-width: 600px) {
            .container {
                padding: 15px;
                margin-top: 10px;
            }
            .title {
                font-size: 24px;
            }
            .button-group {
                flex-direction: column;
                gap: 10px;
            }
            button {
                width: 100%;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="title">Token 管理系统</div>
        <div class="stats">当前 Token 数量: <span id="tokenCount">0</span></div>
        <div class="input-group">
            <textarea id="newTokens" placeholder="请输入Token，每行一个&#10;例如：&#10;token1&#10;token2&#10;token3"></textarea>
            <div class="tips">
                📝 使用说明：
                <br>• 每行输入一个Token
                <br>• 支持一次性添加多个Token
                <br>• 使用 Ctrl + Enter 快捷键添加
            </div>
        </div>
        <div class="button-group">
            <button onclick="addTokens()">批量添加</button>
            <button class="clear-btn" onclick="clearTokens()">清空所有</button>
        </div>
        <div class="token-list" id="tokenList"></div>
        <div id="message" class="message"></div>
    </div>
    <script>
        const password = window.location.pathname.split('/')[1];

        function showMessage(text, isError = false) {
            const msg = document.getElementById('message');
            msg.textContent = text;
            msg.style.display = 'block';
            msg.className = 'message ' + (isError ? 'error' : 'success');
            setTimeout(() => {
                msg.style.opacity = '0';
                msg.style.transition = 'opacity 0.5s ease';
                setTimeout(() => {
                    msg.style.display = 'none';
                    msg.style.opacity = '1';
                    msg.style.transition = '';
                }, 500);
            }, 3000);
        }

        async function loadTokens() {
            try {
                const response = await fetch('/' + password + '/token/list');
                const text = await response.text();
                
                let data;
                try {
                    data = JSON.parse(text);
                } catch (e) {
                    throw new Error(`Invalid JSON response: ${text}`);
                }

                if (data.code === 200 && data.data) {
                    const tokens = data.data.trim().split('\n').filter(t => t);
                    document.getElementById('tokenCount').textContent = tokens.length;
                    const tokenList = document.getElementById('tokenList');
                    tokenList.innerHTML = tokens.map(token =>
                        '<div class="token-item">' +
                        '<span class="token-text">' + token + '</span>' +
                        '</div>'
                    ).join('');
                } else {
                    throw new Error(data.message || 'Unknown error');
                }
            } catch (error) {
                showMessage('加载失败: ' + error.message, true);
            }
        }

        async function addTokens() {
            const textarea = document.getElementById('newTokens');
            const tokens = textarea.value.trim().split('\n').filter(t => t.trim());
            
            if (tokens.length === 0) {
                showMessage('请输入至少一个Token', true);
                return;
            }

            let successCount = 0;
            let failCount = 0;

            for (const token of tokens) {
                try {
                    const response = await fetch('/' + password + '/token/append', {
                        method: 'POST',
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify({token: token.trim()})
                    });
                    const data = await response.json();
                    if (data.code === 200) {
                        successCount++;
                    } else {
                        failCount++;
                    }
                } catch (error) {
                    failCount++;
                }
            }

            if (successCount > 0) {
                showMessage('成功添加 ' + successCount + ' 个Token' + (failCount > 0 ? '，失败 ' + failCount + ' 个' : ''));
                textarea.value = '';
                loadTokens();
            } else {
                showMessage('添加失败', true);
            }
        }

        async function clearTokens() {
            if (!confirm('确定要清空所有 Token 吗？此操作不可恢复！')) return;
            
            try {
                const response = await fetch('/' + password + '/token/clear', {
                    method: 'POST'
                });
                const data = await response.json();
                if (data.code === 200) {
                    showMessage('所有 Token 已清空');
                    loadTokens();
                } else {
                    showMessage(data.message || '清空失败', true);
                }
            } catch (error) {
                showMessage('清空失败: ' + error.message, true);
            }
        }

        document.addEventListener('keydown', function(e) {
            if (e.ctrlKey && e.key === 'Enter') {
                addTokens();
            }
        });

        // 页面加载完成后自动加载 Token 列表
        loadTokens();
    </script>
</body>
</html>
`
    c.Header("Content-Type", "text/html; charset=utf-8")
    c.String(http.StatusOK, html)
}
