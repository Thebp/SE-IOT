import urequests as requests
from lib.LTR329ALS01 import LTR329ALS01
import time
import pycom
lt = LTR329ALS01()

url = 'http://mndkk.dk:50001/light'
#data = {'data':'christian'}
#r = requests.post(url, json=data)
#print(str(r.content))
while True:
    light_data = str(lt.light()[0])
    data = {'data':light_data}
    print(light_data)
    r2 = requests.post(url, json=data)
    print(r2.content)
    print(hex(int(r2.content)))
    pycom.rgbled(int(r2.content))
    time.sleep(1)