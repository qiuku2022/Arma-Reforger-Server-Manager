package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"arsm/config"
	"arsm/models"

	"github.com/gin-gonic/gin"
)

// 官方场景列表
var officialScenarios = []models.Scenario{
	{ID: "{ECC61978EDCC2B5A}Missions/23_Campaign.conf", Name: "Conflict - Everon", Map: "Everon", Mode: "Conflict"},
	{ID: "{C41618FD18E9D714}Missions/23_Campaign_Arland.conf", Name: "Conflict - Arland", Map: "Arland", Mode: "Conflict"},
	{ID: "{C700DB41F0C546E1}Missions/23_Campaign_NorthCentral.conf", Name: "Conflict - Northern Everon", Map: "Everon", Mode: "Conflict"},
	{ID: "{28802845ADA64D52}Missions/23_Campaign_SWCoast.conf", Name: "Conflict - Southern Everon", Map: "Everon", Mode: "Conflict"},
	{ID: "{94992A3D7CE4FF8A}Missions/23_Campaign_Western.conf", Name: "Conflict - Western Everon", Map: "Everon", Mode: "Conflict"},
	{ID: "{FDE33AFE2ED7875B}Missions/23_Campaign_Montignac.conf", Name: "Conflict - Montignac", Map: "Everon", Mode: "Conflict"},
	{ID: "{0220741028718E7F}Missions/23_Campaign_HQC_Everon.conf", Name: "Conflict: HQ Commander - Everon", Map: "Everon", Mode: "Conflict"},
	{ID: "{68D1240A11492545}Missions/23_Campaign_HQC_Arland.conf", Name: "Conflict: HQ Commander - Arland", Map: "Arland", Mode: "Conflict"},
	{ID: "{BB5345C22DD2B655}Missions/23_Campaign_HQC_Cain.conf", Name: "Conflict: HQ Commander - Kolguyev", Map: "Kolguyev", Mode: "Conflict"},
	
	{ID: "{59AD59368755F41A}Missions/21_GM_Eden.conf", Name: "Game Master - Everon", Map: "Everon", Mode: "Game Master"},
	{ID: "{2BBBE828037C6F4B}Missions/22_GM_Arland.conf", Name: "Game Master - Arland", Map: "Arland", Mode: "Game Master"},
	{ID: "{F45C6C15D31252E6}Missions/27_GM_Cain.conf", Name: "Game Master - Kolguyev", Map: "Kolguyev", Mode: "Game Master"},
	
	{ID: "{DAA03C6E6099D50F}Missions/24_CombatOps.conf", Name: "Combat Ops - Arland", Map: "Arland", Mode: "Combat Ops"},
	{ID: "{DFAC5FABD11F2390}Missions/26_CombatOpsEveron.conf", Name: "Combat Ops - Everon", Map: "Everon", Mode: "Combat Ops"},
	{ID: "{CB347F2F10065C9C}Missions/CombatOpsCain.conf", Name: "Combat Ops - Kolguyev", Map: "Kolguyev", Mode: "Combat Ops"},
	
	{ID: "{3F2E005F43DBD2F8}Missions/CAH_Briars_Coast.conf", Name: "Capture & Hold - Briars", Map: "Everon", Mode: "Capture & Hold"},
	{ID: "{F1A1BEA67132113E}Missions/CAH_Castle.conf", Name: "Capture & Hold - Montfort Castle", Map: "Everon", Mode: "Capture & Hold"},
	{ID: "{589945FB9FA7B97D}Missions/CAH_Concrete_Plant.conf", Name: "Capture & Hold - Concrete Plant", Map: "Everon", Mode: "Capture & Hold"},
	{ID: "{9405201CBD22A30C}Missions/CAH_Factory.conf", Name: "Capture & Hold - Almara Factory", Map: "Everon", Mode: "Capture & Hold"},
	{ID: "{1CD06B409C6FAE56}Missions/CAH_Forest.conf", Name: "Capture & Hold - Simon's Wood", Map: "Everon", Mode: "Capture & Hold"},
	{ID: "{7C491B1FCC0FF0E1}Missions/CAH_LeMoule.conf", Name: "Capture & Hold - Le Moule", Map: "Everon", Mode: "Capture & Hold"},
	{ID: "{6EA2E454519E5869}Missions/CAH_Military_Base.conf", Name: "Capture & Hold - Camp Blake", Map: "Everon", Mode: "Capture & Hold"},
	{ID: "{2B4183DF23E88249}Missions/CAH_Morton.conf", Name: "Capture & Hold - Morton", Map: "Everon", Mode: "Capture & Hold"},
	
	{ID: "{002AF7323E0129AF}Missions/Tutorial.conf", Name: "Training", Map: "Arland", Mode: "Tutorial"},
	
	{ID: "{C47A1A6245A13B26}Missions/SP01_ReginaV2.conf", Name: "Elimination", Map: "Arland", Mode: "Singleplayer"},
	{ID: "{0648CDB32D6B02B3}Missions/SP02_AirSupport.conf", Name: "Air Support", Map: "Arland", Mode: "Singleplayer"},
	{ID: "{10B8582BAD9F7040}Missions/Scenario01_Intro.conf", Name: "Operation Omega 01: Over The Hills And Far Away", Map: "Kolguyev", Mode: "Campaign"},
	{ID: "{1D76AF6DC4DF0577}Missions/Scenario02_Steal.conf", Name: "Operation Omega 02: Radio Check", Map: "Kolguyev", Mode: "Campaign"},
	{ID: "{D1647575BCEA5A05}Missions/Scenario03_Villa.conf", Name: "Operation Omega 03: Light In The Dark", Map: "Kolguyev", Mode: "Campaign"},
	{ID: "{6D224A109B973DD8}Missions/Scenario04_Sabotage.conf", Name: "Operation Omega 04: Red Silence", Map: "Kolguyev", Mode: "Campaign"},
	{ID: "{FA2AB0181129CB16}Missions/Scenario05_Hill.conf", Name: "Operation Omega 05: Cliffhanger", Map: "Kolguyev", Mode: "Campaign"},
}

func getConfigPath() string {
	cfg := config.Get()
	return filepath.Join(cfg.ServerPath, "config.json")
}

func getPresetsDir() string {
	cfg := config.Get()
	return filepath.Join(cfg.ServerPath, "presets")
}

// GetConfig 获取服务端配置
func GetConfig(c *gin.Context) {
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		// 返回默认配置
		defaultConfig := models.ServerConfig{
			BindAddress:   "",
			BindPort:      2001,
			PublicAddress: "",
			PublicPort:    2001,
			A2S: models.A2SConfig{
				Address: "",
				Port:    17777,
			},
			RCON: &models.RCONConfig{
				Address:    "",
				Port:       19999,
				Password:   "",
				Permission: "admin",
				Blacklist:  []string{},
				Whitelist:  []string{},
			},
			Game: models.GameConfig{
				Name:          "Arma Reforger Server",
				Password:      "",
				PasswordAdmin: "",
				Admins:        []string{},
				ScenarioID:    "{ECC61978EDCC2B5A}Missions/23_Campaign.conf",
				MaxPlayers:    64,
				Visible:       true,
				CrossPlatform: false,
				SupportedPlatforms: []string{"PLATFORM_PC"},
				GameProperties: models.GameProperties{
					ServerMaxViewDistance:  2500,
					ServerMinGrassDistance: 50,
					NetworkViewDistance:    1500,
					DisableThirdPerson:     false,
					FastValidation:         true,
					BattlEye:               true,
				},
				Mods: []models.ModConfig{},
			},
			Operating: models.OperatingConfig{
				LobbyPlayerSynchronise: true,
				JoinQueue: models.JoinQueueConfig{
					MaxSize: 64,
				},
			},
		}
		success(c, defaultConfig)
		return
	}

	var serverConfig models.ServerConfig
	if err := json.Unmarshal(data, &serverConfig); err != nil {
		fail(c, "配置文件解析失败")
		return
	}

	success(c, serverConfig)
}

// SaveConfig 保存服务端配置
func SaveConfig(c *gin.Context) {
	var serverConfig models.ServerConfig
	if err := c.ShouldBindJSON(&serverConfig); err != nil {
		fail(c, "无效的配置数据")
		return
	}

	// 检查 RCON 是否被禁用
	// 注意：前端如果关闭了 RCON 开关，虽然 serverConfig.RCON 仍然有值（前端发送过来的），
	// 但我们可能需要一种机制来判断用户是否想要禁用它。
	// 然而，models.ServerConfig 中没有 RCONEnabled 字段（这是 config.json 的结构限制）。
	// 前端的做法是：如果 toggle 关闭，就不在 UI 上显示。
	// 为了支持 "删除 rcon 段"，我们约定：如果 password 为空，或者我们增加一个显式的标志？
	// 不，用户说的是 "如果关闭，就把...删掉"。
	// 最简单的办法是前端在 save 时，如果不启用 RCON，就把 rcon 字段置为 null？
	// 但 Go struct 的 JSON unmarshal 如果字段缺省会用零值。
	// 关键是 models.ServerConfig.RCON 已经改为指针 *RCONConfig 了。
	
	// 所以，如果前端传来的数据里 rcon 为 null，serverConfig.RCON 就会是 nil。
	// 或者，如果前端传了值，但我们想在后端根据某种逻辑删掉它？
	// 不，前端应该控制传什么。
	// 但用户是在描述后端逻辑："如果关闭...删掉"。
	// 我们可以简单地检查 RCON.Password 是否为空。如果为空且看起来没配置，就 nil 掉？
	// 不太安全。
	// 最好的方式是前端显式传 null，或者后端判断。
	
	// 这里我们假设：如果 RCON.Password 为空，则视为禁用 RCON？
	// 不，Example.json 里 password 是必须的。
	// 让我们看看前端怎么发。
	
	// 如果 serverConfig.RCON 还是非 nil (因为 gin binding 会实例化它如果 json 有这个 key)，
	// 我们需要一个策略。
	// 既然已经把 models.ServerConfig.RCON 改为 *RCONConfig，
	// 如果前端发送 {"rcon": null} 或者不发 rcon 字段，它就是 nil。
	
	// 验证 RCON 密码
	if serverConfig.RCON != nil {
		if serverConfig.RCON.Password != "" {
			if len(serverConfig.RCON.Password) < 3 || strings.Contains(serverConfig.RCON.Password, " ") {
				fail(c, "RCON密码必须至少3个字符且不能包含空格")
				return
			}
		} else {
			// 如果密码为空，视为禁用？或者强制 nil？
			// 为了满足 "删掉 rcon 段"，我们可以在这里做个判断：
			// 如果密码为空，直接 set 为 nil，这样 json 序列化时就会忽略（omitempty）
			serverConfig.RCON = nil
		}
	}

	configPath := getConfigPath()
	cfg := config.Get()
	os.MkdirAll(cfg.ServerPath, 0755)

	data, err := json.MarshalIndent(serverConfig, "", "    ")
	if err != nil {
		fail(c, "配置序列化失败")
		return
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fail(c, "保存配置失败")
		return
	}

	success(c, nil)
}

// GetPresets 获取预设列表
func GetPresets(c *gin.Context) {
	presetsDir := getPresetsDir()
	os.MkdirAll(presetsDir, 0755)

	entries, err := os.ReadDir(presetsDir)
	if err != nil {
		success(c, []string{})
		return
	}

	var presets []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			presets = append(presets, strings.TrimSuffix(entry.Name(), ".json"))
		}
	}

	success(c, presets)
}

// GetPresetContent 获取指定预设的内容
func GetPresetContent(c *gin.Context) {
	name := c.Param("name")
	presetsDir := getPresetsDir()
	presetPath := filepath.Join(presetsDir, name+".json")

	data, err := os.ReadFile(presetPath)
	if err != nil {
		fail(c, "预设不存在")
		return
	}

	var serverConfig models.ServerConfig
	if err := json.Unmarshal(data, &serverConfig); err != nil {
		fail(c, "预设文件损坏")
		return
	}

	success(c, serverConfig)
}

// SavePreset 保存预设
func SavePreset(c *gin.Context) {
	var req struct {
		Name   string              `json:"name"`
		Config models.ServerConfig `json:"config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "无效的请求数据")
		return
	}

	presetsDir := getPresetsDir()
	os.MkdirAll(presetsDir, 0755)

	presetPath := filepath.Join(presetsDir, req.Name+".json")
	data, _ := json.MarshalIndent(req.Config, "", "    ")

	if err := os.WriteFile(presetPath, data, 0644); err != nil {
		fail(c, "保存预设失败")
		return
	}

	success(c, nil)
}

// DeletePreset 删除预设
func DeletePreset(c *gin.Context) {
	name := c.Param("name")
	presetsDir := getPresetsDir()
	presetPath := filepath.Join(presetsDir, name+".json")

	if err := os.Remove(presetPath); err != nil {
		fail(c, "删除预设失败")
		return
	}

	success(c, nil)
}

// ImportConfig 导入配置
func ImportConfig(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		fail(c, "未找到上传文件")
		return
	}

	f, err := file.Open()
	if err != nil {
		fail(c, "打开文件失败")
		return
	}
	defer f.Close()

	var serverConfig models.ServerConfig
	if err := json.NewDecoder(f).Decode(&serverConfig); err != nil {
		fail(c, "配置文件格式错误")
		return
	}

	success(c, serverConfig)
}

// ExportConfig 导出配置
func ExportConfig(c *gin.Context) {
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		fail(c, "读取配置失败")
		return
	}

	c.Header("Content-Disposition", "attachment; filename=config.json")
	c.Data(200, "application/json", data)
}

// GetScenarios 获取官方场景列表
func GetScenarios(c *gin.Context) {
	success(c, officialScenarios)
}
