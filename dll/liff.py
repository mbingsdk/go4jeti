import requests,json,sys
from Liff.ttypes import LiffChatContext, LiffContext, LiffSquareChatContext, LiffNoneContext, LiffViewRequest
from thrift.transport import THttpClient
from thrift.protocol import TCompactProtocol
from Liff import LiffService
session = requests.session()
token = sys.argv[1]
to = sys.argv[2]
imgurl = sys.argv[3]
def allowLiff():
    data = {'on': ['P', 'CM'], 'off': []}
    headers = {
    'X-Line-Access': token,
    'X-Line-Application': "IOS\t10.10.1\tiPhone OS\t1",
    'X-Line-ChannelId': "1653779160",
    'Content-Type': 'application/json'
    }
    r = session.post("https://access.line.me/dialog/api/permissions", headers=headers, data=json.dumps(data))
    return r.json()
def connectTalk():
    transport = THttpClient.THttpClient("https://legy-jp.line.naver.jp/LIFF1")
    transport.setCustomHeaders({"User-Agent": "Line/10.10.1", "X-Line-Application": "IOS\t10.10.1\tiPhone OS\t1", "X-Line-Access": token})
    protocol = TCompactProtocol.TCompactProtocol(transport)
    transport.open()
    return LiffService.Client(protocol)
def sendLiff(to, data):
    xyz = LiffChatContext(to)
    xyzz = LiffContext(chat=xyz)
    view = LiffViewRequest('1653779160-yw2l2v9d', xyzz)
    token = connectTalk().issueLiffView(view)
    url = 'https://api.line.me/message/v3/share'
    headers = {'Content-Type': 'application/json','Authorization': 'Bearer %s' % token.accessToken}
    data = {"messages":[data]}
    requests.post(url, headers=headers, data=json.dumps(data))
allowLiff()
data = {
    "type": "template",
    "altText": "Big Sticker",
    "template": {
        "type": "image_carousel",
        "columns": [
            {
                "imageUrl":imgurl,
                "size": "full",
                "action": {
                    "type": "uri",
                    "uri": imgurl
                }
            }
        ]
    }
}
sendLiff(to, data)
