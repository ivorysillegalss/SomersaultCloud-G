package ioutil

import (
	"SomersaultCloud/constant/common"
	"io/ioutil"
)

func LoadLuaScript(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return common.ZeroString, err
	}
	return string(data), nil
}
