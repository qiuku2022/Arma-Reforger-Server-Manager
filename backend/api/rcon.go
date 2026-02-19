package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"arsm/config"
	"arsm/models"
	"github.com/gin-gonic/gin"
	"github.com/multiplay/go-battleye"
)

// RCON 日志通道（用于实时推送）
var rconLogChan = make(chan string, 100)

// 从服务端 config.json 读取 RCON 配置
func getServerRCONConfig() (string, int, string, bool) {
	cfg := config.Get()
	configPath := filepath.Join(cfg.ServerPath, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", 0, "", false
	}
	var serverConfig models.ServerConfig
	if err := json.Unmarshal(data, &serverConfig); err != nil {
		return "", 0, "", false
	}
	if serverConfig.RCON == nil {
		return "", 0, "", false
	}
	return serverConfig.RCON.Address, serverConfig.RCON.Port, serverConfig.RCON.Password, true
}

// 创建 RCON 客户端
func getRCONClient() (*battleye.Client, error) {
	address, port, password, enabled := getServerRCONConfig()
	if !enabled {
		return nil, fmt.Errorf("RCON 未启用")
	}
	if password == "" {
		return nil, fmt.Errorf("RCON 密码未设置")
	}
	// 如果地址为空，使用本机
	if address == "" {
		address = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%d", address, port)
	c, err := battleye.NewClient(addr, password)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// 记录 RCON 命令日志
func logRCON(command, response string) {
	timestamp := time.Now().Format("15:04:05")
	logLine := fmt.Sprintf("[%s] > %s\n%s", timestamp, command, response)
	select {
	case rconLogChan <- logLine:
	default:
		// 通道满，丢弃旧日志
	}
}

// 解析 Arma Reforger RCON players 命令输出
// Arma Reforger 格式示例：
// #0 111.111.111.111:12345 <ping:123> <guid:abc123...> PlayerName
// (1 players in total)
func parsePlayers(output string) []models.Player {
	var players []models.Player
	lines := strings.Split(output, "\n")

	// 匹配行格式: #ID IP:Port <ping:X> <guid:Y> Name
	re := regexp.MustCompile(`^#(\d+)\s+([\d.:]+)\s+<ping:(\d+)>\s+<guid:([a-f0-9]+)>\s*(.+)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "players in total") {
			continue
		}
		if strings.HasPrefix(line, "Players on server") {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 6 {
			player := models.Player{
				ID:       matches[1], // Session ID
				Name:     strings.TrimSpace(matches[5]),
				Online:   true,
				OnlineTime: 0, // Arma Reforger 不直接提供在线时长
			}
			players = append(players, player)
		}
	}
	return players
}

// GetPlayers 获取玩家列表
func GetPlayers(c *gin.Context) {
	client, err := getRCONClient()
	if err != nil {
		success(c, []models.Player{})
		return
	}
	defer client.Close()

	resp, err := client.Exec("players")
	if err != nil {
		fail(c, "获取玩家列表失败: "+err.Error())
		return
	}

	logRCON("players", resp)
	players := parsePlayers(resp)
	success(c, players)
}

// KickPlayer 踢出玩家
func KickPlayer(c *gin.Context) {
	id := c.Param("id")
	client, err := getRCONClient()
	if err != nil {
		fail(c, "RCON 连接失败: "+err.Error())
		return
	}
	defer client.Close()

	command := fmt.Sprintf("kick %s Kicked by Admin", id)
	resp, err := client.Exec(command)
	if err != nil {
		fail(c, "踢出失败: "+err.Error())
		return
	}

	logRCON(command, resp)
	success(c, nil)
}

// BanPlayer 封禁玩家
func BanPlayer(c *gin.Context) {
	id := c.Param("id")
	client, err := getRCONClient()
	if err != nil {
		fail(c, "RCON 连接失败: "+err.Error())
		return
	}
	defer client.Close()

	// ban <id> [minutes] [reason] - 0 minutes = permanent
	command := fmt.Sprintf("ban %s 0 Banned by Admin", id)
	resp, err := client.Exec(command)
	if err != nil {
		fail(c, "封禁失败: "+err.Error())
		return
	}

	logRCON(command, resp)
	success(c, nil)
}

// SendRCONCommand 发送RCON命令
func GetRCONStatus(c *gin.Context) {
	client, err := getRCONClient()
	if err != nil {
		success(c, map[string]bool{"connected": false})
		return
	}
	defer client.Close()
	success(c, map[string]bool{"connected": true})
}

// SendRCONCommand 发送RCON命令
func SendRCONCommand(c *gin.Context) {
	var req struct {
		Command string `json:"command"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "无效的命令")
		return
	}

	client, err := getRCONClient()
	if err != nil {
		fail(c, "RCON 连接失败: "+err.Error())
		return
	}
	defer client.Close()

	resp, err := client.Exec(req.Command)
	if err != nil {
		fail(c, "命令执行失败: "+err.Error())
		return
	}

	logRCON(req.Command, resp)
	success(c, map[string]string{"response": resp})
}

// GetRCONLogs 获取最近的 RCON 日志（WebSocket 或轮询）
func GetRCONLogs(c *gin.Context) {
	var logs []string
	size := len(rconLogChan)
	for i := 0; i < size && i < 50; i++ {
		select {
		case log := <-rconLogChan:
			logs = append(logs, log)
		default:
			break
		}
	}
	success(c, logs)
}
