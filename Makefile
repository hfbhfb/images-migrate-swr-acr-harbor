
IMGTAG:=swr.cn-north-4.myhuaweicloud.com/hfbbg4/images-migrate-linux-amd64:v0.3

build:
	@echo "编译和运行镜像迁移工具:支持华为云，阿里云，Harbor"
	GOOS=linux go build -o images-migrate-linux-amd64 ./cmd/main.go

dockerimg: build
	docker build -f ./Dockerfile -t ${IMGTAG} .

k8stest1:
	@echo "测试"
	-kubectl delete -f test1.yaml
	kubectl apply -f test1.yaml
