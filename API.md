# ARSM API 文档

本文档定义了 Arma Reforger Server Manager (ARSM) 的后端接口规范。所有接口均位于 `/api` 路径下（WebSocket 除外）。

## 通用响应结构

所有 JSON 响应遵循以下格式：

```json
{
  "code": 0,          // 0 = 成功, 1 = 失败
  "message": "...",   // 错误信息或 "success"
  "data": { ... }     // 响应数据
}
```

---

## 1. 仪表盘 (Dashboard)

### 系统信息
*   **GET** `/api/system/info`
    *   **描述**: 获取宿主机系统状态（CPU/内存/磁盘）。
    *   **响应**:
        ```json
        {
          "os": "linux",
          "arch": "amd64",
          "hostname": "server-01",
          "cpu_usage": 12.5,      // 百分比
          "memory_total": 16000000000, // 字节
          "memory_used": 8000000000,
          "disk_total": 500000000000,
          "disk_used": 100000000000
        }
        ```

### 服务端控制
*   **GET** `/api/server/status`
    *   **描述**: 获取游戏服务端运行状态。
    *   **响应**:
        ```json
        {
          "installed": true,
          "running": true,
          "pid": 1234
        }
        ```
*   **POST** `/api/server/start`
    *   **描述**: 启动游戏服务端。
*   **POST** `/api/server/stop`
    *   **描述**: 停止游戏服务端。
*   **POST** `/api/server/restart`
    *   **描述**: 重启游戏服务端。

### SteamCMD 管理
*   **GET** `/api/steamcmd/status`
    *   **描述**: 检测 SteamCMD 是否安装。
    *   **响应**: `{"installed": true, "path": "/path/to/steamcmd"}`
*   **POST** `/api/steamcmd/install`
    *   **描述**: 下载并安装 SteamCMD。
*   **POST** `/api/steamcmd/update`
    *   **描述**: 更新 SteamCMD 自身。
*   **DELETE** `/api/steamcmd`
    *   **描述**: 删除 SteamCMD 文件。

### 服务端文件管理
*   **POST** `/api/server/install`
    *   **描述**: 使用 SteamCMD 下载/安装 Arma Reforger Server (AppID 1874900)。
*   **POST** `/api/server/update`
    *   **描述**: 更新游戏服务端。
*   **DELETE** `/api/server`
    *   **描述**: 删除游戏服务端文件。

### 实时日志 (WebSocket)
*   **WS** `/ws/logs`
    *   **描述**: 实时推送服务端控制台日志。
    *   **协议**: WebSocket Text Message。连接后服务器会自动推送新增的日志行。

---

## 2. 服务端配置 (Config)

### 配置文件操作
*   **GET** `/api/config`
    *   **描述**: 读取 `config.json`。
    *   **响应**: 完整的 `ServerConfig` 对象结构。
*   **POST** `/api/config`
    *   **描述**: 保存配置到 `config.json`。
    *   **Body**: 完整的 `ServerConfig` 对象。
*   **POST** `/api/config/import`
    *   **描述**: 上传并导入配置文件。
    *   **Body**: `multipart/form-data` (file)。
*   **GET** `/api/config/export`
    *   **描述**: 下载当前配置文件。

### 预设管理 (Presets)
*   **GET** `/api/config/presets`
    *   **描述**: 获取所有保存的预设名称列表。
*   **POST** `/api/config/presets`
    *   **描述**: 将当前配置保存为新预设。
    *   **Body**: `{"name": "preset_name", "config": {...}}`
*   **DELETE** `/api/config/presets/:name`
    *   **描述**: 删除指定预设。

### 辅助数据
*   **GET** `/api/config/scenarios`
    *   **描述**: 获取官方支持的场景列表（用于下拉选择）。
    *   **响应**: `[{"id": "...", "name": "Conflict - Everon", ...}]`

---

## 3. 模组管理 (Mods)

*   **GET** `/api/mods`
    *   **描述**: 获取所有模组列表（包含启用状态和本地下载状态）。
    *   **响应**:
        ```json
        [
          {
            "id": "5967306236B42D33",
            "name": "Better Loadouts",
            "enabled": true,      // 是否在 config.json 中
            "downloaded": true    // 本地文件是否存在
          }
        ]
        ```
*   **POST** `/api/mods`
    *   **描述**: 添加模组到库（仅记录，不下载）。
    *   **Body**: `{"id": "...", "name": "..."}`
*   **DELETE** `/api/mods/:id`
    *   **描述**: 从库和配置中移除模组。
*   **POST** `/api/mods/:id/enable`
    *   **描述**: 启用模组（写入 config.json）。
*   **POST** `/api/mods/:id/disable`
    *   **描述**: 禁用模组（从 config.json 移除）。
*   **GET** `/api/mods/:id/check`
    *   **描述**: 强制检测模组本地文件状态。

---

## 4. RCON 管理

*   **GET** `/api/rcon/players`
    *   **描述**: 获取在线玩家列表。
    *   **响应**:
        ```json
        [
          {
            "id": "1",
            "name": "PlayerOne",
            "online": true,
            "online_time": 3600
          }
        ]
        ```
*   **POST** `/api/rcon/kick/:id`
    *   **描述**: 踢出指定 ID 的玩家。
*   **POST** `/api/rcon/ban/:id`
    *   **描述**: 封禁指定 ID 的玩家。
*   **POST** `/api/rcon/command`
    *   **描述**: 发送自定义 RCON 命令。
    *   **Body**: `{"command": "#restart"}`

---

## 5. 全局设置 (Settings)

*   **GET** `/api/settings`
    *   **描述**: 获取 ARSM 全局设置（路径、RCON 默认凭据）。
*   **POST** `/api/settings`
    *   **描述**: 更新全局设置。
    *   **Body**:
        ```json
        {
          "steamcmd_path": "/home/user/steamcmd",
          "server_path": "/home/user/arma-reforger-server",
          "rcon_enabled": true,
          "rcon_address": "127.0.0.1",
          "rcon_port": 19999,
          "rcon_password": "..."
        }
        ```
