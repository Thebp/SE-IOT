import urequests as requests

url = 'http://192.168.43.85:5467/light'
data = {'data':'hej marc'}
print('Et eller  andet')
r = requests.post(url, json=data)
print('Nu printer vi JSON\n' + r.text)
print(str(r.content))