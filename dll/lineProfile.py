import requests,sys,os,time,random,tempfile,shutil,base64,json,urllib
from transport import TPersistentHttpClient as thrift_client
from thrift.transport import THttpClient
from thrift.protocol import TCompactProtocol
#from api.talk import Client as anu
#from api.ttypes import TalkException,OpType,Message
from thrift.transport.TTransport import TBufferedTransport as thrift_transport
from thrift.protocol.TCompactProtocol import TCompactProtocolAccelerated as thrift_protocol
import channel
session = requests.session()
if len(sys.argv) == 5:
   pass
else:
    sys.exit()
token = sys.argv[1]
mid = sys.argv[2]
msgid = sys.argv[3]

def deleteFile(path):
    if os.path.exists(path):
        os.remove(path)
        return True
    else:
        return False

def loginChannel():
    client = thrift_client("https://gxx.line.naver.jp/CH4")
    client.setCustomHeaders({"User-Agent": "LLA/2.14.0 F-01H 6.0.1", "X-Line-Application": "ANDROIDLITE 2.14.0 Android OS 6.0.1", "X-Line-Access": token})
    LineTransport_channel = thrift_transport(client)
    LineProtocol_channel = thrift_protocol(LineTransport_channel)
    LineThriftClient_channel = channel.Client(LineProtocol_channel)
    tokens = LineThriftClient_channel.issueChannelToken("1341209850")
    return tokens.channelAccessToken

def uploadObjTalk(path, type='image', objId=""):
    files = {'file': open(path, 'rb')}
    data = {'params': genOBSParams({'oid': objId,'size': len(open(path, 'rb').read()),'type': type,'ver': '1.0', 'name': 'media'})}
    ty = session.post('https://obs-sg.line-apps.com/talk/m/upload.nhn',headers={"User-Agent": "LLA/2.14.0 F-01H 6.0.1", "X-Line-Application": "ANDROIDLITE 2.14.0 Android OS 6.0.1", "X-Line-Access": token}, data=data, files=files)

def sendImage(path, objectId):
    return uploadObjTalk(path=path, type='image', objId=objectId)

def downloadObjectMsg(messageId, returnAs='path', saveAs=''):
    if saveAs == '':
        saveAs = genTempFile('path')
    if returnAs not in ['path','bool','bin']:
        raise Exception('Invalid returnAs value')
    url = "https://obs-sg.line-apps.com/talk/m/download.nhn?oid="+messageId
    r = session.get(url,headers={"User-Agent": "LLA/2.14.0 F-01H 6.0.1", "X-Line-Application": "ANDROIDLITE 2.14.0 Android OS 6.0.1", "X-Line-Access": token},stream=True)
    if r.status_code == 200:
        saveFile(saveAs, r.raw)
        if returnAs == 'path':
            return saveAs
        elif returnAs == 'bool':
         return True
        elif returnAs == 'bin':
            return r.raw
    else:
        raise Exception('Download object failure.')

def genTempFile(returnAs='path'):
    fName, fPath = 'line-%s-%i.bin' % (int(time.time()), random.randint(0, 9)), tempfile.gettempdir()
    if returnAs == 'file':
       return fName
    elif returnAs == 'path':
         return os.path.join(fPath, fName)

def saveFile(path, raw):
    with open(path, 'wb') as f:
        shutil.copyfileobj(raw, f)

def genOBSParams(newList, returnAs='json'):
    oldList = {'name': genTempFile('file'),'ver': '1.0'}
    oldList.update(newList)
    if 'range' in oldList:
       new_range='bytes 0-%s\/%s' % ( str(oldList['range']-1), str(oldList['range']) )
       oldList.update({'range': new_range})
    if returnAs == 'json':
       oldList=json.dumps(oldList)
       return oldList
    elif returnAs == 'b64':
         oldList=json.dumps(oldList)
         return base64.b64encode(oldList.encode('utf-8'))
    elif returnAs == 'default':
         return oldList

def updateProfileCoverById(objId):
    params = {'coverImageId': objId}
    url = 'https://gxx.line.naver.jp/mh/api/v39/home/updateCover.json?'+urllib.parse.urlencode(params)
    r = session.get(url, headers={'Content-Type': 'application/json','User-Agent': "LLA/2.14.0 F-01H 6.0.1",'X-Line-Mid': mid,'X-Line-Carrier': '51089, 1-0','X-Line-Application': "ANDROIDLITE 2.14.0 Android OS 6.0.1",'X-Line-ChannelToken': loginChannel()})
    return r.json()

def uploadObjHome(path, type='image', returnAs='bool', objId=None):
    if returnAs not in ['objId','bool']:
        raise Exception('Invalid returnAs value')
    if type not in ['image','video','audio']:
        raise Exception('Invalid type value')
    if type == 'image':
        contentType = 'image/jpeg'
    elif type == 'video':
        contentType = 'video/mp4'
    elif type == 'audio':
        contentType = 'audio/mp3'
    if not objId:
        objId = int(time.time())
    file = open(path, 'rb').read()
    params = {
        'name': '%s' % str(time.time()*1000),
        'userid': '%s' % mid,
        'oid': '%s' % str(objId),
        'type': type,
        'ver': '1.0'
    }
    r = session.post('https://obs-sg.line-apps.com/myhome/c/upload.nhn', headers={'Content-Type': 'application/json','User-Agent': "LLA/2.14.0 F-01H 6.0.1",'X-Line-Mid': mid,'X-Line-Carrier': '51089, 1-0','X-Line-Application': "ANDROIDLITE 2.14.0 Android OS 6.0.1",'X-Line-ChannelToken': loginChannel(),'Content-Type': contentType,'Content-Length': str(len(file)),'x-obs-params': genOBSParams(params,'b64')}, data=file)
    if r.status_code != 201:
        raise Exception('Upload object home failure.')
    if returnAs == 'objId':
        return objId
    elif returnAs == 'bool':
        return True

def updateCover(path, returnAs='bool'):
    if returnAs not in ['objId','bool']:
        raise Exception('Invalid returnAs value')
    objId = uploadObjHome(path, type='image', returnAs='objId')
    home = updateProfileCoverById(objId)
    if returnAs == 'objId':
        return objId
    elif returnAs == 'bool':
        return True

def updateProfilePicture(path, type='p'):
    files = {'file': open(path, 'rb')}
    params = {'oid': mid,'type': 'image'}
    if type == 'vp':
       params.update({'ver': '2.0', 'cat': 'vp.mp4'})
    data = {'params': genOBSParams(params)}
    r = session.post('https://obs-sg.line-apps.com/talk/p/upload.nhn',headers={"User-Agent": "LLA/2.14.0 F-01H 6.0.1", "X-Line-Application": "ANDROIDLITE 2.14.0 Android OS 6.0.1", "X-Line-Access": token}, data=data, files=files,stream=True)
    if r.status_code != 201:
       raise Exception('Update profile picture failure.')
    return True

try:
   if sys.argv[4] == "picture":
      path = downloadObjectMsg(msgid)
      updateProfilePicture(path)
      deleteFile(path)
   elif sys.argv[4] == "cover":
        path = downloadObjectMsg(msgid)
        updateCover(path)
        deleteFile(path)
   elif sys.argv[4] == "sendimage":
        sendImage("out.png",msgid)
        deleteFile("out.png")
except Exception as e:
    print(e)
