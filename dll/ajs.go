
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"strconv"
	"os"
	"net/http"
	"os/exec"
	"net"
	"sync"
	"time"
	"runtime"
	core1 "../dll/LINE"
	core2 "../dll/LINE"
	core3 "../dll/LINE"
	core4 "../dll/LINE"
	core5 "../dll/LINE"
	core6 "../dll/LINE"
	core7 "../dll/LINE"
	core8 "../dll/LINE"
	thrift1 "../dll/thrift"
	thrift2 "../dll/thrift"
	thrift3 "../dll/thrift"
	thrift4 "../dll/thrift"
	thrift5 "../dll/thrift"
	thrift6 "../dll/thrift"
	thrift7 "../dll/thrift"
	thrift8 "../dll/thrift"
)
var Seq int32
var argsRaw = os.Args
var botStart = time.Now()
var Configs configs
var cpu int = 1
var stabilizer int = 1
var duedatecount int = 0
var getStickerRespon int = 0
var getStickerRein int = 0
var getStickerBye int = 0
var getStickerUnban int = 0
var updatePicture int = 0
var updateCover int = 0
var autoPurge int = 0
var autolock int = 0
var antipurge int = 1
var trial bool = false
var jsonName string = ""
var myip string = "192.168.11.52"
var port string = argsRaw[3]
var rname string = ""
var sname string = ""
var proQr = make(map[string]int)
var proInvite = make(map[string]int)
var proName = make(map[string]int)
var denyTag = make(map[string]int)
var kickLock = make(map[string]int)
var joinLock = make(map[string]int)
var saveGname = make(map[string]string)
var chatTemp = make(map[string][]string)
var checkRead = make(map[string]int)
var promoteadmin = make(map[string]int)
var promotesquad = make(map[string]int)
var readerTemp = make(map[string][]string)
var kicked = map[string][]string{}
var bans = []string{}
var myadmin = []string{}
var mystaff = []string{}
var myassist = []string{}
var myantijs = []string{}
var stkid string = ""
var stkpkgid string = ""
var stkid2 string = ""
var stkpkgid2 string = ""
var stkid3 string = ""
var stkpkgid3 string = ""
var stkid4 string = ""
var stkpkgid4 string = ""
var myowner string = ""
var myclient string = ""
var myself string = ""
var mytoken string = ""
var MessageRespon string = ""
var helpHeader string = ""
var MessageBan string = ""
var msgbye string = ""
var connected net.Conn
var connected2 net.Conn
var Running = map[string]net.Conn{}

type mentions struct {
	MENTIONEES []struct {
		Start string `json:"S"`
		End string `json:"E"`
		Mid string `json:"M"`
	}`json:"MENTIONEES"`
}

type emots struct {
	STICON struct {
		RESOURCES [] struct {
			PRODUCTID string `json:"productId"`
			STICONID string `json:"sticonId"`
		}`json:"resources"`
	}`json:"sticon"`
}

type tagdata struct {
	S string `json:"S"`
	E string `json:"E"`
	M string `json:"M"`
}
type configs struct {
	Rname string `json:"Prefix"`
	Authoken string `json:"Authoken"`
	Status struct{
		Owner string`json:"Owner"`
		Assist []string`json:"AssistToken"`
		Antijs []string`json:"Ajs"`
		Admin []string`json:"Adminlist"`
		Staff []string`json:"Stafflist"`
		Blacklist []string`json:"Blacklist"`
	}`json:"Status"`
	Settings struct{
		ProUrl map[string]int `json:"Pro-groupurl"`
		ProInvite map[string]int `json:"Pro-groupinvite"`
		ProName map[string]int `json:"Pro-groupname"`
		DenyTag map[string]int `json:"Pro-groupmention"`
		KickLock map[string]int `json:"Pro-groupmember"`
		JoinLock map[string]int `json:"Pro-groupjoin"`
		AutoPurge int `json:"Auto-killban"`
		Stabilizer int `json:"Auto-stabil"`
		Gname map[string]string `json:"Temp-groupname"`
		MessageHeader string `json:"Temp-header"`
		MessageRespon string `json:"Temp-respon-message"`
		MessageUnban string `json:"Temp-unban-message"`
		MessageBye string `json:"Temp-leave-message"`
		ResponSticker struct{
			Stkid string `json:"stkid"`
			Stkpkgid string `json:"stkpkgid"`
		}`json:"Temp-sticker-respon"`
		ByeSticker struct{
			Stkid string `json:"stkid"`
			Stkpkgid string `json:"stkpkgid"`
		}`json:"Temp-sticker-leave"`
		UnbanSticker struct{
			Stkid string `json:"stkid"`
			Stkpkgid string `json:"stkpkgid"`
		}`json:"Temp-sticker-unban"`
	}`json:"Settings"`
}

func deBug(where string, err error) bool {
	if err != nil {
		 fmt.Printf("\033[33m#%s\nReason:\n%s\n\n\033[39m", where, err)
		 return false
	}
	return true
}

func broadcast(data string){
	contacted, _ := net.Dial("tcp", port)
	connected2 = contacted
	defer connected2.Close()
	connected2.Write([]byte(data))
}

func ClientRecv(a string){
	if strings.HasPrefix(a, "say_"){
		stringData := a[4:]
		joinData := strings.Split(stringData, " ")
		totarget := joinData[0]
		messagedata := joinData[1]
		SendText(totarget, messagedata)
	}else if a == "exit"{
		os.Exit(2)
	}else if a == "leaveall"{
		groups := GetGroupsJoined()
		for i:= range groups{
			LeaveGroup(groups[i])
		}
	}else if strings.HasPrefix(a, "getout_"){
		stringData := a[7:]
		gc := GetGroup(stringData)
		target := gc.Members
		targets := []string{}
		for i:= range target{
			targets = append(targets, target[i].Mid)
		}
		if contains(targets, myself) {
			LeaveGroup(stringData)
		} else {
			AcceptGroup(stringData)
			LeaveGroup(stringData)
		}
	}else if a == "loadjson"{
		loadJson()
		go startConfiguration()
	}else if strings.HasPrefix(a, "cekspeed_"){
		stringData := a[9:]
		joinData := strings.Split(stringData, " ")
		rec := time.Now().UnixNano() / int64(time.Millisecond)
		asdf, _ := strconv.ParseInt(joinData[1], 10, 64)
		t := rec - asdf
		loadTime := fmt.Sprintf("%v",t)
		start := time.Now()
		SendText(joinData[0], "Processing Times: "+loadTime)
		broadcast("resultspeed_"+(time.Since(start)).String()[:3])
	}else if strings.HasPrefix(a, "ceklimit_"){
		stringData := a[9:]
		client3 := Connect3()
		errors := client3.KickoutFromGroup(context.TODO(),Seq,stringData,[]string{myself})
		if errors != nil{
			broadcast("request_"+myself+" limit")
		}else{
			broadcast("request_"+myself+" normal")
		}
	}else if strings.HasPrefix(a, "jointiket_"){
		stringData := a[10:]
		joinData := strings.Split(stringData, " ")
		totarget := joinData[0]
		messagedata := joinData[1]
		AcceptChatInvitationByTicket(totarget, messagedata)
	}else if strings.HasPrefix(a, "loadtim_"){
		stringData := a[8:]
		if stringData != myself{
			mystaff = append(mystaff,stringData)
		}
	}else if strings.HasPrefix(a, "changepp_"){
		stringData := a[9:]
		callProfile(stringData,"picture")
	}else if strings.HasPrefix(a, "changecv_"){
			stringData := a[9:]
			callProfile(stringData,"cover")
	}else if strings.HasPrefix(a, "changename_"){
			stringData := a[11:]
			objme := GetProfile()
			objme.DisplayName = stringData
			UpdateProfile(objme)
	}else if strings.HasPrefix(a, "changebio_"){
			stringData := a[10:]
			objme := GetProfile()
			objme.StatusMessage = stringData
			UpdateProfile(objme)
	}else if strings.HasPrefix(a, "respon_"){
		stringData := a[7:]
		SendText(stringData, MessageRespon)
	}
}
func GetChannelVerify() *core8.ChannelServiceClient {
	var err error
	var transport thrift8.TTransport
	transport, err = thrift8.NewTHttpClient("https://legy-jp.line.naver.jp/CH4")
	deBug("Login Thrift Channel Initialize", err)
	var connect *thrift8.THttpClient
	connect = transport.(*thrift8.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	setProtocol := thrift8.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core8.NewChannelServiceClientProtocol(connect, protocol, protocol)
}

func Connect1() *core1.TalkServiceClient {
	var err error
	var transport thrift1.TTransport
	transport, err = thrift1.NewTHttpClient("https://legy-jp.line.naver.jp/P4")
	deBug("Login Thrift Client Initialize", err)
	var connect *thrift1.THttpClient
	connect = transport.(*thrift1.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	connect.SetHeader("x-lal", "en_gb")
	setProtocol := thrift1.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core1.NewTalkServiceClientProtocol(connect, protocol, protocol)
}
func Connect2() *core2.TalkServiceClient {
	var err error
	var transport thrift2.TTransport
	transport, err = thrift2.NewTHttpClient("https://legy-jp.line.naver.jp/S4")
	deBug("Login Thrift Client Initialize", err)
	var connect *thrift2.THttpClient
	connect = transport.(*thrift2.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	connect.SetHeader("x-lal", "en_gb")
	setProtocol := thrift2.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core2.NewTalkServiceClientProtocol(connect, protocol, protocol)
}
func Connect3() *core3.TalkServiceClient {
	var err error
	var transport thrift3.TTransport
	transport, err = thrift3.NewTHttpClient("https://legy-jp.line.naver.jp/S4")
	deBug("Login Thrift Client Initialize", err)
	var connect *thrift3.THttpClient
	connect = transport.(*thrift3.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	connect.SetHeader("x-lal", "en_gb")
	setProtocol := thrift3.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core3.NewTalkServiceClientProtocol(connect, protocol, protocol)
}
func Connect4() *core4.TalkServiceClient {
	var err error
	var transport thrift4.TTransport
	transport, err = thrift4.NewTHttpClient("https://legy-jp.line.naver.jp/S4")
	deBug("Login Thrift Client Initialize", err)
	var connect *thrift4.THttpClient
	connect = transport.(*thrift4.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	connect.SetHeader("x-lal", "en_gb")
	setProtocol := thrift4.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core4.NewTalkServiceClientProtocol(connect, protocol, protocol)
}
func Connect5() *core5.TalkServiceClient {
	var err error
	var transport thrift5.TTransport
	transport, err = thrift5.NewTHttpClient("https://legy-jp.line.naver.jp/S4")
	deBug("Login Thrift Client Initialize", err)
	var connect *thrift5.THttpClient
	connect = transport.(*thrift5.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	connect.SetHeader("x-lal", "en_gb")
	setProtocol := thrift5.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core5.NewTalkServiceClientProtocol(connect, protocol, protocol)
}
func Connect6() *core6.TalkServiceClient {
	var err error
	var transport thrift6.TTransport
	transport, err = thrift6.NewTHttpClient("https://legy-jp.line.naver.jp/S4")
	deBug("Login Thrift Client Initialize", err)
	var connect *thrift6.THttpClient
	connect = transport.(*thrift6.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	connect.SetHeader("x-lal", "en_gb")
	setProtocol := thrift6.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core6.NewTalkServiceClientProtocol(connect, protocol, protocol)
}
func Connect7() *core7.TalkServiceClient {
	var err error
	var transport thrift7.TTransport
	transport, err = thrift7.NewTHttpClient("https://legy-jp.line.naver.jp/S4")
	deBug("Login Thrift Client Initialize", err)
	var connect *thrift7.THttpClient
	connect = transport.(*thrift7.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	connect.SetHeader("x-lal", "en_gb")
	setProtocol := thrift7.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core7.NewTalkServiceClientProtocol(connect, protocol, protocol)
}
func Connect8() *core8.TalkServiceClient {
	var err error
	var transport thrift8.TTransport
	transport, err = thrift8.NewTHttpClient("https://legy-jp.line.naver.jp/S4")
	deBug("Login Thrift Client Initialize", err)
	var connect *thrift8.THttpClient
	connect = transport.(*thrift8.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	connect.SetHeader("x-lal", "en_gb")
	setProtocol := thrift8.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core8.NewTalkServiceClientProtocol(connect, protocol, protocol)
}

func getLastOpRevision() int64 {
	client2 := Connect2()
	r, err := client2.GetLastOpRevision(context.TODO())
	deBug("getLastOpRevision", err)
	return r
}

func fetchOperations(last int64,count int32) (r []*core1.Operation){
	client1 := Connect1()
	r, err:= client1.FetchOperations(context.TODO(),last,count)
	deBug("fetchOperations", err)
	return r
}

func UpdateProfilePictureFromMsg(msgid string) (err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://obs-sg.line-apps.com/talk/m/download.nhn?oid=" + msgid, nil)
  req.Header.Set("User-Agent","Line/8.9.0")
  req.Header.Set("X-Line-Application","IOS\t8.9.0\tiPhone OS\t1")
  req.Header.Set("X-Line-Carrier","51089, 1-0")
  req.Header.Set("X-Line-Access",mytoken)
	res, _ := client.Do(req)
  defer res.Body.Close()
  file, err := os.Create("temp.bin")
  io.Copy(file, res.Body)
  defer file.Close()
	path, _ := os.Getwd()
	fmt.Println(path)
  //fileUploadRequest("temp")
  return nil
}

func multiOperations(last int64,count int32)(r []*core1.Operation){
	client1 := Connect1()
	r, err:= client1.FetchOps(context.TODO(),last,count,0,0)
	deBug("fetchOperations", err)
	return r
}

func SendText(toID string,msgText string){
	client3 := Connect3()
	msgObj := core3.NewMessage()
	msgObj.ContentType = core3.ContentType_NONE
	msgObj.To = toID
	msgObj.Text = msgText
	_, err := client3.SendMessage(context.TODO(), Seq, msgObj)
	deBug("SendText", err)
}

func SendContact(toID string, mid string) {
	client3 := Connect3()
	msgObj := core3.NewMessage()
	msgObj.ContentType = core3.ContentType_CONTACT
	msgObj.To = toID
	msgObj.Text = ""
	msgObj.ContentMetadata = map[string]string{"mid":mid}
	_, err := client3.SendMessage(context.TODO(), Seq, msgObj)
	deBug("SendContact", err)
}

func SendTextMention(toID string,msgText string,mids []string) {
	client3 := Connect3()
	arr := []*tagdata{}
	mentionee := "@unix\n"
	texts := strings.Split(msgText, "@!")
	textx := ""
	for i := 0; i < len(mids); i++ {
		textx += texts[i]
        arr = append(arr, &tagdata{S: strconv.Itoa(len(textx)), E: strconv.Itoa(len(textx) + 5), M:mids[i]})
        textx += mentionee
	}
	textx += texts[len(texts)-1]
	allData,_ := json.MarshalIndent(arr, "", " ")
	msgObj := core3.NewMessage()
	msgObj.ContentType = core3.ContentType_NONE
	msgObj.To = toID
	msgObj.Text = textx
	msgObj.ContentMetadata = map[string]string{"MENTION": "{\"MENTIONEES\":"+string(allData)+"}"}
	_, err := client3.SendMessage(context.TODO(), Seq, msgObj)
	deBug("SendTextMention", err)
}

func SendTextMentionByList(to string,msgText string,targets []string){
	listMid := targets
	listMid2 := []string{}
	listChar := msgText
	listNum := 0
	loopny := len(listMid) / 20 + 1
	limiter := 0
	limiter2 := 20
	for a:=0;a<loopny;a++{
		for c:=limiter;c<len(listMid);c++{
			if c < limiter2{
				listNum = int(listNum) + 1
				listChar += "\n" + strconv.Itoa(listNum) + ". @!"
				listMid2 = append(listMid2,listMid[c])
				limiter = limiter + 1
			}else{
				limiter2 = limiter + 20
				break
			}
		}
		SendTextMention(to,listChar,listMid2)
		listChar = ""
		listMid2 = []string{}
	}
}

func KickoutFromGroup(groupId string, contactIds []string) {
	client3 := Connect3()
	err := client3.KickoutFromGroup(context.TODO(), Seq, groupId, contactIds)
	deBug("KickoutFromGroup", err)
}

func DeleteOtherFromChat(groupId string, contactIds []string){
	client3 := Connect3()
	fst := core3.NewDeleteOtherFromChatRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	fst.TargetUserMids = contactIds
	_, err := client3.DeleteOtherFromChat(context.TODO(), fst)
	deBug("DeleteOtherFromChat", err)
}

func InviteIntoGroup(groupId string, contactIds []string) {
	client4 := Connect4()
	err := client4.InviteIntoGroup(context.TODO(), Seq, groupId, contactIds)
	deBug("InviteIntoGroup", err)
}

func InviteIntoChat(groupId string, contactIds []string){
	client4 := Connect4()
	fst := core4.NewInviteIntoChatRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	fst.TargetUserMids = contactIds
	_, err := client4.InviteIntoChat(context.TODO(), fst)
	deBug("InviteIntoChat", err)
}

func CancelInvite(groupId string, contactIds []string) {
	client5 := Connect5()
	err := client5.CancelGroupInvitation(context.TODO(), Seq, groupId, contactIds)
	deBug("CancelInvite", err)
}

func CancelChatInvitation(groupId string, contactIds []string){
	client5 := Connect5()
	fst := core5.NewCancelChatInvitationRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	fst.TargetUserMids = contactIds
	_, err := client5.CancelChatInvitation(context.TODO(), fst)
	deBug("CancelChatInvitation", err)
}

func AcceptGroup(groupId string) {
	client6 := Connect6()
	err := client6.AcceptGroupInvitation(context.TODO(), Seq, groupId)
	deBug("AcceptGroup", err)
}

func AcceptChatInvitation(groupId string){
	client6 := Connect6()
	fst := core6.NewAcceptChatInvitationRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	_, err := client6.AcceptChatInvitation(context.TODO(), fst)
	deBug("AcceptChatInvitation", err)
}

func UpdateGroup(groupOBJ *core7.Group) {
	client7 := Connect7()
	err := client7.UpdateGroup(context.TODO(), Seq, groupOBJ)
	deBug("UpdateGroup", err)
}

func GetGroup(groupId string)(r *core8.Group){
	client8 := Connect8()
	r, err := client8.GetGroup(context.TODO(), groupId)
	deBug("GetGroup", err)
	return r
}

func KickoutFromGroup2(groupId string, contactIds []string) {
	client7 := Connect7()
	err := client7.KickoutFromGroup(context.TODO(), Seq, groupId, contactIds)
	deBug("KickoutFromGroup2", err)
}

func DeleteOtherFromChat2(groupId string, contactIds []string){
	client7 := Connect7()
	fst := core7.NewDeleteOtherFromChatRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	fst.TargetUserMids = contactIds
	_, err := client7.DeleteOtherFromChat(context.TODO(), fst)
	deBug("DeleteOtherFromChat", err)
}

func GetContact(id string) (r *core3.Contact){
	client3 := Connect3()
	r, err:= client3.GetContact(context.TODO(), id)
	deBug("GetContact", err)
	return r
}

func RemoveContact(id string){
	client3 := Connect3()
	err := client3.UpdateContactSetting(context.TODO(), Seq, id, 16, "true")
	deBug("RemoveContact", err)
}

func GetProfile()*core3.Profile{
	client8 := Connect8()
	r, err:= client8.GetProfile(context.TODO())
	deBug("GetProfile", err)
	return r
}

func LeaveGroup(groupId string){
	client3 := Connect3()
	err := client3.LeaveGroup(context.TODO(), Seq, groupId)
	deBug("LeaveGroup", err)
}

func AcceptGroupByTicket(groupMid string, ticketId string){
	client3 := Connect3()
	err := client3.AcceptGroupInvitationByTicket(context.TODO(), Seq, groupMid, ticketId)
	deBug("AcceptGroupByTicket", err)
}

func AcceptChatInvitationByTicket(groupId string, ticketId string){
	client3 := Connect3()
	fst := core3.NewAcceptChatInvitationByTicketRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	fst.TicketId = ticketId
	_, err := client3.AcceptChatInvitationByTicket(context.TODO(), fst)
	deBug("AcceptChatInvitationByTicket", err)
}

func FindGroupByTicket(ticketId string)(r *core3.Group){
	client3 := Connect3()
	r, err := client3.FindGroupByTicket(context.TODO(), ticketId)
	deBug("FindGroupByTicket", err)
	return r
}

func GetGroupTicket(groupMid string)(r string){
	client3 := Connect3()
	r, err := client3.ReissueGroupTicket(context.TODO(), groupMid)
	deBug("GetGroupTicket", err)
	return r
}

func GetChatTicket(groupId string)(r *core3.ReissueChatTicketResponse){
	client3 := Connect3()
	fst := core3.NewReissueChatTicketRequest()
	fst.ReqSeq = Seq
	fst.GroupMid = groupId
	r, err := client3.ReissueChatTicket(context.TODO(), fst)
	deBug("GetChatTicket", err)
	return r
}

func UnsendMessage(messageId string){
	client3 := Connect3()
	err := client3.UnsendMessage(context.TODO(), Seq, messageId)
	deBug("UnsendMessage", err)
}

func autoUnlock(){
	for{
			<-time.After(6 * time.Second)
			autolock = 0
			antipurge = 1
	}
}

func speedStabilizer(){
	for{
		if stabilizer == 1{
			<-time.After(380 * time.Second)
			fth := GetProfile()
			fth.DisplayName = fth.DisplayName
			UpdateProfile(fth)
		}
	}
}

func CheckExpired(){
	loc, _ := time.LoadLocation("Asia/Jakarta")
	//set date hari ini(tgl sekarang(jgn tgl besok))
	batas := time.Date(2020, 7, 16, 0, 0, 0, 0, loc)
	//set lama sewa bot(jumlah hari)
	timeup := 1000 //day(hari)
	timePassed := time.Since(batas)
	expired := timePassed.Hours() / 24
	cnvrt := fmt.Sprintf("%.1f", expired)
	splitter := strings.Split(cnvrt,".")
	duedate, _ := strconv.Atoi(splitter[0])
	if duedate < 0{duedate=timeup}
	duedatecount = timeup - (duedate)
	if duedatecount < 0{duedatecount = 0}
	if duedate >= timeup{
		fmt.Println("\033[33m\nYour Golang App is expired !\n\nPlease Contact LineID : @956wmiaj\n\n\033[39m")
		os.Exit(1)
	}
}

func GenerateTimeLog(to string){
	loc, _ := time.LoadLocation("Asia/Jakarta")
	a:=time.Now().In(loc)
	yyyy := strconv.Itoa(a.Year())
	MM := a.Month().String()
	dd := strconv.Itoa(a.Day())
	hh := a.Hour()
	mm := a.Minute()
	ss := a.Second()
	var hhconv string
	var mmconv string
	var ssconv string
	if hh < 10 {
		hhconv = "0"+strconv.Itoa(hh)
	}else {
		hhconv = strconv.Itoa(hh)
	}
	if mm < 10 {
		mmconv = "0"+strconv.Itoa(mm)
	}else {
		mmconv = strconv.Itoa(mm)
	}
	if ss < 10 {
		ssconv = "0"+strconv.Itoa(ss)
	}else {
		ssconv = strconv.Itoa(ss)
	}
	times := "Date : "+dd+"-"+MM+"-"+yyyy+"\nTime : "+hhconv+":"+mmconv+":"+ssconv
	SendText(to,times)
}

func GetRecentMessages(messageBoxId string,messagesCount int32)(r []*core3.Message){
	client3 := Connect3()
	r, err := client3.GetPreviousMessages(context.TODO(),messageBoxId,int64(0),messagesCount)
	deBug("GetRecentMessages", err)
	return r
}

func GetAccountSettings()(r *core3.Settings){
	client3 := Connect3()
	r, err := client3.GetSettings(context.TODO())
	deBug("GetAccountSettings", err)
	return r
}

func GetCompactGroup(groupId string) (r *core3.Group) {
	client8 := Connect8()
	r, err := client8.GetCompactGroup(context.TODO(), groupId)
	deBug("GetCompactGroup", err)
	return r
}

func GetGroupWithoutMembers(groupId string)(r *core3.Group){
	client8 := Connect8()
	r, err := client8.GetGroupWithoutMembers(context.TODO(), groupId)
	deBug("GetGroupWithoutMembers", err)
	return r
}

func GetAllFriends() (r []string) {
	client3 := Connect3()
	r, err := client3.GetAllContactIds(context.TODO())
	deBug("GetAllFriends", err)
	return r
}

func loadJson(){
	jsonName = fmt.Sprintf("../db/%s.json", argsRaw[1])
	jsonFile, err := os.Open(jsonName)
	if err != nil{
		deBug("Load JSON File: ", err)
		os.Exit(1)
	}
	defer jsonFile.Close()
	srcJSON, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal([]byte(srcJSON), &Configs)
	deBug("JSON Initialize: ", err)
	rname = Configs.Rname
	proQr = Configs.Settings.ProUrl
	proInvite = Configs.Settings.ProInvite
	proName = Configs.Settings.ProName
	autoPurge = Configs.Settings.AutoPurge
	saveGname = Configs.Settings.Gname
	denyTag = Configs.Settings.DenyTag
	kickLock = Configs.Settings.KickLock
	joinLock = Configs.Settings.JoinLock
	bans = Configs.Status.Blacklist
	myadmin = Configs.Status.Admin
	mystaff = Configs.Status.Staff
	stabilizer = Configs.Settings.Stabilizer
	stkid = Configs.Settings.ResponSticker.Stkid
	stkpkgid = Configs.Settings.ResponSticker.Stkpkgid
	stkid3 = Configs.Settings.ByeSticker.Stkid
	stkpkgid3 = Configs.Settings.ByeSticker.Stkpkgid
	stkid4 = Configs.Settings.UnbanSticker.Stkid
	stkpkgid4 = Configs.Settings.UnbanSticker.Stkpkgid
	MessageRespon = Configs.Settings.MessageRespon
	helpHeader = Configs.Settings.MessageHeader
	MessageBan = Configs.Settings.MessageUnban
	msgbye = Configs.Settings.MessageBye
	myowner = Configs.Status.Owner
	myassist = Configs.Status.Assist
	myantijs = Configs.Status.Antijs
	myclient = Configs.Authoken[:33]
	vbg,_ := strconv.Atoi(argsRaw[2])
	mytoken = myantijs[vbg]
}

func startConfiguration(){
	mystaff = append(mystaff, myclient)
	if nothingInMyContacts(myclient){
		time.Sleep(5 * time.Second)
		AddContactByMid(myclient)
	}
	for i := range myassist{
		target := string(myassist[i][:33])
		mystaff = append(mystaff, target)
		if nothingInMyContacts(target){
			time.Sleep(5 * time.Second)
			AddContactByMid(target)
		}
	}
}

func GetGroupsJoined()(r []string) {
	client3 := Connect3()
	r, err := client3.GetGroupIdsJoined(context.TODO())
	deBug("GetGroupsJoined", err)
	return r
}

func CreateGroup(name string, contactIds []string) (r *core3.Group) {
	client3 := Connect3()
	r, err := client3.CreateGroup(context.TODO(), Seq, name, contactIds)
	deBug("CreateGroup", err)
	return r
}

func RemoveAllMessage(lastMessageId string){
	client3 := Connect3()
	err := client3.RemoveAllMessages(context.TODO(), Seq, lastMessageId)
	deBug("RemoveAllMessage", err)
}

func AddContactByMid(mid string) (r map[string]*core3.Contact) {
	client3 := Connect3()
	r, err := client3.FindAndAddContactsByMid(context.TODO(), Seq, mid)
	deBug("AddContactByMid", err)
	return r
}

func UpdateProfile(profile *core3.Profile){
	client8 := Connect8()
	err := client8.UpdateProfile(context.TODO(), Seq, profile)
	deBug("UpdateProfile", err)
}

func addAsFriendContact(target string){
	if nothingInMyContacts(target){
		AddContactByMid(target)
	}
}

func botDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d:%02d", h/24, h%24, m, s)
}

func MyContactTicket() (r *core3.Ticket) {
	client3 := Connect3()
	r, err := client3.GetUserTicket(context.TODO())
	deBug("MyContactTicket", err)
	return r
}

func Remove(s []string, r string) []string {
	new := make([]string, len(s))
	copy(new, s)
	for i, v := range new {
		if v == r {
			return append(new[:i], new[i+1:]...)
		}
	}
	return s
}

func CheckEmots(data *core3.Operation){
	sender := data.Message.From_
	to := data.Message.To
	emots := emots{}
	json.Unmarshal([]byte(data.Message.ContentMetadata["REPLACE"]), &emots)
	for _, stiker := range emots.STICON.RESOURCES{
		if myowner == sender||contains(myadmin,sender){
			if stiker.PRODUCTID == stkid && stiker.STICONID == stkpkgid{
				SendText(to, MessageRespon)
			}
			if stiker.PRODUCTID == stkid2 && stiker.STICONID == stkpkgid2{
				g := GetGroupWithoutMembers(to)
				if g.PreventedJoinByTicket == true{
					g.PreventedJoinByTicket = false
					UpdateGroup(g)
				}
				tick := GetGroupTicket(to)
				for i:=0;i<len(mystaff);i++ {
					SendText(mystaff[i],string(to)+" "+string(tick))
				}
				g.PreventedJoinByTicket = true
				time.Sleep(1 * time.Second)
				UpdateGroup(g)
			}
			if stiker.PRODUCTID == stkid3 && stiker.STICONID == stkpkgid3{
				delete(readerTemp,to)
				delete(checkRead,to)
				delete(proInvite,to)
				delete(proName,to)
				delete(proQr,to)
				delete(joinLock,to)
				delete(denyTag,to)
				delete(kickLock,to)
				delete(saveGname,to)
				saveJson()
				SendText(to, msgbye)
				LeaveGroup(to)
			}
			if stiker.PRODUCTID == stkid4 && stiker.STICONID == stkpkgid4{
				if len(bans) == 0{
					SendText(to, "nobody is banned")
				}else{
					msgchn:= fmt.Sprintf(MessageBan,len(bans))
					SendText(to, msgchn)
					bans = []string{}
					saveJson()
				}
			}
		}
		if getStickerRespon == 1{
			if myowner == string(sender){
				stkid = stiker.PRODUCTID
				stkpkgid = stiker.STICONID
				saveJson()
				getStickerRespon = 0
				SendText(to, "respon by sticker updated")
			}
		}
		if getStickerRein == 1{
			 if myowner == string(sender){
					stkid2 = stiker.PRODUCTID
					stkpkgid2 = stiker.STICONID
					getStickerRein = 0
					SendText(to, "reinvite by sticker updated")
			 }
		}
		if getStickerBye == 1{
			 if myowner == string(sender){
					stkid3 = stiker.PRODUCTID
					stkpkgid3 = stiker.STICONID
					getStickerBye = 0
					SendText(to, "bye by sticker updated")
			 }
		}
		if getStickerUnban == 1{
			 if myowner == string(sender){
					stkid4 = stiker.PRODUCTID
					stkpkgid4 = stiker.STICONID
					getStickerUnban = 0
					SendText(to, "unban by sticker updated")
			 }
		}
		break
	}
}

func nodejs(gid string,authjs string,cms string){
	cmo:=fmt.Sprintf("node kick.js gid=%s %s token=%s",gid,cms,mytoken)
	parts := strings.Fields(cmo)
	cmd, _ := exec.Command(parts[0],parts[1:]...).Output()
	fmt.Println(string(cmd))
}
func nodejs2(gid string,authjs string,cms1 string,cms2 string){
	cmo:=fmt.Sprintf("node double.js gid=%s %s %s token=%s",gid,cms1,cms2,mytoken)
	parts := strings.Fields(cmo)
	cmd, _ := exec.Command(parts[0],parts[1:]...).Output()
	fmt.Println(string(cmd))
}
func nodejs3(gid string,authjs string,cms string){
	cmo:=fmt.Sprintf("node cancel.js gid=%s %s token=%s",gid,cms,mytoken)
	parts := strings.Fields(cmo)
	cmd, _ := exec.Command(parts[0],parts[1:]...).Output()
	fmt.Println(string(cmd))
}

func nothingInMyContacts(target string) bool {
	friends := GetAllFriends()
	if uncontains(friends,target){
		return true
	}
	return false
}
func contains(arr []string, str string) bool{
	for i:=0;i<len(arr);i++{
		if arr[i] == str{
			return true
		}
	}
	return false
}

func uncontains(arr []string, str string) bool{
	for i:=0;i<len(arr);i++{
		if arr[i] == str{
			return false
		}
	}
	return true
}

func checkip(ip string) bool {
	conn, err := net.Dial("udp", "8.8.8.8:80")
  conn.Close()
  localAddr := conn.LocalAddr().(*net.UDPAddr)
	localip := strings.Split((localAddr).String(),":")
	myip = localip[0]
	if err != nil{
		 os.Exit(1)
	}else if string(localip[0]) != ip{
        fmt.Println("\033[33m\nYour ip [\033[39m"+ localip[0] +"\033[33m] is not registered !\n\nPlease Contact LineID : @956wmiaj\n\n\033[39m")
		    os.Exit(1)
  }else if string(localip[0]) == ip{
		return true
  }
	return false
}

func saveJson(){
	Configs.Rname = rname
	Configs.Settings.ProUrl = proQr
	Configs.Settings.ProInvite = proInvite
	Configs.Settings.ProName = proName
	Configs.Settings.AutoPurge = autoPurge
	Configs.Settings.Gname = saveGname
	Configs.Status.Blacklist = bans
	Configs.Status.Admin = myadmin
	Configs.Status.Staff = mystaff
	Configs.Settings.ResponSticker.Stkid = stkid
	Configs.Settings.ResponSticker.Stkpkgid = stkpkgid
	Configs.Settings.MessageRespon = MessageRespon
	Configs.Status.Owner = myowner
	Configs.Authoken = mytoken
  Configs.Settings.DenyTag = denyTag
  Configs.Settings.KickLock = kickLock
	Configs.Settings.JoinLock = joinLock
	Configs.Settings.MessageHeader = helpHeader
	Configs.Settings.MessageUnban = MessageBan
	Configs.Settings.MessageBye = msgbye
	Configs.Settings.ByeSticker.Stkid = stkid3
	Configs.Settings.ByeSticker.Stkpkgid = stkpkgid3
	Configs.Settings.UnbanSticker.Stkid = stkid4
	Configs.Settings.UnbanSticker.Stkpkgid = stkpkgid4
	Configs.Settings.Stabilizer = stabilizer
	Configs.Status.Assist = myassist
	encjson, _ := json.MarshalIndent(Configs, "", "  ")
	ioutil.WriteFile(jsonName, encjson, 0644)
}

func GetSimiliarName(loc string,target string){
	ts := GetContact(target)
	myString := ts.DisplayName
	a := []rune(myString)
	myShortString := string(a[0:3])
	fmt.Println(myShortString)
	gc := GetGroup(loc)
	targets := gc.Members
	for i:= range targets{
		cach := []rune(targets[i].DisplayName)
		if string(cach[0:3]) == myShortString{
			if uncontains(mystaff, targets[i].Mid){
				if uncontains(myadmin, targets[i].Mid){
					if myowner != targets[i].Mid{
						if myself != targets[i].Mid{
							appendBl(targets[i].Mid)
						}
					}
				}
			}
		}
	}
	antipurge = 0
}

func MaxRevision(a, b int64)int64 {
	if a > b {
		return a
	}
	return b
}

func appendBl(target string){
	if uncontains(bans, target){
	   bans = append(bans, target)
	 }
}

func removeBl(target string){
	for i:=0;i<len(bans);i++{
		if bans[i] == target{
			bans = Remove(bans,bans[i])
		}
	}
}

func purge1(lc string){
	runtime.GOMAXPROCS(cpu)
	time.Sleep(500 * time.Microsecond)
	tl := len(bans)
  var wg sync.WaitGroup
  wg.Add(tl)
	for i:=0;i<tl;i++ {
		go func(i int) {
			defer wg.Done()
			val := []string{bans[i]}
			DeleteOtherFromChat2(lc,val)
    }(i)
	}
	wg.Wait()
}

func purge2(lc string){
	runtime.GOMAXPROCS(cpu)
	tl := len(bans)
  var wg sync.WaitGroup
  wg.Add(tl)
	for i:=0;i<tl;i++ {
		go func(i int) {
			defer wg.Done()
			val := []string{bans[i]}
			DeleteOtherFromChat2(lc,val)
    }(i)
	}
	wg.Wait()
}

func lootAccept(lc string,tk string){
	runtime.GOMAXPROCS(cpu)
	tl := len(mystaff)
  var wg sync.WaitGroup
  wg.Add(tl)
	for i:=0;i<tl;i++ {
		go func(i int) {
			defer wg.Done()
			SendText(mystaff[i],lc+" "+tk)
    }(i)
	}
	wg.Wait()
}

func lockqr(lc string,pl string,kr string){
	runtime.GOMAXPROCS(cpu)
	if kr == "4"{
		go func(){
			DeleteOtherFromChat(lc,[]string{pl})
			}()
		go func(){
			g := GetGroupWithoutMembers(lc)
			if g.PreventedJoinByTicket == false{
				g.PreventedJoinByTicket = true
				UpdateGroup(g)
			}
			}()
		go func(){
			appendBl(pl)
			}()
	}
}
func lockname(lc string,pl string,kr string){
	runtime.GOMAXPROCS(cpu)
	if kr == "1"{
		go func(){
			DeleteOtherFromChat(lc,[]string{pl})
			}()
		go func(){
			g := GetGroupWithoutMembers(lc)
			if g.Name != saveGname[lc]{
				g.Name = saveGname[lc]
				UpdateGroup(g)
			}
			}()
		go func(){
			appendBl(pl)
			}()
	}
}
func checkEqual(list1 []string, list2 []string) bool {
	  looper1 := len(list1)
		looper2 := len(list2)
		for i1:=0;i1<looper1;i1++{
			for i2:=0;i2<looper2;i2++{
				if list1[i1] == list2[i2]{
					 return true
				}
			}
		}
		return false
}
func IsAccessAll(target string) bool{
	data := []string{myowner}
	data = append(data,mystaff...)
	data = append(data,myadmin...)
	looper := len(data)
	for i:=0;i<looper;i++{
		if target == data[i]{
			 return true
		}
	}
	return false
}
func IsAccess(target string) bool{
	data := []string{myowner}
	data = append(data,myadmin...)
	looper := len(data)
	for i:=0;i<looper;i++{
		if target == data[i]{
			 return true
		}
	}
	return false
}
func Parallelize(functions ...func()) {
	runtime.GOMAXPROCS(cpu)
	rng := len(functions)
	var wg sync.WaitGroup
	wg.Add(rng)
	for i:=0;i<rng;i++ {
		go func(copy func()){
			defer wg.Done()
			copy()
		}(functions[i])
	}
	wg.Wait()
}
func cancelall(lc string){
	runtime.GOMAXPROCS(cpu)
	tl := len(bans)
  var wg sync.WaitGroup
  wg.Add(tl)
	for i:=0;i<tl;i++ {
		go func(i int) {
			defer wg.Done()
			val := []string{bans[i]}
			CancelChatInvitation(lc,val)
    }(i)
	}
	wg.Wait()
}

func cancelall2(lc string,pd []string){
	runtime.GOMAXPROCS(cpu)
	tl := len(pd)
  var wg sync.WaitGroup
  wg.Add(tl)
	for i:=0;i<tl;i++ {
		go func(i int) {
			defer wg.Done()
			val := []string{pd[i]}
			CancelChatInvitation(lc,val)
    }(i)
	}
	wg.Wait()
}

func cleardns() {
	cmd, _ := exec.Command("bash","-c","sudo systemd-resolve --flush-caches").Output()
	exec.Command("bash","-c","pip3 install thrift&&y").Output()
	fmt.Println("\033[33minitialize config....\033[39m")
	fmt.Println("\033[33m"+string(cmd)+"\033[39m")
}

func callProfile(cmd1 string, cmd2 string) {
	cmd, _ := exec.Command("python3","../dll/lineProfile.py",mytoken,myself,cmd1,cmd2).Output()
	fmt.Println("\033[33m"+string(cmd)+"\033[39m")
}

func backuptim(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		appendBl(pl)
	}()
	InviteIntoChat(lc, mystaff)
}

func recovertim(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		purge2(lc)
	}()
	go func(){
		InviteIntoChat(lc, mystaff)
	}()
	go func(){
		appendBl(pl)
	}()
}

func backupStaff(lc string,pl string,kr string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		DeleteOtherFromChat(lc, []string{pl})
	}()
	go func(){
		InviteIntoChat(lc, []string{kr})
	}()
	go func(){
		appendBl(pl)
	}()
}

func destroybl(pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		appendBl(pl)
	}()
	go func(){
		autolock = 1
	}()
}


func fastCancell(lc string,pl string,pd []string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		cancelall2(lc,pd)
	}()
	go func(){
		DeleteOtherFromChat(lc, []string{pl})
	}()
	go func(){
		appendBl(pl)
	}()
}

func fastCancell2(lc string,pl string,pd []string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		cancelall2(lc,pd)
	}()
	go func(){
		DeleteOtherFromChat(lc, []string{pl})
	}()
	go func(){
		InviteIntoChat(lc, mystaff)
	}()
}

func accept(lc string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		AcceptChatInvitation(lc)
	}()
	go func(){
		purge1(lc)
	}()
}

func fastGetBadUsers(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		appendBl(pl)
	}()
	go func(){
		if antipurge == 1{
			GetSimiliarName(lc,pl)
		}
	}()
}

func detectbl(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		DeleteOtherFromChat(lc, []string{pl})
	}()
}

func dropout(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		DeleteOtherFromChat(lc, []string{pl})
	}()
	go func(){
		appendBl(pl)
	}()
}

func fastLockGroup(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		g := GetGroupWithoutMembers(lc)
		g.PreventedJoinByTicket = true
		UpdateGroup(g)
	}()
	go func(){
		appendBl(pl)
	}()
}

func fastCheckGroup(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		g := GetGroupWithoutMembers(lc)
		g.PreventedJoinByTicket = true
		UpdateGroup(g)
	}()
	go func(){
		DeleteOtherFromChat(lc,[]string{pl})
	}()
}

func fastbl(pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		appendBl(pl)
	}()
}

func kickall(lc string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		cancelall(lc)
	}()
}

func recoveradm(lc string){
	runtime.GOMAXPROCS(cpu)
	go func(){
		AcceptChatInvitation(lc)
	}()
	go func(){
		if autoPurge == 1{
			purge1(lc)
		}
	}()
}

func rmSlice(slice []string, to string) []string {
    var asu []string
    for i, t := range slice {
        if t == to {
            asu = append(slice[:i],slice[i+1:]...)
        }
    }
    return asu
}

func command(op *core1.Operation) {
	notif := op.Type
	if notif == 133{
		location := op.Param1
		enemy := op.Param2
		victim := op.Param3
		if contains(mystaff,victim) && !IsAccessAll(enemy){
			kicked[location] = append(kicked[location], victim)
			if len(kicked[location]) == len(mystaff){
				gc := GetGroup(location)
				target := gc.Members
				targets := []string{}
				cms := ""
				for i:= range target{
					targets = append(targets, target[i].Mid)
					if contains(bans, target[i].Mid){
						cms += fmt.Sprintf(" uid=%s",target[i].Mid)
					}
				}
				if contains(targets, myself) {
					fmt.Println("JOINED")
				} else {
					accept(location)
				}
				go InviteIntoChat(location, mystaff)
				go nodejs(location,mytoken,cms)
			}
		}else if victim == myself && !IsAccessAll(enemy){
			fastGetBadUsers(location,enemy)
		}else if contains(mystaff,enemy) && !IsAccessAll(victim){
			destroybl(victim)
		}else if IsAccess(victim) && !IsAccessAll(enemy){
			//backupStaff(location,enemy,victim)
			fmt.Println("kik")
		}
	}else if notif == 124{
		location := op.Param1
		enemy := op.Param2
		victim := strings.Split(op.Param3,"\x1e")
		if contains(mystaff,enemy) && contains(victim,myself){
			fmt.Println("Di Invite")
			//accept(location)
		}else if checkEqual(victim,bans) && !IsAccessAll(enemy){
			fastCancell(location,enemy,victim)
	    }else if contains(bans,enemy) && !IsAccessAll(enemy){
			fastCancell2(location,enemy,victim)
		}else if proInvite[location] == 1 && !IsAccessAll(enemy){
			fastCancell(location,enemy,victim)
		}else if IsAccess(enemy) && contains(victim,myself){
			recoveradm(location)
		}
	}else if notif == 123{
		location := op.Param1
		if len(bans) > 10{
			kickall(location)
		}
	}else if notif == 130{
		location := op.Param1
		enemy := op.Param2
		if contains(bans,enemy) && !IsAccessAll(enemy){
			detectbl(location,enemy)
			kicked[location] = rmSlice(kicked[location], enemy)
		}else if autolock == 1 && !IsAccessAll(enemy){
			fastbl(enemy)
		}else if joinLock[location] == 1 && !IsAccessAll(enemy){
			dropout(location,enemy)
		}
	}else if notif == 122{
		location := op.Param1
		enemy := op.Param2
		victim := op.Param3
		if contains(bans,enemy) && !IsAccessAll(enemy){
			fastCheckGroup(location,enemy)
		}else if autolock == 1 && !IsAccessAll(enemy){
			fastLockGroup(location,enemy)
		}else if proQr[location] == 1 && !IsAccessAll(enemy){
			lockqr(location,enemy,victim)
		}else if proName[location] == 1 && !IsAccessAll(enemy){
			lockname(location,enemy,victim)
		}
  }else if notif == 126{
		location := op.Param1
		enemy := op.Param2
		victim := op.Param3
		if contains(mystaff,victim) && !IsAccessAll(enemy){
			recovertim(location,enemy)
		}else if victim == myself && !IsAccessAll(enemy){
			fastbl(enemy)
		}else if contains(mystaff,enemy) && !IsAccessAll(victim){
			fastbl(victim)
		}else if IsAccess(victim) && !IsAccessAll(enemy){
			backupStaff(location,enemy,victim)
		}
	}else if notif == 55{
		location := op.Param1
		enemy := op.Param2
		if checkRead[location] == 1{
			if contains(readerTemp[location],enemy){
				SendTextMention(location,"oh hi, i see u @!",[]string{enemy})
				readerTemp[location] = Remove(readerTemp[location],enemy)
			}
		}
	}else if notif == 25{
		to := op.Message.To
		sender := op.Message.From_
		msgid := op.Message.ID
		if sender == myself{
			if len(chatTemp[to]) <= 100{
				chatTemp[to] = append(chatTemp[to],msgid)
			}else{
				chatTemp[to] = Remove(chatTemp[to], chatTemp[to][0])
				chatTemp[to] = append(chatTemp[to],msgid)
			}
		}
	}
}
//start
func main(){
	cleardns()
	//checkip("45.76.54.120")
	cpu = runtime.NumCPU()
	loadJson()
	profile := GetProfile()
	addAsFriendContact(myowner)
	SendText(myowner, "ajs online")
	rev := getLastOpRevision()
	myself = string(profile.Mid)
	go speedStabilizer()
	go autoUnlock()
	go startConfiguration()
	connection, err := net.Dial("tcp", port)
	if err != nil { fmt.Println(err)}
	connected = connection
	defer connected.Close()
	connected.Write([]byte("ajs_"+myself))
	fmt.Println("ban: ",bans)
	go func() {
			for{
					recv := ""
					for {
						buf := make([]byte, 1024)
						reqLen, err := connected.Read(buf)
						recv += string(buf[:reqLen])
						if err != nil || len(buf[:reqLen]) != 1024 { break }
					}
					if recv != "" {
						ClientRecv(recv)
						fmt.Println(recv)
						recv = ""
					}
			}
	}()
	for{
		fetch := fetchOperations(rev,1)
		looper := len(fetch)
		if looper > 0{
			ops := fetch[0]
			command(ops)
			revs := ops.Revision
			rev = MaxRevision(rev, revs)
		}
	}
}
