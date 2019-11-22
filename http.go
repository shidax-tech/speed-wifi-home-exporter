package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

func HTTPGetXML(url string, buf interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(body, buf)
	if err != nil {
		return err
	}

	return nil
}
