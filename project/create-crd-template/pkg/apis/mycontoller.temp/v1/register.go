//+groupName=mycontoller.temp
package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

//将group和version注册到schema
var (
	Schema       = runtime.NewScheme()
	GroupVersion = schema.GroupVersion{Group: "mycontoller.temp", Version: "v1"}
	//解码工具
	Codes = serializer.NewCodecFactory(Schema)
)
