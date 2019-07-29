package main

import (
	"fmt"
	// "io/ioutil"
	"net/http"
	//"net/http/httputil"
)

func main() {
	/*
	resp, err := http.Get("https://www.reddit.com/user/genericSpecimen")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s", body)

	for name, headers := range resp.Header {
		fmt.Printf("%v: %v\n", name, headers)
	}
	*/

	client := &http.Client {}
	resp, err := client.Get("https://www.rlacollege.edu.in/view-all-details.php")
	if err != nil {
		panic(err)
	}
	for name, headers := range resp.Header {
		fmt.Printf("%v: %v\n", name, headers)
	}

	fmt.Println(resp.Header.Get("last-modified"))

	fmt.Println("---------------------------------")

	req, err := http.NewRequest("GET", "https://www.rlacollege.edu.in/view-all-details.php", nil)
	req.Header.Add("If-Modified-Since", "Sun, 28 Jul 2019 16:25:02 GMT")

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	for name, headers := range resp.Header {
		fmt.Printf("%v: %v\n", name, headers)
	}

}