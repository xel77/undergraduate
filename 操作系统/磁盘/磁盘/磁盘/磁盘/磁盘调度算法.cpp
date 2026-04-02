#include<stdio.h>
#include<math.h>
#include<Windows.h>
void fcfs(int wait[], int n) {
    int now = wait[0];//初始化当前磁道数
    int i = 0;
    int move = 0;//初始化移动数
    printf("FCFS:\n");
    while (i != n) {
        move += abs(now - wait[i]);
        now = wait[i];//更新磁道数
        printf("visited track:%d\n", now);
        i++;
    } //循环计算磁道数差
    printf("平均寻道时间：%.1lf\n", 1.0 * move / (i - 1));
    system("pause");
	system("cls");
}
void sstf(int wait[], int n) {
    int now = wait[0], i = 0, a[10], mini, sum = 0;//初始化
    printf("SSTF:\nvisited track:%d\n", now);
    for (int j = 1; j < n; j++) {
        i = 0;
        while (i != n) {
            a[i] = abs(now - wait[i]);
            i++;
        }//遍历数组，通过差值找出距离最近的磁道
        mini = 1;
        for (i = 1; i < n; i++) {
            if (a[mini] > a[i]) {
                mini = i;
            }
        }
        sum += abs(now - wait[mini]);//计算移动数
        now = wait[mini];//更新当前磁道
        wait[mini] = 10000;
        printf("visited track:%d\n", now);
    }
    printf("平均寻道时间：%.1lf\n", (1.0 * sum) / (n - 1));
    system("pause");
	system("cls");
}
void scan(int wait[], int n) {
    int now = wait[0], max = wait[0], t, f, sum = 0;//初始化 
    for (int i = 0; i < n; i++) {
        for (int j = 0; j < n - i; j++) {
            if (wait[j] > wait[j + 1]) {
                t = wait[j + 1];
                wait[j + 1] = wait[j];
                wait[j] = t;
            }//将磁道排序 
        }
    }
    for (int k = 0; k < n; k++) {
        if (wait[k] == now) {
            f = k;
        }
    }//找到当前磁道下标 
    int l = f;//记录下标 
    printf("SCAN:\nvisited track:%d\n", now);
    for (f += 1; f < n + 1; f++) {
        sum += abs(now - wait[f]);
        now = wait[f];
        printf("visited track:%d\n", now);
    }//遍历右边的磁道 
    for (int m = l - 1; m > 0; m--) {
        sum += abs(now - wait[m]);
        now = wait[m];
        printf("visited track:%d\n", now);
    }//遍历左边的磁道 
    printf("平均寻道时间：%.1lf", (1.0 * sum) / (n - 1));
    system("pause");
	system("cls");
}
void init() {
    int a;
    int wait[10] = { 100,55,58,39,18,90,160,150,38,184 };
    int n = sizeof(wait) / sizeof(int);
    do{
    	 printf("  磁盘寻道算法\n");
	    printf("-----------------\n");
	    printf("| 1.FCFS        |\n");
	    printf("| 2.SSTF        |\n");
	    printf("| 3.SCAN        |\n");
	    printf("-----------------\n");
	    printf("请选择算法:\n");
	    scanf("%d", &a);
		system("cls");
	    if (a == 1) {
	        fcfs(wait, n);
	    }
	    else if (a == 2) {
	        sstf(wait, n);
	    }
	    else if (a == 3) {
	        scan(wait, n);
	    }
	    else {
	        printf("错误\n");
	    }
	}while(a);
}
int main() {
    init();
    return 0;
}
