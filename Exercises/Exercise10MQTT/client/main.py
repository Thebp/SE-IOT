# import pycom
# import machine
# import time
# from machine import Pin
# import socket
# import struct
# import ubinascii
# import urequests as requests
# from lib.LTR329ALS01 import LTR329ALS01
# lt = LTR329ALS01()
from mqtt import MQTTClient
import machine
from machine import Pin
import time

url = '192.168.43.85'
p_out = Pin('P19', mode=Pin.OUT)
p_out.value(1)

adc = machine.ADC()             # create an ADC object
apin = adc.channel(pin='P16')   # create an analog pin on P16
val = apin()


def sub_cb(topic, msg):
   print(msg)



# client = MQTTClient("", url,user="", password="", port=12345)
client = MQTTClient("daniel", "io.adafruit.com",user="marcndkk", password="0e78872815dd43e9951bf2915f55761b", port=1883)

# client.set_callback(sub_cb)
client.connect()
# client.subscribe(topic="marcndkk/feeds/tempF")

count = 0
while True:
   # client.connect()
   millivolts = apin.voltage()
   degC = (millivolts - 500.0) / 10.0
   client.publish(topic="marcndkk/feeds/count", msg=str(count))

   #print(degC)
   degC_data = str(degC)
   time.sleep(1)


   print("Sending Data %s" % degC_data)
   client.publish(topic="marcndkk/feeds/temp", msg=degC_data)
   # client.check_msg()
   count += 1

   time.sleep(59)