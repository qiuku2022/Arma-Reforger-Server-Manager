package models

// SystemInfo 系统信息
type SystemInfo struct {
	OS           string  `json:"os"`
	Arch         string  `json:"arch"`
	Hostname     string  `json:"hostname"`
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryTotal  uint64  `json:"memory_total"`
	MemoryUsed   uint64  `json:"memory_used"`
	MemoryFree   uint64  `json:"memory_free"`
	DiskTotal    uint64  `json:"disk_total"`
	DiskUsed     uint64  `json:"disk_used"`
	DiskFree     uint64  `json:"disk_free"`
}

// ServerStatus 服务端状态
type ServerStatus struct {
	Installed bool   `json:"installed"`
	Running   bool   `json:"running"`
	PID       int    `json:"pid,omitempty"`
	Version   string `json:"version,omitempty"`
}

// SteamCMDStatus SteamCMD状态
type SteamCMDStatus struct {
	Installed bool   `json:"installed"`
	Path      string `json:"path"`
	Version   string `json:"version,omitempty"`
}

// Mod 模组信息
type Mod struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Enabled     bool   `json:"enabled"`
	Downloaded  bool   `json:"downloaded"`
}

// Player 玩家信息
type Player struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Online     bool   `json:"online"`
	OnlineTime int64  `json:"online_time"` // 在线时长（秒）
}

// ServerConfig Arma Reforger 服务器配置
type ServerConfig struct {
	BindAddress   string          `json:"bindAddress"`
	BindPort      int             `json:"bindPort"`
	PublicAddress string          `json:"publicAddress"`
	PublicPort    int             `json:"publicPort"`
	A2S           A2SConfig       `json:"a2s"`
	RCON          *RCONConfig     `json:"rcon,omitempty"`
	Game          GameConfig      `json:"game"`
	Operating     OperatingConfig `json:"operating,omitempty"`
}

type A2SConfig struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type RCONConfig struct {
	Address    string   `json:"address"`
	Port       int      `json:"port"`
	Password   string   `json:"password"`
	Permission string   `json:"permission"`
	Blacklist  []string `json:"blacklist"`
	Whitelist  []string `json:"whitelist"`
}

type GameConfig struct {
	Name               string         `json:"name"`
	Password           string         `json:"password"`
	PasswordAdmin      string         `json:"passwordAdmin"`
	Admins             []string       `json:"admins"`
	ScenarioID         string         `json:"scenarioId"`
	MaxPlayers         int            `json:"maxPlayers"`
	Visible            bool           `json:"visible"`
	CrossPlatform      bool           `json:"crossPlatform"`
	SupportedPlatforms []string       `json:"supportedPlatforms"`
	GameProperties     GameProperties `json:"gameProperties"`
	Mods               []ModConfig    `json:"mods"`
}

type GameProperties struct {
	ServerMaxViewDistance   int  `json:"serverMaxViewDistance"`
	ServerMinGrassDistance  int  `json:"serverMinGrassDistance"`
	NetworkViewDistance     int  `json:"networkViewDistance"`
	DisableThirdPerson      bool `json:"disableThirdPerson"`
	FastValidation          bool `json:"fastValidation"`
	BattlEye                bool `json:"battlEye"`
	VONDisableUI            bool `json:"VONDisableUI"`
	VONDisableDirectSpeechUI bool `json:"VONDisableDirectSpeechUI"`
	VONCanTransmitCrossFaction bool `json:"VONCanTransmitCrossFaction"`
}

type ModConfig struct {
	ModID   string `json:"modId"`
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type OperatingConfig struct {
	LobbyPlayerSynchronise  bool            `json:"lobbyPlayerSynchronise"`
	JoinQueue               JoinQueueConfig `json:"joinQueue"`
	DisableNavmeshStreaming []string        `json:"disableNavmeshStreaming,omitempty"`
}

type JoinQueueConfig struct {
	MaxSize int `json:"maxSize"`
}

// Preset 配置预设
type Preset struct {
	Name   string       `json:"name"`
	Config ServerConfig `json:"config"`
}

// Scenario 官方场景
type Scenario struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Map  string `json:"map"`
	Mode string `json:"mode"`
}
