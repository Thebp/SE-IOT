import machine
import time
import ujson
from components import Led, Lightsensor
from mqtt import MQTTClient
from machine import Timer
import ubinascii

MQTT_HOST = "mndkk.dk"
MQTT_USER = "iot"
MQTT_PASSWORD = "uS831ACCL6sZHz4"
MQTT_PORT = 1883

class Board:
    def __init__(self):
        self.id = ubinascii.hexlify(machine.unique_id()).decode("utf-8") 
        print("machine id: {}".format(self.id))
        self.mqtt = MQTTClient(self.id, MQTT_HOST, MQTT_PORT, MQTT_USER, MQTT_PASSWORD)
        self.led = Led()
        self.lightsensor = Lightsensor()

        self.dispatcher = {}
        self.dispatcher["{}/led/rgb".format(self.id)] = lambda rgb: self.led.set_rgb(rgb["red"], rgb["green"], rgb["blue"])
        self.dispatcher["{}/led/ping".format(self.id)] = lambda msg: self.led.ping()

    def process_message(self, topic, msg):
        topic_str = topic.decode("utf-8")
        msg_str = topic.decode("utf-8")
        if topic_str in self.dispatcher:
            self.dispatcher[topic_str](ujson.loads(msg_str))

    def publish_lightlevel(self, alarm):
        self.mqtt.publish(topic="lightdata", msg=ujson.dumps({"lightlevel":self.lightsensor.get_lightlevel(),"board_id":self.id}))

    def run(self):
        self.mqtt.set_callback(self.process_message)
        self.mqtt.connect()

        self.mqtt.subscribe("{}/led/rgb".format(self.id))
        self.mqtt.subscribe("{}/led/ping".format(self.id))

        self.mqtt.publish(topic="board_discovery", msg=ujson.dumps({"id":self.id}))

        alarms = []
        alarms.append(Timer.Alarm(handler=self.publish_lightlevel, s=5, periodic=True))

        try:
            while True:
                self.mqtt.wait_msg()
                machine.idle()
        finally:
            for alarm in alarms:
                alarm.cancel()
            self.mqtt.disconnect()

if __name__ == "__main__":
    Board().run()