package there

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	. "github.com/Gebes/there/there/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
)

type BodyReader struct {
	request http.Request
}

func (read BodyReader) Json(dest interface{}) error {
	return read.format(&dest, json.Unmarshal)
}

func (read BodyReader) Xml(dest interface{}) error {
	return read.format(&dest, xml.Unmarshal)
}

func (read BodyReader) Yaml(dest interface{}) error {
	fmt.Println(ContentTypeApplicationEdiDashX12)
	return read.format(&dest, yaml.Unmarshal)
}

func (read BodyReader) format(dest *interface{}, formatter func(data []byte, v interface{}) error) error {
	body, err := read.Bytes()
	if err != nil {
		return err
	}
	err = formatter(body, dest)
	return err
}

func (read BodyReader) String() (string, error) {
	data, err := read.Bytes()
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (read BodyReader) Bytes() ([]byte, error) {
	data, err := ioutil.ReadAll(read.request.Body)
	defer read.request.Body.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}
