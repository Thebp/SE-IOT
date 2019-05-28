from lib.LTR329ALS01 import LTR329ALS01
import pycom
import time

class Led:
    def __init__(self):
        pycom.heartbeat(False)
        self.value = 0x000000
        pycom.rgbled(self.value)

    def set_rgb(self, red, green, blue):
        self.value = int('0x%02x%02x%02x' % (red, green, blue))
        pycom.rgbled(self.value)
    
    def ping(self):
        pycom.rgbled(0x00FF00)
        time.sleep(2)
        pycom.rgbled(self.value)



class Lightsensor:
    def __init__(self):
        self.lt = LTR329ALS01()
    
    def get_lightlevel(self):
        return self.lt.light()[0]