
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
	"os/exec"
	"net"
	"net/url"
	"net/http"
	"sync"
	"time"
	"runtime"
	"math/rand"
	"github.com/fogleman/gg"
	"github.com/leekchan/timeutil"
	"github.com/shirou/gopsutil/mem"
	core "./dll/LINE"
	thrift "./dll/thrift"
)
var Seq int32
var argsRaw = os.Args
var botStart = time.Now()
var Configs configs
var cpu int = 1
var timRespon int = 0
var targetNumber int = 0
var stabilizer int = 1
var duedatecount int = 0
var getStickerRespon int = 0
var getStickerRein int = 0
var getStickerBye int = 0
var getStickerUnban int = 0
var getStickerKick int = 0
var updatePicture int = 0
var updatePicture2 int = 0
var updateCover int = 0
var updateCover2 int = 0
var autoPurge int = 0
var antipurge int = 1
var limiterset int = 0
var trial bool = false
var jsonName string = ""
var myip string = "192.168.11.52"
var port string = "127.0.0.1:8000"
var rname string = ""
var sname string = ""
var proQr = make(map[string]int)
var proInvite = make(map[string]int)
var proName = make(map[string]int)
var denyTag = make(map[string]int)
var kickLock = make(map[string]int)
var joinLock = make(map[string]int)
var welcome = make(map[string]int)
var saveGname = make(map[string]string)
var chatTemp = make(map[string][]string)
var checkRead = make(map[string]int)
var promoteadmin = make(map[string]int)
var promotestaff = make(map[string]int)
var demoteadmin = make(map[string]int)
var demotestaff = make(map[string]int)
var readerTemp = make(map[string][]string)
var msgTemp = make(map[string][]string)
var groupLock = make(map[string][]string)
var commands = []string{"prefix","upimage","upcover","upname","upbio","promote","demote","addadmin","deladmin","msgleave","msgunban","msgrespon","msgwelcome","stcleave","stcunban","stcrespon","stckick","clearcontacts","clearadmin","leaveall","acceptall","addassist","delassist","addajs","delajs","reload","shutdown","upheader","setlimiter","change","out","bye","stay","stand","limitout","groups","gurl","gnuke","ginvite","joinurl","clearban","clearchat","clearstaff","cleanse","cancelall","unsend","kick","vkick","nk","contacts","status","check","banlist","managers","promote","demote","addstaff","delstaff","notag","reply","welcome","viewcontact","viewpost","ajs","mode","limiter","stabilizer","nukejoin","blockurl","blockjoin","blockgname","blockinvite","lockmember","killban","help","respon","speed","speeds","time","runtime","myuid","mygrade","info","log","set","refresh","checkram","fs","getuid","getsmule","addgif","tagall","ourl","curl","groupinfo","lurk","say","goblokall","duar anjeng"}
var bans = []string{}
var limitoutTemp = []string{}
var mystaff = []string{}
var myadmin = []string{}
var myteam = []string{}
var mycreator = []string{"ue2330fdb6b7db69eb771c3176388d0ff"}
var myassist = []string{}
var myantijs = []string{}
var targetbc = []string{}
var targetViewLimit = []string{}
var stkid string = ""
var stkpkgid string = ""
var stkid2 string = ""
var stkpkgid2 string = ""
var stkid3 string = ""
var stkpkgid3 string = ""
var stkid4 string = ""
var stkpkgid4 string = ""
var myowner string = ""
var myself string = ""
var mytoken string = ""
var MessageRespon string = ""
var helpHeader string = ""
var MessageBan string = ""
var MessageWelcome string = ""
var msgbye string = ""
var targetCache string = ""
var singletarget string = ""
var targetCommand string = ""
var limitStatus string = "Condition:\n"
var speedAll string = "Performance:\n"
var pnd bool = false
var connected net.Conn
var Running = map[string]net.Conn{}
var Ajs = map[string]net.Conn{}

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

type resultConver struct{
	Resultny struct{
		Ipny string`json:"login_ip"`
		Tokeny string`json:"token"`
	}`json:"result"`
}

type resultSimi struct{
	Kata string`json:"result"`
	Status int`json:"status"`
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
		Staylist map[string][]string`json:"Staybot"`
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
		MessageWelcome string `json:"Temp-welcome-message"`
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

func killAssist(){
	for x := range Running{
		val := Running[x]
		val.Close()
    }
}

func killAssist2(target string){
	val := Running[target]
	val.Close()
}

func broadcast(data string){
	for x := range Running{
		val := Running[x]
		val.Write([]byte(data))
    }
}
func broadcast2(target string,data string){
	val := Running[target]
	val.Write([]byte(data))
}

func killAjs(){
	for x := range Running{
		val := Running[x]
		val.Close()
    }
}

func killAjs2(target string){
	val := Running[target]
	val.Close()
}

func broadcastod(data string){
	for x := range Ajs{
		val := Ajs[x]
		val.Write([]byte(data))
    }
}
func broadcastod2(target string,data string){
	val := Ajs[target]
	val.Write([]byte(data))
}

func handleConnection(conn net.Conn) {
	bs := make([]byte, 1024)
	lenk, err := conn.Read(bs)
	if err != nil { fmt.Printf("Connection error: %s\n", err) }
	recv := string(bs[:lenk])
	if strings.HasPrefix(recv, "asis_"){
		mid := recv[5:]
		Running[mid] = conn
	}else if strings.HasPrefix(recv, "ajs_"){
		mid := recv[4:]
		Ajs[mid] = conn
	}else if strings.HasPrefix(recv, "request_"){
		stringData := recv[8:]
		joinData := strings.Split(stringData, " ")
		if joinData[0] == singletarget{
			if joinData[1] == "limit"{
				timRespon = timRespon + 1
				limitoutTemp = append(limitoutTemp,joinData[0])
				limitStatus += "\n"+ strconv.Itoa(timRespon) +". @! : limit"
				targetViewLimit = append(targetViewLimit,joinData[2])
			}else{
				timRespon = timRespon + 1
				limitoutTemp = Remove(limitoutTemp,joinData[0])
				limitStatus += "\n"+ strconv.Itoa(timRespon) +". @! : normal"
				targetViewLimit = append(targetViewLimit,joinData[2])
			}
			targetNumber = targetNumber + 1
			if targetNumber <= len(groupLock[targetCache])-1{
				broadcast2(groupLock[targetCache][targetNumber],"ceklimit_"+targetCache)
				singletarget = groupLock[targetCache][targetNumber]
			}
		}
		if len(groupLock[targetCache]) == targetNumber{
			singletarget = ""
			targetNumber = 0
			SendTextMentionByList(targetCache,limitStatus,targetViewLimit)
			timRespon = 0
			targetViewLimit = []string{}
			limitStatus = "Condition:\n"
		}
	}else if strings.HasPrefix(recv, "resultspeed_"){
		stringData := recv[12:]
		timRespon = timRespon + 1
		speedAll += "\nbot"+ strconv.Itoa(timRespon)+" : "+stringData+" ms"
		if len(myassist) == timRespon{
			SendText(targetCache, speedAll)
			timRespon = 0
			speedAll = "Performance:\n"
		}
	}
}

/* files */
func SaveFile(path string, data []byte) (bool) {
	err := ioutil.WriteFile(path, data, 0644)
	if err!=nil {
		return false
	}
	return true
}
func DeleteFile(path string) (bool) {
	err := os.Remove(path)
	if err!=nil {
		return false
	}
	return true
}
func DownloadObjectMsg(msgid string){
	client := &http.Client{}
	req , err := http.NewRequest("GET","https://obs-sg.line-apps.com/talk/m/download.nhn?oid="+msgid,nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	req.Header.Set("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	req.Header.Set("X-Line-Access", mytoken)
	res, _ := client.Do(req)
	defer res.Body.Close()
	file, err := os.Create("line-result.bin")
	io.Copy(file, res.Body)
	defer file.Close()
}

/* connection */
func ConnectChannel() *core.ChannelServiceClient {
	var err error
	var transport thrift.TTransport
	transport, err = thrift.NewTHttpClient("https://gxx.line.naver.jp/CH4")
	deBug("Login Thrift Channel Initialize", err)
	var connect *thrift.THttpClient
	connect = transport.(*thrift.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	setProtocol := thrift.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core.NewChannelServiceClientProtocol(connect, protocol, protocol)
}

func ConnectPoll() *core.TalkServiceClient {
	var err error
	var transport thrift.TTransport
	transport, err = thrift.NewTHttpClient("https://gxx.line.naver.jp/P4")
	deBug("Login Thrift Client Initialize", err)
	var connect *thrift.THttpClient
	connect = transport.(*thrift.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	connect.SetHeader("x-lal", "en_id")
	setProtocol := thrift.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core.NewTalkServiceClientProtocol(connect, protocol, protocol)
}

func GetlastOp() *core.TalkServiceClient {
	var err error
	var transport thrift.TTransport
	transport, err = thrift.NewTHttpClient("https://gxx.line.naver.jp/S4")
	deBug("Login Thrift Client Initialize", err)
	var connect *thrift.THttpClient
	connect = transport.(*thrift.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	connect.SetHeader("x-lal", "en_id")
	setProtocol := thrift.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core.NewTalkServiceClientProtocol(connect, protocol, protocol)
}

func ConnectTalk() *core.TalkServiceClient {
	var err error
	var transport thrift.TTransport
	transport, err = thrift.NewTHttpClient("https://gxx.line.naver.jp/S4")
	deBug("Login Thrift Client Initialize", err)
	var connect *thrift.THttpClient
	connect = transport.(*thrift.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.15.0 Nexus 5X 10")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.15.0\tAndroid OS\t6.0.1")
	setProtocol := thrift.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core.NewTalkServiceClientProtocol(connect, protocol, protocol)
}

func getLastOpRevision()int64{
	client := GetlastOp()
	r, e := client.GetLastOpRevision(context.TODO())
	deBug("getLastOpRevision", e)
	return r
}

func GetServerTime()int64{
	client := ConnectTalk()
	r, e := client.GetServerTime(context.TODO())
	deBug("GetServerTime", e)
	return r
}

func fetchOperations(last int64,count int32) (r []*core.Operation){
	client := ConnectPoll()
	r, e:= client.FetchOperations(context.TODO(),last,count)
	deBug("fetchOperations", e)
	return r
}

func multiOperations(last int64,count int32,global int64,individu int64)(r []*core.Operation){
	client := ConnectPoll()
	r, e := client.FetchOps(context.TODO(),last,count,global,individu)
	deBug("multiOperations", e)
	return r
}

// **generate secondary token** //
func convertTokenPrymary(){
	client := &http.Client{}
	rndvalue := []string{"pool-1","pool-2"}
	req, _ := http.NewRequest("GET", "https://api.be-team.me/primary2secondary", nil)
	req.Header.Set("apiKey", "Eo68zQtYcjXt")
	req.Header.Set("appName", "CHROMEOS\t2.3.8\tChrome OS\t1")
	req.Header.Set("server", rndvalue[rand.Intn(1)])
	req.Header.Set("sysname", "unixbot")
	req.Header.Set("authToken", mytoken)
	res, _ := client.Do(req)
	defer res.Body.Close()
	resp_body, _ := ioutil.ReadAll(res.Body)
	var responseObject resultConver
	_ = json.Unmarshal([]byte(resp_body), &responseObject)
	fmt.Println(string(resp_body))
	fmt.Println(responseObject.Resultny.Tokeny)
	//mytoken = responseObject.Resultny.Tokeny
}

func chatBot(to string,kata string){
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.be-team.me/simisimi?text="+kata+"&lang=id", nil)
	req.Header.Set("apiKey", "Eo68zQtYcjXt")
	res, err := client.Do(req)
	if err == nil{
		defer res.Body.Close()
		resp_body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resp_body))
		var responseObject resultSimi
		_ = json.Unmarshal([]byte(resp_body), &responseObject)
		if responseObject.Status == 200{
			kata := strings.Replace(responseObject.Kata, "simi", rname, -1)
			SendText(to,kata)
		}
	}
}

func GetIdForImage(toID string)(r *core.Message){
	client := ConnectTalk()
	msgObj := core.NewMessage()
	msgObj.ContentType = core.ContentType_IMAGE
	msgObj.To = toID
	r , e := client.SendMessage(context.TODO(), Seq, msgObj)
	deBug("GetIdForImage", e)
	return r
}

func SendText(toID string,msgText string){
	client := ConnectTalk()
	msgObj := core.NewMessage()
	msgObj.ContentType = core.ContentType_NONE
	msgObj.To = toID
	msgObj.Text = msgText
	_, e := client.SendMessage(context.TODO(), Seq, msgObj)
	deBug("SendText", e)
}

func SendContact(toID string, mid string) {
	client := ConnectTalk()
	msgObj := core.NewMessage()
	msgObj.ContentType = core.ContentType_CONTACT
	msgObj.To = toID
	msgObj.Text = ""
	msgObj.ContentMetadata = map[string]string{"mid":mid}
	_, e := client.SendMessage(context.TODO(), Seq, msgObj)
	deBug("SendContact", e)
}

func SendTextMention(toID string,msgText string,mids []string) {
	client := ConnectTalk()
	arr := []*tagdata{}
	mentionee := "@unix"
	texts := strings.Split(msgText, "@!")
	textx := ""
	for i := 0; i < len(mids); i++ {
		textx += texts[i]
        arr = append(arr, &tagdata{S: strconv.Itoa(len(textx)), E: strconv.Itoa(len(textx) + 5), M:mids[i]})
        textx += mentionee
	}
	textx += texts[len(texts)-1]
	allData,_ := json.MarshalIndent(arr, "", " ")
	msgObj := core.NewMessage()
	msgObj.ContentType = core.ContentType_NONE
	msgObj.To = toID
	msgObj.Text = textx
	msgObj.ContentMetadata = map[string]string{"MENTION": "{\"MENTIONEES\":"+string(allData)+"}"}
	_, e := client.SendMessage(context.TODO(), Seq, msgObj)
	deBug("SendTextMention", e)
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
				//listChar += "\n" + strconv.Itoa(listNum) + ". @!"
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

// **send message multi mentions** //
func SendTextMentionByList2(to string,msgText string,targets []string){
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
				listChar += "\n" + strconv.Itoa(listNum) + ". @!\n"
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

func DeleteOtherFromChat(groupId string, contactIds []string){
	client := ConnectTalk()
	fst := core.NewDeleteOtherFromChatRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	fst.TargetUserMids = contactIds
	_, e := client.DeleteOtherFromChat(context.TODO(), fst)
	deBug("DeleteOtherFromChat", e)
}
func InviteIntoChat(groupId string, contactIds []string){
	client := ConnectTalk()
	fst := core.NewInviteIntoChatRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	fst.TargetUserMids = contactIds
	_, e := client.InviteIntoChat(context.TODO(), fst)
	deBug("InviteIntoChat", e)
}
func CancelChatInvitation(groupId string, contactIds []string){
	client := ConnectTalk()
	fst := core.NewCancelChatInvitationRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	fst.TargetUserMids = contactIds
	_, err := client.CancelChatInvitation(context.TODO(), fst)
	deBug("CancelChatInvitation", err)
}
func AcceptChatInvitation(groupId string){
	client := ConnectTalk()
	v := core.NewAcceptChatInvitationRequest()
	v.ReqSeq = Seq
	v.ChatMid = groupId
	_, e := client.AcceptChatInvitation(context.TODO(), v)
	deBug("AcceptChatInvitation", e)
}
func GetChat(targets []string, opsiMembers bool, opsiPendings bool)(r *core.GetChatsResponse){
	client := ConnectTalk()
	v := core.NewGetChatsRequest()
	v.ChatMids = targets
	v.WithMembers = opsiMembers
	v.WithInvitees = opsiPendings
	r, e := client.GetChats(context.TODO(), v)
	deBug("GetChat", e)
	return r
}
func UpdateQrChat(groupOBJ *core.Chat){
	client := ConnectTalk()
	v := core.NewUpdateChatRequest()
	v.ReqSeq = Seq
	v.Chat = groupOBJ
	v.UpdatedAttribute = core.ChatAttribute_PREVENTED_JOIN_BY_TICKET
	_, e := client.UpdateChat(context.TODO(), v)
	deBug("UpdateChatQr", e)
}
func UpdateNameChat(groupOBJ *core.Chat){
	client := ConnectTalk()
	v := core.NewUpdateChatRequest()
	v.ReqSeq = Seq
	v.Chat = groupOBJ
	v.UpdatedAttribute = core.ChatAttribute_NAME
	_, e := client.UpdateChat(context.TODO(), v)
	deBug("UpdateChatName", e)
}
func AcceptChatInvitationByTicket(groupId string, ticketId string){
	client := ConnectTalk()
	v := core.NewAcceptChatInvitationByTicketRequest()
	v.ReqSeq = Seq
	v.ChatMid = groupId
	v.TicketId = ticketId
	_, e := client.AcceptChatInvitationByTicket(context.TODO(), v)
	deBug("AcceptChatInvitationByTicket", e)
}
func FindChatByTicket(ticketId string)(r *core.FindChatByTicketResponse){
	client := ConnectTalk()
	v := core.NewFindChatByTicketRequest()
	v.TicketId = ticketId
	r, e := client.FindChatByTicket(context.TODO(), v)
	deBug("FindChatByTicket", e)
	return r
}
func GetChatTicket(groupId string)(r *core.ReissueChatTicketResponse){
	client := ConnectTalk()
	v := core.NewReissueChatTicketRequest()
	v.ReqSeq = Seq
	v.GroupMid = groupId
	r, err := client.ReissueChatTicket(context.TODO(), v)
	deBug("GetChatTicket", err)
	return r
}

// talk v1  //
func KickoutFromGroup(groupId string, contactIds []string) {
	client := ConnectTalk()
	e := client.KickoutFromGroup(context.TODO(), Seq, groupId, contactIds)
	deBug("KickoutFromGroup", e)
}
func InviteIntoGroup(groupId string, contactIds []string) {
	client := ConnectTalk()
	e := client.InviteIntoGroup(context.TODO(), Seq, groupId, contactIds)
	deBug("InviteIntoGroup", e)
}
func CancelInvite(groupId string, contactIds []string) {
	client := ConnectTalk()
	e := client.CancelGroupInvitation(context.TODO(), Seq, groupId, contactIds)
	deBug("CancelInvite", e)
}
func AcceptGroup(groupId string) {
	client := ConnectTalk()
	e := client.AcceptGroupInvitation(context.TODO(), Seq, groupId)
	deBug("AcceptGroup", e)
}
func UpdateGroup(groupOBJ *core.Group) {
	client := ConnectTalk()
	e := client.UpdateGroup(context.TODO(), Seq, groupOBJ)
	deBug("UpdateGroup", e)
}
func GetGroup(groupId string)(r *core.Group){
	client := ConnectTalk()
	r, e := client.GetGroup(context.TODO(), groupId)
	deBug("GetGroup", e)
	return r
}
func GetContact(id string) (r *core.Contact){
	client := ConnectTalk()
	r, e := client.GetContact(context.TODO(), id)
	deBug("GetContact", e)
	return r
}
func RemoveContact(id string){
	client := ConnectTalk()
	e := client.UpdateContactSetting(context.TODO(), Seq, id, 16, "true")
	deBug("RemoveContact", e)
}
func GetProfile()*core.Profile{
	client := ConnectTalk()
	r, e := client.GetProfile(context.TODO())
	deBug("GetProfile", e)
	return r
}
func LeaveGroup(groupId string){
	client := ConnectTalk()
	e := client.LeaveGroup(context.TODO(), Seq, groupId)
	deBug("LeaveGroup", e)
}
func AcceptGroupByTicket(groupMid string, ticketId string){
	client := ConnectTalk()
	e := client.AcceptGroupInvitationByTicket(context.TODO(), Seq, groupMid, ticketId)
	deBug("AcceptGroupByTicket", e)
}
func FindGroupByTicket(ticketId string)(r *core.Group){
	client := ConnectTalk()
	r, e := client.FindGroupByTicket(context.TODO(), ticketId)
	deBug("FindGroupByTicket", e)
	return r
}
func GetGroupTicket(groupMid string)(r string){
	client := ConnectTalk()
	r, e := client.ReissueGroupTicket(context.TODO(), groupMid)
	deBug("GetGroupTicket", e)
	return r
}
func UnsendMessage(messageId string){
	client := ConnectTalk()
	e := client.UnsendMessage(context.TODO(), Seq, messageId)
	deBug("UnsendMessage", e)
}

// ........................................... //
func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
func fansign(to string, name string, typ string){
    im, err := gg.LoadImage("dll/images/fs"+typ+".bin")
    if err != nil {fmt.Println(err)}
		size := im.Bounds().Size()
    dc := gg.NewContext(size.X, size.Y)//lebar(x),tinggi(y)
		dc.SetRGB(1, 1, 1)
	  dc.Clear()
	  dc.SetRGB(0, 0, 0)
		if typ == "1"{
			if err := dc.LoadFontFace("dll/font.ttf", 80); err != nil {fmt.Println(err)}
			dc.SetHexColor("#7E4031")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
			dc.DrawStringAnchored(text, 220, 220,   0,0)
			dc.DrawStringAnchored(text, 440, 200,   0,0)
		}else if typ == "2"{
			if err := dc.LoadFontFace("dll/font.ttf", 80); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
			cekcontent := strings.Split((text)," ")
			if len(cekcontent) == 1{
				dc.DrawStringAnchored(text, 100, 390,   0,0)
			}else if len(cekcontent) == 2{
				dc.DrawStringAnchored(text, 100, 390,   0,0)
			}else if len(cekcontent) == 3{
				dc.DrawStringAnchored(cekcontent[0]+" "+cekcontent[1], 100, 390,   0,0)
				dc.DrawStringAnchored(cekcontent[2], 105, 440,   0,0)
			}else{
				dc.DrawStringAnchored(cekcontent[0]+" "+cekcontent[1], 100, 390,   0,0)
				dc.DrawStringAnchored(cekcontent[2]+" "+cekcontent[3], 105, 440,   0,0)
			}
		}else if typ == "3"{
			if err := dc.LoadFontFace("dll/font.ttf", 60); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(8))
		  dc.Stroke()
			cekcontent := strings.Split((text)," ")
			if len(cekcontent) == 1{
				dc.DrawStringAnchored(text, 100, 280,   0,0)
			}else if len(cekcontent) == 2{
				dc.DrawStringAnchored(text, 100, 280,   0,0)
			}else if len(cekcontent) == 3{
				dc.DrawStringAnchored(cekcontent[0]+" "+cekcontent[1], 100, 280,   0,0)
				dc.DrawStringAnchored(cekcontent[2], 115, 320,   0,0)
			}else{
				dc.DrawStringAnchored(cekcontent[0]+" "+cekcontent[1], 100, 280,   0,0)
				dc.DrawStringAnchored(cekcontent[2]+" "+cekcontent[3], 115, 320,   0,0)
			}
		}else if typ == "4"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
			cekcontent := strings.Split((text)," ")
			if len(cekcontent) == 1{
				dc.DrawStringAnchored(text, 80, 690, 0,0)
			}else if len(cekcontent) == 2{
				dc.DrawStringAnchored(text, 80, 690,  0,0)
			}else if len(cekcontent) == 3{
				dc.DrawStringAnchored(cekcontent[0]+" "+cekcontent[1], 80, 690,   0,0)
				dc.DrawStringAnchored(cekcontent[2], 80, 770,   0,0)
			}else{
				dc.DrawStringAnchored(cekcontent[0]+" "+cekcontent[1], 80, 690,   0,0)
				dc.DrawStringAnchored(cekcontent[2]+" "+cekcontent[3], 80, 770,   0,0)
			}
		}else if typ == "5"{
			if err := dc.LoadFontFace("dll/font.ttf", 45); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
			cekcontent := strings.Split((text)," ")
			if len(cekcontent) == 1{
				dc.DrawStringAnchored(text, 160, 270,   0,0)
			}else{
				dc.DrawStringAnchored(cekcontent[0], 160, 270,   0,0)
				dc.DrawStringAnchored(cekcontent[1], 160, 295,   0,0)
			}
		}else if typ == "6"{
			if err := dc.LoadFontFace("dll/font.ttf", 80); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 185, 539,   0,0)
		}else if typ == "7"{
			if err := dc.LoadFontFace("dll/font.ttf", 90); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 165, 440,   0,0)
		}else if typ == "8"{
			if err := dc.LoadFontFace("dll/font.ttf", 160); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(-40))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, -480, 800,   0,0)
		}else if typ == "9"{
			if err := dc.LoadFontFace("dll/font.ttf", 110); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 190, 720,   0,0)
		}else if typ == "10"{
			if err := dc.LoadFontFace("dll/font.ttf", 80); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
			cekcontent := strings.Split((text)," ")
			if len(cekcontent) == 1{
				dc.DrawStringAnchored(text, 290, 430,   0,0)
			}else if len(cekcontent) == 2{
				dc.DrawStringAnchored(text, 290, 430,   0,0)
			}else if len(cekcontent) == 3{
				dc.DrawStringAnchored(cekcontent[0]+" "+cekcontent[1], 290, 430,   0,0)
				dc.DrawStringAnchored(cekcontent[2], 290, 500,   0,0)
			}else{
				dc.DrawStringAnchored(cekcontent[0]+" "+cekcontent[1], 290, 430,   0,0)
				dc.DrawStringAnchored(cekcontent[2]+" "+cekcontent[3], 290, 500,   0,0)
			}
		}else if typ == "11"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(7))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 210, 130,   0,0)
		}else if typ == "12"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(28))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 950, 540,   0,0)
		}else if typ == "13"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
			cekcontent := strings.Split((text)," ")
			if len(cekcontent) == 1{
				dc.DrawStringAnchored(text, 380, 360,   0,0)
			}else{
				dc.DrawStringAnchored(cekcontent[0], 380, 360,   0,0)
				dc.DrawStringAnchored(cekcontent[1], 660, 360,   0,0)
			}
		}else if typ == "14"{
			if err := dc.LoadFontFace("dll/font.ttf", 120); err != nil {fmt.Println(err)}
			dc.SetHexColor("#7E4031")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 350, 1030,   0,0)
		}else if typ == "15"{
			if err := dc.LoadFontFace("dll/font.ttf", 80); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 30, 250,   0,0)
		}else if typ == "16"{
			if err := dc.LoadFontFace("dll/font.ttf", 140); err != nil {fmt.Println(err)}
			dc.SetHexColor("#503335")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 330, 1230,   0,0)
		}else if typ == "17"{
			if err := dc.LoadFontFace("dll/font.ttf", 120); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 120, 540,   0,0)
		}else if typ == "18"{
			if err := dc.LoadFontFace("dll/font.ttf", 80); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 430, 370,   0,0)
		}else if typ == "19"{
			if err := dc.LoadFontFace("dll/font.ttf", 60); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 250, 260,   0,0)
		}else if typ == "20"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 100, 360,   0,0)
		}else if typ == "21"{
			if err := dc.LoadFontFace("dll/font.ttf", 90); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 80, 460,   0,0)
		}else if typ == "22"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 320, 480,   0,0)
		}else if typ == "23"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#503335")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(29))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 720, 650,   0,0)
		}else if typ == "24"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 150, 820,   0,0)
		}else if typ == "25"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#7E4031")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(-25))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, -90, 1430,   0,0)
		}else if typ == "26"{
			if err := dc.LoadFontFace("dll/font.ttf", 120); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 910, 360,   0,0)
		}else if typ == "27"{
			if err := dc.LoadFontFace("dll/font.ttf", 60); err != nil {fmt.Println(err)}
			dc.SetHexColor("#7E4031")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 340, 680,   0,0)
		}else if typ == "28"{
			if err := dc.LoadFontFace("dll/font.ttf", 80); err != nil {fmt.Println(err)}
			dc.SetHexColor("#7E4031")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 220, 880,   0,0)
			dc.DrawStringAnchored(text, 690, 880,   0,0)
		}else if typ == "29"{
			if err := dc.LoadFontFace("dll/font.ttf", 130); err != nil {fmt.Println(err)}
			dc.SetHexColor("#7E4031")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 290, 1260,   0,0)
		}else if typ == "30"{
			if err := dc.LoadFontFace("dll/font.ttf", 70); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(-5))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 180, 430,   0,0)
			dc.DrawStringAnchored(text, 180, 430,   0,0)
		}else if typ == "31"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 210, 1200,   0,0)
		}else if typ == "32"{
			if err := dc.LoadFontFace("dll/font.ttf", 60); err != nil {fmt.Println(err)}
			dc.SetHexColor("#FF1493")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(-10))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 270, 440,   0,0)
		}else if typ == "33"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(8))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 260, 230,   0,0)
		}else if typ == "34"{
			if err := dc.LoadFontFace("dll/font.ttf", 80); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(30))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 800, 580,   0,0)
		}else if typ == "35"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 410, 1000,   0,0)
		}else if typ == "36"{
			if err := dc.LoadFontFace("dll/font.ttf", 100); err != nil {fmt.Println(err)}
			dc.SetHexColor("#503335")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 410, 825,   0,0)
		}else if typ == "37"{
			if err := dc.LoadFontFace("dll/font.ttf", 60); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(-12))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 230, 460,   0,0)
		}else if typ == "38"{
			if err := dc.LoadFontFace("dll/font.ttf", 130); err != nil {fmt.Println(err)}
			dc.SetHexColor("#7E4031")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(8))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 650, 270,   0,0)
		}else if typ == "39"{
			if err := dc.LoadFontFace("dll/font.ttf", 60); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(30))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 280, 200,   0,0)
		}else if typ == "40"{
			if err := dc.LoadFontFace("dll/font.ttf", 30); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(-32))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 140, 340,   0,0)
		}else if typ == "41"{
			if err := dc.LoadFontFace("dll/font.ttf", 40); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(-10))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 207, 370,   0,0)
		}else if typ == "42"{
			if err := dc.LoadFontFace("dll/font.ttf", 60); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(54))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 490, 60,   0,0)
		}else if typ == "43"{
			if err := dc.LoadFontFace("dll/font.ttf", 140); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(-22))
		  dc.Stroke()
			cekcontent := strings.Split((text)," ")
			if len(cekcontent) == 1{
				dc.DrawStringAnchored(text, -50, 780,   0,0)
			}else if len(cekcontent) == 2{
				dc.DrawStringAnchored(text, -50, 780,   0,0)
			}else if len(cekcontent) == 3{
				dc.DrawStringAnchored(cekcontent[0]+" "+cekcontent[1], -50, 780,   0,0)
				dc.DrawStringAnchored(cekcontent[2], -50, 920,   0,0)
			}else{
				dc.DrawStringAnchored(cekcontent[0]+" "+cekcontent[1], -50, 780,   0,0)
				dc.DrawStringAnchored(cekcontent[2]+" "+cekcontent[3], -50, 920,   0,0)
			}
		}else if typ == "44"{
			if err := dc.LoadFontFace("dll/font.ttf", 40); err != nil {fmt.Println(err)}
			dc.SetHexColor("#7E4031")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(0))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 460, 250,   0,0)
		}else if typ == "45"{
			if err := dc.LoadFontFace("dll/font.ttf", 40); err != nil {fmt.Println(err)}
			dc.SetHexColor("#000000")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(-50))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, -90, 660,   0,0)
		}else if typ == "46"{
			if err := dc.LoadFontFace("dll/font.ttf", 70); err != nil {fmt.Println(err)}
			dc.SetHexColor("#7E4031")//font color (hexcode)
	    dc.DrawImage(im, 0, 0)
			text := name
		  dc.MeasureString(text)
		  dc.Rotate(gg.Radians(8))
		  dc.Stroke()
	    dc.DrawStringAnchored(text, 270, 560,   0,0)
		}else{return}
    dc.SavePNG("out.png")
		cntn := GetIdForImage(to).ID
		callProfile(cntn,"sendimage")
}
func autoUnlock(){
	for{
			<-time.After(6 * time.Second)
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
	loc, _ := time.LoadLocation("Asia/Makassar")
	//set date hari ini(tgl sekarang(jgn tgl besok))
	batas := time.Date(2021, 1, 30, 0, 0, 0, 0, loc)
	//set lama sewa bot(jumlah hari)
	timeup := 365 //day(hari)
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
	loc, _ := time.LoadLocation("Asia/Makassar")
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

func GetRecentMessages(messageBoxId string)(r []*core.Message){
	client := ConnectTalk()
	r, err := client.GetRecentMessages(context.TODO(),messageBoxId,1000)
	deBug("GetRecentMessages", err)
	return r
}

func GetAccountSettings()(r *core.Settings){
	client := ConnectTalk()
	r, err := client.GetSettings(context.TODO())
	deBug("GetAccountSettings", err)
	return r
}

func GetCompactGroup(groupId string) (r *core.Group) {
	client := ConnectTalk()
	r, err := client.GetCompactGroup(context.TODO(), groupId)
	deBug("GetCompactGroup", err)
	return r
}

func GetGroupWithoutMembers(groupId string)(r *core.Group){
	client := ConnectTalk()
	r, err := client.GetGroupWithoutMembers(context.TODO(), groupId)
	deBug("GetGroupWithoutMembers", err)
	return r
}

func GetAllFriends() (r []string) {
	client := ConnectTalk()
	r, err := client.GetAllContactIds(context.TODO())
	deBug("GetAllFriends", err)
	return r
}

func loadJson(){
	jsonName = fmt.Sprintf("db/%s.json", argsRaw[1])
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
	mytoken = Configs.Authoken
	myassist = Configs.Status.Assist
	myantijs = Configs.Status.Antijs
	groupLock = Configs.Status.Staylist
	MessageWelcome = Configs.Settings.MessageWelcome
}

func GetGroupsInvited()(r []string) {
	client := ConnectTalk()
	r, err := client.GetGroupIdsInvited(context.TODO())
	deBug("GetGroupsInvited", err)
	return r
}

func GetGroupsJoined()(r []string) {
	client := ConnectTalk()
	r, err := client.GetGroupIdsJoined(context.TODO())
	deBug("GetGroupsJoined", err)
	return r
}

func CreateGroup(name string, contactIds []string) (r *core.Group) {
	client := ConnectTalk()
	r, err := client.CreateGroup(context.TODO(), Seq, name, contactIds)
	deBug("CreateGroup", err)
	return r
}

func RemoveAllMessage(lastMessageId string){
	client := ConnectTalk()
	err := client.RemoveAllMessages(context.TODO(), Seq, lastMessageId)
	deBug("RemoveAllMessage", err)
}

func AddContactByMid(mid string) (r map[string]*core.Contact) {
	client := ConnectTalk()
	r, err := client.FindAndAddContactsByMid(context.TODO(), Seq, mid)
	deBug("AddContactByMid", err)
	return r
}

func UpdateProfile(profile *core.Profile){
	client := ConnectTalk()
	err := client.UpdateProfile(context.TODO(), Seq, profile)
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

func MyContactTicket() (r *core.Ticket) {
	client := ConnectTalk()
	r, err := client.GetUserTicket(context.TODO())
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

func CheckEmots(data *core.Operation){
	sender := data.Message.From_
	to := data.Message.To
	emots := emots{}
	json.Unmarshal([]byte(data.Message.ContentMetadata["REPLACE"]), &emots)
	for _, stiker := range emots.STICON.RESOURCES{
		if myowner == sender||contains(myadmin,sender)||contains(mycreator,sender){
			if getStickerRespon == 1{
				if myowner == string(sender){
					stkid = stiker.PRODUCTID
					stkpkgid = stiker.STICONID
					saveJson()
					getStickerRespon = 0
					SendText(to, "respon by sticker updated")
				}
			}else if getStickerKick == 1{
				 if myowner == string(sender) || contains(mycreator,sender){
						stkid2 = stiker.PRODUCTID
						stkpkgid2 = stiker.STICONID
						saveJson()
						getStickerKick = 0
						SendText(to, "kick by sticker updated")
				 }
			}else if getStickerBye == 1{
				 if myowner == string(sender){
						stkid3 = stiker.PRODUCTID
						stkpkgid3 = stiker.STICONID
						getStickerBye = 0
						SendText(to, "bye by sticker updated")
				 }
			}else if getStickerUnban == 1{
				 if myowner == string(sender){
						stkid4 = stiker.PRODUCTID
						stkpkgid4 = stiker.STICONID
						getStickerUnban = 0
						SendText(to, "unban by sticker updated")
				 }
			}else if stiker.PRODUCTID == stkid && stiker.STICONID == stkpkgid{
				_, found := groupLock[to]
				if found == false{SendText(to, MessageRespon)
				}else{
					if len(groupLock[to]) > 0{
						SendText(to, MessageRespon)
						targetNumber = 0
						targetCommand = "respon"
						singletarget = string(groupLock[to][0])
						broadcast2(groupLock[to][0],"respon_"+to)
					}else{SendText(to, MessageRespon)}
				}
			}else if stiker.PRODUCTID == stkid2 && stiker.STICONID == stkpkgid2{
				_, found := data.Message.ContentMetadata["message_relation_server_message_id"]
				if found == true{
					for i := range msgTemp{
						if contains(msgTemp[i],data.Message.ContentMetadata["message_relation_server_message_id"]){
							if !fullAccess(i){
								DeleteOtherFromChat(to,[]string{i})
								break
							}
						}
					}
				}
			}else if stiker.PRODUCTID == stkid3 && stiker.STICONID == stkpkgid3{
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
			}else if stiker.PRODUCTID == stkid4 && stiker.STICONID == stkpkgid4{
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
		break
	}
}

func cekTypeToken(){
	r := len([]rune(mytoken))
	if r != 184{
		if r == 92{//prymary
			fmt.Println("\033[33mPrymary authoken is not permitted to login\033[39m")
			os.Exit(1)
		}else if r == 54{//secondary
			fmt.Println("\033[33mSecondary authoken is not permitted to login\033[39m")
			os.Exit(1)
		}else{//ngasal
			fmt.Println("\033[33minvalid authoken\033[39m")
			os.Exit(1)
		}
	}
}

func startAssist(){
	for i:= range myassist{
		err := exec.Command("bash", "-c", "cd dll&&./assist "+argsRaw[1]+" "+strconv.Itoa(i)+" "+port).Start()
		if err != nil { fmt.Println(err) }
	}
	getReady := GetGroupsJoined()
	cloudtim := []string{}
	cloudget := []string{}
	for cu := range myassist{
		cloudtim = append(cloudtim,myassist[cu][:33])
		myteam = append(myteam,myassist[cu][:33])
	}
	for in := range getReady{
		gtarget := getReady[in]
		checkin := GetGroup(gtarget).Members
		for ui := range checkin{
			cloudget = append(cloudget,checkin[ui].Mid)
		}
		_, found := groupLock[gtarget]
		if found == false{
		for ix := range checkin{
			if contains(cloudtim,checkin[ix].Mid){
				if uncontains(groupLock[gtarget],checkin[ix].Mid){
					groupLock[gtarget] = append(groupLock[gtarget],checkin[ix].Mid)
				}
			}
		}
	}else{
		if len(groupLock[gtarget]) == 0{
			delete(groupLock,gtarget)
		}else{
			for ixx := range groupLock[gtarget]{
				if uncontains(cloudget,groupLock[gtarget][ixx]){
					groupLock[gtarget] = Remove(groupLock[gtarget],groupLock[gtarget][ixx])
				}
			}
		}
		cloudget = []string{}
	}
	}
	saveJson()
	for i := range myassist{
		target := string(myassist[i][:33])
		fmt.Println(target)
		if target != myself{
			if nothingInMyContacts(target){
				time.Sleep(1 * time.Second)
				fmt.Println(target)
				AddContactByMid(target)
			}
		}
	}
}

func startAjs(){
	for i:= range myantijs{
		err := exec.Command("bash", "-c", "cd dll&&./anjes "+argsRaw[1]+" "+strconv.Itoa(i)+" "+port).Start()
		if err != nil { fmt.Println(err) }
	}
	getReady := GetGroupsJoined()
	cloudtim := []string{}
	cloudget := []string{}
	for cu := range myantijs{
		cloudtim = append(cloudtim,myantijs[cu][:33])
		myteam = append(myteam,myantijs[cu][:33])
	}
	for in := range getReady{
		gtarget := getReady[in]
		checkin := GetGroup(gtarget).Members
		for ui := range checkin{
			cloudget = append(cloudget,checkin[ui].Mid)
		}
		_, found := groupLock[gtarget]
		if found == false{
		for ix := range checkin{
			if contains(cloudtim,checkin[ix].Mid){
				if uncontains(groupLock[gtarget],checkin[ix].Mid){
					groupLock[gtarget] = append(groupLock[gtarget],checkin[ix].Mid)
				}
			}
		}
	}else{
		if len(groupLock[gtarget]) == 0{
			delete(groupLock,gtarget)
		}else{
			for ixx := range groupLock[gtarget]{
				if uncontains(cloudget,groupLock[gtarget][ixx]){
					groupLock[gtarget] = Remove(groupLock[gtarget],groupLock[gtarget][ixx])
				}
			}
		}
		cloudget = []string{}
	}
	}
	saveJson()
	for i := range myantijs{
		target := string(myantijs[i][:33])
		fmt.Println(target)
		if target != myself{
			if nothingInMyContacts(target){
				time.Sleep(1 * time.Second)
				fmt.Println(target)
				AddContactByMid(target)
			}
		}
	}
}

func nodejs(gid string,authjs string,cms string){
	cmo:=fmt.Sprintf("node dll/kick.js gid=%s %s token=%s",gid,cms,authjs)
	parts := strings.Fields(cmo)
	cmd, _ := exec.Command(parts[0],parts[1:]...).Output()
	fmt.Println(string(cmd))
}
func nodejs2(gid string,authjs string,cms1 string,cms2 string){
	cmo:=fmt.Sprintf("node dll/double.js gid=%s %s %s token=%s",gid,cms1,cms2,mytoken)
	parts := strings.Fields(cmo)
	cmd, _ := exec.Command(parts[0],parts[1:]...).Output()
	fmt.Println(string(cmd))
}
func nodejs3(gid string,authjs string,cms string){
	cmo:=fmt.Sprintf("node dll/cancel.js gid=%s %s token=%s",gid,cms,mytoken)
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
	//port = ":"+localip[1]
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
	Configs.Settings.MessageWelcome = MessageWelcome
	Configs.Settings.ByeSticker.Stkid = stkid3
	Configs.Settings.ByeSticker.Stkpkgid = stkpkgid3
	Configs.Settings.UnbanSticker.Stkid = stkid4
	Configs.Settings.UnbanSticker.Stkpkgid = stkpkgid4
	Configs.Settings.Stabilizer = stabilizer
	Configs.Status.Assist = myassist
	Configs.Status.Antijs = myantijs
	Configs.Status.Staylist = groupLock
	encjson, _ := json.MarshalIndent(Configs, "", "  ")
	ioutil.WriteFile(jsonName, encjson, 0644)
}

func GetSimiliarName(loc string,target string){
	ts := GetContact(target)
	myString := ts.DisplayName
	a := []rune(myString)
	myShortString := string(a[0:3])
	gc := GetGroup(loc)
	targets := gc.Members
	for i:= range targets{
		cach := []rune(targets[i].DisplayName)
		if string(cach[0:3]) == myShortString{
			if !fullAccess(targets[i].Mid){
				if myself != targets[i].Mid{
					appendBl(targets[i].Mid)
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

func poolKickBanWhenAccept(lc string){
	runtime.GOMAXPROCS(cpu)
	time.Sleep(1000 * time.Microsecond)
  var wg sync.WaitGroup
  wg.Add(len(bans))
	for i:=0;i<len(bans);i++ {
		go func(i int) {
			defer wg.Done()
			DeleteOtherFromChat(lc,[]string{bans[i]})
    }(i)
	}
	wg.Wait()
}

func poolKickBans(lc string){
	runtime.GOMAXPROCS(cpu)
  var wg sync.WaitGroup
  wg.Add(len(bans))
	for i:=0;i<len(bans);i++ {
		go func(i int) {
			defer wg.Done()
			DeleteOtherFromChat(lc,[]string{bans[i]})
    }(i)
	}
	wg.Wait()
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

// **cek target all access** //
func fullAccess(target string) bool{
	data := []string{myowner}
	data = append(data,myteam...)
	data = append(data,mycreator...)
	data = append(data,myadmin...)
	data = append(data,mystaff...)
	looper := len(data)
	for i:=0;i<looper;i++{
		if target == data[i]{
			 return true
		}
	}
	return false
}

// **cek target some access** //
func access(target string) bool{
	data := []string{myowner}
	data = append(data,mycreator...)
	data = append(data,myadmin...)
	data = append(data,mystaff...)
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

func poolCancel(lc string,pd []string){
	runtime.GOMAXPROCS(cpu)
  var wg sync.WaitGroup
  wg.Add(len(pd))
	for i:=0;i<len(pd);i++ {
		go func(i int) {
			defer wg.Done()
			CancelChatInvitation(lc,[]string{pd[i]})
    }(i)
	}
	wg.Wait()
}

func cleardns() {
	cmd, _ := exec.Command("bash","-c","sudo systemd-resolve --flush-caches").Output()
	fmt.Println("\033[33m"+string(cmd)+"\033[39m")
}
func installer() {
	cmd, _ := exec.Command("bash","-c","sudo systemd-resolve --flush-caches").Output()
	exec.Command("bash","-c","pip3 uninstall thrift&&y").Output()
	fmt.Println("\033[33m"+string(cmd)+"\033[39m")
}

func callProfile(cmd1 string, cmd2 string) {
	cmd, _ := exec.Command("python3","dll/lineProfile.py",mytoken,myself,cmd1,cmd2).Output()
	fmt.Println("\033[33m"+string(cmd)+"\033[39m")
}

func sendBigImage(grp string, imglink string) {
	cmd, _ := exec.Command("python3","dll/liff.py",mytoken,grp,imglink).Output()
	fmt.Println("\033[33m"+string(cmd)+"\033[39m")
}

//clear
func backup(lc string,pl string,kr string){
	runtime.GOMAXPROCS(cpu)
	cek := GetChat([]string{lc}, true, false)
	midMembers := cek.Chat[0].Extra.GroupExtra.MemberMids
	_, foundMe := midMembers[myself]
	_, foundEnemy := midMembers[pl]
	_, foundVictim := midMembers[kr]
	if foundMe == true{
		if foundEnemy == true{
			if foundVictim == false{
				go func(){DeleteOtherFromChat(lc,[]string{pl})}()
				go func(){InviteIntoChat(lc, []string{kr})}()
				go func(){appendBl(pl)}()
			}
		}
	}
}

//clear
func acceptManagers(lc string){
	runtime.GOMAXPROCS(cpu)
	go func(){AcceptChatInvitation(lc)}()
	go func(){if autoPurge == 1{poolKickBanWhenAccept(lc)}}()
}
//clear
func cancelAllEnemy(lc string,pl string,pd []string){
	runtime.GOMAXPROCS(cpu)
	go func(){poolCancel(lc,pd)}()
	go func(){if !fullAccess(pl){DeleteOtherFromChat(lc, []string{pl})}}()
	go func(){for i:=range pd{appendBl(pd[i])}}()
	go func(){if !fullAccess(pl){appendBl(pl)}}()
}
//clear
func scanNameTarget(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){if antipurge == 1{GetSimiliarName(lc,pl)}}()
	go func(){appendBl(pl)}()
}
//clear
func proJoin(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){DeleteOtherFromChat(lc, []string{pl})}()
	go func(){appendBl(pl)}()
}
//clear
func proQrGroup(lc string,pl string,kr string){
	runtime.GOMAXPROCS(cpu)
	if kr == "4"{
		go func(){DeleteOtherFromChat(lc,[]string{pl})}()
		go func(){g := GetGroupWithoutMembers(lc);if g.PreventedJoinByTicket == false{g.PreventedJoinByTicket = true;UpdateGroup(g)}}()
		go func(){appendBl(pl)}()
	}
}

func closeQrAndKick(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){DeleteOtherFromChat(lc,[]string{pl})}()
	go func(){g := GetGroupWithoutMembers(lc);if g.PreventedJoinByTicket == false{g.PreventedJoinByTicket = true;UpdateGroup(g)}}()
}

func proNameGroup(lc string,pl string,kr string){
	runtime.GOMAXPROCS(cpu)
	if kr == "1"{
		go func(){DeleteOtherFromChat(lc,[]string{pl})}()
		go func(){g := GetGroupWithoutMembers(lc);if g.Name != saveGname[lc]{g.Name = saveGname[lc];UpdateGroup(g)}}()
		go func(){appendBl(pl)}()
	}
}

func pend(to string){
	if pnd == false{
		pnd = true
		proQr[to] = 1
		time.Sleep(10 * time.Second)
		delete(proQr,to)
	}
}

func command(op *core.Operation) {
	if op.Type == 133{
		location := op.Param1
		enemy := op.Param2
		victim := op.Param3
		if contains(myteam,victim) && !fullAccess(enemy){
			go pend(location)
			backup(location,enemy,victim)
		}else if victim == myself && !fullAccess(enemy){
			go pend(location)
			scanNameTarget(location,enemy)
		}else if contains(myteam,enemy) && !fullAccess(victim){
			appendBl(victim)
		}else if access(victim) && !fullAccess(enemy){
			go pend(location)
		  	backup(location,enemy,victim)
		}
	}else if op.Type == 124{
		location := op.Param1
		enemy := op.Param2
		victim := strings.Split(op.Param3,"\x1e")
		if contains(myteam,enemy) && contains(victim,myself){
			AcceptChatInvitation(location)
		}else if checkEqual(victim,bans){
			cancelAllEnemy(location,enemy,victim)
	  }else if contains(bans,enemy){
			cancelAllEnemy(location,enemy,victim)
		}else if proInvite[location] == 1 && !fullAccess(enemy){
			cancelAllEnemy(location,enemy,victim)
		}else if access(enemy) && contains(victim,myself){
			acceptManagers(location)
		}
	}else if op.Type == 130{
		location := op.Param1
		enemy := op.Param2
		if contains(bans,enemy){
			closeQrAndKick(location,enemy)
		}else if joinLock[location] == 1 && !fullAccess(enemy){
			proJoin(location,enemy)
		}else if welcome[location] == 1 && !fullAccess(enemy){
			SendText(location,MessageWelcome)
		}
	}else if op.Type == 122{
		location := op.Param1
		enemy := op.Param2
		victim := op.Param3
		if contains(bans,enemy){
			closeQrAndKick(location,enemy)
		}else if proQr[location] == 1 && !fullAccess(enemy){
			proQrGroup(location,enemy,victim)
		}else if proName[location] == 1 && !fullAccess(enemy){
			proNameGroup(location,enemy,victim)
		}
  }else if op.Type == 126{
		location := op.Param1
		enemy := op.Param2
		victim := op.Param3
		if contains(myteam,victim) && !fullAccess(enemy){
			backup(location,enemy,victim)
		}else if victim == myself && !fullAccess(enemy){
			appendBl(enemy)
		}else if contains(myteam,enemy) && !fullAccess(victim){
			appendBl(victim)
		}else if access(victim) && !fullAccess(enemy){
			backup(location,enemy,victim)
		}
	}else if op.Type == 55{
		location := op.Param1
		enemy := op.Param2
		if checkRead[location] == 1{
			if contains(readerTemp[location],enemy){
				SendTextMention(location,"oh hi, i see u @!",[]string{enemy})
				readerTemp[location] = Remove(readerTemp[location],enemy)
			}
		}
	}else if op.Type == 25{
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
	}else if op.Type == 26{
		go CheckEmots(op)
		sender := op.Message.From_
		text := op.Message.Text
		to := op.Message.To
		msgid := op.Message.ID
		if len(msgTemp[sender]) <= 200{
			msgTemp[sender] = append(msgTemp[sender],msgid)
		}else{
			msgTemp[sender] = Remove(msgTemp[sender], msgTemp[sender][0])
			msgTemp[sender] = append(msgTemp[sender],msgid)
		}
		if (op.Message.ToType).String() == "GROUP"||(op.Message.ToType).String() == "USER"{
			if singletarget == sender{
				if targetCommand == "respon"{
					targetNumber = targetNumber + 1
					if targetNumber <= len(groupLock[to])-1{
						broadcast2(groupLock[to][targetNumber],"respon_"+to)
						singletarget = groupLock[to][targetNumber]
						if len(groupLock[to])-1 == targetNumber{
							singletarget = ""
							targetNumber = 0
							targetCommand = ""
							singletarget = ""
						}
					}
				}
			}
			if contains(myadmin,sender)||sender == myowner||contains(mycreator,sender)||contains(mystaff,sender){
				if promoteadmin[to] == 1{
					if uncontains(mystaff,sender){
						if strings.ToLower(text) == "stop"{
							delete(promoteadmin,to)
							SendText(to, "promote contact canceled")
						}
					}
				}
				if promotestaff[to] == 1{
					if uncontains(mystaff,sender){
						if strings.ToLower(text) == "stop"{
							delete(promotestaff,to)
							SendText(to, "promote contact canceled")
						}
					}
				}
				if demoteadmin[to] == 1{
					if uncontains(mystaff,sender){
						if strings.ToLower(text) == "stop"{
							delete(demoteadmin,to)
							SendText(to, "demote contact canceled")
						}
					}
				}
				if demotestaff[to] == 1{
					if uncontains(mystaff,sender){
						if strings.ToLower(text) == "stop"{
							delete(demotestaff,to)
							SendText(to, "demote contact canceled")
						}
					}
				}
				if (op.Message.ContentType).String() == "NONE"{
					if strings.ToLower(text) == "name"{SendText(to, rname)}
					if strings.HasPrefix(strings.ToLower(text), rname){
						tobject := strings.Replace(strings.ToLower(text), rname+" ", "", -1)
						tobject = strings.Replace(tobject, rname, "", -1)
						looping := strings.Split(tobject,",")
						for cmd := range looping{
							if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[0]+":"){//prefix
								if uncontains(myadmin,sender) || uncontains(mystaff,sender){
									result := strings.Split((op.Message.Text),":")
									rname = result[1]
									saveJson()
									SendText(to, "ᴘʀᴇꜰɪx ᴜᴘᴅᴀᴛᴇᴅ")
								}
							}else if strings.ToLower(looping[cmd]) == commands[1]{//upimage
								if uncontains(myadmin,sender) || uncontains(mystaff,sender){
									SendText(to, "ᴘʟᴇᴀsᴇ sᴇɴᴅ ᴀɴ ɪᴍᴀɢᴇ")
									updatePicture = 1
								}
							}else if strings.ToLower(looping[cmd]) == commands[2]{//upcover
								if uncontains(myadmin,sender) || uncontains(mystaff,sender){
									SendText(to, "ᴘʟᴇᴀsᴇ sᴇɴᴅ ᴀɴ ɪᴍᴀɢᴇ")
									updateCover = 1
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[3]+":"){//upname
								if uncontains(myadmin,sender) || uncontains(mystaff,sender){
									result := strings.Split((op.Message.Text),":")
									objme := GetProfile()
									objme.DisplayName = result[1]
									UpdateProfile(objme)
									SendText(to, "ᴘʀᴏꜰɪʟᴇ ɴᴀᴍᴇ ᴜᴘᴅᴀᴛᴇᴅ")
	 							}
	 						}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[4]+":"){//upbio
	 							if uncontains(myadmin,sender) || uncontains(mystaff,sender){
	 								result := strings.Split((op.Message.Text),":")
	 								objme := GetProfile()
	 								objme.StatusMessage = result[1]
	 								UpdateProfile(objme)
	 								SendText(to, "ᴘʀᴏꜰɪʟᴇ sᴛᴀᴛᴜs ᴜᴘᴅᴀᴛᴇᴅ")
	 							}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[5]+":"){//promote
									result := strings.Split(text,":")
									if strings.ToLower(result[1]) == "admin"{
										if uncontains(myadmin,sender) || uncontains(mystaff,sender){
											promoteadmin[to] = 1
											SendText(to, "ᴘʟᴇᴀsᴇ sᴇɴᴅ ᴀ ᴄᴏɴᴛᴀᴄᴛ")
										}
									}else if strings.ToLower(result[1]) == "staff"{
										if uncontains(mystaff,sender){
											promotestaff[to] = 1
											SendText(to, "ᴘʟᴇᴀsᴇ sᴇɴᴅ ᴀ ᴄᴏɴᴛᴀᴄᴛ")
										}
									}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[6]+":"){//demote
												result := strings.Split(text,":")
												if strings.ToLower(result[1]) == "admin"{
													if myowner == string(sender) || contains(mycreator,sender){
														demoteadmin[to] = 1
														SendText(to, "ᴘʟᴇᴀsᴇ sᴇɴᴅ ᴀ ᴄᴏɴᴛᴀᴄᴛ")
													}
												}else if strings.ToLower(result[1]) == "staff"{
													if uncontains(mystaff,sender){
														demotestaff[to] = 1
														SendText(to, "ᴘʟᴇᴀsᴇ sᴇɴᴅ ᴀ ᴄᴏɴᴛᴀᴄᴛ")
													}
												}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[7]){//addadmin
															if uncontains(myadmin,sender) || contains(mystaff,sender){
																go func(){
																	mentions := mentions{}
																	json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
																	for _, mention := range mentions.MENTIONEES{
																		if uncontains(myadmin, mention.Mid){
																			if mention.Mid != myself{
																				myadmin = append(myadmin, mention.Mid)
																				if nothingInMyContacts(mention.Mid){
																					time.Sleep(1 * time.Millisecond)
																					AddContactByMid(mention.Mid)
																				}
																			}
																		}
																	}
																	saveJson()
																	broadcast("loadjson")
																	broadcastod("loadjson")
																}()
																SendText(to, "ᴀᴅᴅᴇᴅ ᴛᴏ ᴀᴅᴍɪɴʟɪsᴛ")
															}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[8]){//deladmin
															if uncontains(myadmin,sender) || contains(mystaff,sender){
																mentions := mentions{}
																json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
																for _, mention := range mentions.MENTIONEES{
																	for i := 0; i < len(myadmin); i++ {
																		if myadmin[i] == mention.Mid{
																			myadmin = Remove(myadmin,myadmin[i])
																		}
																	}
																}
																saveJson()
																broadcast("loadjson")
																broadcastod("loadjson")
																SendText(to,"ʀᴇᴍᴏᴠᴇᴅ ꜰʀᴏᴍ ᴀᴅᴍɪɴʟɪsᴛ")
															}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[9]+":"){//msgleave
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									result := strings.Split((text),":")
									msgbye = result[1]
									saveJson()
									broadcast("loadjson")
									broadcastod("loadjson")
									SendText(to, "ʟᴇᴀᴠᴇ ᴍᴇssᴀɢᴇ ᴜᴘᴅᴀᴛᴇᴅ")
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[10]+":"){//msgunban
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									result := strings.Split((text),":")
									MessageBan = result[1]
									saveJson()
									broadcast("loadjson")
									broadcastod("loadjson")
									SendText(to, "ᴜɴʙᴀɴ ᴍᴇssᴀɢᴇ ᴜᴘᴅᴀᴛᴇᴅ")
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[11]+":"){//msgrespon
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									result := strings.Split((text),":")
									MessageRespon = result[1]
									saveJson()
									broadcast("loadjson")
									broadcastod("loadjson")
									SendText(to, "ʀᴇsᴘᴏɴ ᴍᴇssᴀɢᴇ ᴜᴘᴅᴀᴛᴇᴅ")
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[12]+":"){//msgwelcome
								if uncontains(myadmin,sender) || contains(mystaff,sender){
										result := strings.Split((text),":")
										MessageWelcome = result[1]
										saveJson()
										broadcast("loadjson")
										broadcastod("loadjson")
										SendText(to, "ᴡᴇʟᴄᴏᴍᴇ ᴍᴇssᴀɢᴇ ᴜᴘᴅᴀᴛᴇᴅ")
								}
							}else if strings.ToLower(looping[cmd]) == commands[13]{//stcleave
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									getStickerBye = 1
									SendText(to, "ᴘʟᴇᴀsᴇ sᴇɴᴅ ᴀ sᴛɪᴄᴋᴇʀ")
								}
							}else if strings.ToLower(looping[cmd]) == commands[14]{//stcunban
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									getStickerUnban = 1
									SendText(to, "ᴘʟᴇᴀsᴇ sᴇɴᴅ ᴀ sᴛɪᴄᴋᴇʀ")
								}
							}else if strings.ToLower(looping[cmd]) == commands[15]{//stcrespon
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									getStickerRespon = 1
									SendText(to, "ᴘʟᴇᴀsᴇ sᴇɴᴅ ᴀ sᴛɪᴄᴋᴇʀ")
								}
							}else if strings.ToLower(looping[cmd]) == commands[16]{//stckick
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									getStickerKick = 1
									SendText(to, "ᴘʟᴇᴀsᴇ sᴇɴᴅ ᴀ sᴛɪᴄᴋᴇʀ")
								}
							}else if strings.ToLower(looping[cmd]) == commands[17]{//clearfriends
									if uncontains(myadmin,sender) || contains(mystaff,sender){
										friends := GetAllFriends()
										count := 0
										for i:= range friends{
											if !fullAccess(friends[i]){
												RemoveContact(friends[i])
												count = count + 1
											}
										}
										SendText(to, strconv.Itoa(count)+" ᴄᴏɴᴛᴀᴄᴛ ᴄʟᴇᴀʀᴇᴅ")
									}
							}else if strings.ToLower(looping[cmd]) == commands[18]{//clearadmin
										if uncontains(myadmin,sender) || contains(mystaff,sender){
											SendText(to, strconv.Itoa(len(myadmin))+" ᴀᴅᴍɪɴ ᴄʟᴇᴀʀᴇᴅ")
											myadmin = []string{}
											saveJson()
											broadcast("loadjson")
											broadcastod("loadjson")
										}
							}else if strings.ToLower(looping[cmd]) == commands[19]{//leaveall
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									groups := GetGroupsJoined()
									SendText(to, "ʟᴇᴀᴠɪɴɢ "+strconv.Itoa(len(groups))+" ɢʀᴏᴜᴘ")
									for i:= range groups{
										LeaveGroup(groups[i])
										_, found := groupLock[groups[i]]
										if found == true{
											delete(groupLock,groups[i])
										}
									}
									saveJson()
									broadcast("loadjson")
									broadcast("leaveall")
								}
							}else if strings.ToLower(looping[cmd]) == commands[20]{//acceptall
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									groups := GetGroupsInvited()
									if len(groups) != 0{
										for i:= range groups{
											AcceptChatInvitation(groups[i])
										}
										SendText(to, "ᴀᴄᴄᴇᴘᴛ "+strconv.Itoa(len(groups))+" ɢʀᴏᴜᴘs")
									}else{SendText(to, "ɴᴏ ɢʀᴏᴜᴘs")}
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[21]+":"){//addassist
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									result := strings.Split((text),":")
									result2 := result[1] +":"+ result[2]
									if uncontains(myassist,result2){
										myassist = append(myassist,result2)
										saveJson()
										if nothingInMyContacts(result[1]){
											time.Sleep(1 * time.Second)
											AddContactByMid(result[1])
										}
										for i := range myassist{
											if myassist[i] == result2{
												err := exec.Command("bash", "-c", "cd dll&&./assist "+argsRaw[1]+" "+strconv.Itoa(i)+" "+port).Start()
												if err!=nil{fmt.Println(err)}
												SendText(to, "ᴀssɪsᴛ "+strconv.Itoa(i+1)+" ʀᴜɴ sᴜᴄᴄᴇsꜰᴜʟʟʏ")
												break
											}
										}
										getReady := GetGroupsJoined()
										for in := range getReady{
											gtarget := getReady[in]
											checkin := GetGroup(getReady[in]).Members
											for ix := range checkin{
												if contains(myteam,checkin[ix].Mid){
													if uncontains(groupLock[gtarget],checkin[ix].Mid){
														groupLock[gtarget] = append(groupLock[gtarget],checkin[ix].Mid)
													}
												}
											}
										}
									}else{SendText(to, "ᴀʟʀᴇᴀᴅʏ ᴇxɪsᴛ")}
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[22]+":"){//delassist:number
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									result := strings.Split((text),":")
									if result[1] != "0"{
										result2, err := strconv.Atoi(result[1])
										if err != nil{SendText(to, "ɪɴᴠᴀʟɪᴅ ɴᴜᴍʙᴇʀ")
											}else{
												for i := range myassist{
													if result2 > 0 && result2-1 < len(myassist){
														if i == result2-1{
															broadcast2(myassist[i][:33],"exit")
															killAssist2(myassist[i][:33])
															for ii := range groupLock{
																if contains(groupLock[ii],myassist[i][:33]){
																	groupLock[ii] = Remove(groupLock[ii],myassist[i][:33])
																}
															}
															myassist = Remove(myassist,myassist[i])
															SendText(to, "ᴀssɪsᴛ "+result[1]+"ʀᴇᴍᴏᴠᴇᴅ")
															break
														}
														}else{SendText(to, "ᴏᴜᴛ ᴏꜰ ʀᴀɴɢᴇ")}
													}
												}
										}else{SendText(to, "ɪɴᴠᴀʟɪᴅ ʀᴀɴɢᴇ")}
										saveJson()
										broadcast("loadjson")
										broadcastod("loadjson")
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[23]+":"){//addajs
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									result := strings.Split((text),":")
									result2 := result[1] +":"+ result[2]
									if uncontains(myantijs,result2){
										myantijs = append(myantijs,result2)
										saveJson()
										broadcast("loadjson")
										broadcastod("loadjson")
										for i := range myantijs{
											if myantijs[i] == result2{
												err := exec.Command("bash", "-c", "cd dll&&./anjes "+argsRaw[1]+" "+strconv.Itoa(i)+" "+port).Start()
												if err != nil { fmt.Println(err) }
												SendText(to, "ᴀᴊs "+strconv.Itoa(i+1)+" ʀᴜɴ sᴜᴄᴄᴇsꜰᴜʟʟʏ")
												break
											}
										}
									}else{SendText(to, "ᴀʟʀᴇᴀᴅʏ ᴇxɪsᴛ")}
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[24]+":"){//delajs
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									result := strings.Split((text),":")
									result2, _ := strconv.Atoi(result[1])
									if len(myantijs) != 0{
										if result2 > 0 && result2 <= len(myantijs){
											broadcast2(myantijs[result2-1][:33],"exit")
											myantijs = Remove(myantijs,myantijs[result2-1])
											SendText(to,"ᴀᴊs ʀᴇᴍᴏᴠᴇᴅ")
										}else{SendText(to,"ᴏᴜᴛ ᴏꜰ ʀᴀɴɢᴇ")}
									}else{SendText(to,"ᴀᴊs ɪs ᴇᴍᴘᴛʏ")}
									saveJson()
								}
							}else if strings.ToLower(looping[cmd]) == commands[25]{//reload
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									myteam = []string{}
									broadcast("exit")
									broadcastod("exit")
									killAssist()
									killAjs()
									startAssist()
									startAjs()
									SendText(to, strconv.Itoa(len(myassist))+" ᴀssɪsᴛ ʀᴇʟᴏᴀᴅᴇᴅ")
								}
							}else if strings.ToLower(looping[cmd]) == commands[26]{//shutdown
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									if myowner == string(sender){
										broadcast("exit")
										broadcastod("exit")
										killAssist()
										killAjs()
										SendText(to, "sᴇʟꜰʙᴏᴛ ᴏᴘᴇʀᴀᴛɪᴏɴ sᴛᴏᴘᴘᴇᴅ")
										os.Exit(1)
									}
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[27]+":"){//upheader
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									result := strings.Split((text),":")
									helpHeader = result[1]
									SendText(to, "ʜᴇᴀᴅᴇʀ ᴜᴘᴅᴀᴛᴇᴅ")
									saveJson()
									broadcast("loadjson")
									broadcastod("loadjson")
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[28]+":"){//setlimiter
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									result := strings.Split((text),":")
									r,_ := strconv.Atoi(result[1])
									limiterset = r
									SendText(to, "ʟɪᴍɪᴛᴀᴛɪᴏɴ ᴄʜᴀɴɢᴇᴅ")
									broadcast("limiterset_"+result[1])
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[29]+":"){//change
								if uncontains(myadmin,sender) || contains(mystaff,sender){
									result := strings.Split((text),":")
									msgjoin := strings.Split((result[1])," ")
									numb, _ := strconv.Atoi(msgjoin[0])
									if numb > 0&&numb <= len(commands){
										numb = numb - 1
										commands[numb] = msgjoin[1]
										SendText(to, "ᴄᴏᴍᴍᴀɴᴅ ᴄʜᴀɴɢᴇᴅ")
									}else{SendText(to,"ᴏᴜᴛ ᴏꜰ ʀᴀɴɢᴇ")}
								}
							}else if strings.ToLower(looping[cmd]) == commands[30]{//out
								if uncontains(mystaff,sender){
									_, found := groupLock[to]
									if found == true{
										for i := range groupLock[to]{
											time.Sleep(60 * time.Millisecond)
											broadcast2(groupLock[to][i],"getout_"+to)
										}
										getmem := GetGroup(to)
										target := getmem.Members
										for i:= range target{
											if contains(myteam,target[i].Mid){
												time.Sleep(60 * time.Millisecond)
												broadcast2(target[i].Mid,"getout_"+to)
											}
										}
									}
									delete(groupLock,to)
									saveJson()
									broadcast("loadjson")
								}
							}else if strings.ToLower(looping[cmd]) == commands[31]{//bye
								if uncontains(mystaff,sender){
									delete(readerTemp,to)
									delete(checkRead,to)
									delete(proInvite,to)
									delete(proName,to)
									delete(proQr,to)
									delete(joinLock,to)
									delete(denyTag,to)
									delete(kickLock,to)
									delete(saveGname,to)
									SendText(to, msgbye)
									delete(groupLock,to)
									LeaveGroup(to)
									saveJson()
									broadcast("loadjson")
									broadcastod("loadjson")
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[32]+":"){//stay:number
								if uncontains(mystaff,sender){
									result := strings.Split((text),":")
									result2, _ := strconv.Atoi(result[1])
									result2 = result2 - 1
									tick := GetChatTicket(to)
									if result2 > 0&&result2 <= len(myassist){
										getmem := GetGroup(to)
										target := getmem.Members
										targets:= []string{}
										batastim := 0
										batastim2 := 0
										//proses simpan smua mid member grup
										for i:= range target{
											targets = append(targets, target[i].Mid)
										}
										_, found := groupLock[to]
										if found == false{
											//if contains(targets,myself){
											//	groupLock[to] = append(groupLock[to],myself)
											//}
											for utt := range myassist{
												if contains(targets,myassist[utt][:33]){
													groupLock[to] = append(groupLock[to],myassist[utt][:33])
												}
											}
										}
										//proses cek squad yg sudah ada di grup, jika kebanyaan kluar
										for i:= range targets{
											if contains(groupLock[to],targets[i]){
												if batastim < result2{
													batastim = batastim + 1
												}else{
													broadcast2(targets[i],"getout_"+to)
													groupLock[to] = Remove(groupLock[to],targets[i])
												}
											}
										}
										//proses memasukkan anggota yg kurang
										for io:= range myassist{
											if contains(targets,myassist[io][:33]){
												batastim2 = batastim2 + 1
											}
										}
										if batastim2 < result2{
											if getmem.PreventedJoinByTicket == true{
												getmem.PreventedJoinByTicket = false
												UpdateGroup(getmem)
											}
										}
										for i:= range myassist{
											if batastim2 < result2{
												if uncontains(targets,myassist[i][:33]){
													batastim2 = batastim2 + 1
													broadcast2(myassist[i][:33],"jointiket_"+to+" "+tick.TicketId)
													if uncontains(groupLock[to],myassist[i][:33]){
														groupLock[to] = append(groupLock[to],myassist[i][:33])
													}
												}
											}
										}
										if batastim2 == result2{
											time.Sleep(1 * time.Second)
											if getmem.PreventedJoinByTicket == false{
												getmem.PreventedJoinByTicket = true
												UpdateGroup(getmem)
											}
										}
									}else{SendText(to, "out of range")}
									saveJson()
									broadcast("loadjson")
									broadcastod("loadjson")
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[33]+":"){//stand:number
								if uncontains(mystaff,sender){
									result := strings.Split((text),":")
									result2, _ := strconv.Atoi(result[1])
									result2 = result2 - 1
									if result2 > 0&&result2 <= len(myassist){
										getmem := GetGroup(to)
										target := getmem.Members
										targets := []string{}
										tempInv := []string{}
										batastim := 0
										batastim2 := 0
										//proses simpan smua mid member grup
										for i:= range target{
											targets = append(targets, target[i].Mid)
										}
										_, found := groupLock[to]
										if found == false{
											for utt := range myassist{
												if contains(targets,myassist[utt][:33]){
													groupLock[to] = append(groupLock[to],myassist[utt][:33])
												}
											}
										}
										//proses cek squad yg sudah ada di grup, jika kebanyaan kluar
										for i:= range targets{
											if contains(groupLock[to],targets[i]){
												if batastim < result2{
													batastim = batastim + 1
												}else{
													broadcast2(targets[i],"getout_"+to)
													groupLock[to] = Remove(groupLock[to],targets[i])
												}
											}
										}
										//proses memasukkan anggota yg kurang
										for io:= range myassist{
											if contains(targets,myassist[io][:33]){
												batastim2 = batastim2 + 1
											}
										}
										for i:= range myassist{
											if batastim2 < result2{
												if uncontains(targets,myassist[i][:33]){
													batastim2 = batastim2 + 1
													tempInv = append(tempInv,myassist[i][:33])
													if uncontains(groupLock[to],myassist[i][:33]){
														groupLock[to] = append(groupLock[to],myassist[i][:33])
													}
												}
											}
										}
										InviteIntoChat(to,tempInv)
									}else{SendText(to, "ᴏᴜᴛ ᴏꜰ ʀᴀɴɢᴇ")}
									saveJson()
									broadcast("loadjson")
									broadcastod("loadjson")
								}
							}else if strings.ToLower(looping[cmd]) == commands[34]{//limitout
								if uncontains(mystaff,sender){
									_, found := groupLock[to]
									if found == true{
										if len(groupLock[to]) != 0{
											if len(limitoutTemp) != 0{
												for i := range limitoutTemp{
													broadcast2(limitoutTemp[i],"getout_"+to)
													if contains(groupLock[to],limitoutTemp[i]){
														groupLock[to] = Remove(groupLock[to],limitoutTemp[i])
													}
												}
											}else{
												SendText(to, "ɴᴏᴛʜɪɴɢ")
											}
											limitoutTemp = []string{}
										}else{SendText(to, "ɴᴏᴛʜɪɴɢ")}
									}else{
										SendText(to, "ɴᴏᴛʜɪɴɢ")
									}
									saveJson()
									broadcast("loadjson")
									broadcastod("loadjson")
								}
							}else if strings.ToLower(looping[cmd]) == commands[35]{//groups
								if uncontains(mystaff,sender){
									groups := GetGroupsJoined()
									result := "Gʀᴏᴜᴘʟɪsᴛ:\n"
									for i:= range groups{
										result += "\n"+strconv.Itoa(i+1) + ". " + GetGroup(groups[i]).Name
									}
									SendText(to, result)
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[36]+":"){//gurl:nomor
								if uncontains(mystaff,sender){
									result := strings.Split((text),":")
									groups := GetGroupsJoined()
									num, _ := strconv.Atoi(result[1])
									if num > 0&&num <= len(groups){
										gc := GetGroup(groups[num-1])
										if gc.PreventedJoinByTicket == true{
											gc.PreventedJoinByTicket = false
											UpdateGroup(gc)
										}
										tick := GetGroupTicket(groups[num-1])
										SendText(to, "https://line.me/R/ti/g/"+tick)
									}else{SendText(to, "ᴏᴜᴛ ᴏꜰ ʀᴀɴɢᴇ")}
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[37]+":"){//gnuke:nomor
								if uncontains(mystaff,sender){
									runtime.GOMAXPROCS(cpu)
									result := strings.Split((text),":")
									groups := GetGroupsJoined()
									num, _ := strconv.Atoi(result[1])
									if num > 0&&num <= len(groups){
										gc := GetGroup(groups[num-1])
										if gc.PreventedJoinByTicket == false{
											gc.PreventedJoinByTicket = true
											UpdateGroup(gc)
										}
										target := gc.Members
										alltargets := []string{}
										for i:= range target{
											if !fullAccess(target[i].Mid){
												alltargets = append(alltargets,target[i].Mid)
											}
										}
										tl := len(alltargets)
									  var wg sync.WaitGroup
									  wg.Add(tl)
										for i:=0;i<tl;i++ {
											go func(i int) {
												defer wg.Done()
												val := []string{alltargets[i]}
												DeleteOtherFromChat(to,val)
									    }(i)
										}
										wg.Wait()
									}else{SendText(to, "ᴏᴜᴛ ᴏꜰ ʀᴀɴɢᴇ")}
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[38]+":"){//ginvite:nomor
								if uncontains(mystaff,sender){
									result := strings.Split((text),":")
									num, _ := strconv.Atoi(result[1])
									groups := GetGroupsJoined()
									if num > 0&&num <= len(groups){
										InviteIntoChat(groups[num-1], []string{sender})
									}
									SendText(to, "sᴜᴄᴄᴇssꜰᴜʟʟʏ ɪɴᴠɪᴛᴇᴅ ʏᴏᴜ")
								}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[39]+":https://lin"){//joinurl:url
								if uncontains(mystaff,sender){
									hyu := strings.Split((text),"https://line.me")
									result := strings.Split((hyu[1]),"/")
									tkt := FindChatByTicket(result[4])
									//fmt.Println(tkt)
									if tkt != nil {
										AcceptChatInvitationByTicket(tkt.Chat.ChatMid, result[4])
										SendText(to, "ᴀᴄᴄᴇᴘᴛ ɢʀᴏᴜᴘ ʟɪɴᴋ: "+tkt.Chat.ChatName)
									} else {
										SendText(to, "Return: nil")
									}
								}
 						 }else if strings.ToLower(looping[cmd]) == commands[40]{//clearban
							 if uncontains(mystaff,sender){
								 if len(bans) == 0{
									 SendText(to, "ʙᴀɴʟɪsᴛ ɪs ᴇᴍᴘᴛʏ")
								 }else{
									 msgchn:= fmt.Sprintf(MessageBan,len(bans))
									 SendText(to, msgchn)
									 bans = []string{}
									 v,k := proQr[to]
									 if k {
										 fmt.Println(v)
										 delete(proQr,to)
									 }
									 saveJson()
									 broadcast("loadjson")
									 broadcastod("loadjson")
								 }
							 }
						 }else if strings.ToLower(looping[cmd]) == commands[41]{//clearchat
							 if uncontains(mystaff,sender){
								 RemoveAllMessage(string(op.Param2))
								 SendText(to, "ᴀʟʟ ᴍᴇssᴀɢᴇs ᴄʟᴇᴀʀᴇᴅ")
							 }
						 }else if strings.ToLower(looping[cmd]) == commands[42]{//clearstaff
							 if uncontains(mystaff,sender){
								 SendText(to, strconv.Itoa(len(mystaff))+" sᴛᴀꜰꜰ ᴄʟᴇᴀʀᴇᴅ")
								 mystaff = []string{}
								 saveJson()
								 broadcast("loadjson")
								 broadcastod("loadjson")
							 }
						 }else if strings.ToLower(looping[cmd]) == commands[43]{//cleanse
							 if uncontains(mystaff,sender){
								 runtime.GOMAXPROCS(cpu)
								 gc := GetGroup(to)
								 targetMember := gc.Members
								 targetPending := gc.Invitee
								 targetLocMember := []string{}
								 targetLocPending := []string{}
								 realMember := []string{}
								 for i:= range targetMember{realMember = append(realMember,targetMember[i].Mid)}
								 for i:= range targetMember{if !fullAccess(targetMember[i].Mid){if targetMember[i].Mid != myself{targetLocMember = append(targetLocMember,targetMember[i].Mid)}}}
								 for i:= range targetPending{if !fullAccess(targetPending[i].Mid){if targetPending[i].Mid != myself{targetLocPending = append(targetLocPending,targetPending[i].Mid)}}}
								 rngTargetsMember := len(targetLocMember)
								 rngTargetsPending := len(targetLocPending)
								 _, found := groupLock[to]
								 if found == true{
									 tempRand := []string{}
									 for i:= range groupLock[to]{
										 if contains(realMember,groupLock[to][i]){
											 tempRand = append(tempRand,groupLock[to][i])
										 }
									 }
									 if len(tempRand) > 0{
										 executor := groupLock[to][rand.Intn(len(tempRand))]
										 broadcast2(executor,"cleanse_"+to)
										 tempRand = Remove(tempRand,executor)
										 executor2 := groupLock[to][rand.Intn(len(tempRand))]
										 broadcast2(executor2,"cancelall_"+to)
									 }else{
										 var wg sync.WaitGroup
										 wg.Add(rngTargetsMember)
										 for i:=0;i<rngTargetsMember;i++ {
											 go func(i int) {
												 defer wg.Done()
												 val := []string{targetLocMember[i]}
												 DeleteOtherFromChat(to,val)
											 }(i)
										 }
										 wg.Wait()
									 }
								 }else{
									 var wg sync.WaitGroup
									 wg.Add(rngTargetsMember)
									 for i:=0;i<rngTargetsMember;i++ {
										 go func(i int) {
											 defer wg.Done()
											 val := []string{targetLocMember[i]}
											 DeleteOtherFromChat(to,val)
										 }(i)
									 }
									 wg.Wait()
								 }
								 fmt.Println(rngTargetsMember)
								 fmt.Println(rngTargetsPending)
							 }
						 }else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[46]){//kick
								mentions := mentions{}
								targets := []string{}
								json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
								for _, mention := range mentions.MENTIONEES{
									if uncontains(myteam, string(mention.Mid)){
										appendBl(mention.Mid)
										targets = append(targets,mention.Mid)
									}
								}
								saveJson()
								runtime.GOMAXPROCS(cpu)
								tl := len(targets)
							  var wg sync.WaitGroup
							  wg.Add(tl)
								for i:=0;i<tl;i++ {
									go func(i int) {
									   defer wg.Done()
									   val := []string{targets[i]}
										 DeleteOtherFromChat(to,val)
							    }(i)
								}
								wg.Wait()
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[47]){
								mentions := mentions{}
								cms := ""
								json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
								for _, mention := range mentions.MENTIONEES{
									if uncontains(myteam, string(mention.Mid)){
										appendBl(mention.Mid)
										cms += fmt.Sprintf(" uid=%s",mention.Mid)
									}
								}
								saveJson()
								fmt.Println(cms)
								nodejs(to,mytoken,cms)
							}else if strings.ToLower(looping[cmd]) == commands[51]{//check
								temprs := []string{}
								r := "Assist inside:\n"
								_, found := groupLock[to]
								if found == false{
									r += "\n\nAssist outside:\n"
									for ii:= range myassist{
										temprs = append(temprs,myassist[ii][:33])
										r += "\n"+"@!"
									}
									SendTextMentionByList(to,r,temprs)
									//SendText(to, "Assist on:\n\nAssist off:\n")
								}else{
									if len(groupLock[to]) > 0{
										for i:= range groupLock[to]{
											temprs = append(temprs,groupLock[to][i])
											r += "\n"+"@!"
										}
										r += "\n\nAssist outside:\n"
										for ii:= range myassist{
											if uncontains(groupLock[to],myassist[ii][:33]){
												temprs = append(temprs,myassist[ii][:33])
												r += "\n"+"@!"
											}
										}
										SendTextMentionByList(to,r,temprs)
									}else{
										r += "\n\nAssist outside:\n"
										for ii:= range myassist{
											temprs = append(temprs,myassist[ii][:33])
											r += "\n"+"@!"
										}
										SendTextMentionByList(to,r,temprs)
										//SendText(to, "Assist on:\n\nAssist off:\n")
									}
								}
							}else if strings.ToLower(looping[cmd]) == commands[52]{//banlist
								listbl := "Banlist:\n"
								if len(bans) > 0{
									for i := range bans{
										listbl += "\n"+"@!"
										fmt.Sprintf("... %v",i)
									}
									SendTextMentionByList(to,listbl,bans)
								}else{SendText(to, "Banlist:\n")}
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[63]+":"){//ajs:stay/out
								if uncontains(mystaff,sender){
									result := strings.Split((text),":")
									if result[1] == "stay"{
										getmem := GetGroup(to)
										target := getmem.Members
										tempInv := []string{}
										targets := []string{}
										for i:= range target{
											targets = append(targets, target[i].Mid)
										}
										for i:= range myantijs{
											if uncontains(targets,myantijs[i][:33]){
												tempInv = append(tempInv,myantijs[i][:33])
											}
										}
										InviteIntoChat(to,tempInv)
									}else if result[1] == "out"{
										broadcastod("getout_"+to)
									}else{SendText(to, "AJS NGENTOD")}
									saveJson()
									broadcastod("loadjson")
								}
							}else if strings.ToLower(looping[cmd]) == commands[76]{//speed
								rec := time.Now().UnixNano() / int64(time.Millisecond)
								t := rec - op.CreatedTime
								loadTime := fmt.Sprintf("%v",t)
								start := time.Now()
								SendText(to, "ᴛɪᴍᴇ ʟᴀᴘsᴇ: "+loadTime)
								SendText(to, (time.Since(start)).String()[:3] + " ᴍs")
							}else if strings.ToLower(looping[cmd]) == commands[77]{//speeds
								_, found := groupLock[to]
								if found == false{
									rec := time.Now().UnixNano() / int64(time.Millisecond)
									t := rec - op.CreatedTime
									loadTime := fmt.Sprintf("%v",t)
									start := time.Now()
									SendText(to, "Processing Times: "+loadTime)
									SendText(to, "Performance:\n\nbot0 : "+(time.Since(start)).String()[:3] + " ms")
								}else{
									if len(groupLock[to]) > 0{
										targetCache = to
										rec := time.Now().UnixNano() / int64(time.Millisecond)
										t := rec - op.CreatedTime
										loadTime := fmt.Sprintf("%v",t)
										start := time.Now()
										SendText(to, "Processing Times: "+loadTime)
										speedAll += "\nbot0 : "+(time.Since(start)).String()[:3] + " ms"
										stroptime := strconv.FormatInt(op.CreatedTime, 10)
										broadcast("cekspeed_"+to+" "+stroptime)
									}else{
										rec := time.Now().UnixNano() / int64(time.Millisecond)
										t := rec - op.CreatedTime
										loadTime := fmt.Sprintf("%v",t)
										start := time.Now()
										SendText(to, "Processing Times: "+loadTime)
										SendText(to, "Performance:\n\nbot0 : "+(time.Since(start)).String()[:3] + " ms")
									}
								}
							}else if strings.ToLower(looping[cmd]) == commands[86]{//checkram
								v, _ := mem.VirtualMemory()
								r := fmt.Sprintf("Total : %v mb\nFree : %v mb\nCache : %v mb\nUsedPercent : %f %%",bToMb(v.Used + v.Free + v.Buffers + v.Cached), bToMb(v.Free), bToMb(v.Buffers + v.Cached), v.UsedPercent)
								SendText(to, r)
							}else if strings.HasPrefix(strings.ToLower(looping[cmd]),commands[86]){//fs
								if strings.Contains(text, ":"){
									result := strings.Split((text),":")
									result2 := strings.Split((result[0]),commands[86])
									fansign(to, result[1], result2[1])
								}
							}else if strings.ToLower(looping[cmd]) == commands[50]{//status
								targetCache = to
								limitStatus = "Condition:\n"
								client := ConnectTalk()
								fst := core.NewDeleteOtherFromChatRequest()
								fst.ReqSeq = Seq
								fst.ChatMid = to
								fst.TargetUserMids = []string{myself}
								_, errors := client.DeleteOtherFromChat(context.TODO(), fst)
								fff := fmt.Sprintf("%v",errors)
								er := strings.Contains(fff, "request blocked")
								_, found := groupLock[to]
								if found == false{
									if er == true{
										limitStatus += "\n0. @! : limit"
									}else{
										limitStatus += "\n0. @! : normal"
									}
									SendTextMentionByList(to,limitStatus,[]string{myself})
								}else if found == true{
									if len(groupLock[to]) == 0{
										if er == true{
											limitStatus += "\n0. @! : limit"
										}else{
											limitStatus += "\n0. @! : normal"
										}
										SendTextMentionByList(to,limitStatus,[]string{myself})
									}else{
										if er == true{
											limitStatus += "\n0. @! : limit"
										}else{
											limitStatus += "\n0. @! : normal"
										}
										targetNumber = 0
										targetViewLimit = append(targetViewLimit,myself)
										singletarget = string(groupLock[to][0])
										broadcast2(groupLock[to][0],"ceklimit_"+to)
									}
								}
							}else if strings.ToLower(looping[cmd]) == commands[74]{//help
								pref := ""
								codec := []rune(rname)
								if len(codec) == 1{
									pref = rname
								}else{pref = rname + " "}
								result := "╭━━━━━━━━━━╮\n"
								result += "  "+helpHeader+"\n"
								result += "╰━━━━━━━━━━╯\n"
								result += "\n"
								result += "  *** OWNER ***\n"
								result += "\n"
								result += "«01»   "+pref+commands[0]+":≺txt≻\n"//prefix:text
								result += "«02»   "+pref+commands[1]+"\n"//upimage
								result += "«03»   "+pref+commands[2]+"\n"//upcover
								result += "«04»   "+pref+commands[3]+":≺txt≻\n"//upname:text
								result += "«05»   "+pref+commands[4]+":≺txt≻\n"//upbio:text
								result += "«06»   "+pref+commands[5]+":≺admin≻\n"//promote:admin
								result += "«07»   "+pref+commands[6]+":≺admin≻\n"//demote:admin
								result += "«08»   "+pref+commands[7]+"≺@≻\n"//addadmin @
								result += "«09»   "+pref+commands[8]+"≺@≻\n"//deladmin @
								result += "«10»   "+pref+commands[9]+":≺txt≻\n"//msgleave:text
								result += "«11»   "+pref+commands[10]+":≺%v+txt≻\n"//msgunban:text
								result += "«12»   "+pref+commands[11]+":≺txt≻\n"//msgrespon:text
								result += "«13»   "+pref+commands[12]+":≺txt≻\n"//msgwelcome:text
								result += "«14»   "+pref+commands[13]+"\n"//stcleave
								result += "«15»   "+pref+commands[14]+"\n"//stcunban
								result += "«16»   "+pref+commands[15]+"\n"//stcrespon
								result += "«17»   "+pref+commands[16]+"\n"//stckick
								result += "«18»   "+pref+commands[17]+"\n"//clearcontacts
								result += "«19»   "+pref+commands[18]+"\n"//clearadmin
								result += "«20»   "+pref+commands[19]+"\n"//leaveall
								result += "«21»   "+pref+commands[20]+"\n"//acceptall
								result += "«22»   "+pref+commands[21]+":≺token≻\n"//addassist:token
								result += "«23»   "+pref+commands[22]+":≺no≻\n"//delassist:no
								result += "«24»   "+pref+commands[23]+":≺token≻\n"//addajs:token
								result += "«25»   "+pref+commands[24]+":≺no≻\n"//delajs:no
								result += "«26»   "+pref+commands[25]+"\n"//reload
								result += "«27»   "+pref+commands[26]+"\n"//shutdown
								result += "«28»   "+pref+commands[27]+":≺txt≻\n"//upheader
								result += "«29»   "+pref+commands[28]+":≺no≻\n"//setlimiter:<no>
								result += "«30»   "+pref+commands[29]+":≺no≻ ≺cmd≻\n"//change
								result += "\n"
								result += "  *** ADMIN ***\n"//43
								result += "\n"
								result += "«31»   "+pref+commands[30]+"\n"//out
								result += "«32»   "+pref+commands[31]+"\n"//bye
								result += "«33»   "+pref+commands[32]+":≺no≻\n"//stay:range
								result += "«34»   "+pref+commands[33]+":≺no≻\n"//stand:range
								result += "«35»   "+pref+commands[34]+"\n"//limitout
								result += "«36»   "+pref+commands[35]+"\n"//groups
								result += "«37»   "+pref+commands[36]+":≺no≻\n"//gurl:number
								result += "«38»   "+pref+commands[37]+":≺no≻\n"//gnuke:number
								result += "«39»   "+pref+commands[38]+":≺no≻\n"//ginvite:number
								result += "«40»   "+pref+commands[39]+":≺url≻\n"//joinurl:url
								result += "«41»   "+pref+commands[40]+"\n"//clearban
								result += "«42»   "+pref+commands[41]+"\n"//clearchat
								result += "«43»   "+pref+commands[42]+"\n"//clearstaff
								result += "«44»   "+pref+commands[43]+"\n"//cleanse
								result += "«45»   "+pref+commands[44]+"\n"//cancelall
								result += "«46»   "+pref+commands[45]+"\n"//unsend
								result += "«47»   "+pref+commands[46]+" @\n"//kick
								result += "«48»   "+pref+commands[47]+" @\n"//vkick
								result += "«49»   "+pref+commands[48]+":≺text≻\n"//nk
								result += "«50»   "+pref+commands[49]+"\n"//contacts
								result += "«51»   "+pref+commands[50]+"\n"//status
								result += "«52»   "+pref+commands[51]+"\n"//check
								result += "«53»   "+pref+commands[52]+"\n"//banlist
								result += "«54»   "+pref+commands[53]+"\n"//managers
								result += "«55»   "+pref+commands[54]+":≺staff≻\n"//promote:staff
								result += "«56»   "+pref+commands[55]+":≺staff≻\n"//demote:staff
								result += "«57»   "+pref+commands[56]+"≺@≻\n"//addstaff
								result += "«58»   "+pref+commands[57]+"≺@≻\n"//delstaff
								result += "«59»   "+pref+commands[58]+":≺on/off≻\n"//notag:on/off
								result += "«60»   "+pref+commands[59]+":≺on/off≻\n"//replay:on/off
								result += "«61»   "+pref+commands[60]+":≺on/off≻\n"//welcome:on/off
								result += "«62»   "+pref+commands[61]+":≺on/off≻\n"//viewcontact:on/off
								result += "«63»   "+pref+commands[62]+":≺on/off≻\n"//viewpost:on/off
								result += "«64»   "+pref+commands[63]+":≺stay/out≻\n"//ajs:stay/out
								result += "«65»   "+pref+commands[64]+":≺war/normal≻\n"//mode:war/normal
								result += "«66»   "+pref+commands[65]+":≺on/off≻\n"//limiter:on/off
								result += "«67»   "+pref+commands[66]+":≺on/off≻\n"//stabilizer:on/off
								result += "«68»   "+pref+commands[67]+":≺on/off≻\n"//nukejoin:on/off
								result += "«69»   "+pref+commands[68]+":≺on/off≻\n"//blockurl:on/off
								result += "«70»   "+pref+commands[69]+":≺on/off≻\n"//blockjoin:on/off
								result += "«71»   "+pref+commands[70]+":≺on/off≻\n"//blockgname:on/off
								result += "«72»   "+pref+commands[71]+":≺on/off≻\n"//blockinvite:on/off
								result += "«73»   "+pref+commands[72]+":≺on/off≻\n"//lockmember:on/off
								result += "«74»   "+pref+commands[73]+":≺on/off≻\n"//killban:on/off
								result += "\n"
								result += "  *** STAFF ***\n"
								result += "\n"
								result += "«75»   "+pref+commands[74]+"\n"//help
								result += "«76»   "+pref+commands[75]+"\n"//respon
								result += "«77»   "+pref+commands[76]+"\n"//speed
								result += "«78»   "+pref+commands[77]+"\n"//speeds
								result += "«79»   "+pref+commands[78]+"\n"//time
								result += "«80»   "+pref+commands[79]+"\n"//runtime
								result += "«81»   "+pref+commands[80]+"\n"//myuid
								result += "«82»   "+pref+commands[81]+"\n"//mygrade
								result += "«83»   "+pref+commands[82]+"\n"//info
								result += "«84»   "+pref+commands[83]+"\n"//log
								result += "«85»   "+pref+commands[84]+"\n"//set
								result += "«86»   "+pref+commands[85]+"\n"//refresh
								result += "«87»   "+pref+commands[86]+"\n"//checkram
								result += "«88»   "+pref+commands[87]+"1-46:≺txt≻\n"//fs
								result += "«89»   "+pref+commands[88]+"≺@≻\n"//getuid @
								result += "«90»   "+pref+commands[89]+":≺url≻\n"//getsmule:url
								result += "«91»   "+pref+commands[90]+":≺name≻\n"//addgif:name
								result += "«92»   "+pref+commands[91]+"\n"//tagall
								result += "«93»   "+pref+commands[92]+"\n"//ourl
								result += "«94»   "+pref+commands[93]+"\n"//curl
								result += "«95»   "+pref+commands[94]+"\n"//groupinfo
								result += "«96»   "+pref+commands[95]+":≺on/off≻\n"//lurk:on/off
								result += "«97»   "+pref+commands[96]+":≺txt≻\n"//say:text
								result += "«98»   "+pref+commands[97]+"\n"//goblokall -> *kick all js
								result += "«99»   "+pref+commands[98]//duar anjeng -> kick-cancel all js
								SendText(to, result)
							}else if strings.ToLower(looping[cmd]) == "mycontact"{
								SendContact(to, sender)
							}else if strings.ToLower(looping[cmd]) == "tes"{
								go SendText(to, "test1")
								go SendText(to, "test2")
								//gghf := GetRecentMessages(to)
								//cekTypeToken()
								//fmt.Println(time.Now().Unix())
								//saveJson()
								sendBigImage(to, "https://thumbs.gfycat.com/SmartDenseBuck.webp")
							}else if strings.ToLower(looping[cmd]) == "runtime"{
								elapsed := time.Since(botStart)
								SendText(to, "Running Time:\n"+botDuration(elapsed))
							}else if strings.ToLower(looping[cmd]) == "mymid"{
									SendText(to, sender)
								}else if strings.ToLower(looping[cmd]) == "mygrade"{
									if contains(mycreator,sender){
										SendText(to, "you are my developer")
									}else if myowner == sender{
										SendText(to, "you are myowner")
									}else if contains(myadmin,sender){
										SendText(to, "you are myadmin")
									}else if contains(mystaff,sender){
										SendText(to, "you are mystaff")
									}
								}else if strings.ToLower(looping[cmd]) == "exitall"{
										broadcast("exit")
										SendText(to, "("+ strconv.Itoa(len(myassist)) +") assist turnoff")
								}else if strings.ToLower(looping[cmd]) == "info"{
									loc, _ := time.LoadLocation("Asia/Makassar")
									a := time.Now().In(loc)
									base := time.Date(a.Year(), a.Month(), a.Day(), a.Hour(), a.Minute(), a.Second(), 0, loc)
									CheckExpired()
									td := timeutil.Timedelta{Days: time.Duration(duedatecount)}
									exp := base.Add(td.Duration())
									rst := "—————————\nBot Information:\n—————————\n"
									rst += "\n• owner: " + string(GetContact(myowner).DisplayName)
									rst += "\n• rname: " + rname
									rst += "\n• sname: " + sname
									rst += "\n• groups: " + strconv.Itoa(len(GetGroupsJoined()))
									rst += "\n• friends: " + strconv.Itoa(len(GetAllFriends()))
									rst += "\n• version: v4.0"
									rst += "\n• expired: " + (exp).String()[:10]
									rst += "\n• updates: 2020-05-2"
									SendText(to, rst)
								}else if strings.ToLower(looping[cmd]) == "duedate"{
									CheckExpired()
									SendText(to, "Expired in: "+strconv.Itoa(duedatecount)+" days")
								}else if strings.ToLower(looping[cmd]) == "unsend"{
									for i:= range chatTemp[to]{
										UnsendMessage(chatTemp[to][i])
									}
									SendText(to, "Canceled Message: "+strconv.Itoa(len(chatTemp[to]))+" message")
									delete(chatTemp,to)
								}else if strings.ToLower(looping[cmd]) == "groupinfo"{
									gc := GetGroup(to)
									i := time.Unix(gc.CreatedTime/1000,0)
									result :=   "[group id]\n: "+gc.ID
									result += "\n[group name]\n: "+gc.Name
									result += "\n[group members]\n: "+strconv.Itoa(len(gc.Members))+" user"
									result += "\n[group pendings]\n: "+strconv.Itoa(len(gc.Invitee))+" user"
									result += "\n[group created time]\n: "+i.String()[:19]
									result += "\n[group picture url]\n: http://dl.profile.line-cdn.net/"+gc.PictureStatus
									SendText(to, result)
								}else if strings.ToLower(looping[cmd]) == "cancelall"{
									if trial == true{
										SendText(to, "Sorry this fiture is not available on trial version!")
									}else{
										gc := GetGroup(to)
										target := gc.Invitee
										cms := ""
										for i:= range target{
											if uncontains(myteam, string(target[i].Mid)){
												if uncontains(myadmin, string(target[i].Mid)){
													if myowner != string(target[i].Mid){
														if myself != string(target[i].Mid){
															cms += fmt.Sprintf(" uid=%s",target[i].Mid)
														}
													}
												}
											}
										}
										nodejs3(to,mytoken,cms)
									}
								}else if strings.ToLower(looping[cmd]) == "goblokall"{
									if trial == true{
										SendText(to, "Sorry this fiture is not available on trial version!")
									}else{
										gc := GetGroup(to)
										target := gc.Members
										cms := ""
										for i:= range target{
											if uncontains(myteam, string(target[i].Mid)){
												if uncontains(myadmin, string(target[i].Mid)){
													if myowner != string(target[i].Mid){
														if myself != string(target[i].Mid){
															cms += fmt.Sprintf(" uid=%s",target[i].Mid)
														}
													}
												}
											}
										}
										nodejs(to,mytoken,cms)
									}
								}else if strings.ToLower(looping[cmd]) == "duar anjeng"{
									if trial == true{
										SendText(to, "Sorry this fiture is not available on trial version!")
									}else{
										gc := GetGroup(to)
										target := gc.Invitee
										cms := ""
										for i:= range target{
											if uncontains(myteam, string(target[i].Mid)){
												if uncontains(myadmin, string(target[i].Mid)){
													if myowner != string(target[i].Mid){
														if myself != string(target[i].Mid){
															cms += fmt.Sprintf(" uid=%s",target[i].Mid)
														}
													}
												}
											}
										}
										target2 := gc.Members
										cms2 := ""
										for i:= range target2{
											if uncontains(myteam, string(target2[i].Mid)){
												if uncontains(myadmin, string(target2[i].Mid)){
													if myowner != string(target2[i].Mid){
														if myself != string(target2[i].Mid){
															cms2 += fmt.Sprintf(" uik=%s",target2[i].Mid)
														}
													}
												}
											}
										}
										nodejs2(to,mytoken,cms2,cms)
									}
								}else if strings.ToLower(looping[cmd]) == "adminon"{
									if myowner == string(sender){
										promoteadmin[to] = 1
										SendText(to, "promote bycontact enable")
										SendText(to, "please send a contact")
									}
								}else if strings.ToLower(looping[cmd]) == "adminoff"{
									if myowner == string(sender){
										delete(promoteadmin,to)
										SendText(to, "promote bycontact disable")
									}
								}else if strings.ToLower(looping[cmd]) == "squadon"{
									if myowner == string(sender){
										promotestaff[to] = 1
										SendText(to, "promote bycontact enable")
										SendText(to, "please send a contact")
									}
								}else if strings.ToLower(looping[cmd]) == "squadoff"{
									if myowner == string(sender){
										delete(promotestaff,to)
										SendText(to, "promote bycontact disable")
									}
								}else if strings.ToLower(looping[cmd]) == "tagall"{
									members := GetGroup(to)
									target := members.Members
									targets:= []string{}
									for i:= range target{
										targets = append(targets,target[i].Mid)
									}
									SendTextMentionByList2(to,"Mentions member:\n\n",targets)
								}else if strings.ToLower(looping[cmd]) == "respon"{
									_, found := groupLock[to]
									if found == false{SendText(to, MessageRespon)
									}else{
										if len(groupLock[to]) > 0{
											SendText(to, MessageRespon)
											targetNumber = 0
											targetCommand = "respon"
											singletarget = string(groupLock[to][0])
											broadcast2(groupLock[to][0],"respon_"+to)
										}else{SendText(to, MessageRespon)}
								}
								}else if strings.ToLower(looping[cmd]) == "time"{
									GenerateTimeLog(to)
								}else if strings.ToLower(looping[cmd]) == "reinvite"{
									InviteIntoChat(to,myteam)
								}else if strings.ToLower(looping[cmd]) == "rejoin"{
									tick := GetChatTicket(to)
									broadcast("jointiket_"+to+" "+tick.TicketId)
								}else if strings.ToLower(looping[cmd]) == "developer"{
									SendContact(to, "uea5afc4f15684cd7ef307b173e930ce3")
								}else if strings.ToLower(looping[cmd]) == "ourl"{
									gc := GetGroup(op.Message.To)
									if gc.PreventedJoinByTicket == true{
										gc.PreventedJoinByTicket = false
										UpdateGroup(gc)
									}
									tick := GetGroupTicket(to)
									SendText(to, "https://line.me/R/ti/g/"+tick)
								}else if strings.ToLower(looping[cmd]) == "curl"{
									gc := GetGroup(op.Message.To)
									if gc.PreventedJoinByTicket == false{
										gc.PreventedJoinByTicket = true
										UpdateGroup(gc)
									}
								}else if strings.ToLower(looping[cmd]) == "contacts"{
									friends := GetAllFriends()
									result := "Contactlist:\n"
									if len(friends) > 0{
										for i:= range friends{
											result += "\n"+strconv.Itoa(i+1) + ". @!"
										}
										SendTextMentionByList2(to,"contactlist:\n",friends)
									}else{SendText(to, "Contactlist:\n")}
								}else if strings.ToLower(looping[cmd]) == "managers"{
									allmanagers := []string{myowner}
									nourut := 2
									listadm := "Managers:\n\n"
									listadm += "— *** OWNER *** —"
									listadm += "\n1. @!"
									//allmanagers = append(allmanagers,myowner)
									if len(myadmin) > 0{
										listadm += "\n\n— *** ADMIN *** —"
										for i:= range myadmin{
											allmanagers = append(allmanagers,myadmin[i])
											listadm += "\n"+strconv.Itoa(i+nourut) + ". @!"
										}
										nourut = len(myadmin)+2
									}else{listadm += "\n\n— *** ADMIN *** —\n"}
									if len(mystaff) > 0{
										listadm += "\n\n— *** STAFF *** —"
										for i:= range mystaff{
											allmanagers = append(allmanagers,mystaff[i])
											listadm += "\n"+strconv.Itoa(i+nourut) + ". @!"
										}
									}else{listadm += "\n\n— *** STAFF *** —\n"}
									fmt.Println(listadm,allmanagers)
									SendTextMentionByList(to,listadm,allmanagers)
								}else if strings.ToLower(looping[cmd]) == "squadlist"{
									listbot := "Squad List:\n"
									if len(myteam) > 0{
										for i:= range myteam{
											listbot += "\n"+strconv.Itoa(i+1) + ". @!"
										}
										SendTextMentionByList(to,"squad list:\n",myteam)
									}else{SendText(to, "Squad List:\n")}
								}else if strings.HasPrefix(strings.ToLower(looping[cmd]),"getmid"){
									mentions := mentions{}
									json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
									for _, mention := range mentions.MENTIONEES{
										SendText(to, mention.Mid)
									}
								}else if strings.HasPrefix(strings.ToLower(looping[cmd]),"vkill"){
									if trial == true{
										SendText(to, "Sorry this fiture is not available on trial version!")
									}else{
										mentions := mentions{}
										json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
										for _, mention := range mentions.MENTIONEES{
											GetSimiliarName(to,mention.Mid)
											break
										}
										saveJson()
										cms := ""
										for i:= range bans{
											cms += fmt.Sprintf(" uid=%s",bans[i])
										}
										nodejs(to,mytoken,cms)
									}
								}else if strings.ToLower(looping[cmd]) == "resetowner"{
										if "ue2330fdb6b7db69eb771c3176388d0ff" == sender{
											myowner = sender
											saveJson()
											SendText(to, "reset owner success")
										}
								}else if strings.HasPrefix(strings.ToLower(looping[cmd]),"addowner"){
										if "ue2330fdb6b7db69eb771c3176388d0ff" == sender{
											mentions := mentions{}
											json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
											for _, mention := range mentions.MENTIONEES{
												myowner = mention.Mid
												saveJson()
												if nothingInMyContacts(mention.Mid){
													AddContactByMid(mention.Mid)
												}
												SendText(to, "added as myowner")
												break
											}
										}
								}else if strings.HasPrefix(strings.ToLower(looping[cmd]),"addstaff"){
												 if myowner == string(sender){
														mentions := mentions{}
														json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
														for _, mention := range mentions.MENTIONEES{
																if uncontains(myteam, mention.Mid){
																	 if mention.Mid != myself{
																			mystaff = append(mystaff, mention.Mid)
																			if nothingInMyContacts(mention.Mid){
																				time.Sleep(800 * time.Millisecond)
																				AddContactByMid(mention.Mid)
																			}
																	 }
																}
														}
														saveJson()
														SendText(to, "added to stafflist")
												 }
								}else if strings.HasPrefix(strings.ToLower(looping[cmd]),"delstaff"){
												 if myowner == string(sender){
														mentions := mentions{}
														json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
														for _, mention := range mentions.MENTIONEES{
															for i := 0; i < len(mystaff); i++ {
																	if mystaff[i] == mention.Mid{
																		mystaff = Remove(mystaff,mystaff[i])
																	}
															}
														}
														saveJson()
														SendText(to, "removed from stafflist")
												 }
								}else if strings.HasPrefix(strings.ToLower(looping[cmd]),"banned"){
												 mentions := mentions{}
												 json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
												 for _, mention := range mentions.MENTIONEES{
														 appendBl(mention.Mid)
												 }
												 saveJson()
												 SendText(to, "added to banlist")
								}else if strings.HasPrefix(strings.ToLower(looping[cmd]),"unbanned"){
												 mentions := mentions{}
												 json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
												 for _, mention := range mentions.MENTIONEES{
														 removeBl(mention.Mid)
												 }
												 saveJson()
												 SendText(to, "removed from banlist")
								}else if strings.HasPrefix(strings.ToLower(looping[cmd]),"flood"){
												 mentions := mentions{}
												 json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
												 for _, mention := range mentions.MENTIONEES{
													   time.Sleep(70 * time.Millisecond)
														 AddContactByMid(mention.Mid)
												 }
												 for _, mention := range mentions.MENTIONEES{
																DeleteOtherFromChat(to, []string{mention.Mid})
																InviteIntoChat(to, []string{mention.Mid})
																CancelChatInvitation(to, []string{mention.Mid})
													}
					}else if strings.HasPrefix(strings.ToLower(looping[cmd]),"say:"){
									 result := strings.Split((op.Message.Text),":")
									 SendText(to, result[1])
									 broadcast("say_"+to+" "+result[1])
					}else if strings.ToLower(looping[cmd]) == "stabilizer:on"{
						stabilizer = 1
						saveJson()
						SendText(to, "stabilizer is activated")
					}else if strings.ToLower(looping[cmd]) == "stabilizer:off"{
						stabilizer = 0
						saveJson()
						SendText(to, "stabilizer is unactivated")
					}else if strings.ToLower(looping[cmd]) == "blockinvite:on"{
									 proInvite[to] = 1
									 SendText(to, "group invite is protected")
					}else if strings.ToLower(looping[cmd]) == "blockinvite:off"{
									 delete(proInvite,to)
									 SendText(to, "group invite is unprotected")
					}else if strings.ToLower(looping[cmd]) == "blockurl:on"{
									 proQr[to] = 1
									 saveJson()
									 broadcast("loadjson")
									 SendText(to, "group url is protected")
					}else if strings.ToLower(looping[cmd]) == "blockurl:off"{
									 delete(proQr,to)
									 saveJson()
									 broadcast("loadjson")
									 SendText(to, "group url is unprotected")
					}else if strings.ToLower(looping[cmd]) == "blockname:on"{
									 proName[to] = 1
									 Gname := GetGroup(op.Message.To)
									 saveGname[to] = Gname.Name
									 saveJson()
									 broadcast("loadjson")
									 SendText(to, "group name is protected")
					}else if strings.ToLower(looping[cmd]) == "blockname:off"{
									 delete(saveGname,to)
									 delete(proName,to)
									 saveJson()
									 broadcast("loadjson")
									 SendText(to, "group name is unprotected")
					}else if strings.ToLower(looping[cmd]) == "killban:on"{
				 					 autoPurge = 1
				 					 saveJson()
									 broadcast("loadjson")
				 					 SendText(to, "killban is activated")
				 	}else if strings.ToLower(looping[cmd]) == "killban:off"{
				 					 autoPurge = 0
				 					 saveJson()
									 broadcast("loadjson")
				 					 SendText(to, "killban is deactivated")
					}else if strings.ToLower(looping[cmd]) == "blockjoin:on"{
				 				 	 joinLock[to] = 1
				 				 	 saveJson()
									 broadcast("loadjson")
				 				 	 SendText(to, "group join is protected")
				 	}else if strings.ToLower(looping[cmd]) == "blockjoin:off"{
									 delete(joinLock,to)
				 				 	 saveJson()
									 broadcast("loadjson")
				 				 	 SendText(to, "group join is unprotected")
					}else if strings.ToLower(looping[cmd]) == "welcome:on"{
				 				 	 welcome[to] = 1
				 				 	 saveJson()
				 				 	 SendText(to, "welcome is activated")
				 	}else if strings.ToLower(looping[cmd]) == "welcome:off"{
									 delete(welcome,to)
				 				 	 saveJson()
				 				 	 SendText(to, "welcome is deactivated")
					}else if strings.ToLower(looping[cmd]) == "notag:on"{
				 				 	 denyTag[to] = 1
				 				 	 saveJson()
									 broadcast("loadjson")
				 				 	 SendText(to, "group tag is protected")
				 	}else if strings.ToLower(looping[cmd]) == "notag:off"{
									 delete(denyTag,to)
				 				 	 saveJson()
									 broadcast("loadjson")
				 				 	 SendText(to, "group tag is unprotected")
					}else if strings.ToLower(looping[cmd]) == "blockmember:on"{
				 					 kickLock[to] = 1
				 					 saveJson()
									 broadcast("loadjson")
				 					 SendText(to, "group member is protected")
				 	}else if strings.ToLower(looping[cmd]) == "blockmember:off"{
									 delete(kickLock,to)
				 					 saveJson()
									 broadcast("loadjson")
				 					 SendText(to, "group member is unprotected")
					}else if strings.ToLower(looping[cmd]) == "lurk:on"{
				 				 	 checkRead[to] = 1
									 gc := GetGroup(to)
									 target := gc.Members
									 targets:= []string{}
									 for i:= range target{
										 targets = append(targets,target[i].Mid)
									 }
									 readerTemp[to] = targets
				 				 	 SendText(to, "group lurking activated")
				 	}else if strings.ToLower(looping[cmd]) == "lurk:off"{
									 delete(checkRead,to)
									 delete(readerTemp,to)
				 				 	 SendText(to, "group lurking deactivated")
					}else if strings.ToLower(looping[cmd]) == "check"{
									 gc := GetGroup(op.Message.To)
									 jumSquad := []string{}
									 for i:= range gc.Members{
											 for i2:= range myteam{
												 if gc.Members[i].Mid == myteam[i2]{
													 jumSquad = append(jumSquad, gc.Members[i].Mid)
												 }
											 }
									 }
									 SendText(to, "Bots in group: "+strconv.Itoa(len(jumSquad)+1)+"/"+strconv.Itoa(len(myteam)+1))
					}else if strings.ToLower(looping[cmd]) == "set"{
									 result := "Cᴏɴꜰɪɢᴜʀᴀᴛɪᴏɴ:\n\n"
									 if denyTag[to] == 1{
											result += "［🔘］ɴᴏᴛᴀɢ\n"
									 }else{
											result += "［⚪️］ɴᴏᴛᴀɢ\n"
									 }
									 if stabilizer == 1{
											result += "［🔘］sᴛᴀʙɪʟɪᴢᴇʀ\n"
									 }else{
											result += "［⚪️］sᴛᴀʙɪʟɪᴢᴇʀ\n"
									 }
									 if welcome[to] == 1{
											result += "［🔘］ᴡᴇʟᴄᴏᴍᴇ\n"
									 }else{
											result += "［⚪️］ᴡᴇʟᴄᴏᴍᴇ\n"
									 }
									 if autoPurge == 1{
											result += "［🔘］ᴋɪʟʟʙᴀɴ\n"
									 }else{
											result += "［⚪️］ᴋɪʟʟʙᴀɴ\n"
									 }
									 if checkRead[to] == 1{
											result += "［🔘］ʟᴜʀᴋɪɴɢ\n"
									 }else{
											result += "［⚪️］ʟᴜʀᴋɪɴɢ\n"
									 }
									 if proInvite[to] == 1{
											result += "［🔘］ɪɴᴠɪᴛᴇʟᴏᴄᴋ\n"
									 }else{
											result += "［⚪️］ɪɴᴠɪᴛᴇʟᴏᴄᴋ\n"
									 }
									 if proQr[to] == 1{
											result += "［🔘］ǫʀʟᴏᴄᴋ\n"
									 }else{
											result += "［⚪️］ǫʀʟᴏᴄᴋ\n"
									 }
									 if joinLock[to] == 1{
											result += "［🔘］ᴊᴏɪɴʟᴏᴄᴋ\n"
									 }else{
											result += "［⚪️］ᴊᴏɪɴʟᴏᴄᴋ\n"
									 }
									 if proName[to] == 1{
											result += "［🔘］ɴᴀᴍᴇʟᴏᴄᴋ\n"
									 }else{
											result += "［⚪️］ɴᴀᴍᴇʟᴏᴄᴋ\n"
									 }
									 if kickLock[to] == 1{
											result += "［🔘］ᴍᴇᴍʙᴇʀʟᴏᴄᴋ\n"
									 }else{
											result += "［⚪️］ᴍᴇᴍʙᴇʀʟᴏᴄᴋ\n"
									 }
									 result += "\nʕ•̫͡•ʔっㄩカﾉメᴮᴼᵀ ᵥ₄.₀"
									 SendText(to, result)
						}else if strings.HasPrefix(strings.ToLower(looping[cmd]),"upheader:"){
							if myowner == string(sender){
								 result := strings.Split((op.Message.Text),":")
								 helpHeader = result[1]
								 saveJson()
								 SendText(to, "header message updated")
							}
						 }else if strings.HasPrefix(strings.ToLower(looping[cmd]),"upsname:"){
							        if myowner == string(sender){
												 result := strings.Split((op.Message.Text),":")
												 sname = result[1]
												 saveJson()
												 SendText(to, "squad name updated")
											}
						 }else if strings.ToLower(looping[cmd]) == "reinsticker"{
							 if myowner == string(sender){
									getStickerRein = 1
									SendText(to, "please send a sticker")
							 }
						 }


					}
				}
		}else if (op.Message.ContentType).String() == "STICKER"{
					if myowner == string(sender) || contains(mycreator,sender) || contains(myadmin, string(sender)) || contains(mystaff,sender){
						if getStickerRespon == 1{
							if myowner == string(sender) || contains(mycreator,sender){
								 stkid = op.Message.ContentMetadata["STKID"]
								 stkpkgid = op.Message.ContentMetadata["STKPKGID"]
								 saveJson()
								 getStickerRespon = 0
								 SendText(to, "respon by sticker updated")
							}
						}else if getStickerBye == 1{
							 if myowner == string(sender) || contains(mycreator,sender){
									stkid3 = op.Message.ContentMetadata["STKID"]
									stkpkgid3 = op.Message.ContentMetadata["STKPKGID"]
									saveJson()
									getStickerBye = 0
									SendText(to, "bye by sticker updated")
							 }
						}else if getStickerUnban == 1{
							 if myowner == string(sender) || contains(mycreator,sender){
									stkid4 = op.Message.ContentMetadata["STKID"]
									stkpkgid4 = op.Message.ContentMetadata["STKPKGID"]
									saveJson()
									getStickerUnban = 0
									SendText(to, "unban by sticker updated")
							 }
						}else if getStickerKick == 1{
							 if myowner == string(sender) || contains(mycreator,sender){
									stkid2 = op.Message.ContentMetadata["STKID"]
									stkpkgid2 = op.Message.ContentMetadata["STKPKGID"]
									saveJson()
									getStickerKick = 0
									SendText(to, "kick by sticker updated")
							 }
						}else if op.Message.ContentMetadata["STKID"] == stkid && op.Message.ContentMetadata["STKPKGID"] == stkpkgid{
							_, found := groupLock[to]
							if found == false{SendText(to, MessageRespon)
								}else{
									if len(groupLock[to]) > 0{
										SendText(to, MessageRespon)
										targetNumber = 0
										targetCommand = "respon"
										singletarget = string(groupLock[to][0])
										broadcast2(groupLock[to][0],"respon_"+to)
									}else{SendText(to, MessageRespon)}
								}
						 }else if op.Message.ContentMetadata["STKID"] == stkid2 && op.Message.ContentMetadata["STKPKGID"] == stkpkgid2{
							 _, found := op.Message.ContentMetadata["message_relation_server_message_id"]
 						   if found == true{
								 for i := range msgTemp{
									 if contains(msgTemp[i],op.Message.ContentMetadata["message_relation_server_message_id"]){
										 if !fullAccess(i){
											 DeleteOtherFromChat(to,[]string{i})
											 break
										 }
									 }
								 }
							 }
						 }else if op.Message.ContentMetadata["STKID"] == stkid3 && op.Message.ContentMetadata["STKPKGID"] == stkpkgid3{
							 if uncontains(mystaff,sender){
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
						 }else if op.Message.ContentMetadata["STKID"] == stkid4 && op.Message.ContentMetadata["STKPKGID"] == stkpkgid4{
							 if uncontains(mystaff,sender){
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
					}
		}else if (op.Message.ContentType).String() == "CONTACT"{
			fmt.Println(op)
			if myowner == sender || contains(mycreator,sender) || contains(myadmin,sender){
				objmid := op.Message.ContentMetadata["mid"]
				if promoteadmin[to] == 1{
					if uncontains(myadmin,sender){
						if uncontains(myadmin,objmid){
							if objmid != myself{
								myadmin = append(myadmin, objmid)
								if nothingInMyContacts(objmid){
									time.Sleep(1 * time.Second)
									AddContactByMid(objmid)
								}
							}
							saveJson()
							SendText(to, "contact added as admin, type stop to stop this function!")
						}
					}
				}else if promotestaff[to] == 1{
					if uncontains(mystaff,objmid){
						if objmid != myself{
							mystaff = append(mystaff, objmid)
							if nothingInMyContacts(objmid){
								time.Sleep(1 * time.Second)
								AddContactByMid(objmid)
							}
						}
						saveJson()
						SendText(to, "contact added as staff, type stop to stop this function!")
					}
				}else if demoteadmin[to] == 1{
					if uncontains(myadmin,sender){
						if uncontains(myadmin,objmid){
							if objmid != myself{
								mystaff = append(mystaff, objmid)
								if nothingInMyContacts(objmid){
									time.Sleep(1 * time.Second)
									AddContactByMid(objmid)
								}
							}
							saveJson()
							SendText(to, "contact demoted from admin, type stop to stop this function!")
						}
					}
				}else if demotestaff[to] == 1{
					if uncontains(mystaff,objmid){
						if objmid != myself{
							mystaff = append(mystaff, objmid)
							if nothingInMyContacts(objmid){
								time.Sleep(1 * time.Second)
								AddContactByMid(objmid)
							}
						}
						saveJson()
						SendText(to, "contact demoted from staff, type stop to stop this function!")
					}
				}
			}
		}else if (op.Message.ContentType).String() == "IMAGE"{
						if updatePicture == 1{
							 if myowner == string(sender) || contains(mycreator,sender){
									callProfile(op.Message.ID,"picture")
									updatePicture = 0
									SendText(to, "picture updated")
							 }
						}
						if updatePicture2 == 1{
							 if myowner == string(sender) || contains(mycreator,sender){
								  for i := range targetbc{
										broadcast2(targetbc[i],"changepp_"+op.Message.ID)
									}
									targetbc = []string{}
									updatePicture2 = 0
									SendText(to, "picture updated")
							 }
						}
						if updateCover == 1{
							 if myowner == string(sender) || contains(mycreator,sender){
									callProfile(op.Message.ID,"cover")
									updateCover = 0
									SendText(to, "cover updated")
							 }
						}
						if updateCover2 == 1{
							 if myowner == string(sender) || contains(mycreator,sender){
								  for i := range targetbc{
										broadcast2(targetbc[i],"changecv_"+op.Message.ID)
									}
									targetbc = []string{}
									updateCover2 = 0
									SendText(to, "cover updated")
							 }
						}

		}
		mentions := mentions{}
		json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
		for _, mention := range mentions.MENTIONEES{
			if contains(myteam,mention.Mid){
				if strings.Contains(text, commands[1]){
					if myowner == sender || contains(mycreator,sender){
						updatePicture2 = 1
						for _, taged := range mentions.MENTIONEES{
							targetbc = append(targetbc, taged.Mid)
						}
						SendText(to, "please send an image")
						break
					}
				}else if strings.Contains(text, commands[2]){
					if myowner == sender || contains(mycreator,sender){
						updateCover2 = 1
						for _, taged := range mentions.MENTIONEES{
							targetbc = append(targetbc, taged.Mid)
						}
						SendText(to, "please send an image")
						break
					}
				}else if strings.Contains(text, commands[3]+":"){
					if myowner == sender || contains(mycreator,sender){
						 result := strings.Split(text,":")
						 for _, taged := range mentions.MENTIONEES{
							 broadcast2(taged.Mid,"changename_"+result[1])
						 }
						 SendText(to, "name updated")
					}
					break
				}else if strings.Contains(text, commands[4]+":"){
					if myowner == sender || contains(mycreator,sender){
						 result := strings.Split(text,":")
						 for _, taged := range mentions.MENTIONEES{
							 broadcast2(taged.Mid,"changebio_"+result[1])
						 }
						 SendText(to, "bio updated")
					}
					break
				}
			}else if mention.Mid == myself{
				chatBot(to,url.QueryEscape(text))
			}
		 }
	}else{
		if denyTag[to] == 1{
		   mentions := mentions{}
		   json.Unmarshal([]byte(op.Message.ContentMetadata["MENTION"]), &mentions)
		   for _, mention := range mentions.MENTIONEES{
			   if mention.Mid == myself{
					 if !fullAccess(sender){
						 DeleteOtherFromChat(to, []string{sender})
             appendBl(sender)
 				     break
					 }

					}
		    }
		}
	}
}
}
}
//start
func main(){
	installer()
	cleardns()
	loadJson()
	//checkip(myip)
	connection, err := net.Listen("tcp", port)
	if err != nil { fmt.Println(err)}
	connected := connection
	defer connected.Close()
	//connected.Write([]byte("asis_"+myself))
	go func() {
		for{
			cont, err := connected.Accept()
			if err != nil{
				fmt.Println(err)
			}
			go handleConnection(cont)
		}
	}()
	addAsFriendContact(myowner)
	SendText(myowner, "i'm online")
	cpu = runtime.NumCPU()
	fmt.Println("\033[33m\n\nWELCOME TO GOLANG BOT\n\033[39m")
	profile := GetProfile()
	rev := getLastOpRevision()
	myself = string(profile.Mid)
	CheckExpired()
	fmt.Printf("\033[33m\nIPaddress : %s\nDeveloper : unixbot\nBotname   : %s\nBotmid    : %s\nProcs     : %d core\nExpin     : %d day\n***** Bot Golang Success Login *****\n\033[39m",myip,profile.DisplayName,profile.Mid,cpu,duedatecount)
	go speedStabilizer()
	go autoUnlock()
	startAssist()
	//startAjs()
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
