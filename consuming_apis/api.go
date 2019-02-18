package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	Status string                 `json:"status"`
	Dogs   map[string]interface{} `json:"message"`
}

type FoodResponse struct {
	Receipe []Receipe `json:"results"`
}

type Receipe struct {
	Title       string `json:"title"`
	Ingredients string `json:"ingredients"`
}

func main() {
	fmt.Println("---------")
	fmt.Println("read structured data")
	getStructuredData()
	fmt.Println("---------")
	fmt.Println("read unstructured data")
	fmt.Println("---------")
	c := make(chan string)
	go getUnstructuredData(c)
	fmt.Println(<-c)
}

func getStructuredData() {
	response, err := http.Get("http://www.recipepuppy.com/api/?i=onions,garlic&q=omelet&p=3")
	if err != nil {
		fmt.Println(err.Error())
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseObject FoodResponse
	json.Unmarshal(responseData, &responseObject)
	for _, val := range responseObject.Receipe {
		fmt.Println(val.Title)
	}

}
func getUnstructuredData(c chan string) {
	response, err := http.Get("https://dog.ceo/api/breeds/list/all")

	if err != nil {
		fmt.Println(err.Error())
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	dogs := responseObject.Dogs
	var dogList []string
	for key := range dogs {
		dogList = append(dogList, key)
	}
	c <- strings.Join(dogList, "\n")
}
