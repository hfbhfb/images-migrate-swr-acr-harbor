package main

import (
	"errors"
	"flag"
	"fmt"
	"images-migrate/pkg"
	"time"

	"github.com/AliyunContainerService/image-syncer/pkg/client"
	"github.com/AliyunContainerService/image-syncer/pkg/utils"
)

var (
	logPath, configFile, namespaces string

	authFile = "auth.json"

	imageFile = "images.json"

	procNum, retries int

	forceUpdate bool
)

var (
	config *pkg.Config
)

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "config file path")
	flag.StringVar(&logPath, "log", "", "log file path (default in os.Stderr)")
	flag.StringVar(&namespaces, "namespaces", "", "the namespaces that you want to sync,multiple values use ','. by default, all are sync")
	flag.IntVar(&procNum, "proc", 5, "numbers of working goroutines")
	flag.IntVar(&retries, "retries", 2, "times to retry failed task")
	flag.BoolVar(&forceUpdate, "force", false, "force update manifest whether the destination manifest exists")
}

func preCheckFile(cfg *pkg.Config) error {
	if cfg == nil {
		fmt.Println("配置文件没有正确配置。文件内容为空？yaml格式不对？")
		return errors.New("配置文件没有正确配置。文件内容为空？")
	}
	if cfg.AccessKey == "" && cfg.FromHwFlag == false && cfg.FromAliFlag == false && cfg.FromHarborFlag == false {
		fmt.Println("从配置文件校验失败: 源目镜像registry flag没设置")
		return errors.New("从配置文件校验失败: 源目镜像registry flag没设置")
	}
	return nil

}

func v01Opt() {

	if err := pkg.GenAuthFile(authFile, config); err != nil {
		fmt.Println("生成认证文件失败: ", err)
		return
	}
	if err := pkg.GenImagesFile(imageFile, namespaces, config); err != nil {
		fmt.Println("获取镜像列表失败: ", err)
		return
	}
	client, err := client.NewSyncClient("", authFile, imageFile, logPath, procNum, retries,
		utils.RemoveEmptyItems([]string{}), utils.RemoveEmptyItems([]string{}), forceUpdate)
	if err != nil {
		fmt.Println("init sync client error: ", err)
		return
	}
	client.Run()
}

func loopRun() {
	flag.Parse()
	var err error
	config, err = pkg.ReadConfigFromFile(configFile)
	if err != nil {
		fmt.Println("从配置文件读取失败: ", err)
		return
	}

	if preCheckFile(config) != nil {
		return
	}
	if config != nil && (config.FromHwFlag || config.FromAliFlag || config.FromHarborFlag) {
		// fmt.Println("新的逻辑处理")
		pkg.StartOptNew(config)

	} else {
		// v0.1 版本 保留从阿里迁移到华为云的逻辑
		v01Opt()
	}

}

func main() {

	for {
		var st time.Duration
		st = 20 * 365 * 24 * time.Hour
		// fmt.Println("运行")
		loopRun()
		// fmt.Println("sleep 等待下一次循环")
		if config != nil && config.LoopGap != "" {
			// durationString := "3h30m15s"

			// Parse the duration string into a time.Duration
			duration, err := time.ParseDuration(config.LoopGap)
			if err != nil {
				fmt.Printf("sleep 等待下一次循环: %v s\n", int(st/time.Second))
				time.Sleep(st)
			} else {
				if duration < time.Second {
					return // 0s自动退出程序
				}
				fmt.Printf("sleep 等待下一次循环: %v s\n", int(duration/time.Second))
				time.Sleep(duration)
			}

		} else {
			fmt.Printf("sleep 等待下一次循环: %v s\n", int(st/time.Second))
			time.Sleep(st)
		}

	}
}
