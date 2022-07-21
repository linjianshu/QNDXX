package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

var m = make(map[string]string)
var Folder string

func main() {
	//做个字典 后续方便根据姓名拿到学号
	//利用学号进行命名
	m["杜勇敢"] = "2020170229"
	m["孟凡森"] = "2020170230"
	m["李文博"] = "2020170231"
	m["杨江帆"] = "2020170232"
	m["边尚琪"] = "2020170233"
	m["王雷"] = "2020170234"
	m["杨庚"] = "2020170235"
	m["涂启添"] = "2020170236"
	m["冯也"] = "2020170237"
	m["窦浩天"] = "2020170238"
	m["权家友"] = "2020170239"
	m["贵兴锋"] = "2020170240"
	m["赵欢"] = "2020170241"
	m["孙一鸣"] = "2020170242"
	m["杜志杰"] = "2020170243"
	m["帅深龙"] = "2020170244"
	m["徐宏坤"] = "2020170245"
	m["陈鑫"] = "2020170246"
	m["鹿牧野"] = "2020170247"
	m["周帆"] = "2020170248"
	m["虢锐"] = "2020170249"
	m["尹群群"] = "2020170250"
	m["白杨"] = "2020170251"
	m["孟丽"] = "2020170252"
	m["陈瀚"] = "2020170253"
	m["李润"] = "2020170254"
	m["李掌珠"] = "2020170255"
	m["杨子健"] = "2020170256"
	m["李晓龙"] = "2020170257"
	m["赵麟寒"] = "2020170258"
	m["蒋晓明"] = "2020170259"
	m["冯家乐"] = "2020170260"
	m["石行"] = "2020170261"
	m["赵仁陈"] = "2020170262"
	m["邓峰"] = "2020170263"
	m["王艺超"] = "2020170264"
	m["赵智健"] = "2020170265"
	m["韩大力"] = "2020170266"
	m["陆国庆"] = "2020170267"
	m["陈涛"] = "2020170268"
	m["周紫剑"] = "2020170269"
	m["刘亦丰"] = "2020170270"
	m["姜晨皓"] = "2020170271"
	m["丁颖"] = "2020170272"
	m["黄以豪"] = "2020170273"
	m["雷俊"] = "2020170274"
	m["贺竞娇"] = "2020170275"
	m["范永正"] = "2020170276"
	m["肖锴"] = "2020170277"
	m["张子豪"] = "2020170278"
	m["钱坤"] = "2020170279"
	m["孙凯"] = "2020170280"
	m["林健树"] = "2020170281"
	m["江文涛"] = "2020170282"
	m["李继恩"] = "2020170283"
	m["牛谊博"] = "2020170284"
	m["赵梦梓"] = "2020170285"

	fmt.Println("starting ... ")
	//得到程序运行的路径
	getwd, err2 := os.Getwd()
	if err2 != nil {
		fmt.Println(err2.Error())
		os.Exit(1)
		return
	}

	//路径下的public
	folder := path.Join(getwd, "public")
	_, err2 = os.Stat(folder)
	//没有public文件夹就创建一个
	if os.IsNotExist(err2) {
		err2 = os.Mkdir(folder, 0755)
		if err2 != nil {
			fmt.Println(err2.Error())
			os.Exit(1)
			return
		}
	}

	Folder = folder

	// FileServer返回一个使用FileSystem接口root提供文件访问服务的HTTP处理器
	//所以使用 ip地址:端口号/  就可以查看所有的文件(在public路径下的)
	http.Handle("/", http.FileServer(http.Dir(Folder)))

	//ip地址:端口号/file  上传文件操作
	http.HandleFunc("/file", uploadFile)

	//ip地址:端口号/uploadPage  打开上传界面
	http.HandleFunc("/uploadPage", uploadPage)
	http.ListenAndServe(":8766", nil)
}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("test.html"))
	temp.Execute(w, nil)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	var data string
	data = "success!"

	//获得文件
	r.ParseMultipartForm(32 << 20)
	fmt.Println(r.MultipartForm.File)
	file := r.MultipartForm.File["uploadify"][0]
	r.ParseForm()

	//获得姓名
	name := r.PostFormValue("name")
	fmt.Println(name)

	//打开文件
	open, err := file.Open()
	if err != nil {
		fmt.Println("open")
		fmt.Println(err)
		return
	}
	defer open.Close()

	//读取文件
	all, err := ioutil.ReadAll(open)
	if err != nil {
		fmt.Println("readAll")
		fmt.Println(err)
		return
	}

	//根据姓名匹配学号  得到新的文件名称
	if sno, ok := m[name]; ok {
		fmt.Println(ok)
		newname := fmt.Sprintf("%s-%s", sno, name)
		index := strings.LastIndex(file.Filename, ".")
		file.Filename = newname + file.Filename[index:]
	} else {
		fmt.Println(ok)
		data = "error ,您输入的姓名不正确,请重新上传!"
	}

	//文件名称(带路径 /public/2020170281-林健树.jpg)
	filename := path.Join(Folder, file.Filename)

	//创建文件
	openFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("openFile")
		fmt.Println(err)
		return
	}

	//把刚刚读取到的东西存进去
	write, err := openFile.Write(all)
	fmt.Println(write)
	if err != nil {
		fmt.Println("write ")
		fmt.Println(err)
	}
	defer openFile.Close()

	//成功  返回状态码200
	w.WriteHeader(200)
	//跳转到绿绿的 success网页
	funcMap := template.FuncMap{"check": checkSuccess}
	tmpl := template.New("result.html").Funcs(funcMap)
	template := template.Must(tmpl.ParseFiles("result.html"))
	template.Execute(w, data)
}

func checkSuccess(s string) bool {
	if s == "success!" {
		return true
	}
	return false
}
