import serial
import time

file_location = 'data.txt'


def write(ser, text):
	ser.write((str(text) + '\r\n').encode())


def read(ser):
	return ser.readline().decode().rstrip('\r\n')


def write_to_device1():
	with serial.Serial("/dev/ttyACM0", 115200, timeout=None) as ser:
		start = time.time()
		print('Writing to device 1 (time: ' + str(start) + ')')
		write(ser, 'on')

		read(ser) #skip echo of what was just written
		
		if ser.in_waiting > 0:
			line = read(ser)
			print(line)
		return start

		#print('Device 1 response: ' + ser.readline().decode())


def read_from_device2(start):
	with serial.Serial("/dev/ttyACM1", 115200, timeout=None) as ser:
		print('Waiting for device 2')
		light = read(ser)
		end = time.time()
		print(light)
		print('Device 2 response at time: ' + str(end))

		return end, light


def log_to_file(*args):
	text = '\t'.join(map(str, args))
	with open(file_location, 'w+') as data_file:
		data_file.write(text)


def main():
	for i in range(100):
		start = write_to_device1()
		end, light = read_from_device2(start)

		difference = abs(start - end)

		log_to_file(start, end, difference, light)
		with serial.Serial("/dev/ttyACM0", 115200, timeout=None) as ser:
			write(ser, 'off')
			read(ser)
			line = read(ser)
			print(line)
			time.sleep(0.5)


if __name__ == '__main__':
    main()