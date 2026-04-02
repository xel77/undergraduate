import threading
import time

buffer = []  # 共享缓冲区
condition = threading.Condition()  # 同步控制对象

def producer(num_items):
    for i in range(num_items):
        with condition:
            item = f"产品{i}"
            buffer.append(item)
            print(f"生产者生产了 {item}")
            condition.notify()  # 通知消费者
        time.sleep(0.1)  # 模拟生产耗时

def consumer(num_items):
    for _ in range(num_items):
        with condition:
            while not buffer:  # 缓冲区空时等待
                condition.wait()
            item = buffer.pop(0)
            print(f"消费者消费了 {item}")
        time.sleep(0.2)  # 模拟消费耗时

if __name__ == "__main__":
    num = 5
    p = threading.Thread(target=producer, args=(num,))
    c = threading.Thread(target=consumer, args=(num,))
    
    p.start()
    c.start()
    
    p.join()
    c.join()
    print("所有任务完成")