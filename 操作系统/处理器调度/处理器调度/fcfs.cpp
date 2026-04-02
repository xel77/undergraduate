#include <stdio.h>
#include <stdlib.h>
typedef struct Job {
    char name[10]; // 作业名称
    int arrival_time; // 到达时间
    int processing_time; // 运行时间
    int c; // 主存需求量
    int x; // 磁带机需求量
    int r;//未完成为0，已完成为1;
} Job;

int judge(Job* jobs, int* C, int* X, int current_time) { //优先级判断
    int m;
    for (int i = 1; i < 5; i++) {
        if (jobs[i].c <= *C && jobs[i].x <= *X && jobs[i].r == 0 && jobs[i].arrival_time < current_time) {
            *X = 4 - jobs[i].x;
            *C = 100 - jobs[i].c;
            m = i;
            jobs[i].r = 1;
            return m;
        }
    }
    if (jobs[1].r == 1 && jobs[2].r == 1 && jobs[3].r == 1 && jobs[4].r == 1) {
        m = 6;
        return m;
    }
    else {
        for (int i = 1; i < 5; i++) {
            if (jobs[i].r == 0) {
                *X = 4 - jobs[i].x;
                *C = 100 - jobs[i].c;
                jobs[i].r = 1;
                m = i;
                return m;
            }
        }
    }
}

void FCFS(Job* jobs) {
    int current_time = jobs[0].arrival_time; // 当前的时间
    int total_turnaround_time = 0; // 总的周转时间
    double total_weighted_turnaround_time = 0; // 总的带权周转时间
    int C = 100; //磁盘大小
    int X = 4; //磁带机需求量
    int a = 0;
    int A[5];//用来存最终的输出顺序，记录数字
    int m = 0;
    X -= jobs[0].x;
    C -= jobs[0].c;
    jobs[0].r = 1;
    A[a] = 0;
    a++;
    while (m < 5) {
        int start_time = jobs[m].arrival_time;
        current_time += jobs[m].processing_time;
        int turnaround_time = current_time - jobs[m].arrival_time;
        total_turnaround_time += turnaround_time;
        double weighted_turnaround_time = (double)turnaround_time / jobs[m].processing_time;
        total_weighted_turnaround_time += weighted_turnaround_time;
        printf("%s: %d:%02d-%d:%02d 周转时间：%d min 带权周转时间：%.1f min\n", jobs[m].name, start_time / 60, start_time % 60, current_time / 60, current_time % 60, turnaround_time, weighted_turnaround_time);
        m = judge(jobs, &C, &X, current_time);
        if (m == 6) {
            break;
        }
        A[a] = m;
        a++;
    }
    printf("作业执行顺序：\n");
    for (int i = 0; i < a; i++) {
        if (A[i] != 6) {
            printf("JOB%d ", A[i] + 1);
        }
    }
    printf("\n平均周转时间：%dmin\n", total_turnaround_time / a);
    printf("平均带权周转时间：%.1f\n", total_weighted_turnaround_time / a);
}

int main() {
    Job* jobs = (Job*)malloc(5 * sizeof(Job));
    jobs[0] = Job{ "JOB1", 10 * 60, 40, 35, 3, 0 };
    jobs[1] = Job{ "JOB2", 10 * 60 + 10, 30, 70, 1, 0 };
    jobs[2] = Job{ "JOB3", 10 * 60 + 15, 20, 50, 3, 0 };
    jobs[3] = Job{ "JOB4", 10 * 60 + 35, 10, 25, 2, 0 };
    jobs[4] = Job{ "JOB5", 10 * 60 + 40, 5, 20, 2, 0 };
    FCFS(jobs);
    free(jobs);
    return 0;
}
