import socket

client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.connect(('127.0.0.1', 8000))

while True:
    msg = input("请输入消息：")
    client.send(msg.encode())
    if msg.lower() == 'bye':
        break
    response = client.recv(1024).decode()
    print(response)

client.close()
