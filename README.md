
## 项目目标：
- 应用镜像在 华为云，阿里云，Harbor 中便利的相互迁移
- 支持正则配置过滤和替换

## 项目使用的库介绍: 
- 使用华为云[huaweicloud-sdk-go](https://github.com/huaweicloud/huaweicloud-sdk-go-v3)
- 使用阿里云[alibaba-cloud-sdk-go](https://github.com/aliyun/alibaba-cloud-sdk-go) 获取镜像列表并生成auth.json和images.json文件
- 使用harbor [github.com/goharbor/go-client](https://github.com/goharbor/go-client)
- 使用阿里云开源镜像同步工具[image-syncer](https://github.com/AliyunContainerService/image-syncer) 进行镜像同步


#### 在k8s中部署
kubectl apply -f test1.yaml

#### config.yaml 配置说明
```yaml


# dry_run: true  # 不实际执行，只打印需要创建哪些 （组织，命名空间，项目） 和 迁移哪些镜像

# 配置从华为源【from】相关信息
from_hw_flag: true
from_hw_access_key: xxx
from_hw_secret_key: xxx
from_hw_region: cn-north-4
from_hw_docker_user: xxx
from_hw_docker_passwd: xxx

# 配置从华为目【to】相关信息
to_hw_flag: true
to_hw_region: cn-north-4
to_hw_docker_user: xxx
to_hw_docker_passwd: xxx


# 配置从阿里源【from】相关信息
# from_ali_flag: true
# from_ali_access_key: xxx
# from_ali_secret_key: xxx
# from_ali_region: cn-hangzhou
# from_ali_docker_user: xxx
# from_ali_docker_passwd: xxx

# 配置从阿里目【to】相关信息
# to_ali_flag: true
# to_ali_region: cn-shenzhen
# to_ali_docker_user: xxx
# to_ali_docker_passwd: xxx

# 配置从Harbor源【from】相关信息
# from_harbor_flag: true
# from_harbor_url: myharbor1.com:32443
# from_harbor_user: admin
# from_harbor_passwd: Harbor12345

# 配置从Harbor目【from】相关信息
# to_harbor_flag: true
# to_harbor_url: myharbor1.com:32443
# to_harbor_user: user1
# to_harbor_passwd: user1AAA

# 自动创建组织【华为】,命名空间【阿里】，项目(harbor)
auto_create_to_flag: true  
auto_create_hw_ak: xxx
auto_create_hw_sk: xxx
auto_create_ali_ak: xxx
auto_create_ali_sk: xxx

# 当源from，目to完全一样的域名时，需要使用中转registry镜像仓库
# transition_registry: registry:5000 
transition_registry: 192.168.255.246:32005  


# loop_gap: 30m15s #默认是在容器化中执行：默认从不再执行，保留日志，0s自动退出（二进制运行时）
# loop_gap: 15s 
# loop_gap: 0s 

#配置镜像过滤规则 drop keep replacement
opt_map:
# - action: drop
#   regex: (t1/.*) #过滤掉自己原本的
# - action: drop
#   regex: (t2/.*) #过滤掉自己原本的
# - action: drop
#   regex: (t3my328uay/.*) #过滤掉自己原本的
# - action: drop
#   regex: (t3mynewaaa/.*) #过滤掉自己原本的
# - action: drop
#   regex: (group3/.*)
# - action: keep
#   regex: (.*)
# - action: keep
#   regex: (.*myfilter.*)
# - action: replacement # 默认只替换组织
#   regex: (.*)
#   replacement: t3my328uay/$1  # 所有（组织，命名空间，项目）增加一层t1
# - action: replacement # 
#   regex: (.*/)(.*)
#   replacement: t3my328uay/$2  # 
# - action: replacement # 
#   regex: ([^/]*)/(.*)
#   replacement: t3mynewaaa/$1-$2  # 
# - action: replacement # 
#   regex: ([^/]*)/(.*)
#   replacement: t3mynewaaa/$1/$2  #


```

---

## [image-syncer](https://github.com/AliyunContainerService/image-syncer) Features

- Support for many-to-many registry synchronization
- Supports docker registry services based on Docker Registry V2 (e.g., Alibaba Cloud Container Registry Service, Docker Hub, Quay.io, Harbor, etc.)
- **Network & Memory Only, doesn't rely on any large disk storage, fast synchronization**
- Incremental Synchronization, ignore unchanged images automatically
- BloB-Level Concurrent Synchronization, adjustable goroutine numbers
- Automatic Retries of Failed Sync Tasks, to resolve the network problems (rate limit, etc.) while synchronizing
- Doesn't rely on Docker daemon or other programs



