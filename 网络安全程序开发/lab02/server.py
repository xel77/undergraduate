# 服务器端代码 (server.py)
import multiprocessing
from multiprocessing.managers import BaseManager
import time

# (1) 导入必要的模块已完成

# (2) 定义共享数据结构
class SharedData:
    def __init__(self):
        self.number = 0
        self.data_list = []
        self.data_dict = {}

    def set_number(self, value):
        self.number = value
    
    def get_number(self):
        return self.number
    
    def append_list(self, item):
        self.data_list.append(item)
    
    def get_list(self):
        return self.data_list
    
    def update_dict(self, key, value):
        self.data_dict[key] = value
    
    def get_dict(self):
        return self.data_dict

# 定义Manager类
class MyManager(BaseManager):
    pass

def main():
    # (3) 设置共享数据
    shared_data = SharedData()
    
    # (4) 创建并注册管理器服务
    MyManager.register('SharedData', callable=lambda: shared_data)
    manager = MyManager(address=('0.0.0.0', 5000), authkey=b'password')
    
    # 启动服务器
    server = manager.get_server()
    print("服务器启动在 0.0.0.0:5000")
    
    # 保持服务器运行
    server.serve_forever()

if __name__ == '__main__':
    main()