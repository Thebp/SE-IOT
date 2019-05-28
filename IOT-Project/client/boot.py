from network import WLAN
import machine
import pycom
import time

pycom.hearbeat(False)
wlan = WLAN(mode=WLAN.STA)

ssid = "OnePlus 5"
password = "eef2a2620ab5"

access_points = wlan.scan()
for ap in access_points:
    if ap.ssid == ssid:
        print('Network found!')
        wlan.connect(ap.ssid, auth=(ap.sec, password))
        while not wlan.isconnected():
            machine.idle() # save power while waiting
        
        print('WLAN connection succeeded!')
        # 5 second blue flash to show successful connection
        pycom.rgbled(0x0000FF)
        time.sleep(5)
        pycom.rgbled(0x000000)

        machine.main('main.py')
        break