# ARSM - Arma Reforger Server Manager

一款轻量级的 Arma Reforger 专用服务器管理工具，提供 Web 界面管理服务器配置、模组、RCON 控制台等功能。

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows-lightgrey.svg)

## ✨ 功能特性

- 📊 **仪表盘** - 实时查看系统资源（CPU/内存/磁盘）、服务器状态、实时日志流
- ⚙️ **服务端配置** - 可视化配置服务器参数（网络、RCON、游戏设置、场景等）
- 💾 **预设管理** - 保存多套配置预设，支持快速切换和默认预设自动加载
- 📦 **模组管理** - 本地模组库管理，支持版本号，一键启用/禁用
- 🔐 **安全认证** - 支持多用户登录、JWT 认证、会话隔离（关闭页面自动登出）
- 👥 **RCON 控制台** - 实时玩家列表、踢出/封禁玩家、自定义命令控制台
- 🔧 **SteamCMD 自动化** - 一键安装/更新/删除 SteamCMD 和游戏服务端
- 📜 **日志持久化** - 日志自动保存到 localStorage，刷新页面不丢失

## 📋 系统要求

| 平台 | 要求 |
|------|------|
| **Linux** | x64 架构，建议 4GB+ 内存 |
| **Windows** | x64 架构，Windows 10/11/Server 2019+ |
| **通用** | 10GB+ 磁盘空间，连接 Steam 的网络 |

## 🚀 快速开始

### 默认账户

| 用户名 | 密码 | 权限 |
|------|------|------|
| `admin` | `admin` | 管理员 |

> 首次登录后请务必在 **🔧 设置** 页面修改默认用户名和密码。

### 下载 Release

| 平台 | 下载 |
|------|------|
| Linux | `arsm-linux-amd64` |
| Windows | `arsm-windows-amd64.exe` |

### Linux 部署

```bash
# 下载并运行
wget https://github.com/qiuku2022/Arma-Reforger-Server-Manager/releases/latest/download/arsm-linux-amd64
chmod +x arsm-linux-amd64
./arsm-linux-amd64

# 或使用 systemd
sudo nano /etc/systemd/system/arsm.service
```

### Windows 部署

```powershell
# 下载后双击运行，或使用命令行
arsm-windows-amd64.exe

# 指定端口
set PORT=8081
arsm-windows-amd64.exe
```

### 访问 Web 界面

打开浏览器访问：`http://localhost:8080`

## 🔧 配置说明

### 首次配置

1. 点击左侧菜单 **🔧 设置**
2. 配置 **SteamCMD 路径** 和 **游戏服务端路径**

**Linux 推荐路径**
```
/home/user/
├── steamcmd/
└── arma-reforger-server/
```

**Windows 推荐路径**
```
C:\
├── steamcmd\
└── ArmaReforgerServer\
```

### RCON 配置

RCON 配置存储在服务端 `config.json`：

```json
{
  "rcon": {
    "address": "",
    "port": 19999,
    "password": "your_password",
    "permission": "admin",
    "blacklist": [],
    "whitelist": []
  }
}
```

## 🛠️ 构建指南

### 环境要求

- Go 1.21+
- Node.js 18+

### 构建脚本

```bash
# 构建当前平台
bash build.sh

# 构建 Windows 版本
bash build.sh windows

# 构建所有平台
bash build.sh all
```

### 手动构建

```bash
# 前端
cd frontend
npm install
npm run build

# 后端 (Linux)
cd ../backend
GOOS=linux GOARCH=amd64 go build -o arsm-linux-amd64

# 后端 (Windows)
GOOS=windows GOARCH=amd64 go build -o arsm-windows-amd64.exe
```

## 📁 项目结构

```
arsm/
├── backend/          # Go 后端
│   ├── api/         # API 处理器
│   ├── config/      # 配置管理
│   └── main.go      # 入口
├── frontend/         # Vue 3 前端
│   └── src/
│       └── views/   # 页面组件
├── build.sh          # 构建脚本
└── README.md
```

## 🔌 API 文档

### REST API

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/info` | GET | 系统信息 |
| `/api/server/start` | POST | 启动服务器 |
| `/api/server/stop` | POST | 停止服务器 |
| `/api/config` | GET/POST | 配置管理 |
| `/api/rcon/players` | GET | 玩家列表 |

### WebSocket

```
ws://host/ws/logs
```

实时推送服务端日志。

## 🐛 故障排查

### RCON 页面空白

1. 检查 `config.json` 中 RCON 配置
2. 确认 RCON 密码已设置（≥3字符，无空格）
3. 确认服务器已启动

### Windows 进程无法终止

- 程序使用 `taskkill` 优雅终止进程
- 如遇顽固进程，可手动在任务管理器结束

### 日志不显示

1. 检查 `profile/logs/` 目录
2. 确认 WebSocket 连接正常

## 📄 许可证

MIT License

## 📞 支持

- GitHub: [https://github.com/qiuku2022/Arma-Reforger-Server-Manager](https://github.com/qiuku2022/Arma-Reforger-Server-Manager)

---

*ARSM - 让 Arma Reforger 服务器管理更简单*
