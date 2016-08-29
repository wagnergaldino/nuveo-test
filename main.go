package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type outputdata []map[string]string

func main() {

	var data outputdata
	var err error
	urlStr1 := "http://www.galdino.net.br/nuveotest.csv" //csv
	urlStr2 := "http://www.galdino.net.br/nuveotest.php" //json

	err = getData(urlStr1, &data)
	if err != nil {
		fmt.Println("Erro ao processar CSV:", err)
	}

	fmt.Println("")
	fmt.Println("CSV DATA from", urlStr1)
	fmt.Println("")
	fmt.Println(data)
	fmt.Println("")
	fmt.Println("-------------------------------------------------------------------------------")
	fmt.Println("")

	data = nil
	err = getData(urlStr2, &data)
	if err != nil {
		fmt.Println("Erro ao processar CSV:", err)
	}
	fmt.Println("JSON DATA from", urlStr2)
	fmt.Println("")
	fmt.Println(data)

}

func getData(url string, data *outputdata) (err error) {

	resp, _ := http.Get(url)
	defer resp.Body.Close()
	ctype := resp.Header.Get("Content-Type")

	if strings.Contains(ctype, "text/csv") {
		err = getCsvData(resp, data)
		if err != nil {
			return
		}
	} else if strings.Contains(ctype, "application/json") {
		err = getJsonData(resp, data)
		if err != nil {
			return
		}
		putMandatoryData(data)
	} else {
		return errors.New("\nContent-Type incompatível (" + ctype + ").\nUse apenas 'text/csv' ou 'application/json'")
	}
	return
}

func getCsvData(resp *http.Response, data *outputdata) error {

	var mapKeys []string
	var mapValues []string
	reader := csv.NewReader(resp.Body)
	lineCount := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		basemap := getBaseMap()

		if lineCount == 1 {
			keys := record[0]
			mapKeys = strings.Split(keys, ";")
		} else {
			values := record[0]
			mapValues = strings.Split(values, ";")
			for i := 0; i < len(mapKeys); i++ {
				basemap[mapKeys[i]] = mapValues[i]
			}
			*data = append(*data, basemap)
		}
		lineCount += 1
	}
	return nil
}

func getJsonData(resp *http.Response, data *outputdata) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(body, data)
	return nil
}

func getBaseMap() map[string]string {
	m := make(map[string]string)
	m["nome"] = ""
	m["email"] = ""
	m["sexo"] = ""
	m["idade"] = ""
	return m
}

func putMandatoryData(data *outputdata) {
	for _, m := range *data {
		if _, ok := m["nome"]; !ok {
			m["nome"] = ""
		}
		if _, ok := m["email"]; !ok {
			m["email"] = ""
		}
		if _, ok := m["sexo"]; !ok {
			m["sexo"] = ""
		}
		if _, ok := m["idade"]; !ok {
			m["idade"] = ""
		}
	}
	return
}
