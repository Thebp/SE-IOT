import pycom
import time

def on():
    pycom.rgbled(0xFFFFFF)#whiter than marc

def off():
    pycom.rgbled(0x000000)

commands = {
    'on': on, 
    'off': off
}


def dispatch(command):
    try:
        commands[command]()
        print('Command \'' + command + '\' dispatched')
    except KeyError:
        print('Command \'' + command + '\' not recognized')


def main():
    pycom.heartbeat(False)

    while(True):
        line = input()
        dispatch(line)


if __name__ == '__main__':
    main()