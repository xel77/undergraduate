import socket
import threading
from struct import pack, unpack

def middle(client_socket, target_host, target_port, client_addr):
    try:
        # 连接目标服务器
        server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        server.connect((target_host, target_port))
        print(f"[*] 路由建立: {client_addr} ↔ {target_host}:{target_port}")

        while True:
            # 接收客户端数据并显示
            data = client_socket.recv(4096)
            if not data:
                break
            print(f"\n[→] 转发客户端 {client_addr} 至服务器 {target_host}:{target_port}")
            print(f"  原始数据 (HEX): {data[:100].hex()}")  # 显示前100字节HEX
            print(f"  ASCII可读部分: {repr(data.decode('utf-8', errors='replace')[:50])}")  # 显示前50字符

            # 转发到目标服务器
            server.send(data)

            # 接收服务器响应并显示
            response = server.recv(4096)
            print(f"\n[←] 服务器 {target_host}:{target_port} 响应至客户端 {client_addr}")
            print(f"  响应数据 (HEX): {response[:100].hex()}")
            print(f"  ASCII可读部分: {repr(response.decode('utf-8', errors='replace')[:50])}")

            # 转发响应到客户端
            client_socket.send(response)

    except Exception as e:
        print(f"[!] 代理异常: {e}")
    finally:
        client_socket.close()
        server.close()

def start_proxy(local_port, target_host, target_port):
    proxy = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    proxy.bind(('0.0.0.0', local_port))
    proxy.listen(200)
    print(f"[*] 代理启动: 0.0.0.0:{local_port} → {target_host}:{target_port}")

    while True:
        client, addr = proxy.accept()
        print(f"[+] 客户端连接: {addr}")
        threading.Thread(target=middle, args=(client, target_host, target_port, addr)).start()

if __name__ == '__main__':
    start_proxy(8000, '127.0.0.1', 8888)  # 将本机8080映射到目标8888端口
