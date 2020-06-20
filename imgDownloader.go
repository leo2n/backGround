package main

import (
	"encoding/json"
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

// 用于拿到待下载的image信息
func getImgInfo() (*ImgInfo, error) {
	client := http.Client{Timeout: 10 * time.Second}
	baseUrl := "https://wallpaper.wispx.cn/api/find"
	request, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	q := request.URL.Query()
	q.Add("rand", "1")
	request.URL.RawQuery = q.Encode()
	request.Header.Add("User-Agent", "PostmanRuntime/7.24.0")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Referer", "https://wallpaper.wispx.cn/random")
	request.Header.Add("TE", "Trailers")
	request.Header.Add("X-Requested-With", "XMLHttpRequest")
	// send request
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Error when client.Do")
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // binary stream
	if err != nil {
		log.Println("io.Reader to []byte error happend!")
		return nil, err
	}
	r := new(ImgInfo)
	err = json.Unmarshal(body, r)
	if err != nil {
		log.Println("Error happen when json.Unmarshal")
		return nil, err
	}
	return r, nil
}

// store img file in specfic directory
func downloadImg() (status bool, err error)  {
	// If 'image' folder exists in current directory
	currentDir, err := os.Getwd()
	if err!=nil {
		return false, err
	}
	if _, err := os.Stat(currentDir + "/image"); os.IsNotExist(err) {
		os.MkdirAll(currentDir+"/image", 0777)
	}
	img, err := getImgInfo()
	//defer recoverFrom()
	if err!=nil {
		log.Printf("getImgInfo error happend!")
		return false, nil
	}
	imgURL := img.Url
	imgName := img.WorksName
	log.Printf("The image's name U will downloading is: %s, URL is %s", imgName, imgURL)
	// This step target is get image extension
	imgExtendName := (strings.Split(img.SaveName, "."))[len(strings.Split(img.SaveName, "."))-1]
	if strings.Contains(imgExtendName, "&") {
		imgExtendName = (strings.Split(imgExtendName, "&"))[0]
	}
	log.Printf("This image's extendName is: %s", imgExtendName)
	//fileName := uuid.New().String()
	//fileNamePlusExtend := fileName+"."+imgExtendName
	fName := "\"" + imgName + "\"" + "." + imgExtendName
	fName = strings.ReplaceAll(fName, "/", "-")
	log.Printf("This image's Name is: %s", fName)
	client := http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest("GET", imgURL, nil)
	if err != nil {
		log.Println("when make http request, error happen!", err)
		return false, nil
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
		log.Printf("client do request:")
		return false, err
	}
	defer resp.Body.Close()

	// io.Reader convert to []byte
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("get resp's bytes:")
		return false, err
	}

	// save image locally
	currentDir, err = os.Getwd()
	if err != nil {
		log.Printf("get currentDir:")
		return false, err
	}

	f, err := os.Create(currentDir + "/image/" + fName)
	if err != nil {
		log.Printf("create file error")
		return false, err
	}

	n, err := f.Write(body)
	if err != nil {
		log.Printf("Write bytes error happend!")
		return false, err
	}
	log.Printf("%s Writed %.2f KB", fName, float64(n)/1024)
	return true, nil
}
