package main

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"stones/comment"
	"stones/dynamic"
	"stones/exp"
	"stones/login"
	"stones/post"
	"stones/tempsuid"
	"time"
)

func main() {

	User := login.NewUser()

	err := User.Login()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = sign(User)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\n----开始点赞----")
	sum := 0
	page := 1
	for {
		time.Sleep(time.Second)
		list, err := post.GetPostList(User, page)
		if err != nil {
			fmt.Println(err)
		}
		likeNum, err := list.Like()
		if err != nil {
			fmt.Println(err)
		}
		sum += likeNum
		if sum >= 5 {
			fmt.Println("----点赞完成----")
			break
		} else {
			page++
		}
	}

	time.Sleep(time.Second)
	fmt.Println("\n----开始发评论----")
	err = comment.Comment(User)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("----评论完成----")

	fmt.Println("\n----将在5秒后开始发动态----")
	time.Sleep(5 * time.Second)
	_, err = dynamic.Dynamic(User)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("----发动态完成----")

	fmt.Println("\n----将在5秒后开始获取经验----")
	fmt.Println("\n----如果不需要刷经验请按任意键退出----")
	fmt.Println()
	go func() {
		time.Sleep(5 * time.Second)
		err = exp.Exp(User)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	b := make([]byte, 1)
	os.Stdin.Read(b)
}

func sign(user *login.UserData) error {
	signPath := `https://apiff14risingstones.web.sdo.com/api/home/sign/signIn?tempsuid=`
	tempsUid, err := tempsuid.Get()
	if err != nil {
		return err
	}
	signPath = signPath + tempsUid
	req, err := http.NewRequest("POST", signPath, nil)
	if err != nil {
		return err
	}
	resp, err := user.GetClient().Do(req)
	if err != nil {
		return err
	}
	type resultBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			SqMsg          string `json:"sqMsg"`
			ContinuousDays int    `json:"continuousDays"`
			TotalDays      string `json:"totalDays"`
			SqExp          int    `json:"sqExp"`
			ShopExp        int    `json:"shopExp"`
		} `json:"data"`
	}
	re := new(resultBody)
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	result = bytes.TrimSpace(result)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(result, re)
	if err != nil {
		return err
	}
	fmt.Println(re.Msg)
	if re.Code == 10000 {
		fmt.Println(re.Data.SqMsg)
	}
	return nil
}
