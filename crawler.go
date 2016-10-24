package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

var dir string
var wg sync.WaitGroup

func main() {
	dir = "food/"
	os.Mkdir(dir, 0777)

	for i := 1; i < 50; i++ {

		data, err := goquery.NewDocument("http://www.meishij.net/list.php?lm=43&page=" + strconv.Itoa(i))
		if err != nil {
			fmt.Println("url is error!")
		}

		data.Find("div.listtyle1_list div.listtyle1 a").Each(func(i int, s *goquery.Selection) {

			href, hrefexist := s.Attr("href")
			img, imgexist := s.Find("img").Attr("src")
			alt, altexist := s.Find("img").Attr("alt")

			if !hrefexist {
				fmt.Println("href not exist!")
				return
			}
			if !imgexist {
				fmt.Println("img src not exist!")
				return
			}
			if !altexist {
				fmt.Println("alt not exist!")
				return
			}
			wg.Add(1)
			go result(href, img, alt)

		})
	}
	wg.Wait()
}

func result(href, img, alt string) {
	data := getUrl(img)
	os.Mkdir(dir+alt, 0777)
	dirname := dir + alt + "/"

	txtfile, txterr := os.Create(dirname + alt + ".txt")
	defer txtfile.Close()
	imgfile, imgerr := os.Create(dirname + path.Base(img))
	defer imgfile.Close()

	if txterr != nil {
		fmt.Println("txt create is error!")
	}
	txtfile.WriteString(href)

	if imgerr != nil {
		fmt.Println("img file create is error!")
	}
	imgfile.Write(data)
	wg.Done()
}

func getUrl(url string) []byte {
	resp, resperr := http.Get(url)
	if resperr != nil {
		fmt.Println("imgurl error!")
	}
	data, _ := ioutil.ReadAll(resp.Body)

	return data
}
