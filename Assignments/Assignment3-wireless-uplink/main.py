import pycom
import machine
import time
from machine import Pin
import socket
import struct
import ubinascii
import urequests as requests
from lib.LTR329ALS01 import LTR329ALS01
lt = LTR329ALS01()

url = 'http://mndkk.dk:50001/data'

headers = {'IGmJuizNPufozspY4Xry':'NuB19BA6TV8bmvTpgOgo'}

id = machine.unique_id()
board = ubinascii.hexlify(id).decode("utf-8") 
print(board)

p_out = Pin('P19', mode=Pin.OUT)
p_out.value(1)

adc = machine.ADC()             # create an ADC object
apin = adc.channel(pin='P16')   # create an analog pin on P16
val = apin()

count = 0

pycom.heartbeat(False)

while True:
  millivolts = apin.voltage()
  degC = (millivolts - 500.0) / 10.0
  
  print(degC)
  degC_data = str(degC)
  light_data = str(lt.light()[0])
  print(light_data)
  data = {'light':light_data, 'temp':degC_data, 'board':board, 'count':count}
  count += 1
  try:
    r2 = requests.post(url, json=data, headers=headers)
    r2.close()
    if(r2.status_code == 200):
      pycom.rgbled(0x001100)
      time.sleep(0.5)
      pycom.rgbled(0x000000)
      print("Code 200")
    else:
      pycom.rgbled(0xFF0000)
      time.sleep(0.5)
      pycom.rgbled(0x000000)
      print("something is shit")
  except Exception as e:
    pycom.rgbled(0xFFFF00)
    print(e)

  #print("Message sent, sleeping for 15 minutes.")
  time.sleep(60)