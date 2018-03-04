package main

import (
	"github.com/ChimeraCoder/anaconda"
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
	"net/url"
	"github.com/parnurzeal/gorequest"
)

const DEFAULT_CONFIG = "config.json"
const DEFAULT_MEMBERINFO = "member.json"
const DEFAULT_OUTPUT = "output.json"

func main() {
	config := loadConfig()
	memberFollowCount := make(map[string]int)

	api := anaconda.NewTwitterApiWithCredentials(
		config["access-token"].(string),
		config["your-access-token-secret"].(string),
		config["your-consumer-key"].(string),
		config["your-consumer-secret"].(string))

	members := loadMemberInfo()

	for _, member := range members{
		fmt.Println("processing", member.MemberName)
		users, err := api.GetUserSearch(member.PageName, nil)
		if err != nil {
			panic(err)
		}
		if len(users) > 0{
			listUser := getAllFollowersList(api, users[0].Id)
			addUnique(&member, listUser)
		}
		memberFollowCount[member.MemberName] = memberFollowCount[member.MemberName] + len(member.Follower)
		time.Sleep(1 * time.Second)
	}

	for k, v := range memberFollowCount{
		t := map[string]interface{}{
			"BnkName" : k,
			"Followers" : v,
			"Timestamp" : time.Now().Unix()*1000,
		}

		j, err := json.Marshal(t)
		if err != nil{
			panic(err)
		}
		publishJsonMetric(config, j)
	}

}

func publishJsonMetric(config map[string]interface{}, json []byte){
	req := gorequest.New()
	req.Post(fmt.Sprintf("http://localhost:9200/%s/%s", config["metric-index"], config["metric-type"])).Send(string(json)).End()
}
func addUnique(member *MemberInfo, followers []int64){
	for _, f := range followers{
		if inTheList(member.Follower, f){
			member.Follower = append(member.Follower, f)
		}
	}
}
func inTheList(l1 []int64, item int64) bool{
	for _, fid := range l1{
		if fid == item{
			return false
		}
	}
	return true
}
func getAllFollowersList(api *anaconda.TwitterApi, id int64) []int64{
	users := make([]int64, 0)
	v := url.Values{}
	for {
		c, err := api.GetFollowersUser(id, v)
		if err != nil{
			break
		}
		users = append(users, c.Ids...)
		if err != nil || c.Next_cursor_str == "0" {
			break
		}
		time.Sleep(100*time.Millisecond)
		v.Set("cursor", c.Next_cursor_str)
	}

	return users
}
func loadMemberInfo() []MemberInfo{
	b := loadFile(DEFAULT_MEMBERINFO)
	members := make([]MemberInfo, 28)
	if err := json.Unmarshal(b, &members); err != nil{
		panic(err)
	}

	return members
}
func loadConfig() map[string]interface{}{
	configFileName := getConfigFileName()
	checkConfigExist(configFileName)
	b := loadFile(configFileName)
	out := make(map[string]interface{})
	if err := json.Unmarshal(b, &out); err != nil{
		panic(err)
	}

	return out
}
func loadFile(filename string) ([]byte) {
	b, err := ioutil.ReadFile(filename);
	if err != nil {
		panic(err)
	}
	return b
}
func checkConfigExist(configName string) {
	if !checkFileExist(configName) {
		panic(fmt.Sprintln("Couldn't load script properly, please check", configName, "filename"))
	}
}
func checkFileExist(f string) bool{
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}
	return true
}
func getConfigFileName() string{
	args := os.Args[1:]
	if len(args) > 0{
		return args[len(args) -1]
	}else{
		return DEFAULT_CONFIG
	}
}