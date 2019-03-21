import urequests as requests
from lib.LTR329ALS01 import LTR329ALS01

lt = LTR329ALS01()

url = 'http://192.168.43.85:5467/light'
data = {'data':'hej marc'}
print('Et eller  andet')
r = requests.post(url, json=data)
print('Nu printer vi JSON\n' + r.text)
print(str(r.content))