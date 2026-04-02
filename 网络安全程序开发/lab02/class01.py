import threading
import time
import random

# 创建一个并发访问上限为3的 BoundedSemaphore 对象
semaphore = threading.BoundedSemaphore(value=3)

def access_resource(thread_id):
    print(f"线程 {thread_id} 正在尝试获取信号量...")
    # 获取信号量，若当前已达上限，则线程将阻塞等待
    semaphore.acquire()
    try:
        print(f"线程 {thread_id} 已获取信号量，正在访问资源...")
        # 模拟访问共享资源所需的耗时操作
        time.sleep(random.uniform(1, 3))
    finally:
        # 访问结束后释放信号量
        print(f"线程 {thread_id} 释放信号量。")
        semaphore.release()

threads = []

# 创建并启动10个线程
for i in range(10):
    t = threading.Thread(target=access_resource, args=(i,))
    threads.append(t)
    t.start()

# 等待所有线程执行完毕
for t in threads:
    t.join()

print("所有线程均已完成。")
