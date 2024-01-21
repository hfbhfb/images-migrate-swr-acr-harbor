package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	// "github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type AliACR struct {
	c *Config
}

func NewAliACR(cfg *Config) *AliACR {
	return &AliACR{c: cfg}
}

func (h *AliACR) GetAllFromImages() []string {
	client, err := cr.NewClientWithAccessKey(h.c.FromAliRegion, h.c.FromAliAccessKey, h.c.FromAliSecretKey)
	// fmt.Println("11412312")
	if err != nil {
		return []string{}
	}
	// fmt.Println("33333")

	var imageArr []string
	i := 1
	for {
		req := cr.CreateGetRepoListRequest()
		req.PageSize = requests.Integer(strconv.Itoa(99))
		req.Page = requests.Integer(strconv.Itoa(i))
		resp, err := client.GetRepoList(req)
		if err != nil {
			return []string{}
		}

		var res RespStruct
		bs := []byte(resp.GetHttpContentString())

		if err = json.Unmarshal(bs, &res); err != nil {
			return []string{}
		}
		for _, v := range res.Data.Repos {
			// img := fmt.Sprintf("registry.%s.aliyuncs.com", h.c.FromAliRegion)
			// imageArr = append(imageArr, img+"/"+v.RepoNamespace+"/"+v.RepoName)
			imageArr = append(imageArr, v.RepoNamespace+"/"+v.RepoName)
		}

		if len(res.Data.Repos) != 99 {
			break
		}
		i++
	}
	return imageArr
}

type NsRepos struct {
	Namespace string `json:"namespace"`
}
type NsRespData struct {
	Namespaces []NsRepos
}
type NsRespStruct struct {
	Data NsRespData
}

func (h *AliACR) GetToNamespaces() []string {

	client, err := cr.NewClientWithAccessKey(h.c.ToAliRegion, h.c.AutoToAliAK, h.c.AutoToAliSK)
	// fmt.Println("11412312")
	if err != nil {
		return []string{}
	}
	// fmt.Println("4444")

	req := cr.CreateGetNamespaceListRequest()

	resp, err := client.GetNamespaceList(req)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	var res NsRespStruct
	bs := []byte(resp.GetHttpContentString())
	// fmt.Println("33333")

	if err = json.Unmarshal(bs, &res); err != nil {
		fmt.Println(err)
		return []string{}
	}
	// fmt.Println(string(bs))

	var arr []string
	for _, v := range res.Data.Namespaces {
		arr = append(arr, v.Namespace)
	}

	return arr
}

func (h *AliACR) CreateToNameSpace(name string) error {
	return errors.New("阿里不支持自动创建命名空间，需要先提前创建好")
}

func (h *AliACR) tmp2(name string) error {
	client, err := cr.NewClientWithAccessKey(h.c.ToAliRegion, h.c.AutoToAliAK, h.c.AutoToAliSK)
	// fmt.Println("11412312")
	if err != nil {
		return err
	}
	// fmt.Println("33333")

	req := cr.CreateCreateNamespaceRequest()

	req.QueryParams["NamespaceName"] = name
	req.QueryParams["InstanceId"] = ""

	resp, err := client.CreateNamespace(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var res RespStruct
	bs := []byte(resp.GetHttpContentString())
	fmt.Println(string(bs))
	if err = json.Unmarshal(bs, &res); err != nil {
		return err
	}
	return nil
}

func (h *AliACR) tmp1(name string) error {
	config := sdk.NewConfig()

	// https://api.alibabacloud.com/api/cr/2018-12-01/CreateNamespace?tab=DEMO&lang=GO&params={%22NamespaceName%22:%22aaajkj%22,%22InstanceId%22:%22%22}&sdkStyle=old

	// Please ensure that the environment variables ALIBABA_CLOUD_ACCESS_KEY_ID and ALIBABA_CLOUD_ACCESS_KEY_SECRET are set.
	credential := credentials.NewAccessKeyCredential(h.c.AutoToAliAK, h.c.AutoToAliSK)
	/* use STS Token
	credential := credentials.NewStsTokenCredential(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"), os.Getenv("ALIBABA_CLOUD_SECURITY_TOKEN"))
	*/
	client, err := sdk.NewClientWithOptions(h.c.ToAliRegion, config, credential)
	if err != nil {
		panic(err)
	}

	request := requests.NewCommonRequest()

	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "cr.cn-qingdao.aliyuncs.com"
	request.Domain = "cr." + h.c.ToAliRegion + ".aliyuncs.com"
	request.Version = "2018-12-01"
	request.ApiName = "CreateNamespace"
	request.QueryParams["NamespaceName"] = name
	request.QueryParams["InstanceId"] = ""

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
	return nil
}
