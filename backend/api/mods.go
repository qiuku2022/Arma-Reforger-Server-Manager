package api

import (
	"encoding/json"
	"os"
	"path/filepath"

	"arsm/config"
	"arsm/models"
	"github.com/gin-gonic/gin"
)

// 本地模组库文件（存储所有添加过的模组信息）
func getLocalModsLibraryPath() string {
	cfg := config.Get()
	return filepath.Join(cfg.ServerPath, "arsm_mods_library.json")
}

func loadLibraryMods() ([]models.Mod, error) {
	path := getLocalModsLibraryPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return []models.Mod{}, nil
	}
	var mods []models.Mod
	if err := json.Unmarshal(data, &mods); err != nil {
		return []models.Mod{}, nil
	}
	return mods, nil
}

func saveLibraryMods(mods []models.Mod) error {
	path := getLocalModsLibraryPath()
	data, err := json.MarshalIndent(mods, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// 从 config.json 加载当前启用的模组
func loadEnabledMods() ([]models.ModConfig, error) {
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return []models.ModConfig{}, nil
	}
	var serverConfig models.ServerConfig
	if err := json.Unmarshal(data, &serverConfig); err != nil {
		return []models.ModConfig{}, nil
	}
	return serverConfig.Game.Mods, nil
}

// 保存启用模组到 config.json
func saveEnabledMods(modConfigs []models.ModConfig) error {
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	var serverConfig models.ServerConfig
	if err := json.Unmarshal(data, &serverConfig); err != nil {
		return err
	}
	serverConfig.Game.Mods = modConfigs
	newData, _ := json.MarshalIndent(serverConfig, "", "  ")
	return os.WriteFile(configPath, newData, 0644)
}

// GetMods 获取模组列表 (合并 Library 和 Config)
func GetMods(c *gin.Context) {
	libMods, _ := loadLibraryMods()
	enabledMods, _ := loadEnabledMods()

	// 构建启用模组的 Map 用于快速查找
	enabledMap := make(map[string]bool)
	for _, m := range enabledMods {
		enabledMap[m.ModID] = true
	}

	// 检查本地文件状态
	cfg := config.Get()
	addonsDir := filepath.Join(cfg.ServerPath, "addons")

	for i := range libMods {
		// 标记启用状态
		if enabledMap[libMods[i].ID] {
			libMods[i].Enabled = true
		} else {
			libMods[i].Enabled = false
		}
		// 检查是否下载
		modPath := filepath.Join(addonsDir, libMods[i].ID)
		if _, err := os.Stat(modPath); err == nil {
			libMods[i].Downloaded = true
		} else {
			libMods[i].Downloaded = false
		}
	}

	success(c, libMods)
}

// AddMod 添加模组到 Library
func AddMod(c *gin.Context) {
	var mod models.Mod
	if err := c.ShouldBindJSON(&mod); err != nil {
		fail(c, "无效的模组数据")
		return
	}

	libMods, _ := loadLibraryMods()

	// 检查是否已存在
	for _, m := range libMods {
		if m.ID == mod.ID {
			fail(c, "模组已存在")
			return
		}
	}

	// 初始化状态
	mod.Enabled = false
	mod.Downloaded = false
	libMods = append(libMods, mod)

	if err := saveLibraryMods(libMods); err != nil {
		fail(c, "保存失败")
		return
	}

	success(c, nil)
}

// DeleteMod 从 Library 和 Config 中删除模组
func DeleteMod(c *gin.Context) {
	id := c.Param("id")
	libMods, _ := loadLibraryMods()

	var newLibMods []models.Mod
	for _, m := range libMods {
		if m.ID != id {
			newLibMods = append(newLibMods, m)
		}
	}
	saveLibraryMods(newLibMods)

	// 同时从 config.json 移除
	enabledMods, _ := loadEnabledMods()
	var newEnabledMods []models.ModConfig
	for _, m := range enabledMods {
		if m.ModID != id {
			newEnabledMods = append(newEnabledMods, m)
		}
	}
	saveEnabledMods(newEnabledMods)

	success(c, nil)
}

// EnableMod 启用模组 (添加到 config.json)
func EnableMod(c *gin.Context) {
	id := c.Param("id")
	// 从 Library 获取模组信息
	libMods, _ := loadLibraryMods()
	var targetMod *models.Mod
	for i := range libMods {
		if libMods[i].ID == id {
			targetMod = &libMods[i]
			break
		}
	}
	if targetMod == nil {
		fail(c, "模组不存在")
		return
	}

	enabledMods, _ := loadEnabledMods()
	// 检查是否已启用
	for _, m := range enabledMods {
		if m.ModID == id {
			success(c, nil)
			return
		}
	}

	// 添加到 config.json，使用 Library 中存储的 version
	enabledMods = append(enabledMods, models.ModConfig{
		ModID:   targetMod.ID,
		Name:    targetMod.Name,
		Version: targetMod.Version,
	})
	if err := saveEnabledMods(enabledMods); err != nil {
		fail(c, "启用失败: "+err.Error())
		return
	}
	success(c, nil)
}

// DisableMod 禁用模组 (从 config.json 移除)
func DisableMod(c *gin.Context) {
	id := c.Param("id")
	enabledMods, _ := loadEnabledMods()
	var newEnabledMods []models.ModConfig
	for _, m := range enabledMods {
		if m.ModID != id {
			newEnabledMods = append(newEnabledMods, m)
		}
	}
	if err := saveEnabledMods(newEnabledMods); err != nil {
		fail(c, "禁用失败: "+err.Error())
		return
	}
	success(c, nil)
}

// CheckModFiles 检查模组文件
func CheckModFiles(c *gin.Context) {
	id := c.Param("id")
	cfg := config.Get()
	modPath := filepath.Join(cfg.ServerPath, "addons", id)
	if _, err := os.Stat(modPath); err == nil {
		success(c, map[string]bool{"downloaded": true})
	} else {
		success(c, map[string]bool{"downloaded": false})
	}
}
