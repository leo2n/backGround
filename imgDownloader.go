package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// define image info struct
type ImgInfo struct {
	WorksName   string   `json:"works_name"`
	NickName    string   `json:"nick_name"`
	Link        string   `json:"link"`
	SaveName    string   `json:"save_name"`
	Md5         string   `json:"md5"`
	Size        float64  `json:"size"`
	Type        string   `json:"type"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Cat         string   `json:"cat"`
	Tags        []string `json:"tags"`
	ReadNum     int      `json:"read_num"`
	Fabulous    int      `json:"fabulous"`
	CreateTime  string   `json:"create_time"`
	Url         string   `json:"url"`
	DownloadUrl string   `json:"download_url"`
}

func getImgInfo() *ImgInfo {
	client := http.Client{Timeout: 5 * time.Second}
	baseUrl := "https://wallpaper.wispx.cn/api/find"
	request, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}
	q := request.URL.Query()
	q.Add("rand", "1")
	request.URL.RawQuery = q.Encode()
	//request.Header.Add("Host", "https://wallpaper.wispx.cn")
	request.Header.Add("User-Agent", "PostmanRuntime/7.24.0")
	request.Header.Add("Accept", "*/*")
	//request.Header.Add("Accept-Encoding", "gzip, deflate, br")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Referer", "https://wallpaper.wispx.cn/random")
	request.Header.Add("TE", "Trailers")
	request.Header.Add("X-Requested-With", "XMLHttpRequest")
	// send request
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	r := new(ImgInfo)
	body, err := ioutil.ReadAll(resp.Body) // binary stream
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, r)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(r)
	return r
}

// store img file in specfic directory
func downloadImg() {
	img := getImgInfo()
	imgURL := img.Url
	imgName := img.WorksName
	fmt.Println("imgName: ", imgName)
	imgExtendName := (strings.Split(img.SaveName, "."))[len(strings.Split(img.SaveName, "."))-1]
	if strings.Contains(imgExtendName, "&"){
		imgExtendName = (strings.Split(imgExtendName, "&"))[0]
	}
	fmt.Println("extendName: ", imgExtendName)
	fileName := uuid.New().String()
	fileNamePlusExtend := fileName+"."+imgExtendName
	fName := "\""+imgName+"\"" + "." + imgExtendName
	fName = strings.ReplaceAll(fName, " ", "-")
	fmt.Println("fName: ",fName)
	client := http.Client{Timeout: 30 * time.Second}
	request, err := http.NewRequest("GET", imgURL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Add("User-Agent", "PostmanRuntime/7.24.0")
	request.Header.Add("Accept", "*/*")
	//request.Header.Add("Accept-Encoding", "gzip, deflate, br")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Referer", "https://wallpaper.wispx.cn/random")
	request.Header.Add("TE", "Trailers")
	request.Header.Add("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(request)
	if err != nil {
		//log.Printf("%s", "have not finished!")
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// save img locally
	currentDir, err := os.Getwd()
	fmt.Println("currentDir is:", currentDir)
	if err != nil {
		log.Fatalln(err)
	}
	f, err := os.Create(currentDir + "/image/" + fileNamePlusExtend)
	if err != nil {
		log.Fatalln(err)
	}
	n, err := f.Write(body)
	log.Printf("Writed %d bytes", n)
	if err != nil {
		log.Fatalln("err", err)
	}

}
