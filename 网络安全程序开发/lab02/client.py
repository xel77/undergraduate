# 客户端代码 (client.py)
from multiprocessing.managers import BaseManager

# 定义Manager类
class MyManager(BaseManager):
    pass

def main():
    # (5) 在接收端获取共享数据
    # 注册共享对象
    MyManager.register('SharedData')
    
    # 连接到服务器
    manager = MyManager(address=('localhost', 5000), authkey=b'password')
    manager.connect()
    
    # 获取共享数据对象
    shared_data = manager.SharedData()
    
    # 示例操作
    # 修改数字
    shared_data.set_number(42)
    print(f"Number: {shared_data.get_number()}")
    
    # 添加列表项
    shared_data.append_list("Hello")
    shared_data.append_list("World")
    print(f"List: {shared_data.get_list()}")
    
    # 更新字典
    shared_data.update_dict("key1", "value1")
    shared_data.update_dict("key2", "value2")
    print(f"Dict: {shared_data.get_dict()}")

if __name__ == '__main__':
    main()