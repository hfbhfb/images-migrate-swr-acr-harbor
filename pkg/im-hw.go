package pkg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	swr "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/region"
)

type HwSWR struct {
	c *Config
}

func NewHwSWR(cfg *Config) *HwSWR {
	return &HwSWR{c: cfg}
}

func removePrefixBeforeFirstSlash(input string) string {
	index := strings.Index(input, "/")
	if index != -1 {
		// 找到第一个斜杠
		return input[index+1:]
	}
	// 如果没有找到斜杠，返回原始字符串
	return input
}

func (h *HwSWR) GetAllFromImages() []string {

	auth := basic.NewCredentialsBuilder().
		WithAk(h.c.FromHwAccessKey).
		WithSk(h.c.FromHwSecretKey).
		Build()

	client := swr.NewSwrClient(
		swr.SwrClientBuilder().
			WithRegion(region.ValueOf(h.c.FromHwRegion)).
			WithCredential(auth).
			Build())

	var arrs []string
	request := &model.ListReposDetailsRequest{}
	strlimit := "99999" //
	request.Limit = &strlimit
	response, err := client.ListReposDetails(request)
	if err == nil {
		// var res HwRepoRespStruct
		// bs := []byte(response.GetHttpContentString())
		// // fmt.Println("33333")

		// if err = json.Unmarshal([]byte(response.String()), &res); err != nil {
		// 	fmt.Println(err)
		// 	return []string{}
		// }

		if response.Body != nil {
			for _, v := range *response.Body {
				arrs = append(arrs, removePrefixBeforeFirstSlash(v.Path))
			}
		}
		// fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
	return arrs
}

func (h *HwSWR) GetToNamespaces() []string {
	auth := basic.NewCredentialsBuilder().
		WithAk(h.c.AutoToHwAK).
		WithSk(h.c.AutoToHwSK).
		Build()

	client := swr.NewSwrClient(
		swr.SwrClientBuilder().
			WithRegion(region.ValueOf(h.c.ToHwRegion)).
			WithCredential(auth).
			Build())

	var arrs []string
	request := &model.ListNamespacesRequest{}
	response, err := client.ListNamespaces(request)
	if err == nil {
		if response.Namespaces != nil {
			for _, v := range *response.Namespaces {
				arrs = append(arrs, v.Name)
			}
		}
		// fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)

	}
	return arrs
}

func (h *HwSWR) CreateToNameSpace(name string) error {
	auth := basic.NewCredentialsBuilder().
		WithAk(h.c.AutoToHwAK).
		WithSk(h.c.AutoToHwSK).
		Build()

	client := swr.NewSwrClient(
		swr.SwrClientBuilder().
			WithRegion(region.ValueOf(h.c.ToHwRegion)).
			WithCredential(auth).
			Build())

	request := &model.CreateNamespaceRequest{}
	request.Body = &model.CreateNamespaceRequestBody{
		Namespace: name,
	}
	response, err := client.CreateNamespace(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
		return nil
	} else {
		fmt.Println(err)
		return err
	}
	return errors.New("创建华为组织错误")
}
