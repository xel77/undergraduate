import socket
import ssl
import re
from urllib.parse import urlsplit

BUFFER_SIZE = 4096
cxt = ssl.create_default_context()

def get_response(url):
    # 解析URL
    parsed = urlsplit(url)
    scheme = parsed.scheme
    host = parsed.hostname
    port = parsed.port or (443 if scheme == 'https' else 80)
    path = parsed.path if parsed.path else '/'
    
    # 拼接查询参数
    if parsed.query:
        path += '?' + parsed.query
    
    # 构造HTTP请求
    request_text = (
        f"GET {path} HTTP/1.1\r\n"
        f"Host: {host}\r\n"
        "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)\r\n"
        "Connection: close\r\n"
        "\r\n"
    ).encode('utf-8')
    
    # 创建套接字连接
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        # HTTPS处理
        if scheme == 'https':
            sock = cxt.wrap_socket(sock, server_hostname=host)
        
        sock.connect((host, port))
        sock.sendall(request_text)
        
        # 接收响应数据
        response = bytearray()
        while True:
            chunk = sock.recv(BUFFER_SIZE)
            if not chunk:
                break
            response.extend(chunk)
        
        return response.decode('utf-8', errors='replace')

# 测试URL列表
test_urls = [
    "http://www.sdtbu.edu.cn/info/1043/24641.htm",
    "https://docs.python.org/3/library/index.html",
    "https://mp.weixin.qq.com/s/u9FeqoBaA3Mr0fPCUMbpqA",
    "https://mp.weixin.qq.com/s/x7eakXGyonwA4M1BTV4JIg"
]

for url in test_urls:
    try:
        # 获取网页内容
        content = get_response(url)
        
        # 生成安全文件名
        filename = re.sub(r'[:/&?#=]', '_', url)[:60]
        if not filename.endswith(('.html', '.htm')):
            filename += '.txt'
        
        # 保存文件
        with open(filename, 'w', encoding='utf-8') as f:
            f.write(content)
            print(f"成功保存: {filename} ({len(content)} 字符)")
            
    except Exception as e:
        print(f"处理 {url} 失败: {str(e)}")