package pkg

import (
	"context"
	"fmt"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

type Harbor struct {
	c *Config
}

func NewHarbor(cfg *Config) *Harbor {

	return &Harbor{c: cfg}
}

func (h *Harbor) GetAllFromImages() []string {
	cs, _ := harbor.NewClientSet(&harbor.ClientSetConfig{
		URL:      "https://" + h.c.FromHarborUrl,
		Insecure: true,
		Username: h.c.FromHarborUser,
		Password: h.c.FromHarborPasswd})
	pi := repository.NewListAllRepositoriesParams()

	ret, err := cs.V2().Repository.ListAllRepositories(context.Background(), pi) // v2 client
	if err != nil {
		fmt.Println(err)
		// fmt.Println("出错了")
		return []string{}
	}
	var arr []string
	for _, value := range ret.Payload {
		// fmt.Printf("Index: %d, Value: %v\n", index, value.Name)
		arr = append(arr, value.Name)
	}

	return arr
}

func (h *Harbor) GetToNamespaces() []string {

	cs, _ := harbor.NewClientSet(&harbor.ClientSetConfig{
		URL:      "https://" + h.c.ToHarborUrl,
		Insecure: true,
		Username: h.c.ToHarborUser,
		Password: h.c.ToHarborPasswd})
	pi := project.NewListProjectsParams()
	// fmt.Println("11111111111nnnnn")

	ret, err := cs.V2().Project.ListProjects(context.Background(), pi) // v2 client
	if err != nil {
		fmt.Println(err)
		fmt.Println("harbor GetToNamespaces 出错了")
		return []string{}
	}
	var arr []string
	// fmt.Println(len(ret.Payload))
	for _, value := range ret.Payload {
		// fmt.Printf("Index: %d, Value: %v\n", index, value.Name)
		arr = append(arr, value.Name)
	}
	// fmt.Println(arr)

	return arr
}

func (h *Harbor) CreateToNameSpace(name string) error {
	cs, _ := harbor.NewClientSet(&harbor.ClientSetConfig{
		URL:      "https://" + h.c.ToHarborUrl,
		Insecure: true,
		Username: h.c.ToHarborUser,
		Password: h.c.ToHarborPasswd})
	pi := project.NewCreateProjectParams()
	pi.Project = &models.ProjectReq{
		ProjectName: name,
	}

	// fmt.Println("11111")
	_, err := cs.V2().Project.CreateProject(context.Background(), pi) // v2 client
	if err != nil {
		fmt.Println(err)
		// fmt.Println("出错了")
	}
	// fmt.Println("2222")
	return nil
}
