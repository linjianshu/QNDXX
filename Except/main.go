package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
)

func main() {
	dir, err := ioutil.ReadDir("C:\\Users\\YourTreeDad\\Desktop\\解封核酸报告20220612")
	if err != nil {
		panic(err)
		return
	}

	start := 2020170229
	end := 2020170285
	m := make(map[string]bool)
	for i := start; i <= end; i++ {
		itoa := strconv.Itoa(i)
		m[itoa] = true
	}

	sort.Slice(dir, func(i, j int) bool {
		return dir[i].Name() < dir[j].Name()
	})

	for _, info := range dir {
		if m[info.Name()[:10]] {
			delete(m, info.Name()[:10])
		}
	}

	m1 := make(map[string]struct{})
	m1["2020170269"] = struct{}{} //周紫剑
	m1["2020170272"] = struct{}{} //丁颖
	m1["2020170274"] = struct{}{} //雷俊
	m1["2020170276"] = struct{}{} //范永正
	m1["2020170278"] = struct{}{} //张子豪
	m1["2020170279"] = struct{}{} //钱坤
	m1["2020170280"] = struct{}{} //孙凯
	m1["2020170282"] = struct{}{} //江文涛

	for key := range m {
		if _, ok := m1[key]; !ok {
			fmt.Println(key)
		}
	}
}
