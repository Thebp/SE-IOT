import machine
import time
import ujson
from components import Led, Lightsensor
from mqtt import MQTTClient
from machine import Timer

MQTT_HOST = "mndkk.dk"
MQTT_USER = "iot"
MQTT_PASSWORD = "newpass12345"
MQTT_PORT = 1883

class Board:
    def __init__(self):
        self.id = str(machine.unique_id())
        print(self.id)
        self.mqtt = MQTTClient(self.id, MQTT_HOST, MQTT_USER, MQTT_PASSWORD, MQTT_PORT)
        self.led = Led()
        self.lightsensor = Lightsensor()

        self.dispatcher = {}
        self.dispatcher["{}/led/rgb".format(self.id), lambda rgb: self.led.set_rgb(rgb["red"], rgb["green"], rgb["blue"])]

    def process_message(self, topic, msg):
        topic_str = topic.decode("utf-8")
        msg_str = topic.decode("utf-8")
        self.dispatcher[topic_str](ujson.loads(msg_str))

    def publish_lightlevel(self, alarm):
        self.mqtt.publish(topic="{}/lightsensor/lightlevel".format(self.id), msg=str(self.lightsensor.get_lightlevel()))

    def run(self):
        self.mqtt.set_callback(self.process_message)
        self.mqtt.connect()

        self.mqtt.subscribe("{}/led/rgb".format(self.id))

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