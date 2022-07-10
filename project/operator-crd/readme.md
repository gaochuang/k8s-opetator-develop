#Operator-CRD demo
#工程参考链接：https://github.com/kubernetes/sample-controller

###一、生成DeepCopyObject() Object方法
code-generator链接:https://github.com/kubernetes/code-generator.git \
需要使用deepcopy-gen生成, 为了方便，借助code-generator工程中提供的generate-groups.sh脚本生成 \
./generate-groups.sh \
Usage: generate-groups.sh <generators> <output-package> <apis-package> <groups-versions> ... \
\
&ensp;&ensp;<generators>        the generators comma separated to run (deepcopy,defaulter,client,lister,informer) or "all". \
&ensp;&ensp; <output-package>    the output package name (e.g. github.com/example/project/pkg/generated). \
&ensp;&ensp;<apis-package>      the external types dir (e.g. github.com/example/api or github.com/example/project/pkg/apis).\
&ensp;&ensp;<groups-versions>   the groups and their versions in the format "groupA:v1,v2 groupB:v1 groupC:v2", relative \
&ensp;&ensp;&ensp;&ensp;&ensp;&ensp;to <api-package>. \
&ensp;&ensp;...                 arbitrary flags passed to all generator binaries.\
\
\
Examples:\
&ensp;&ensp; generate-groups.sh all             github.com/example/project/pkg/client github.com/example/project/pkg/apis "foo:v1 bar:v1alpha1,v1beta1" \
&ensp;&ensp; generate-groups.sh deepcopy,client github.com/example/project/pkg/client github.com/example/project/pkg/apis "foo:v1 bar:v1alpha1,v1beta1" \
生成命令如下: \
&ensp;&ensp;/root/code/code-generator/generate-groups.sh  all operator-crd/pkg/genetated operator-crd/pkg/apis crd.example.com:v1 --go-header-file=/root/code/code-generator/hack/boilerplate.go.txt

