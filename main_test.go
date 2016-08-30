package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

var m = map[string]string{"nome":"","email":"","sexo":"","idade":"",}

var csvData = []map[string]string{
	map[string]string{"nome":"wagal","email":"wagal@hotmail.com","sexo":"m","idade":"18","telefone":"","endereco":"",},
	map[string]string{"nome":"wagner","email":"wagner@galdino.net.br","sexo":"m","idade":"50","telefone":"993986604","endereco":"av celso garcia 5398 ap 4",},
	map[string]string{"nome":"galdino","email":"","sexo":"m","idade":"","telefone":"993986604","endereco":"av celso garcia 5398 ap 4",},
}

var jsonData = []map[string]string{
	map[string]string{"nome":"wagal","email":"wagal@hotmail.com","sexo":"m","idade":"18",},
	map[string]string{"nome":"wagner","email":"wagner@galdino.net.br","sexo":"m","idade":"50","telefone":"993986604","endereco":"av celso garcia 5398 ap 4",},
	map[string]string{"nome":"galdino","email":"","sexo":"m","idade":"","telefone":"993986604","endereco":"av celso garcia 5398 ap 4",},
}

func TestIsCsv(t *testing.T){
	resp, _ := http.Get("http://www.galdino.net.br/nuveotest.csv")
	defer resp.Body.Close()
	expected := true
	actual := strings.Contains(resp.Header.Get("Content-Type"), "text/csv")
	if actual != expected {
		t.Error("TestIsCsv failed")
	}
}

func TestIsJson(t *testing.T){
	resp, _ := http.Get("http://www.galdino.net.br/nuveotest.php")
	defer resp.Body.Close()
	expected := true
	actual := strings.Contains(resp.Header.Get("Content-Type"), "application/json")
	if actual != expected {
		t.Error("TestIsJson failed")
	}
}

func TestContentTypeNotSupported(t *testing.T){
	resp, _ := http.Get("http://www.galdino.net.br/wagner/index.html")
	defer resp.Body.Close()
	expected := true
	actual := (!strings.Contains(resp.Header.Get("Content-Type"), "application/json") && !strings.Contains(resp.Header.Get("Content-Type"), "text/csv"))
	if actual != expected {
		t.Error("TestContentTypeNotSupported failed")
	}
}

func TestGetBaseMap(t *testing.T){
	expected := true
	actual := checkBaseMap()
	if actual != expected {
		t.Error("TestGetBaseMap failed")
	}
}

func checkBaseMap() bool {
	basemap := getBaseMap()
	if len(m) != len(basemap) {
		return false
	}
	for k, _ := range basemap {
		if _, ok := m[k]; !ok {
			return false
		} else if basemap[k] != m[k] {
			return false
		}
	}
	return true
}

func TestGetCsvData(t *testing.T){
	var data outputdata
	var err error
	resp, _ := http.Get("http://www.galdino.net.br/nuveotest.csv")
	defer resp.Body.Close()
	err = getCsvData(resp, &data)
	if err != nil {
		fmt.Println("Erro ao processar CSV:", err)
	}
	expected := true
	actual := checkCsvData(data)
	if actual != expected {
		t.Error("TestGetCsvData failed")
	}
}

func checkCsvData(data outputdata) bool {
	if len(data) != len(csvData) {
		return false
	}
	for i := 0; i < len(data); i++ {
		retMap := data[i]
		csvMap := csvData[i]
		if len(retMap) != len(csvMap) {
			return false
		}
		for k, _ := range csvMap {
			if _, ok := retMap[k]; !ok {
				return false
			} else if csvMap[k] != retMap[k] {
				return false
			}
		}
	}
	return true
}

func TestGetJsonData(t *testing.T){
	var data outputdata
	var err error
	resp, _ := http.Get("http://www.galdino.net.br/nuveotest.php")
	defer resp.Body.Close()
	err = getJsonData(resp, &data)
	if err != nil {
		fmt.Println("Erro ao processar JSON:", err)
	}
	putMandatoryData(&data)
	expected := true
	actual := checkJsonData(data)
	if actual != expected {
		t.Error("TestGetJsonData failed")
	}
}

func checkJsonData(data outputdata) bool {
	if len(data) != len(jsonData) {
		return false
	}
	for i := 0; i < len(data); i++ {
		retMap := data[i]
		jsonMap := jsonData[i]
		if len(retMap) != len(jsonMap) {
			return false
		}
		for k, _ := range jsonMap {
			if _, ok := retMap[k]; !ok {
				return false
			} else if jsonMap[k] != retMap[k] {
				return false
			}
		}
	}
	return true
}
