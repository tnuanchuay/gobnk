package main

type MemberInfo struct {
	MemberName	string `json:"memberName"`
	PageName  	string `json:"pageName"`
	Follower	[]int64
}
