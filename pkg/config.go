package pkg

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Opt struct { // 设计参考
	Action      string `yaml:"action"`      // drop keep replacement
	Regex       string `yaml:"regex"`       //
	Replacement string `yaml:"replacement"` //
}

type Config struct {
	//Ak/sk of cloud providers,Obtain the mirror list on the source side
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	//Cloud Provider Region
	RegionAli string `yaml:"region_ali"`
	RegionHw  string `yaml:"region_hw"`
	//Aliyun ACR Account
	UserAli   string `yaml:"user_ali"`
	PasswdAli string `yaml:"passwd_ali"`
	//Huawei SWR Account
	UserHw   string            `yaml:"user_hw"`
	PasswdHw string            `yaml:"passwd_hw"`
	NsOrgMap map[string]string `yaml:"ns_org_map,omitempty"`

	// loop gap 默认10分钟
	LoopGap string `yaml:"loop_gap"`

	// from setting
	FromHwFlag         bool   `yaml:"from_hw_flag"`
	FromHwAccessKey    string `yaml:"from_hw_access_key"`
	FromHwSecretKey    string `yaml:"from_hw_secret_key"`
	FromHwRegion       string `yaml:"from_hw_region"`
	FromHwDockerUser   string `yaml:"from_hw_docker_user"`
	FromHwDockerPasswd string `yaml:"from_hw_docker_passwd"`

	FromAliFlag         bool   `yaml:"from_ali_flag"`
	FromAliAccessKey    string `yaml:"from_ali_access_key"`
	FromAliSecretKey    string `yaml:"from_ali_secret_key"`
	FromAliRegion       string `yaml:"from_ali_region"`
	FromAliDockerUser   string `yaml:"from_ali_docker_user"`
	FromAliDockerPasswd string `yaml:"from_ali_docker_passwd"`

	FromHarborFlag   bool   `yaml:"from_harbor_flag"`
	FromHarborUrl    string `yaml:"from_harbor_url"`
	FromHarborUser   string `yaml:"from_harbor_user"`
	FromHarborPasswd string `yaml:"from_harbor_passwd"`

	// To setting
	ToHwFlag         bool   `yaml:"to_hw_flag"`
	ToHwRegion       string `yaml:"to_hw_region"`
	ToHwDockerUser   string `yaml:"to_hw_docker_user"`
	ToHwDockerPasswd string `yaml:"to_hw_docker_passwd"`

	ToAliFlag         bool   `yaml:"to_ali_flag"`
	ToAliRegion       string `yaml:"to_ali_region"`
	ToAliDockerUser   string `yaml:"to_ali_docker_user"`
	ToAliDockerPasswd string `yaml:"to_ali_docker_passwd"`

	ToHarborFlag   bool   `yaml:"to_harbor_flag"`
	ToHarborUrl    string `yaml:"to_harbor_url"`
	ToHarborUser   string `yaml:"to_harbor_user"`
	ToHarborPasswd string `yaml:"to_harbor_passwd"`

	// 打印会做什么，不实际执行
	DryRunFlag bool `yaml:"dry_run"`

	// transition registry: 当配置为相同的局点时，from 和 to 属于不同帐号，则需要配置。
	TransitionRegistry string `yaml:"transition_registry"`

	// auto create organization（hw） namespace(ali) project(harbor) 自动创建组织，命名空间，项目
	AutoCreateFlag bool   `yaml:"auto_create_to_flag"`
	AutoToHwAK     string `yaml:"auto_create_hw_ak"` // 目标是哪，填上对应的值
	AutoToHwSK     string `yaml:"auto_create_hw_sk"`
	AutoToAliAK    string `yaml:"auto_create_ali_ak"` // 华为SWR，阿里ACR 需要。harbor不需要
	AutoToAliSK    string `yaml:"auto_create_ali_sk"`

	// 映射操作规则 drop keep replace
	OptMap []Opt `yaml:"opt_map,omitempty"`
}

func ReadConfigFromFile(path string) (*Config, error) {
	var config *Config
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(file, &config)
	return config, nil
}
