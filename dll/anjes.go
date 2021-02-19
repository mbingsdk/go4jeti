
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"os"
	"os/exec"
	"net"
	"sync"
	"time"
	"runtime"
	core "../dll/LINE"
	thrift "../dll/thrift"
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
var limiterset int = 0
var trial bool = false
var jsonName string = ""
var myip string = "172.24.76.153"
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
var groupLock = make(map[string][]string)
var bans = []string{}
var myteam = []string{}
var myadmin = []string{}
var mystaff = []string{}
var mycreator = []string{"ue2330fdb6b7db69eb771c3176388d0ff"}
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
var myself string = ""
var myclient string = ""
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
			AcceptChatInvitation(stringData)
			LeaveGroup(stringData)
		}
	}else if a == "loadjson"{
		loadJson()
	}else if strings.HasPrefix(a, "limiterset_"){
		stringData := a[11:]
		r,_ := strconv.Atoi(stringData)
		limiterset = r
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
	}else if strings.HasPrefix(a, "cleanse_"){
		stringData := a[8:]
		runtime.GOMAXPROCS(cpu)
		gc := GetGroup(stringData)
		targetMember := gc.Members
		targetLocMember := []string{}
		for i:= range targetMember{if !fullAccess(targetMember[i].Mid){if targetMember[i].Mid != myself{targetLocMember = append(targetLocMember,targetMember[i].Mid)}}}
		rngTargetsMember := len(targetLocMember)
		var wg sync.WaitGroup
		wg.Add(rngTargetsMember)
		for i:=0;i<rngTargetsMember;i++ {
			go func(i int) {
				defer wg.Done()
				val := []string{targetLocMember[i]}
				DeleteOtherFromChat(stringData,val)
			}(i)
		}
		wg.Wait()
	}else if strings.HasPrefix(a, "ceklimit_"){
		stringData := a[9:]
		client := ConnectTalk()
		fst := core.NewDeleteOtherFromChatRequest()
		fst.ReqSeq = Seq
		fst.ChatMid = stringData
		fst.TargetUserMids = []string{myself}
		_, errors := client.DeleteOtherFromChat(context.TODO(), fst)
		fff := fmt.Sprintf("%v",errors)
		er := strings.Contains(fff, "request blocked")
		if er == true{
			broadcast("request_"+myself+" limit "+myself)
		}else{
			broadcast("request_"+myself+" normal "+myself)
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

func ConnectPoll() *core.TalkServiceClient {
	var err error
	var transport thrift.TTransport
	transport, err = thrift.NewTHttpClient("https://gxx.line.naver.jp/P4")
	deBug("Login Thrift Client Initialize", err)
	var connect *thrift.THttpClient
	connect = transport.(*thrift.THttpClient)
	connect.SetHeader("X-Line-Access", mytoken)
	connect.SetHeader("User-Agent", "LLA/2.16.0 SMJ730G 6.0.1")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.16.0\tAndroid OS\t6.0.1")
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
	connect.SetHeader("User-Agent", "LLA/2.16.0 SMJ730G 6.0.1")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.16.0\tAndroid OS\t6.0.1")
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
	connect.SetHeader("User-Agent", "LLA/2.16.0 SMJ730G 6.0.1")
	connect.SetHeader("X-Line-Application", "ANDROIDLITE\t2.16.0\tAndroid OS\t6.0.1")
	setProtocol := thrift.NewTCompactProtocolFactory()
	protocol := setProtocol.GetProtocol(connect)
	return core.NewTalkServiceClientProtocol(connect, protocol, protocol)
}

func getLastOpRevision() int64 {
	client := GetlastOp()
	r, e := client.GetLastOpRevision(context.TODO())
	deBug("getLastOpRevision", e)
	return r
}

func fetchOperations(last int64,count int32) (r []*core.Operation){
	client := ConnectPoll()
	r, e := client.FetchOperations(context.TODO(),last,count)
	deBug("fetchOperations", e)
	return r
}

func SendText(toID string,msgText string){
	client := ConnectTalk()
	v := core.NewMessage()
	v.ContentType = core.ContentType_NONE
	v.To = toID
	v.Text = msgText
	_, e := client.SendMessage(context.TODO(), Seq, v)
	deBug("SendText", e)
}

func DeleteOtherFromChat(groupId string, contactIds []string){
	client := ConnectTalk()
	v := core.NewDeleteOtherFromChatRequest()
	v.ReqSeq = Seq
	v.ChatMid = groupId
	v.TargetUserMids = contactIds
	_, e := client.DeleteOtherFromChat(context.TODO(), v)
	deBug("DeleteOtherFromChat", e)
}

func InviteIntoChat(groupId string, contactIds []string){
	client := ConnectTalk()
	v := core.NewInviteIntoChatRequest()
	v.ReqSeq = Seq
	v.ChatMid = groupId
	v.TargetUserMids = contactIds
	_, e := client.InviteIntoChat(context.TODO(), v)
	deBug("InviteIntoChat", e)
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

func CancelChatInvitation(groupId string, contactIds []string){
	client := ConnectTalk()
	f := core.NewCancelChatInvitationRequest()
	f.ReqSeq = Seq
	f.ChatMid = groupId
	f.TargetUserMids = contactIds
	_, err := client.CancelChatInvitation(context.TODO(), f)
	deBug("CancelChatInvitation", err)
}

func AcceptChatInvitation(groupId string){
	client := ConnectTalk()
	f := core.NewAcceptChatInvitationRequest()
	f.ReqSeq = Seq
	f.ChatMid = groupId
	_, err := client.AcceptChatInvitation(context.TODO(), f)
	deBug("AcceptChatInvitation", err)
}

// **update_group_qr** //
func UpdateQrChat(groupOBJ *core.Chat){
	client := ConnectTalk()
	f := core.NewUpdateChatRequest()
	f.ReqSeq = Seq
	f.Chat = groupOBJ
	f.UpdatedAttribute = core.ChatAttribute_PREVENTED_JOIN_BY_TICKET
	_, err := client.UpdateChat(context.TODO(), f)
	deBug("UpdateChat", err)
}

// **update_group_name** //
func UpdateNameChat(groupOBJ *core.Chat){
	client := ConnectTalk()
	f := core.NewUpdateChatRequest()
	f.ReqSeq = Seq
	f.Chat = groupOBJ
	f.UpdatedAttribute = core.ChatAttribute_NAME
	_, err := client.UpdateChat(context.TODO(), f)
	deBug("UpdateChat", err)
}

func GetGroup(groupId string)(r *core.Group){
	client := ConnectTalk()
	r, err := client.GetGroup(context.TODO(), groupId)
	deBug("GetGroup", err)
	return r
}

func GetContact(id string) (r *core.Contact){
	client := ConnectTalk()
	r, err:= client.GetContact(context.TODO(), id)
	deBug("GetContact", err)
	return r
}

func GetProfile()*core.Profile{
	client := ConnectTalk()
	r, err:= client.GetProfile(context.TODO())
	deBug("GetProfile", err)
	return r
}

func LeaveGroup(groupId string){
	client := ConnectTalk()
	err := client.LeaveGroup(context.TODO(), Seq, groupId)
	deBug("LeaveGroup", err)
}

func AcceptGroupByTicket(groupMid string, ticketId string){
	client := ConnectTalk()
	err := client.AcceptGroupInvitationByTicket(context.TODO(), Seq, groupMid, ticketId)
	deBug("AcceptGroupByTicket", err)
}

func AcceptChatInvitationByTicket(groupId string, ticketId string){
	client := ConnectTalk()
	fst := core.NewAcceptChatInvitationByTicketRequest()
	fst.ReqSeq = Seq
	fst.ChatMid = groupId
	fst.TicketId = ticketId
	_, err := client.AcceptChatInvitationByTicket(context.TODO(), fst)
	deBug("AcceptChatInvitationByTicket", err)
}

func FindGroupByTicket(ticketId string)(r *core.Group){
	client := ConnectTalk()
	r, err := client.FindGroupByTicket(context.TODO(), ticketId)
	deBug("FindGroupByTicket", err)
	return r
}

func GetChatTicket(groupId string)(r *core.ReissueChatTicketResponse){
	client := ConnectTalk()
	fst := core.NewReissueChatTicketRequest()
	fst.ReqSeq = Seq
	fst.GroupMid = groupId
	r, err := client.ReissueChatTicket(context.TODO(), fst)
	deBug("GetChatTicket", err)
	return r
}

func UnsendMessage(messageId string){
	client := ConnectTalk()
	err := client.UnsendMessage(context.TODO(), Seq, messageId)
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
	groupLock = Configs.Status.Staylist
	myclient = Configs.Authoken[:33]
	vbg,_ := strconv.Atoi(argsRaw[2])
	mytoken = myantijs[vbg]
}

func startConfiguration(){
	time.Sleep(800 * time.Millisecond)
	if nothingInMyContacts(myclient){
		time.Sleep(5 * time.Second)
		AddContactByMid(myclient)
	}
	if uncontains(myteam,myclient){
			myteam = append(myteam, myclient)
	}
	for i := range myassist{
		target := string(myassist[i][:33])
		if target != myself{
			if uncontains(myteam,target){
				myteam = append(myteam, target)
				if nothingInMyContacts(target){
					time.Sleep(5 * time.Second)
					AddContactByMid(target)
				}
			}
		}
	}
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

func UpdateGroup(groupOBJ *core.Group) {
	client := ConnectTalk()
	e := client.UpdateGroup(context.TODO(), Seq, groupOBJ)
	deBug("UpdateGroup", e)
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

func callProfile(cmd1 string, cmd2 string) {
	cmd, _ := exec.Command("python3","../dll/lineProfile.py",mytoken,myself,cmd1,cmd2).Output()
	fmt.Println("\033[33m"+string(cmd)+"\033[39m")
}

func backupAssist(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	_, found := groupLock[lc]
	if found == true{
		anu := []string{myclient}
		for i := range groupLock[lc] {
			anu = append(anu, groupLock[lc][i])
		}
		appendBl(pl)
		acceptAssist(lc)
		//go func(){DeleteOtherFromChat(lc,[]string{pl})}()
		go func(){InviteIntoChat(lc, anu)}()
		//go func(){appendBl(pl)}()
	}
}

func scanNameTarget(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	if antipurge == 1{GetSimiliarName(lc,pl)}
	go func(){appendBl(pl)}()
}

func backupManagers(lc string,pl string,kr string){
	runtime.GOMAXPROCS(cpu)
	cek := GetChat([]string{lc}, true, false)
	midMembers := cek.Chat[0].Extra.GroupExtra.MemberMids
	_, foundGC := midMembers[myself]
	if foundGC == true{
		go func(){DeleteOtherFromChat(lc, []string{pl})}()
		go func(){InviteIntoChat(lc, []string{kr})}()
		go func(){appendBl(pl)}()
	} else {
		acceptAssist(lc)
	}
}

func cancelAllEnemy(lc string,pl string,pd []string){
	runtime.GOMAXPROCS(cpu)
	go func(){poolCancel(lc,pd)}()
	go func(){DeleteOtherFromChat(lc, []string{pl})}()
	go func(){appendBl(pl)}()
}

func acceptAssist(lc string){
	runtime.GOMAXPROCS(cpu)
	AcceptChatInvitation(lc)
	go func(){if autoPurge == 1{poolKickBanWhenAccept(lc)}}()
}

func acceptManagers(lc string){
	runtime.GOMAXPROCS(cpu)
	go func(){AcceptChatInvitation(lc)}()
	go func(){
		if autoPurge == 1{poolKickBanWhenAccept(lc)}
	}()
}

func proJoin(lc string,pl string){
	runtime.GOMAXPROCS(cpu)
	go func(){DeleteOtherFromChat(lc, []string{pl})}()
	go func(){appendBl(pl)}()
}

func proQrGroup(lc string,pl string,kr string){
	runtime.GOMAXPROCS(cpu)
	if kr == "4"{
		go func(){DeleteOtherFromChat(lc,[]string{pl})}()
		go func(){g := GetGroupWithoutMembers(lc);if g.PreventedJoinByTicket == false{g.PreventedJoinByTicket = true;UpdateGroup(g)}}()
		go func(){appendBl(pl)}()
	}
}

func closeQrAndKickBan(lc string,pl string,kr string){
	runtime.GOMAXPROCS(cpu)
	if kr == "4"{
		go func(){g := GetGroupWithoutMembers(lc);if g.PreventedJoinByTicket == false{g.PreventedJoinByTicket = true;UpdateGroup(g)}}()
		go func(){DeleteOtherFromChat(lc,[]string{pl})}()
	}
}

func proNameGroup(lc string,pl string,kr string){
	runtime.GOMAXPROCS(cpu)
	if kr == "1"{
		go func(){DeleteOtherFromChat(lc,[]string{pl})}()
		go func(){g := GetGroupWithoutMembers(lc);if g.Name != saveGname[lc]{g.Name = saveGname[lc];UpdateGroup(g)}}()
		go func(){appendBl(pl)}()
	}
}

func closeQrAndKickJoiner(lc string){
	runtime.GOMAXPROCS(cpu)
	go func(){poolKickBans(lc)}()
	go func(){g := GetGroupWithoutMembers(lc);if g.PreventedJoinByTicket == false{g.PreventedJoinByTicket = true;UpdateGroup(g)}}()
}

func command(op *core.Operation) {
	if op.Type == 133{
		location := op.Param1
		enemy := op.Param2
		victim := op.Param3
		if contains(myteam,victim) && !fullAccess(enemy){
			backupAssist(location,enemy)
		}else if victim == myself && !fullAccess(enemy){
			scanNameTarget(location,enemy)
		}else if contains(myteam,enemy) && !fullAccess(victim){
			appendBl(victim)
		}else if access(victim) && !fullAccess(enemy){
		    backupManagers(location,enemy,victim)
		}
	}else if op.Type == 124{
		location := op.Param1
		enemy := op.Param2
		victim := strings.Split(op.Param3,"\x1e")
		if contains(myteam,enemy) && contains(victim,myself){
			fmt.Println("Di Jepit")
		}else if checkEqual(victim,bans){
			fmt.Println("Di Jepit")
	  }else if contains(bans,enemy){
		fmt.Println("Di Jepit")
		}else if proInvite[location] == 1 && !fullAccess(enemy){
			fmt.Println("Di Jepit")
		}else if access(enemy) && contains(victim,myself){
			fmt.Println("Di Jepit")
		}
	}else if op.Type == 123{
		location := op.Param1
		poolKickBans(location)
  }else if op.Type == 126{
		location := op.Param1
		enemy := op.Param2
		victim := op.Param3
		if contains(myteam,victim) && !fullAccess(enemy){
			backupAssist(location,enemy)
		}else if victim == myself && !fullAccess(enemy){
			appendBl(enemy)
		}else if contains(myteam,enemy) && !fullAccess(victim){
			appendBl(victim)
		}else if access(victim) && !fullAccess(enemy){
			backupManagers(location,enemy,victim)
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
	}
}
//start
func main(){
	cpu = runtime.NumCPU()
	cleardns()
	loadJson()
	profile := GetProfile()
	SendText(myowner, "ajs online")
	addAsFriendContact(myowner)
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
		fmt.Println(fetch)
		looper := len(fetch)
		if looper > 0{
			ops := fetch[0]
			command(ops)
			revs := ops.Revision
			rev = MaxRevision(rev, revs)
		}
	}
}
