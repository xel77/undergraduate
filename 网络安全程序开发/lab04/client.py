import socket
from struct import pack, unpack

def read_exact(sock, n):
    """确保读取指定长度的数据"""
    data = b''
    while len(data) < n:
        packet = sock.recv(n - len(data))
        if not packet:  # 连接断开
            raise ConnectionError("连接已关闭")
        data += packet
    return data

def client_main():
    # 创建TCP套接字
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        try:
            # 连接到服务器
            sock.connect(('localhost', 5005))
            print("已连接服务器，输入数字进行素数验证（输入 'exit' 退出）")
            
            while True:
                # 获取用户输入
                user_input = input("请输入数字 >> ").strip()
                
                # 退出机制
                if user_input.lower() in ('exit', 'quit', 'q'):
                    print("客户端退出")
                    break
                
                # 输入验证
                try:
                    number = int(user_input)
                except ValueError:
                    print("错误：请输入有效整数")
                    continue
                
                # 数字序列化
                number_str = str(number)
                data = number_str.encode('utf-8')
                data_len = pack('!I', len(data))  # 4字节长度头

                # 发送数据
                sock.sendall(data_len + data)
                
                try:
                    # 读取响应长度
                    resp_len = unpack('!I', read_exact(sock, 4))[0]
                    # 读取响应内容
                    response = read_exact(sock, resp_len).decode('utf-8')
                    print(f"服务器响应：{response}")
                except ConnectionError as e:
                    print(f"连接异常：{e}")
                    break
                    
        except ConnectionRefusedError:
            print("无法连接服务器，请确认服务器已启动")
        except KeyboardInterrupt:
            print("\n客户端主动断开")
        except Exception as e:
            print(f"客户端错误: {e}")

if __name__ == "__main__":
    client_main()