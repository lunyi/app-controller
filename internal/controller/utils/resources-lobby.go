package utils

import (
	v1 "app-controller/api/v1"
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/go-logr/logr"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	setupLog = ctrl.Log.WithName("template")
)

func parseTemplate(templateName string, app *v1.Lobby, log logr.Logger) []byte {
	exePath, err := os.Executable()
	if err != nil {
		log.Info("获取执行程序路径失败:")
		log.Info(err.Error())
	}

	absPath, err := filepath.Abs(exePath)
	if err != nil {
		log.Info("转换为绝对路径失败:", err)
		log.Info(err.Error())
	}

	log.Info("===current folder===")
	log.Info(absPath)

	tmpl, err := template.ParseFiles("internal/controller/template/" + templateName + ".yml")
	if err != nil {
		panic(err)
	}
	b := new(bytes.Buffer)
	err = tmpl.Execute(b, app)
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func NewDeployment(app *v1.Lobby, log logr.Logger) *appv1.Deployment {
	d := &appv1.Deployment{}
	err := yaml.Unmarshal(parseTemplate("deployment", app, log), d)
	if err != nil {
		panic(err)
	}
	return d
}

func NewIngress(app *v1.Lobby, log logr.Logger) *netv1.Ingress {
	i := &netv1.Ingress{}
	err := yaml.Unmarshal(parseTemplate("ingress", app, log), i)
	if err != nil {
		panic(err)
	}
	return i
}

func NewService(app *v1.Lobby, log logr.Logger) *corev1.Service {
	s := &corev1.Service{}
	err := yaml.Unmarshal(parseTemplate("service", app, log), s)
	if err != nil {
		panic(err)
	}
	return s
}
