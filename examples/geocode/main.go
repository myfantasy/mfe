package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/myfantasy/mfe"
)

func main() {
	file, er := os.Open("data.txt")
	if er != nil {
		log.Fatal(er)
	}
	defer file.Close()

	fileOut, erro := os.Create("res.txt")
	if erro != nil {
		log.Fatal(erro)
	}
	defer fileOut.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		sa := strings.Split(s, "\t")
		fmt.Println(sa[1])

		resp, err := http.Get("https://geocode-maps.yandex.ru/1.x/?format=json&geocode=" + sa[1])
		if err != nil {
			log.Fatalln(err)
		}

		if resp.StatusCode != 200 {
			log.Fatalln("Not Good Query")
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		v, err := mfe.VariantNewFromJSON(string(body))
		if err != nil {
			log.Fatalln(err)
		}

		its := v.GE("response", "GeoObjectCollection", "featureMember")
		if its.Count() > 0 {
			vp := its.GI(0).GE("GeoObject", "Point", "pos")

			if vp.IsNull() {
				fileOut.WriteString("" + sa[0] + "\t" + sa[1] + "\tNot Found2")
			} else {

				fileOut.WriteString("" + sa[0] + "\t" + sa[1] + "\t" + vp.Str() + "")
			}
		} else {

			fileOut.WriteString("" + sa[0] + "\t" + sa[1] + "\tNot Found")
		}
		fileOut.WriteString("\n")
	}

}
