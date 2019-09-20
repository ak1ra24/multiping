package utils

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Data struct {
	Addresses []Address `yaml:addresses`
}

type Address struct {
	HostName string `yaml:hostname`
	Address  string `yaml:address`
}

func ReadYaml() Data {
	buf, err := ioutil.ReadFile("pinglist.yaml")
	if err != nil {
		panic(err)
	}

	var d Data

	err = yaml.Unmarshal(buf, &d)
	if err != nil {
		panic(err)
	}

	return d

}
