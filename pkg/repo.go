package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
)

type Repos struct {
	RepoNamespace string `json:"repoNamespace"`
	RepoName      string `json:"repoName"`
	RegionId      string `json:"regionId,omitempty"`
	RegionType    string `json:"regionType,omitempty"`
}
type RespData struct {
	Repos []Repos
	Total int
}
type RespStruct struct {
	Data RespData
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Insecure bool   `json:"insecure"` // https://pkg.go.dev/github.com/cndoit18/image-syncer#section-readme
}

func GenAuthFile(authFile string, c *Config) error {
	repoMap := make(map[string]interface{})

	aliRepoName := fmt.Sprintf("registry.%s.aliyuncs.com", c.RegionAli)
	hwRepoName := fmt.Sprintf("swr.%s.myhuaweicloud.com", c.RegionHw)
	repoMap[hwRepoName] = Auth{
		Username: c.UserHw,
		Password: c.PasswdHw,
	}
	repoMap[aliRepoName] = Auth{
		Username: c.UserAli,
		Password: c.PasswdAli,
	}
	if err := WriteFile(authFile, repoMap); err != nil {
		return err
	}
	return nil
}

func GenImagesFile(imageFile, namespaces string, c *Config) error {
	client, err := cr.NewClientWithAccessKey(c.RegionAli, c.AccessKey, c.SecretKey)
	if err != nil {
		return err
	}
	aliRepoName := fmt.Sprintf("registry.%s.aliyuncs.com", c.RegionAli)
	hwRepoName := fmt.Sprintf("swr.%s.myhuaweicloud.com", c.RegionHw)
	imagesMap := make(map[string]interface{})
	i := 1
	for {
		req := cr.CreateGetRepoListRequest()
		req.PageSize = requests.Integer(strconv.Itoa(99))
		req.Page = requests.Integer(strconv.Itoa(i))
		resp, err := client.GetRepoList(req)
		if err != nil {
			return err
		}

		var res RespStruct
		bs := []byte(resp.GetHttpContentString())

		if err = json.Unmarshal(bs, &res); err != nil {
			return err
		}

		for _, v := range res.Data.Repos {
			if namespaces != "" && !strings.Contains(namespaces, v.RepoNamespace) {
				continue
			}
			aliImage := fmt.Sprintf("%s/%s/%s", aliRepoName, v.RepoNamespace, v.RepoName)
			hwImage := fmt.Sprintf("%s/%s/%s", hwRepoName, v.RepoNamespace, v.RepoName)

			//If you specify a mapping relationship between namespace and organization, the original name is overwritten
			if r, ok := c.NsOrgMap[v.RepoNamespace]; ok {
				hwImage = fmt.Sprintf("%s/%s/%s", hwRepoName, r, v.RepoName)
			}
			imagesMap[aliImage] = hwImage
		}

		if len(res.Data.Repos) != 99 {
			break
		}
		i++
	}
	if err := WriteFile(imageFile, imagesMap); err != nil {
		return err
	}
	return nil
}

func WriteFile(fileName string, m map[string]interface{}) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(m)
	if err != nil {
		return err
	}
	return nil
}

func LoadFile(fileName string) (error, map[string]string) {

	file, err := os.Open(fileName)
	if err != nil {
		return err, nil
	}
	defer file.Close()

	m := make(map[string]string)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&m)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, nil
	}

	// Access the map[string]string
	fmt.Println("Decoded map:", m)

	return nil, m
}
