# Claude2Api
将Claude的网页服务转为Api服务，支持识图，文件上传，流式传输, 思考输出……

Api支持访问格式为 openai 格式

# Claude2API
[![Go Report Card](https://goreportcard.com/badge/github.com/yushangxiao/claude2api)](https://goreportcard.com/report/github.com/yushangxiao/claude2api)
[![License](https://img.shields.io/github/license/yushangxiao/claude2api)](LICENSE)
|[英文](https://github.com/yushangxiao/claude2api/edit/main/README.md)

提醒： 只有 PRO 用户可以使用所有模型。免费用户只能使用 claude-sonnet-4-20250514

## ✨ 特性
- 🖼️ **图像识别** - 发送图像给Claude进行分析
- 📝 **自动对话管理** - 对话可在使用后自动删除
- 🌊 **流式响应** - 获取Claude实时流式输出
- 📁 **文件上传支持** - 上传长文本内容
- 🧠 **思考过程** - 访问Claude的逐步推理，自动输出`<think>`标签
 - 🔄 **聊天历史管理** - 控制对话上下文长度，超出将上传为文件
 - 🌐 **代理支持** - 通过您首选的代理请求
 - 🔐 **API密钥认证** - 保护您的API端点
 - 🔁 **自动重试** - 请求失败时，自动切换下一个账号
  - 🌐 **直接代理** -  使用 sk-ant-* 直接作为key使用
 ## 📋 前提条件
 - Go 1.23+（从源代码构建）
 - Docker（用于容器化部署）
 
 ## 🚀 部署选项
 ### Docker
 ```bash
 docker run -d \
   -p 8080:8080 \
   -e SESSIONS=sk-ant-sid01-xxxx,sk-ant-sid01-yyyy \
   -e APIKEY=123 \
   -e CHAT_DELETE=true \
   -e MAX_CHAT_HISTORY_LENGTH=10000 \
   -e NO_ROLE_PREFIX=false \
   -e PROMPT_DISABLE_ARTIFACTS=false \
   -e ENABLE_MIRROR_API=false \
   -e MIRROR_API_PREFIX=/mirror \
   --name claude2api \
   ghcr.io/yushangxiao/claude2api:latest
 ```
 
 ### Docker Compose
 创建一个`docker-compose.yml`文件：
 ```yaml
 version: '3'
 services:
   claude2api:
     image: ghcr.io/yushangxiao/claude2api:latest
     container_name: claude2api
     ports:
       - "8080:8080"
     environment:
       - SESSIONS=sk-ant-sid01-xxxx,sk-ant-sid01-yyyy
       - ADDRESS=0.0.0.0:8080
       - APIKEY=123
       - PROXY=http://proxy:2080  # 可选
       - CHAT_DELETE=true
       - MAX_CHAT_HISTORY_LENGTH=10000
       - NO_ROLE_PREFIX=false
       - PROMPT_DISABLE_ARTIFACTS=true
       - ENABLE_MIRROR_API=false
       - MIRROR_API_PREFIX=/mirror
     restart: unless-stopped
 ```
 然后运行：
 ```bash
 docker-compose up -d
 ```
 
 ### Hugging Face Spaces
 您可以使用Docker将此项目部署到Hugging Face Spaces：
 1. Fork Hugging Face Space：[https://huggingface.co/spaces/rclon/claude2api](https://huggingface.co/spaces/rclon/claude2api)
 2. 在设置选项卡中配置您的环境变量
 3. Space将自动部署Docker镜像
 
 注意：在Hugging Face中，/v1可能被屏蔽，您可以使用/hf/v1代替。
 
 ### 直接部署
```bash
# Clone the repository
git clone https://github.com/yushangxiao/claude2api.git
cd claude2api
cp .env.example .env  
vim .env  
# Build the binary
go build -o claude2api .

./claude2api
```
 
 ## ⚙️ 配置
 
### YAML 配置

你可以在应用程序的根目录下使用 config.yaml 文件来配置 Claude2API。如果此文件存在，将会使用它而不是环境变量。

config.yaml 示例：

```yaml
# Sessions configuration
sessions:
  - sessionKey: "sk-ant-sid01-xxxx"
    orgID: ""
  - sessionKey: "sk-ant-sid01-yyyy"
    orgID: ""

# Server address
address: "0.0.0.0:8080"

# API authentication key
apiKey: "123"

# Other configuration options...
chatDelete: true
maxChatHistoryLength: 10000
noRolePrefix: false
promptDisableArtifacts: false
enableMirrorApi: false
mirrorApiPrefix: ""
```

仓库中提供了一个名为 config.yaml.example 的示例配置文件。

 
 | 环境变量 | 描述 | 默认值 |
 |----------------------|-------------|---------|
 | `SESSIONS` | 逗号分隔的Claude API会话密钥列表 | 必填 |
 | `ADDRESS` | 服务器地址和端口 | `0.0.0.0:8080` |
 | `APIKEY` | 用于认证的API密钥 | 必填 |
 | `PROXY` | HTTP代理URL | 可选 |
 | `CHAT_DELETE` | 是否在使用后删除聊天会话 | `true` |
 | `MAX_CHAT_HISTORY_LENGTH` | 超出此长度将文本转为文件 | `10000` |
 | `NO_ROLE_PREFIX` |不在每条消息前添加角色 | `false` |
 | `PROMPT_DISABLE_ARTIFACTS` | 添加提示词尝试禁用 ARTIFACTS| `false` |
 | `ENABLE_MIRROR_API` | 允许直接使用 sk-ant-* 作为 key 使用 | `false` |
 | `MIRROR_API_PREFIX` | 对直接使用增加接口前缀，开启ENABLE_MIRROR_API时必填 | `` |
 
 ## 📝 API使用
 ### 认证
 在请求头中包含您的API密钥：
 ```
 Authorization: Bearer YOUR_API_KEY
 ```
 
 ### 聊天完成
 ```bash
 curl -X POST http://localhost:8080/v1/chat/completions \
   -H "Content-Type: application/json" \
   -H "Authorization: Bearer YOUR_API_KEY" \
   -d '{
     "model": "claude-sonnet-4-5-20250929",
     "messages": [
       {
         "role": "user",
         "content": "你好，Claude！"
       }
     ],
     "stream": true
   }'
 ```
 
 ### 图像分析
 ```bash
 curl -X POST http://localhost:8080/v1/chat/completions \
   -H "Content-Type: application/json" \
   -H "Authorization: Bearer YOUR_API_KEY" \
   -d '{
     "model": "claude-sonnet-4-5-20250929",
     "messages": [
       {
         "role": "user",
         "content": [
           {
             "type": "text",
             "text": "这张图片里有什么？"
           },
           {
             "type": "image_url",
             "image_url": {
               "url": "data:image/jpeg;base64,..."
             }
           }
         ]
       }
     ]
   }'
 ```
 
 ## 🤝 贡献
 欢迎贡献！请随时提交Pull Request。
 1. Fork仓库
 2. 创建特性分支（`git checkout -b feature/amazing-feature`）
 3. 提交您的更改（`git commit -m '添加一些惊人的特性'`）
 4. 推送到分支（`git push origin feature/amazing-feature`）
 5. 打开Pull Request
 
 ## 📄 许可证
 本项目采用MIT许可证 - 详见[LICENSE](LICENSE)文件。
 
 ## 🙏 致谢
 - 感谢[Anthropic](https://www.anthropic.com/)创建Claude
 - 感谢Go社区提供的优秀生态系统

 ## 🎁 项目支持

如果你觉得这个项目对你有帮助，可以考虑通过 [爱发电](https://afdian.com/a/iscoker) 支持我😘
 ---
 由[yushangxiao](https://github.com/yushangxiao)用❤️制作
</details
