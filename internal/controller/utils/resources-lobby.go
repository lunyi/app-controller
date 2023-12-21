package utils

import (
	v1 "app-controller/api/v1"
	"bytes"
	"fmt"
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

func searchFile(rootDir, targetFile string, log logr.Logger) (string, error) {
	var foundFilePath string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // skip directories
		}

		log.Info(info.Name())
		if info.Name() == targetFile {
			foundFilePath = path
			return fmt.Errorf("file found") // stop walking
		}

		return nil
	})

	if err != nil && err.Error() != "file found" {
		return "", err
	}

	return foundFilePath, nil
}

func parseTemplate(templateName string, app *v1.Lobby, log logr.Logger) []byte {

	rootDir := "/"
	targetFile := "deployment.yml"

	foundFilePath, err := searchFile(rootDir, targetFile, log)
	if err != nil {
		log.Error(err, "Error:")
	}

	if foundFilePath != "" {
		log.Info("File found")
		log.Info(foundFilePath)
	} else {
		log.Info("File not found")
	}

	tmpl, err := template.ParseFiles("/workspace/internal/controller/template/" + templateName + ".yml")
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
