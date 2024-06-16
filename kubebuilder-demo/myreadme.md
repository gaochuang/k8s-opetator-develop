#kubuilder-demo project

##1.创建工程
1.1 创建kubebuilder-demo文件夹\
1.2 go mod init github.com/kubebuilder-demo 生成包依赖文件go.mod
## 2.初始化 kubebuilder operator功能,并定义groups为my.domin
2.1 kubebuilder init --domain my.domain\
2.2 创建api-- eg. kubebuilder create api --group ship --version v1beta1 --kind Frigate \
    

