import urllib.request
import os
import urllib.parse

# 定义文件上传函数，接收文件名作为参数
def file_upload(url, filename):
    with open(filename, 'rb') as file:
        file_data = file.read()
        file_size = os.path.getsize(filename)
        headers = {'Content-Length': str(file_size)}
        req = urllib.request.Request(url, data=file_data, headers=headers, method='POST')
        try:
            with urllib.request.urlopen(req) as response:
                if response.status == 200:
                    print("文件上传成功")
                else:
                    print("文件上传失败")
        except urllib.error.URLError as e:
            print("文件上传失败:", e.reason)

# 获取脚本所在的目录路径
script_dir = os.path.dirname(os.path.realpath(__file__))

# 获取同级目录下的 info.txt 的绝对路径
filename = os.path.join(script_dir, 'info.txt')

# 将文件路径转换为 Windows 风格的路径
filename = filename.replace('\\', '/')  # 使用双反斜杠

print(filename)

# 调用文件上传函数，传递文件名作为参数
url = 'http://192.168.159.143:8000/%E6%A1%8C%E9%9D%A2/'
url = urllib.parse.quote(url, safe='/:')  # 对URL进行编码
file_upload(url, filename)
