package utils

import (
	v1 "app-controller/api/v1"
	"bytes"
	"os"
	"text/template"

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	setupLog = ctrl.Log.WithName("template")
)

func parseTemplate(templateName string, app *v1.Lobby) []byte {
	wd, err := os.Getwd()

	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	setupLog.Info("current folder:", wd)

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

func NewDeployment(app *v1.Lobby) *appv1.Deployment {
	d := &appv1.Deployment{}
	err := yaml.Unmarshal(parseTemplate("deployment", app), d)
	if err != nil {
		panic(err)
	}
	return d
}

func NewIngress(app *v1.Lobby) *netv1.Ingress {
	i := &netv1.Ingress{}
	err := yaml.Unmarshal(parseTemplate("ingress", app), i)
	if err != nil {
		panic(err)
	}
	return i
}

func NewService(app *v1.Lobby) *corev1.Service {
	s := &corev1.Service{}
	err := yaml.Unmarshal(parseTemplate("service", app), s)
	if err != nil {
		panic(err)
	}
	return s
}
