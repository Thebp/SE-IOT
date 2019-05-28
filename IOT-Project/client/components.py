from lib.LTR329ALS01 import LTR329ALS01
import pycom

class Led:
    def set_rgb(self, red, green, blue):
        value = int('0x%02x%02x%02x' % (red, green, blue))
        pycom.rgbled(value)

class Lightsensor:
    def __init__(self):
        self.lt = LTR329ALS01()
    
    def get_lightlevel(self):
        return self.lt.light()[0]