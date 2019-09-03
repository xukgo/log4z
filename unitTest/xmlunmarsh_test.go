package unitTest

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log4z"
	"testing"
)

func TestXmlDeserializer(t *testing.T) {
	bs, err := ioutil.ReadFile("./conf/log4z.xml")
	if err != nil {
		t.Fail()
	}
	configRoot := new(log4z.ConfXmlRoot)
	err = xml.Unmarshal(bs, configRoot)
	if err != nil {
		t.Fail()
	}
	fmt.Println(configRoot)
}
