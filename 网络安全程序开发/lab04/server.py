import socket
import msvcrt
from itertools import count
import threading
from struct import pack, unpack

primes = {2, 3, 5}
primes_lock = threading.Lock()
server_socket = None

def is_prime(n):
    """判断一个数是否为素数"""
    if n <= 1:
        return False
    with primes_lock:
        if n in primes:
            return True
    if n % 6 not in (1, 5):
        return False
    sqrt_n = int(n**0.5) + 1
    for i in range(3, sqrt_n, 2):
        if n % i == 0:
            return False
    return True

def generate_primes():
    """持续生成素数并添加到集合"""
    try:
        for n in count(7, 2):
            if is_prime(n):
                with primes_lock:
                    primes.add(n)
    except Exception as e:
        print(f"素数生成异常: {e}")

def handle_client(conn):
    """处理客户端请求"""
    try:
        # 接收数字长度
        length_data = conn.recv(4)
        if not length_data:
            return
        length = unpack('!I', length_data)[0]
        
        # 接收数字内容
        number_bytes = conn.recv(length)
        if not number_bytes:
            return
        number = int(number_bytes.decode('utf-8'))
        
        # 判断素数并发送响应
        if is_prime(number):
            response = f"{number} 是素数。".encode('utf-8')
        else:
            response = f"{number} 不是素数。".encode('utf-8')
        
        conn.send(pack('!I', len(response)))
        conn.send(response)
    except Exception as e:
        print(f"处理客户端时出错: {e}")
    finally:
        conn.close()

def start_server():
    """启动TCP服务器"""
    global server_socket
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind(('localhost', 5005))
    server_socket.listen(50)
    
    try:
        while True:
            conn, addr = server_socket.accept()
            threading.Thread(target=handle_client, args=(conn,), daemon=True).start()
    except Exception as e:
        print(f"服务器异常: {e}")
    finally:
        if server_socket:
            server_socket.close()

def shutdown_server():
    """关闭服务器"""
    global server_socket
    if server_socket:
        server_socket.close()
    print("\n服务器已关闭")

def monitor_input():
    """监听用户输入"""
    print("按 'q' 退出服务器...")
    while True:
        char = msvcrt.getwche().lower()
        if char == 'q':
            shutdown_server()
            break

if __name__ == "__main__":
    # 启动素数生成线程
    threading.Thread(target=generate_primes, daemon=True).start()
    
    # 启动服务器线程
    server_thread = threading.Thread(target=start_server, daemon=True)
    server_thread.start()
    
    # 启动输入监控线程
    input_thread = threading.Thread(target=monitor_input)
    input_thread.start()
    
    # 主线程等待输入线程结束
    input_thread.join()