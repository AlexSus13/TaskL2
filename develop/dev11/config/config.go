package config

import (
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"

	"io/ioutil"
)

type Conf struct {
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
}

func Get() (*Conf, error) {

	var dconf Conf
	//Reading the yaml file
	yamlFile, err := ioutil.ReadFile("/home/ubuntu/TaskL2/wbschool_exam_L2/develop/dev11/etc/etc.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "Read .yaml File, func Get")
	}
	//Save the received data in the Conf structure
	err = yaml.Unmarshal(yamlFile, &dconf)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal .yaml File, func Get")
	}

	return &dconf, nil
}
