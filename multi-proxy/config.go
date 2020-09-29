package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ProxySource ProxySource `yaml:"proxy_source"`
}

//type ProxySource struct {
//	Nz string `yaml:"nz"`
//	Sa string `yaml:"sa"`
//	Mz string `yaml:"mz"`
//	Tw string `yaml:"tw"`
//	Ve string `yaml:"ve"`
//	Kg string `yaml:"kg"`
//	Pe string `yaml:"pe"`
//	Uz string `yaml:"uz"`
//	Mg string `yaml:"mg"`
//	Np string `yaml:"np"`
//	Pk string `yaml:"pk"`
//	Ph string `yaml:"ph"`
//	Vn string `yaml:"vn"`
//	Mm string `yaml:"mm"`
//	Th string `yaml:"th"`
//	My string `yaml:"my"`
//	Id string `yaml:"id"`
//	Tl string `yaml:"tl"`
//	Kw string `yaml:"kw"`
//	Cn string `yaml:"cn"`
//	Hn string `yaml:"hn"`
//}
type ProxySource struct {
	Nz string `yaml:"nz"`
	Sa string `yaml:"sa"`
	Mz string `yaml:"mz"`
	Tw string `yaml:"tw"`
	Ve string `yaml:"ve"`
	Kg string `yaml:"kg"`
	Pe string `yaml:"pe"`
	Uz string `yaml:"uz"`
	Mg string `yaml:"mg"`
	Np string `yaml:"np"`
	Pk string `yaml:"pk"`
	Ph string `yaml:"ph"`
	Vn string `yaml:"vn"`
	Mm string `yaml:"mm"`
	Th string `yaml:"th"`
	My string `yaml:"my"`
	Id string `yaml:"id"`
	Tl string `yaml:"tl"`
	Kw string `yaml:"kw"`
	Cn string `yaml:"cn"`
	Hn string `yaml:"hn"`
}

var (
	defaultPath   = "."
	DefaultConfig Config

	ConfigMap = make(map[string]map[string]string)
)

func init() {
	err := loadConfigFile("app.yaml", &ConfigMap)
	if err != nil {
		panic(err)
	}
}

func loadConfigFile(name string, cfg interface{}) error {

	confPath := fmt.Sprintf("%s/%s", defaultPath, name)
	if !fileExists(confPath) {
		return fmt.Errorf("config file `%s' is not exists", confPath)
	}

	file, err := os.Open(confPath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(buf, cfg)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
