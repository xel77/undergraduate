#include <stddef.h>
#include <stdio.h>
#include <unistd.h>
#include <pthread.h>
#include <stdlib.h>
#include <semaphore.h>

sem_t plate_mutex;       // 定义盘子互斥量
sem_t fruit_mutex;       // 定义水果互斥量
int apple_count = 0;     // 苹果计数器
int pear_count = 0;      // 梨计数器
int execution_count = 0; // 执行计数器

void* Dad(void* arg) {
    while (execution_count < 5) {
        sem_wait(&plate_mutex); // 进入临界区
        if (apple_count + pear_count < 5) {
            if (rand() % 2 == 0) {
                printf("Dad puts an apple.\n");
                apple_count++;
            } else {
                printf("Dad puts a peer.\n");
                pear_count++;
            }	
        }
        sem_post(&plate_mutex); // 离开临界区
        sleep(2);
        execution_count++; // 增加计数器
    }
    pthread_exit(NULL); // 退出线程
}

void* Son(void* arg) {
    while (execution_count < 5) {
        sem_wait(&fruit_mutex); // 进入临界区
        if (pear_count > 0) {
            printf("Son got a peer.\n");
            pear_count--;
        }
        sem_post(&fruit_mutex); // 离开临界区
        sleep(1);
    }
    pthread_exit(NULL); // 退出线程
}

void* Daughter(void* arg) {
    while (execution_count < 5) {
        sem_wait(&fruit_mutex); // 进入临界区
        if (apple_count > 0) {
            printf("Daughter got an apple.\n");
            apple_count--;
        }
        sem_post(&fruit_mutex); // 离开临界区
        sleep(1);
    }
    pthread_exit(NULL); // 退出线程
}

int main() {
    sem_init(&plate_mutex, 0, 1); // 初始化盘子互斥量
    sem_init(&fruit_mutex, 0, 1); // 初始化水果互斥量
    pthread_t father_thread, son_thread, daughter_thread; // 初始化线程
    pthread_create(&father_thread, NULL, Dad, NULL); // 创建线程
    pthread_create(&son_thread, NULL, Son, NULL);
    pthread_create(&daughter_thread, NULL, Daughter, NULL);
    pthread_join(father_thread, NULL); // 等待线程结束
    pthread_join(son_thread, NULL);
    pthread_join(daughter_thread, NULL);
    sem_destroy(&plate_mutex); // 销毁盘子互斥量
    sem_destroy(&fruit_mutex); // 销毁水果互斥量
    return 0;
}

