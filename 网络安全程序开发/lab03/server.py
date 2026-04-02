import socket

def handle_client(client_socket):
    while True:
        data = client_socket.recv(1024).decode()
        if not data or data.lower() == 'bye':
            break
        response = f"AI回复：{data[::-1]}"  # 示例逻辑（反转消息）
        client_socket.send(response.encode())
    client_socket.close()

server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.bind(('0.0.0.0', 8888))
server.listen(5)
print("聊天服务器已启动，等待连接...")

while True:
    client, addr = server.accept()
    print(f"客户端 {addr} 已连接")
    handle_client(client)