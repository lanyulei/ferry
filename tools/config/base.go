package config

import (
    "github.com/spf13/viper"
    "os"
    "path/filepath"
)

var BaseDir string
var ScriptPath string

func InitBase() {
    execPath, _ := os.Executable()
    BaseDir = filepath.Dir(execPath)
    ScriptPath = viper.GetString("script.path")
    if !filepath.IsAbs(ScriptPath) {
        ScriptPath = filepath.Join(BaseDir, ScriptPath)
    }
}
