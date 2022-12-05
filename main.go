package main

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var m = make(map[string]string)

var Folder string

// 学号=>姓名的映射集合
var mSnoToSName = make(map[string]string)

// Config 定义配置文件 从配置文件中读取
var Config = struct {
	Hour   int
	Minute int
	Sync   bool
}{}

// Report 定义结构体 行号 学号 姓名 是否上传 上传时间
type Report struct {
	LineId     int
	Sno        string
	Sname      string
	IsUpload   string
	UploadTime string
}

// 时区修正 东八区
var zone = time.FixedZone("CST", 8*3600)

// 集合 m[2020170281-林健树]=...
var mReport = make(map[string]*Report, 57)

func main() {
	//日志初始化
	logrus.SetFormatter(&logrus.TextFormatter{})
	logHandler, err := os.Create("debug.log")
	if err != nil {
		logrus.Error("debugLog init failed :", err)
		panic(err.Error())
	}
	logrus.SetOutput(logHandler)

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

	for key, value := range m {
		mSnoToSName[value] = key
	}

	logrus.Info(time.Now().In(zone).Format("2006-01-02") + "starting ...")
	//得到程序运行的路径
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Error("getPwd failed , err :", err)
		panic(err)
	}

	//路径下的public
	folder := path.Join(pwd, "public")
	_, err = os.Stat(folder)

	//没有public文件夹就创建一个
	if os.IsNotExist(err) {
		err = os.Mkdir(folder, 0755)
		if err != nil {
			logrus.Error("mkdir "+folder+" failed", err)
			panic(err)
		}
	}

	//静态变量
	Folder = folder

	//读文件夹 查看已经上传的文件
	dir, err := os.ReadDir(folder)
	if err != nil {
		logrus.Error("read dir fatherPath:"+folder+" failed , err : ", err)
		panic(err)
	}

	//初始化map mReport
	for _, entry := range dir {
		report := new(Report)
		report.IsUpload = "已上传🎉🎉🎉"
		info, err := entry.Info()
		if err != nil {
			logrus.Error(err.Error())
			report.UploadTime = time.Now().In(zone).Format("2006-01-02 15:04:05")
		} else {
			report.UploadTime = info.ModTime().In(zone).Format("2006-01-02 15:04:05")
		}

		prefix := entry.Name()[:strings.LastIndex(entry.Name(), ".")]
		split := strings.Split(prefix, "-")
		sno := split[0]
		sName := split[1]
		report.Sno = sno
		report.Sname = sName
		mReport[sno+"-"+sName] = report
	}

	//读取配置文件
	err = configor.Load(&Config, "config.json")
	if err != nil {
		logrus.Error(err.Error())
		panic(err.Error())
	}

	//开启一个子协程 各跑各的 到配置的时间了就移动文件位置到yyyy-mm-dd文件夹中 重新初始化一遍public文件夹为空
	go func() {
		for {
			nowTime := time.Now().In(zone)
			//时间达到指定的配置文件中要求的时间的时候
			if nowTime.Hour() == Config.Hour && nowTime.Minute() == Config.Minute {
				//初始化 重置mReport
				mReport = map[string]*Report{}

				//移动当前文件夹的内容到历史文件夹 yyyy-mm-dd文件夹
				fatherPath, _ := path.Split(folder)
				historyFolderName := nowTime.Format("2006-01-02")
				newPath := path.Join(fatherPath, historyFolderName)

				//移动文件的方式使用的是重命名 然后新建一个文件夹的操作
				err = os.Rename(folder, newPath)
				if err != nil {
					logrus.Error("rename folder ", folder, " to newFolder ", historyFolderName, " failed ,err: ", err)
					time.Sleep(60 * time.Second)
					continue
				}
				err := os.MkdirAll(folder, 0755)
				if err != nil {
					logrus.Error("mkdir ", folder, " failed , err: ", err)
					time.Sleep(60 * time.Second)
					continue
				}

				time.Sleep(60 * time.Second)
			}
		}
	}()

	//开启子协程 每个15分钟 重新加载一遍配置文件 可以热启动 不需要重启程序获得最新配置参数
	go func() {
		tick := time.Tick(15 * time.Minute)
		for {
			select {
			case <-tick:
				configor.Load(&Config, "config.json")
				logrus.Info("tick 循环")
			}
		}
	}()

	// FileServer返回一个使用FileSystem接口root提供文件访问服务的HTTP处理器
	//所以使用 ip地址:端口号/  就可以查看所有的文件(在public路径下的)
	http.Handle("/", http.FileServer(http.Dir(Folder)))

	//ip地址:端口号/file  上传文件操作
	http.HandleFunc("/file", uploadFile)

	//ip地址:端口号/uploadPage  打开上传界面
	http.HandleFunc("/uploadPage", uploadPage)
	http.ListenAndServe(":8766", nil)
}

// 打开上传文件的界面
func uploadPage(w http.ResponseWriter, r *http.Request) {
	//从mReport集合中加载出已经上传的同学 其他人都是未上传 按照学号
	reports := make([]Report, 0, len(mReport))
	start := 2020170229
	end := 2020170285
	count := 1
	for start <= end {
		report := new(Report)
		report.LineId = count
		sno := strconv.Itoa(start)
		sName := mSnoToSName[sno]
		report.Sno = sno
		report.Sname = sName
		if value, ok := mReport[sno+"-"+sName]; ok {
			value.LineId = count
			reports = append(reports, *value)
			count++
			start++
			continue
		}
		report.IsUpload = "未上传🧬🧬🧬"
		report.UploadTime = ""
		reports = append(reports, *report)
		start++
		count++
	}
	//渲染test.html 的模板
	temp := template.Must(template.ParseFiles("test.html"))
	//将待渲染的数据 丢到模板中展示
	temp.Execute(w, reports)
}

// 上传文件的操作
func uploadFile(w http.ResponseWriter, r *http.Request) {
	//var data string
	//data = "success!"

	report := new(Report)

	//获得文件
	r.ParseMultipartForm(32 << 20)
	fmt.Println(r.MultipartForm.File)

	file := r.MultipartForm.File["uploadify"][0]
	r.ParseForm()

	//获得姓名
	name := r.PostFormValue("name")
	report.Sname = name
	logrus.Info("姓名 ", name)

	//打开文件
	open, err := file.Open()
	defer open.Close()

	if err != nil {
		logrus.Error("file ", file.Filename, " open fail , err : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("file " + file.Filename + " open fail , err : " + err.Error() + " , 请联系管理员"))
		return
	}

	//读取文件
	all, err := io.ReadAll(open)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("readAll fail" + err.Error() + " , 请联系管理员"))
		return
	}

	var newName string
	//根据姓名匹配学号  得到新的文件名称
	if sno, ok := m[name]; ok {
		logrus.Info("匹配学号是否成功:", ok)
		report.Sno = sno

		newName = fmt.Sprintf("%s-%s", sno, name)
		index := strings.LastIndex(file.Filename, ".")
		file.Filename = newName + file.Filename[index:]
	} else {
		logrus.Info("匹配学号是否成功:", ok)
		//data = "error ,您输入的姓名不正确,请重新上传!"
	}

	//文件名称(带路径 /public/2020170281-林健树.jpg)
	filename := path.Join(Folder, file.Filename)

	//创建文件
	openFile, err := os.Create(filename)
	if err != nil {
		logrus.Error("create file ", filename, " fail , err : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("create file " + filename + " fail , err : " + err.Error()))
		return
	}

	//把刚刚读取到的东西存进去
	_, err = openFile.Write(all)
	defer openFile.Close()
	if err != nil {
		logrus.Error("write to newFile ", openFile.Name(), " fail , err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("write to newFile " + openFile.Name() + " fail , err: " + err.Error()))
		return
	}

	report.IsUpload = "已上传🎉🎉🎉"
	report.UploadTime = time.Now().In(zone).Format("2006-01-02 15:04:05")

	//更新map
	mReport[newName] = report

	//是否同步到核酸打卡
	if Config.Sync {
		url := "http://localhost:9999/fuckVote"
		method := "POST"

		//携带的数据
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("name", name)
		err := writer.Close()
		if err != nil {
			logrus.Error("close writer fail , err :	", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("close writer fail , err :	" + err.Error()))
			return
		}

		//创建新客户端
		client := &http.Client{}
		//创建新的请求
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			logrus.Error("newRequest fail , method: ", method, " url : "+url+" payload ", payload, " err: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("newRequest fail , method: " + method + " url : " + url + " payload " + payload.String() + " err: " + err.Error()))
			return
		}
		//设置header
		req.Header.Set("Content-Type", writer.FormDataContentType())
		//发送请求 打卡完成
		res, err := client.Do(req)
		if err != nil {
			logrus.Error("send request fail , err :", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("send request fail , err :" + err.Error()))
			return
		}
		defer res.Body.Close()

		_, err = io.ReadAll(res.Body)
		if err != nil {
			logrus.Error("read body fail , err :", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("read body fail , err :" + err.Error()))
			return
		}
	}

	//成功  返回状态码200
	//w.WriteHeader(200)

	//重定向到首页
	http.Redirect(w, r, "/uploadPage", http.StatusFound)

	//跳转到绿绿的 success网页
	//funcMap := template.FuncMap{"check": checkSuccess}
	//tmpl := template.New("result.html").Funcs(funcMap)
	//template := template.Must(tmpl.ParseFiles("result.html"))
	//template.Execute(w, data)
}

func checkSuccess(s string) bool {
	if s == "success!" {
		return true
	}
	return false
}
