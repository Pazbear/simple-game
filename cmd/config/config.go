package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

func RewriteConfigForTest(testConfig Config) {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	appconfig = testConfig
}

func AppConfig() Config {

	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	return appconfig
}

// || prod
func initProfile() string {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0], pair[1])
	}
	var profile string
	profile = os.Getenv("GO_PROFILE")
	fmt.Println("GOLANG_PROFILE: " + profile)
	if len(profile) <= 0 {
		profile = ""
	}
	fmt.Println("GOLANG_PROFILE: " + profile)
	return profile
}

func setRuntimeConfig(profile string) {
	if profile != "" {
		profile = fmt.Sprintf(".%s", profile)
	}
	viper.SetConfigName(fmt.Sprintf("backup-config%s", profile)) //without extension
	viper.SetConfigType("json")
	viper.AddConfigPath(etcdirConfPath)
	viper.AddConfigPath(workingdirConfPath)
	viper.AddConfigPath(testConfPath)
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = viper.Unmarshal(&appconfig)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func InitConfig() {
	profile := initProfile()
	setRuntimeConfig(profile)
	fmt.Println("init config completed")
}