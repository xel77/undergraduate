#include <stdio.h>
#include <stdlib.h>
#include <windows.h>
typedef struct Job {
	char name[9]; // 作业名称
	int arrival_time; // 到达时间
	int processing_time; // 运行时间
	int c; // 主存需求量
	int x; // 磁带机需求量
	int r;//未到达为0，调入外存为1，调入内存为2，已完成为3;
} Job;
//计算未进入内存或者未到达的作业数
int count(Job* jobs) {
	int m = 0;
	for (int i = 0; i < 5; i++) {
		if (jobs[i].r != 3 && jobs[i].r != 2) {
			m++;
		}
	}
	return m;
}
//检查到达的作业
int time(int* current_time, Job* jobs) {
	int m = 0;
	int t = 0;
	int k = count(jobs);
	while (1) {
		for (int i = 1; i < 5; i++) {
			//剩余的作业都到达 
			if (jobs[i].r != 3 && jobs[i].r != 2 && jobs[i].arrival_time <= *current_time) {
				jobs[i].r = 1;
				m++;
			}
			//剩余作业都没到达 
			else if (jobs[i].r != 3 && jobs[i].r != 2 && jobs[i].arrival_time > *current_time) {
				t++;
			}
			if (t == 4) {
				break;
			}
			//部分到达或则全部到达结束的条件 
			if (jobs[i].arrival_time > *current_time || m == k) {
				return *current_time;
			}
		}
		//若全部没到达则时间增加 
		*current_time++;
	}
}
//判断是否所有作业已经做完
int boole(Job* jobs) {
	int i = 0;
	while (i != 5) {
		if (jobs[i].r != 3) {
			return 0;
		}
		i++;
	}
	return 1;
}
//在内存中没有作业时，调入外存中的作业执行
int FCFSJudge(Job* jobs, int* C, int* X, int current_time, int star_time[]) {
	int flag = 0;
	for (int i = 1; i < 5; i++) {
		if (jobs[i].r == 1 && jobs[i].c <= *C && jobs[i].x <= *X) {//判断满足条件的作业
			jobs[i].r = 2;
			star_time[i] = current_time;//记录执行时间
			*X = 4 - jobs[i].x;//资源减少
			*C = 100 - jobs[i].c;
			//检查能调入内存的作业
			for (int j = 0; j < 5; j++) {
				if (jobs[j].r == 1 && jobs[j].c <= *C && jobs[j].x <= *X) {
					jobs[j].r = 2;
					flag++;
				}
				if (flag) {
					break;
				}
			}
			jobs[i].r = 3;//标记为完成
			return i;
		}
	}
	return 0;
}
//执行内存中的作业
int FCFSJudge2(Job* jobs, int* C, int* X, int current_time, int star_time[]) {
	int flag = 0;
	for (int i = 1; i < 5; i++) {
		if (jobs[i].r == 2 && jobs[i].c <= *C && jobs[i].x <= *X) {
			star_time[i] = current_time;//记录执行时间
			*X = *X - jobs[i].x;
			*C = *C - jobs[i].c;
			for (int j = 0; j < 5; j++) {
				if (jobs[j].r == 1 && jobs[j].c <= *C && jobs[j].x <= *X) {
					jobs[j].r = 2;
					flag++;
				}
				if (flag) {
					break;
				}
			}
			jobs[i].r = 3;
			return i;
		}
	}
	return 0;
}
//作业执行结束返回资源
void end(Job* jobs, int* C, int* X, int current_time, int star_time[]) {
	for (int k = 0; k < 5; k++) {
		if (jobs[k].r == 3 && star_time[k] + jobs[k].processing_time == current_time) {
			*X += jobs[k].x;
			*C += jobs[k].c;
		}
	}
}
void FCFS(Job* jobs) {
	int current_time = jobs[0].arrival_time; // 当前的时间
	int total_turnaround_time = 0; // 总的周转时间
	double total_weighted_turnaround_time = 0.0; // 总的带权周转时间
	int C = 100; //磁盘大小
	int X = 4; //磁带机需求量
	int a = 0;
	int A[5];//记录执行顺序数组
	int m = 1;
	X -= jobs[0].x;//默认第一个作业执行
	C -= jobs[0].c;//默认第一个作业执行
	current_time += jobs[0].processing_time;
	jobs[0].r = 3;//第一个默认完成
	A[a] = 0;
	a++;
	int B[5][3];//记录各种时间
	double weigth_time[5];//记录带权周转时间
	int star_time[5];//记录开始执行时间
	star_time[0] = jobs[0].arrival_time;
	while (!boole(jobs)) {
		//进行运算
		end(jobs, &C, &X, current_time, star_time);
		time(&current_time, jobs);
		int i = FCFSJudge2(jobs, &C, &X, current_time, star_time);
		if (i) {
			A[a] = i;
			current_time += jobs[A[a]].processing_time;
			a++;
		}
		else {
			i = FCFSJudge(jobs, &C, &X, current_time, star_time);
			if (i) {
				A[a] = i;
				current_time += jobs[A[a]].processing_time;
				a++;
			}
			else{
				current_time++;
			}
		}
	}
	for (int n = 0; n < 5; n++) {
		B[n][0] = jobs[n].arrival_time;
		B[n][1] = jobs[n].processing_time + star_time[n];
		B[n][2] = B[n][1] - B[n][0];
		weigth_time[n] = 1.0 * (jobs[n].processing_time + star_time[n] - jobs[n].arrival_time) / jobs[n].processing_time;
		total_turnaround_time += B[n][2];
		total_weighted_turnaround_time += weigth_time[n];
	}
	printf("作业执行顺序: ");
	//输出执行顺序
	for (int i = 0; i < a; i++) {
		if (A[i] != 6) {
			printf("JOB%d ", A[i] + 1);
		}
	}
	printf("\n");
	for (int i = 0; i < a; i++) {
		printf("JOB%d: %d:%02d-%d:%02d 周转时间：%d  带权周转时间：%.1lf \n", i + 1, B[i][0] / 60, B[i][0] % 60, B[i][1] / 60, B[i][1] % 60, B[i][2], weigth_time[i]);
	}
	printf("平均周转时间：%dmin\n", total_turnaround_time / a);
	printf("平均带权周转时间：%.1f\n", total_weighted_turnaround_time / a);
	system("pause");
	system("cls");
}
//找出外存中能执行的作业
int SJFJudge(Job* jobs, int* C, int* X, int current_time, int star_time[]) {
	int index = 0, index1 = 0, min = 9999, min1 = 9999;
	//判断满足条件的作业
	for (int i = 1; i < 5; i++) {
		if (jobs[i].r == 1 && jobs[i].c <= *C && jobs[i].x <= *X) {
			jobs[i].r = 2;
		}
	}
	//找出满足条件作业中最短的作业
	for (int m = 1; m < 5; m++) {
		if (jobs[m].r == 2 && jobs[m].c <= *C && jobs[m].x <= *X) {
			if (min > jobs[m].processing_time) {
				min = jobs[m].processing_time;
				index = m;
			}
		}
	}
	//将其他的作业调回外存-
	for (int k = 0; k < 5; k++) {
		if (jobs[k].r == 2) {
			if (k != index) {
				jobs[k].r = 1;
			}
		}
	}
	//执行
	if (index) {
		star_time[index] = current_time;//记录执行时间
		*X = 4 - jobs[index].x;
		*C = 100 - jobs[index].c;
		jobs[index].r = 3;
		//查找在当前作业运行中能执行的作业
		for (int j = 0; j < 5; j++) {
			if (jobs[j].r == 1 && jobs[j].c <= *C && jobs[j].x <= *X) {
				jobs[j].r = 2;
			}
		}
		//如果有多个作业满足则找出最短的作业调入内存缓冲区
		for (int m = 1; m < 5; m++) {
			if (jobs[m].r == 2 && jobs[m].c <= *C && jobs[m].x <= *X) {
				if (min1 > jobs[m].processing_time) {//找到最小的运行时间的
					min1 = jobs[m].processing_time;
					index1 = m;//找到最小的下标
				}
			}
		}
		//将其他的调回外存中
		for (int k = 0; k < 5; k++) {
			if (jobs[k].r == 2) {
				if (k != index1) {
					jobs[k].r = 1;
				}
			}
		}
		return index;
	}
	return 0;
}
//将内存缓冲区的作业调入执行
int SJFJudge2(Job* jobs, int* C, int* X, int current_time, int star_time[]) {
	int index = 0, t = 0;
	int min = 9999;
	for (int i = 1; i < 5; i++) {
		if (jobs[i].r == 2) {
			star_time[i] = current_time;//记录执行时间
			*X = *X - jobs[i].x;
			*C = *C - jobs[i].c;
			//从外存中找到在当前作业运行的时候能执行的作业
			for (int j = 0; j < 5; j++) {
				if (jobs[j].r == 1 && jobs[j].c <= *C && jobs[j].x <= *X) {
					jobs[j].r = 2;
				}
			}
			jobs[i].r = 3;
			t = i;
		}
	}
	//如果有多个作业满足则找出最短的作业调入内存缓冲区
	for (int m = 1; m < 5; m++) {
		if (jobs[m].r == 2 && jobs[m].c <= *C && jobs[m].x <= *X) {
			if (min > jobs[m].processing_time) {//找到最小的运行时间的
				min = jobs[m].processing_time;
				index = m;//找到最小的下标
			}
		}
	}
	//将其他的调回外存中
	for (int k = 0; k < 5; k++) {
		if (jobs[k].r == 2) {
			if (k != index) {
				jobs[k].r = 1;
			}
		}
	}
	return t;
}
void SJF(Job* jobs) {
	int current_time = jobs[0].arrival_time; // 当前的时间
	int total_turnaround_time = 0; // 总的周转时间
	double total_weighted_turnaround_time = 0.0; // 总的带权周转时间
	int C = 100; //磁盘大小
	int X = 4; //磁带机需求量
	int a = 0;
	int A[5];//记录执行顺序数组
	X -= jobs[0].x;//默认第一个作业执行
	C -= jobs[0].c;//默认第一个作业执行
	current_time += jobs[0].processing_time;
	jobs[0].r = 3;//第一个默认完成
	A[a] = 0;
	a++;
	int B[5][3];//记录各种时间
	double weigth_time[5];//记录带权周转时间
	int star_time[5];//记录开始执行时间
	star_time[0] = jobs[0].arrival_time;
	while (!boole(jobs)) {//进行运算
		end(jobs, &C, &X, current_time, star_time);
		time(&current_time, jobs);
		int i = SJFJudge2(jobs, &C, &X, current_time, star_time);
		if (i) {
			A[a] = i;
			current_time += jobs[A[a]].processing_time;
			a++;
		}
		else {
			i = SJFJudge(jobs, &C, &X, current_time, star_time);
			if (i) {
				A[a] = i;
				current_time += jobs[A[a]].processing_time;
				a++;
			}
			else{
				current_time++;
			}
		}
	}
	for (int n = 0; n < 5; n++) {
		B[n][0] = jobs[n].arrival_time;
		B[n][1] = jobs[n].processing_time + star_time[n];
		B[n][2] = B[n][1] - B[n][0];
		weigth_time[n] = 1.0 * (jobs[n].processing_time + star_time[n] - jobs[n].arrival_time) / jobs[n].processing_time;
		total_turnaround_time += B[n][2];
		total_weighted_turnaround_time += weigth_time[n];
	}
	printf("作业执行顺序: ");
	//输出执行顺序
	for (int i = 0; i < a; i++) {
		if (A[i] != 6) {
			printf("JOB%d ", A[i] + 1);
		}
	}
	printf("\n");
	for (int i = 0; i < a; i++) {
		printf("JOB%d: %d:%02d-%d:%02d 周转时间：%3d  带权周转时间：%.1lf \n", i + 1, B[i][0] / 60, B[i][0] % 60, B[i][1] / 60, B[i][1] % 60, B[i][2], weigth_time[i]);
	}
	printf("平均周转时间：%dmin\n", total_turnaround_time / a);
	printf("平均带权周转时间：%.1f\n", total_weighted_turnaround_time / a);
	system("pause");
	system("cls");
}
int init() {
	int t;
	//初始化
	Job* jobs = (Job*)malloc(5 * sizeof(Job));
	jobs[0] = Job{ "JOB1", 10 * 60, 40, 35, 3, 0 };
	jobs[1] = Job{ "JOB2", 10 * 60 + 10, 30, 70, 1, 0 };
	jobs[2] = Job{ "JOB3", 10 * 60 + 15, 20, 50, 3, 0 };
	jobs[3] = Job{ "JOB4", 10 * 60 + 35, 10, 25, 2, 0 };
	jobs[4] = Job{ "JOB5", 10 * 60 + 40, 5, 20, 2, 0 };
	do {
		printf("作业调度算法\n");
		printf("  ----------\n");
		printf("  |1.FCFS  |\n");
		printf("  |2.SJF   |\n");
		printf("  |0.EXIT  |\n");
		printf("  ----------\n");
		printf("请进行选择\n");
		scanf("%d", &t);
		if (t == 1) {
			FCFS(jobs);
			return 1;
		}
		else if (t == 2) {
			SJF(jobs);
			return 1;
		}
		else {
			return 0;
			system("cls");
		}
	} while (t);
	
}
int main() {
	while (init()) {
	}
}
