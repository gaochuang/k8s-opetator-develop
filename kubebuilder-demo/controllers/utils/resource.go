package utils

import (
	"bytes"
	"github.com/kubebuilder-demo/api/v1beta1"
	"html/template"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	v12 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func parseTemplate(templateName string, app *v1beta1.App) []byte {
	tmpl, err := template.ParseFiles("controllers/template/" + templateName + ".yaml")
	if err != nil {
		panic(err)
	}
	bytebuffer := new(bytes.Buffer)
	//将app传入，会将template yaml文件中指定的占位符替换，然后写入到bytebuffer
	err = tmpl.Execute(bytebuffer, app)
	if err != nil {
		panic(err)
	}
	return bytebuffer.Bytes()
}

func NewDeployment(app *v1beta1.App) *v1.Deployment {
	d := &v1.Deployment{}
	err := yaml.Unmarshal(parseTemplate("deployment", app), d)
	if err != nil {
		panic(err)
	}
	return d
}

func NewIngress(app *v1beta1.App) *v12.Ingress {
	i := &v12.Ingress{}

	err := yaml.Unmarshal(parseTemplate("ingress", app), i)
	if err != nil {
		panic(err)
	}
	return i
}

func NewService(app *v1beta1.App) *v13.Service {
	s := &v13.Service{}
	err := yaml.Unmarshal(parseTemplate("service", app), s)
	if err != nil {
		panic(err)
	}
	return s
}
