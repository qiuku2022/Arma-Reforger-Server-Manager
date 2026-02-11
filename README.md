# ARSM - Arma Reforger Server Manager

一款轻量级的 Arma Reforger 专用服务器管理工具，提供 Web 界面管理服务器配置、模组、RCON 控制台等功能。

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows-lightgrey.svg)

## ✨ 功能特性

- 📊 **仪表盘** - 实时查看系统资源（CPU/内存/磁盘）、服务器状态、实时日志流
- ⚙️ **服务端配置** - 可视化配置服务器参数（网络、RCON、游戏设置、场景等）
- 💾 **预设管理** - 保存多套配置预设，支持快速切换和默认预设自动加载
- 📦 **模组管理** - 本地模组库管理，支持版本号，一键启用/禁用
- 👥 **RCON 控制台** - 实时玩家列表、踢出/封禁玩家、自定义命令控制台
- 🔧 **SteamCMD 自动化** - 一键安装/更新/删除 SteamCMD 和游戏服务端
- 📜 **日志持久化** - 日志自动保存到 localStorage，刷新页面不丢失

## 📋 系统要求

- **操作系统**: Linux x64 (Windows 需自行编译)
- **内存**: 建议 4GB+ (Arma Reforger 服务端需要)
- **磁盘空间**: 10GB+ 可用空间
- **网络**: 需要连接 Steam 和互联网

## 🚀 快速开始

### 1. 下载程序

```bash
# 下载到本地目录
git clone https://github.com/yourusername/arsm.git
cd arsm

# 或者仅下载 release 包
cd /your/path
wget https://your-release-url/arsm-linux-amd64
chmod +x arsm-linux-amd64
```

### 2. 运行程序

```bash
# 直接运行
./arsm-linux-amd64

# 后台运行
nohup ./arsm-linux-amd64 &

# 使用 systemd (推荐)
sudo nano /etc/systemd/system/arsm.service
```

systemd 服务示例：
```ini
[Unit]
Description=ARSM - Arma Reforger Server Manager
After=network.target

[Service]
Type=simple
User=your-username
WorkingDirectory=/home/your-username/arsm
ExecStart=/home/your-username/arsm/arsm-linux-amd64
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### 3. 访问 Web 界面

打开浏览器访问：`http://服务器IP:8080`

默认监听 `0.0.0.0:8080`，如需修改请在运行前设置环境变量：
```bash
export ARSM_PORT=8081
./arsm-linux-amd64
```

## 🔧 配置说明

### 首次配置

1. 点击左侧菜单 **🔧 设置**
2. 配置 **SteamCMD 路径** 和 **游戏服务端路径**
3. 点击 **保存设置**，然后点击 **检测** 验证路径

推荐路径结构：
```
/home/user/
├── steamcmd/           # SteamCMD 安装目录
│   └── steamcmd.sh
└── arma-reforger-server/  # 游戏服务端目录
    ├── ArmaReforgerServer
    └── config.json
```

### 服务端配置

1. 进入 **⚙️ 服务端配置** 页面
2. 配置以下关键参数：
   - **绑定地址**: 服务器内网 IP 或 0.0.0.0
   - **绑定端口**: 游戏端口（默认 2001）
   - **公网地址**: 服务器公网 IP（用于服务器列表）
   - **RCON**: 用于远程管理（建议启用）
   - **场景**: 选择任务地图
   - **模组**: 从模组管理页面添加到配置

3. 保存预设或直接启动服务器

### RCON 配置

RCON 配置现在存储在服务端 `config.json` 中：

```json
{
  "rcon": {
    "address": "",
    "port": 19999,
    "password": "your_secure_password",
    "permission": "admin",
    "blacklist": [],
    "whitelist": []
  }
}
```

**注意**: RCON 内容通过服务端 `config.json` 配置，不在 ARSM 设置中。

## 📝 重要说明

### 关于模组管理

模组数据分为两层存储：

1. **模组库** (`arsm_mods_library.json`) - 所有添加过的模组
2. **服务端配置** (`config.json`) - 仅存储启用的模组列表

**JSON 格式**:
```json
{
  "mods": [
    {
      "modId": "59727DAE364DEADB",
      "name": "WeaponSwitching",
      "version": "1.0.1"
    }
  ]
}
```

### 关于日志持久化

- 日志自动保存到浏览器 localStorage
- 最多保留 1000 条日志
- 存储满时会自动清理一半旧日志
- 点击 **🗑️ 清除** 可手动清除所有日志

### 配置文件位置

ARSM 自身配置存储在：
```
Linux: ~/.config/arsm/config.json
Windows: %APPDATA%/arsm/config.json
```

内容示例：
```json
{
  "steamcmd_path": "/home/user/steamcmd",
  "server_path": "/home/user/arma-reforger-server",
  "default_preset": "Conflict-Everon"
}
```

## 🛠️ 开发指南

### 项目结构

```
arsm/
├── backend/              # Go 后端代码
│   ├── api/             # API 处理器
│   ├── config/          # 配置管理
│   ├── models/          # 数据模型
│   ├── ws/              # WebSocket
│   ├── static/          # 前端静态文件（嵌入）
│   └── main.go          # 入口
├── frontend/            # Vue 3 前端
│   ├── src/
│   │   ├── views/      # 页面组件
│   │   ├── api/        # API 封装
│   │   └── ...
│   └── package.json
├── build.sh             # 构建脚本
└── README.md
```

### 构建开发环境

```bash
# 克隆代码
cd arsm

# 安装前端依赖
cd frontend
npm install

# 安装后端依赖
cd ../backend
go mod tidy

# 使用构建脚本
cd ..
bash build.sh
```

### 单独构建

```bash
# 仅构建前端
cd frontend && npm run build

# 仅构建后端 (当前平台)
cd backend && go build -o arsm

# 交叉编译 (Linux)
GOOS=linux GOARCH=amd64 go build -o arsm-linux-amd64

# 交叉编译 (Windows)
GOOS=windows GOARCH=amd64 go build -o arsm-windows-amd64.exe
```

## 🔌 API 文档

### REST API

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/info` | GET | 获取系统信息 |
| `/api/server/status` | GET | 获取服务器状态 |
| `/api/server/start` | POST | 启动服务器 |
| `/api/server/stop` | POST | 停止服务器 |
| `/api/config` | GET/POST | 获取/保存配置 |
| `/api/mods` | GET/POST | 获取/添加模组 |
| `/api/rcon/players` | GET | 获取玩家列表 |
| `/api/rcon/command` | POST | 发送 RCON 命令 |

### WebSocket

```
ws://host/ws/logs
```

实时推送服务端日志输出（SteamCMD、服务器启动等）。

## 🐛 故障排查

### RCON 页面空白

1. 检查服务端 `config.json` 中 RCON 配置是否存在
2. 检查 RCON 密码是否设置（不少于3个字符，无空格）
3. 检查服务器是否已启动
4. 查看浏览器开发者工具 -> Console 的错误信息

### 无法启动服务器

1. 确认 SteamCMD 和游戏服务端已安装
2. 检查配置文件 JSON 格式是否正确（无注释）
3. 确认端口未被占用（2001, 19999）
4. 检查是否有运行权限

### 模组不生效

1. 确认模组已在模组库中
2. 在服务端配置中启用模组
3. 保存配置后重启服务器
4. 检查模组是否已下载到 `addons/` 目录

### 日志不显示

1. 检查 `/ws/logs` WebSocket 连接是否正常
2. 查看 `profile/logs/` 目录是否有日志文件
3. 刷新页面查看 localStorage 中的 `arsm_logs`

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request。

## 📞 支持

- GitHub Issues: [https://github.com/yourusername/arsm/issues](https://github.com/yourusername/arsm/issues)
- Discord: [Your Discord Link]

---

*ARSM - 让 Arma Reforger 服务器管理更简单*
