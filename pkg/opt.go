package pkg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AliyunContainerService/image-syncer/pkg/client"
	"github.com/AliyunContainerService/image-syncer/pkg/utils"
)

var (
	authStep1Json = "step1auth.json"
	authStep2Json = "step2auth.json"

	imagesStep1Json = "imagesStep1.json"
	imagesStep2Json = "imagesStep2.json"
)

type optItf interface {
	GetAllFromImages() []string // 获取from所有的镜像

	CreateToNameSpace(string) error // 准备迁移to的后端的镜像，组织，命名空间，项目
	GetToNamespaces() []string      //

}

func CleanFileEnv() error {

	err := os.Remove(authStep1Json)
	if err != nil {
		fmt.Printf("删除文件失败: %v\n", err)
		// return errors.New("删除文件失败")
	}
	err = os.Remove(authStep2Json)
	if err != nil {
		fmt.Printf("删除文件失败: %v\n", err)
		// return errors.New("删除文件失败")
	}
	err = os.Remove(imagesStep1Json)
	if err != nil {
		fmt.Printf("删除文件失败: %v\n", err)
		// return errors.New("删除文件失败")
	}
	err = os.Remove(imagesStep2Json)
	if err != nil {
		fmt.Printf("删除文件失败: %v\n", err)
		// return errors.New("删除文件失败")
	}
	return nil
}

/*
func dropMetricIfMatch(name string) string {
	// 定义匹配的正则表达式
	re := regexp.MustCompile(`^my_metric_prefix_.+$`)

	// 检查是否匹配
	if re.MatchString(name) {
		// 如果匹配，则返回空字符串表示删除该指标
		return ""
	}

	// 如果不匹配，则返回原始指标名称
	return name
}
*/

func check_tcp(url string) error {

	var address string
	// fmt.Println(url)
	if strings.Contains(url, ":") {
		address = url
	} else {
		address = url + ":" + "443"
	}

	// Set a timeout for the connection attempt
	timeout := 5 * time.Second

	// Attempt to establish a TCP connection
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		fmt.Printf("TCP connection to %s failed: %v\n", address, err)
		return err
	}
	defer conn.Close()

	return nil
}

func genAuthFileNew(c *Config) error {

	flagCheckTransitionRegistry := false // 同区域则分成两步，先推到中间镜像仓库
	// docker用的用户名密码
	// repoMap := make(map[string]interface{})
	// 判断是否需要使用registry中转一下
	flagCheckTransitionRegistry = checkMoreStep(c)

	if flagCheckTransitionRegistry {
		err := check_tcp(c.TransitionRegistry)
		if err != nil {
			fmt.Println("error check_tcp: ", err)
			return err
		}
	}

	fromUrl := ""
	fromUser := ""
	fromPw := ""

	toUrl := ""
	toUser := ""
	toPw := ""

	if c.FromHwFlag {
		fromUrl = fmt.Sprintf("swr.%s.myhuaweicloud.com", c.FromHwRegion)
		fromUser = c.FromHwDockerUser
		fromPw = c.FromHwDockerPasswd
	}
	if c.FromAliFlag {
		fromUrl = fmt.Sprintf("registry.%s.aliyuncs.com", c.FromAliRegion)
		fromUser = c.FromAliDockerUser
		fromPw = c.FromAliDockerPasswd
	}
	if c.FromHarborFlag {
		fromUrl = c.FromHarborUrl
		fromUser = c.FromHarborUser
		fromPw = c.FromHarborPasswd
	}

	if c.ToHwFlag {
		toUrl = fmt.Sprintf("swr.%s.myhuaweicloud.com", c.ToHwRegion)
		toUser = c.ToHwDockerUser
		toPw = c.ToHwDockerPasswd
	}
	if c.ToAliFlag {
		toUrl = fmt.Sprintf("registry.%s.aliyuncs.com", c.ToAliRegion)
		toUser = c.ToAliDockerUser
		toPw = c.ToAliDockerPasswd
	}
	if c.ToHarborFlag {
		toUrl = c.ToHarborUrl
		toUser = c.ToHarborUser
		toPw = c.ToHarborPasswd
	}

	if flagCheckTransitionRegistry {

		repoMapStep1 := make(map[string]interface{})
		repoMapStep1[fromUrl] = Auth{
			Username: fromUser,
			Password: fromPw,
			Insecure: true,
		}
		repoMapStep1[c.TransitionRegistry] = Auth{
			Username: "",
			Password: "",
			Insecure: true,
		}

		if err := WriteFile(authStep1Json, repoMapStep1); err != nil {
			return err
		}

		repoMapStep2 := make(map[string]interface{})
		repoMapStep2[toUrl] = Auth{
			Username: toUser,
			Password: toPw,
			Insecure: true,
		}
		repoMapStep2[c.TransitionRegistry] = Auth{
			Username: "",
			Password: "",
			Insecure: true,
		}

		if err := WriteFile(authStep2Json, repoMapStep2); err != nil {
			return err
		}

	} else {
		repoMapStep1 := make(map[string]interface{})

		repoMapStep1[fromUrl] = Auth{
			Username: fromUser,
			Password: fromPw,
			Insecure: true,
		}
		repoMapStep1[toUrl] = Auth{
			Username: toUser,
			Password: toPw,
			Insecure: true,
		}

		if err := WriteFile(authStep1Json, repoMapStep1); err != nil {
			return err
		}

	}

	return nil
}

func regReplace(org, regexPattern, Replacement string) (bool, string) {
	// str := "Golang is awesome, 123 and 456 too."

	// // 定义一个正则表达式，包含两个捕获组，分别匹配数字部分
	// regexPattern := `(\d+) and (\d+)`

	// 编译正则表达式
	regex := regexp.MustCompile(regexPattern)

	// fmt.Println(org, "   ", regexPattern)
	// 检查是否匹配
	if !regex.MatchString(org) {
		return false, ""
	}

	// 使用 ReplaceAllStringFunc 函数进行替换
	Replacement = regex.ReplaceAllStringFunc(org, func(match string) string {
		// 获取捕获组的值
		matches := regex.FindStringSubmatch(match)
		// fmt.Println(len(matches))
		for i := 0; i < len(matches); i++ {
			// fmt.Println(matches[i])
			Replacement = strings.ReplaceAll(Replacement, "$"+strconv.FormatInt(int64(i), 10), matches[i])
		}

		// if len(matches) >= 3 {
		// 	// 返回替换后的字符串，使用 $2 $1 的顺序
		// 	return matches[2] + " " + matches[1]
		// }
		// 没有匹配到或者捕获组数量不够，返回原始字符串
		return Replacement
	})

	return true, Replacement
}

func checkMoreStep(c *Config) bool {

	if c.FromHwFlag && c.ToHwFlag && c.FromHwRegion == c.ToHwRegion {
		// 华为同区域
		return true
	}

	if c.FromAliFlag && c.ToAliFlag && c.FromAliRegion == c.ToAliRegion {
		// 阿里同区域
		return true
	}
	// fmt.Printf("%v", c)

	if c.FromHarborFlag && c.ToHarborFlag && c.FromHarborUrl == c.ToHarborUrl {
		// harbor同区域
		return true
	}
	return false
}

func genImagesFileNew(c *Config) error {

	// 先准备好目标的命名空间
	var optHw optItf
	var optAli optItf
	var optHarbor optItf

	var optFrom optItf

	optHw = NewHwSWR(c)
	optAli = NewAliACR(c)
	optHarbor = NewHarbor(c)

	// var opt interface{}
	flagCheckTransitionRegistry := false // 同区域则分成两步，先推到中间镜像仓库
	// docker用的用户名密码
	// repoMap := make(map[string]interface{})
	// 判断是否需要使用registry中转一下

	flagCheckTransitionRegistry = checkMoreStep(c)

	if c.FromHwFlag {
		optFrom = optHw
	}

	if c.FromAliFlag {
		optFrom = optAli
	}

	if c.FromHarborFlag {
		optFrom = optHarbor
	}

	var allImages []string

	allImages = optFrom.GetAllFromImages()
	// fmt.Println(allImages)

	var fromImages []string
	var toImages []string
	for _, v1 := range allImages {
		flagDrop := false
		tos := v1
		for _, v2 := range c.OptMap {
			if strings.ToUpper(v2.Action) == "DROP" {
				if f, _ := regReplace(v1, v2.Regex, ""); f == true {
					flagDrop = true
				}
			}

			if strings.ToUpper(v2.Action) == "KEEP" {
				if f, _ := regReplace(v1, v2.Regex, ""); f == false {
					flagDrop = true
				}
			}

			if strings.ToUpper(v2.Action) == "REPLACEMENT" {
				if f, s := regReplace(v1, v2.Regex, v2.Replacement); f == true {
					// flagDrop = true
					tos = s
				}
			}
		}
		// fmt.Println("11231231  :", flagDrop, " ", v1)
		if flagDrop {

			fmt.Println("镜像被过滤 ", v1)
		} else {
			fromImages = append(fromImages, v1)
			toImages = append(toImages, tos)
		}
	}
	fmt.Println(fromImages)
	fmt.Println(toImages)

	if flagCheckTransitionRegistry {

		imagesMap1 := make(map[string]interface{})
		var tmpA []string
		for i, v := range fromImages {
			s := c.TransitionRegistry + "/" + toImages[i]
			imagesMap1[getFromUrl(c)+v] = s
			tmpA = append(tmpA, s)
		}
		if err := WriteFile(imagesStep1Json, imagesMap1); err != nil {
			return err
		}

		imagesMap2 := make(map[string]interface{})
		for i, v := range tmpA {
			imagesMap2[v] = getToUrl(c) + "/" + toImages[i]
		}
		if err := WriteFile(imagesStep2Json, imagesMap2); err != nil {
			return err
		}

		return nil
	} else {

		imagesMap1 := make(map[string]interface{})
		for i, v := range fromImages {
			imagesMap1[getFromUrl(c)+v] = getToUrl(c) + "/" + toImages[i]
		}
		if err := WriteFile(imagesStep1Json, imagesMap1); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func getFromUrl(c *Config) string {
	if c.FromHwFlag {
		return fmt.Sprintf("swr.%s.myhuaweicloud.com/", c.FromHwRegion)
	}

	if c.FromAliFlag {
		return fmt.Sprintf("registry.%s.aliyuncs.com/", c.FromAliRegion)
	}

	if c.FromHarborFlag {
		return c.FromHarborUrl + "/"
	}

	return ""
}

func getToUrl(c *Config) string {
	if c.ToHwFlag {
		return fmt.Sprintf("swr.%s.myhuaweicloud.com", c.ToHwRegion)
	}

	if c.ToAliFlag {
		return fmt.Sprintf("registry.%s.aliyuncs.com", c.ToAliRegion)
	}

	if c.ToHarborFlag {
		return c.ToHarborUrl
	}
	return ""
}

func runTranslate(authJson, imagesJson string) error {
	client, err := client.NewSyncClient("", authJson, imagesJson, "", 5, 2,
		utils.RemoveEmptyItems([]string{}), utils.RemoveEmptyItems([]string{}), false)
	if err != nil {
		fmt.Println("init sync client error: ", err)
		return err
	}
	client.Run()
	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func printFile(filePath string) {
	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("无法读取文件:", err)
		return
	}

	// 打印文件内容
	fmt.Println("文件名: ", filePath, "\n", "文件内容:")
	fmt.Println(string(content))
}

func prepareNamespace(c *Config) error {
	// 先准备好目标的命名空间
	var optHw optItf
	var optAli optItf
	var optHarbor optItf

	var optTo optItf

	optHw = NewHwSWR(c)
	optAli = NewAliACR(c)
	optHarbor = NewHarbor(c)

	if c.ToHwFlag {
		optTo = optHw
	}

	if c.ToAliFlag {
		optTo = optAli
	}

	if c.ToHarborFlag {
		optTo = optHarbor
	}

	var allToNs []string

	fileNameJson := ""
	if checkMoreStep(c) {
		fileNameJson = imagesStep2Json
	} else {
		fileNameJson = imagesStep1Json
	}
	if fileExists(fileNameJson) {
		if err, mapimages := LoadFile(fileNameJson); err != nil {
			// fmt.Println("11111111111")
			return err
		} else {
			// fmt.Println("11111111111222")
			for _, v := range mapimages {
				if ns := getNsFromImageString(v); ns != "" {
					alreadyflag := false
					for _, v2 := range allToNs {
						if v2 == ns {
							alreadyflag = true
						}

					}
					if !alreadyflag {
						allToNs = append(allToNs, ns)
					}
				}
				// fmt.Println(i)
				// fmt.Println(v)
			}
		}

	}
	// fmt.Println("000000")
	// fmt.Println(allToNs)
	// fmt.Println("000000aaa")

	// 判断哪些是要新创建 组织，命名空间，项目

	var needCreateNs []string
	if c.AutoCreateFlag == true {

		// 获取to所有命名空间
		arrs := optTo.GetToNamespaces()
		// fmt.Println("jdflsjfl1111111")
		fmt.Println(arrs)

		for _, v1 := range allToNs {
			alreadyNsflag := false
			for _, v2 := range arrs {
				if v2 == v1 {
					alreadyNsflag = true
				}

			}
			if !alreadyNsflag {
				needCreateNs = append(needCreateNs, v1)
			}

		}

	}

	if c.DryRunFlag {
		if len(needCreateNs) > 0 {
			fmt.Println("需要到目标创建的 组织，命名空间，项目 ")
			for _, v2 := range needCreateNs {
				fmt.Println(v2)
			}

		} else {
			fmt.Println("不需要在对端创建命名空间")
		}
		return nil
	}

	if c.AutoCreateFlag != true && len(needCreateNs) > 0 {
		fmt.Println("需要创建目标的 组织，命名空间，项目： 请配置", "auto_create_to_flag 为true")
		return errors.New("需要创建目标的 组织，命名空间，项目： 请配置" + "auto_create_to_flag 为true")
	}
	for _, v2 := range needCreateNs {
		fmt.Println("需要创建命名空间:", v2)
	}
	for _, v2 := range needCreateNs {
		// fmt.Println("需要创建命名空间:", v2)
		if err := optTo.CreateToNameSpace(v2); err != nil {
			fmt.Println("创建命名空间失败: " + err.Error())
			return errors.New("创建命名空间失败: " + err.Error())
		}
	}
	// fmt.Println("000000MM")
	// fmt.Println(needCreateNs)
	// fmt.Println("000000aaann")
	return nil
}

func getNsFromImageString(str string) string {
	// 您的字符串
	// str := "192.168.255.246:32005/t1/library/ubuntu/nginx"

	// 使用 strings.Split 函数拆分字符串
	parts := strings.Split(str, "/")

	var t1 string
	// 如果有足够的部分，则提取 t1
	if len(parts) >= 2 {
		t1 = parts[1]
		// fmt.Println("提取的 t1:", t1)
	} else {
		fmt.Println("无法提取 t1，字符串格式不正确")
		return ""
	}
	return t1

}

func beginImageSync(c *Config) {

	// 先准备好目标的命名空间
	if err := prepareNamespace(c); err != nil {
		fmt.Errorf("%v %v", "ns没有准备好，出错: ", err)
		return
	}
	// return

	if c.DryRunFlag {

		if fileExists(authStep1Json) && fileExists(imagesStep1Json) {
			printFile(authStep1Json)
			printFile(imagesStep1Json)
		}
		if fileExists(authStep2Json) && fileExists(imagesStep2Json) {
			printFile(authStep2Json)
			printFile(imagesStep2Json)
		}

		return
	}

	if fileExists(authStep1Json) && fileExists(imagesStep1Json) {
		runTranslate(authStep1Json, imagesStep1Json)
	}

	if fileExists(authStep2Json) && fileExists(imagesStep2Json) {
		runTranslate(authStep2Json, imagesStep2Json)
	}

}

func StartOptNew(c *Config) {

	CleanFileEnv()
	// return

	// 生成auth文件
	if genAuthFileNew(c) != nil {
		fmt.Println("auth 生成文件出错！！")
		return
	}

	// 生成images配置文件
	genImagesFileNew(c)

	// 执行 ImageSync
	beginImageSync(c)

	//
	/*
		var opt opt
		opt = NewHarbor(c)
			opt.GetAllFromImages()
	*/
	// from.DryRunFrom(c)
}
