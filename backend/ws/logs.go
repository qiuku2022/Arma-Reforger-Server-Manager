package ws

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"arsm/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	logBuffer []string
	bufferMu  sync.Mutex
	maxBuffer = 500 // 最大缓存日志行数
)

func init() {
	go watchLogFile()
}

// 获取最新的日志文件
func getLatestLogFile(logDir string) string {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return ""
	}

	var logFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".log") || strings.HasSuffix(entry.Name(), ".rpt")) {
			logFiles = append(logFiles, filepath.Join(logDir, entry.Name()))
		}
	}

	if len(logFiles) == 0 {
		return ""
	}

	// 按修改时间排序，最新的在最后
	sort.Slice(logFiles, func(i, j int) bool {
		infoI, _ := os.Stat(logFiles[i])
		infoJ, _ := os.Stat(logFiles[j])
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	return logFiles[len(logFiles)-1]
}

func watchLogFile() {
	var currentLogPath string
	var file *os.File
	var reader *bufio.Reader

	for {
		cfg := config.Get()
		logDir := filepath.Join(cfg.ServerPath, "profile", "logs") // Reforger with -profile puts logs in profile/logs
		// 某些配置可能在 profile 目录下
		// 这里假设用户配置的 ServerPath 是根目录，logs 在其中
		// 实际上 Reforger server 运行位置可能产生 profile 目录
		
		latest := getLatestLogFile(logDir)
		
		// 如果找到新日志文件，或者是第一次启动
		if latest != "" && latest != currentLogPath {
			if file != nil {
				file.Close()
			}
			
			var err error
			file, err = os.Open(latest)
			if err == nil {
				currentLogPath = latest
				reader = bufio.NewReader(file)
				// 如果是新切换的文件，读取最后部分，或者从头读？
				// 通常 tail 是从末尾开始。为了简单，新文件从头读（如果是刚启动），
				// 但如果是中途切换，可能不想读旧历史。
				// 这里策略：如果是刚启动监控，跳到末尾；如果是检测到新文件（rotate），从头读。
				if len(logBuffer) > 0 { // 已经有缓存，说明是运行中切换
					file.Seek(0, io.SeekStart)
				} else {
					file.Seek(0, io.SeekEnd)
				}
			}
		}

		if file != nil && reader != nil {
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						break // 读到末尾，等待新内容
					}
					// 其他错误，重置
					file.Close()
					file = nil
					currentLogPath = ""
					break
				}

				// 处理日志行
				line = strings.TrimRight(line, "\r\n")
				
				bufferMu.Lock()
				logBuffer = append(logBuffer, line)
				if len(logBuffer) > maxBuffer {
					logBuffer = logBuffer[len(logBuffer)-maxBuffer:]
				}
				bufferMu.Unlock()

				broadcast(line)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

// Broadcast 发送日志到所有连接的客户端
func Broadcast(message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

func broadcast(message string) {
	Broadcast(message)
}

// HandleLogs WebSocket 日志处理
func HandleLogs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// 注册客户端
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	// 发送缓存的日志
	bufferMu.Lock()
	for _, line := range logBuffer {
		conn.WriteMessage(websocket.TextMessage, []byte(line))
	}
	bufferMu.Unlock()

	// 保持连接
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			clientsMu.Lock()
			delete(clients, conn)
			clientsMu.Unlock()
			break
		}
	}
}
