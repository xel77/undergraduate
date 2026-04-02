import threading
import time

# 普通函数
def regular_function(name):
    for i in range(3):
        print(f"{name} 执行步骤 {i}")
        time.sleep(0.1)

# 线程启动函数
def start_thread(name):
    def task():
        for i in range(3):
            print(f"{name} 执行步骤 {i}")
            time.sleep(0.1)
    thread = threading.Thread(target=task)
    thread.start()
    return thread

if __name__ == "__main__":
    print("=== 调用普通函数 ===")
    regular_function("函数调用1")
    regular_function("函数调用2")

    print("\n=== 启动线程 ===")
    t1 = start_thread("线程1")
    t2 = start_thread("线程2")
    t1.join()
    t2.join()

























