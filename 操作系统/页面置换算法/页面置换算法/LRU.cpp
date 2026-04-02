#include<stdio.h>
struct stack {
	int* top;
	int data[3];
	int size;
};//定义栈

void init(stack* s){
	s->size = 3;
	s->top=s->data;
}//初始化栈

void push(stack* s,int e) {
	if (s->top -s->data>s->size) {
		printf("栈已满\n");
		return;
	}
	*(s->top) = e;
	(s->top)++;
}//压栈

int pop(stack* p) {
	int t;
	if (p->top != p->data) {//栈非空  
		p->top--;
		t = *(p->top);
		return t;
	}
}//出栈

void judge(stack* s,int e) {
	stack* x = s;
	int n,a[3] = {};
	for (int i = 0; i < 3; i++) {
		if (s->data[i] == e) {
			switch (i) {
				case 0:
					n = 3;
					break;
				case 1:
					n = 2;
					break;
				case 2:
					n = 1;
					break;
				default:
					printf("eror\n");
			}
		}
	}//查找需要置顶的元素下标，并计算需要弹出的元素个数
	for (int j = 0; j < n; j++) {
		a[j] = pop(x);
	}//弹出元素（包括置顶元素
	for (int k = n - 2; k >-1; k--) {
		push(x, a[k]);
	}//压入置顶元素之前的元素
	push(x, a[n-1]);//压入置顶元素
}
int caculate(int wait[], stack* s, int i) {
	int b[3] = {};//定义物理块信号量，查找到为1，否则为0
	int count=0,t;
	for (i -= 1; i > -1; i--) {//从缺页元素的前一个元素从右往左遍历页面号
		for (int j = 0; j < 3; j++) {
			if (wait[i] == s->data[j]) {
				b[j] = 1;
			}
		}
		count = 0;//重置计数器
		for (int k = 0; k < 3; k++) {
			if (b[k] == 1) {
				count++;
			}
			else {
				t = k;
			}
		}
		if (count == 2) {//计数器为2，即找到两个元素，停止遍历，返回信号量为0的下标
			return t;
		}
	}
	int flag = 0;
	for (int l = 0; l < 3; l++) {//若计数器不为2，即遍历结束只有一个元素或没有元素出现的情况
		if (b[l] == 1) {
			flag++;
		}
	}
	switch (flag) {
		case 0://为零时，更新栈底元素
			return 0;
		case 1://为1时，如果信号量1为栈底元素更新后一个元素，否则更新栈底元素
			for (int v = 0; v < 3; v++) {
				if (b[v] == 1) {
					if (v == 0) {
						return 1;
					}
					else {
						return 0;
					}
				}
			}
	}
		
}
void lru() {
	int count=0,sum=0,wait[10] = { 1,2,4,2,6,2,1,5,6,1 };
	stack s;
	init(&s);
	for (int i = 1; i <= 3; i++) {
		push(&s, i);
	}//初始化物理块元素
	for (int j = 0; j < (sizeof(wait) / sizeof(int)); j++) {
		count = 0;//计数器更新
		for (int k = 2; k >-1; k--) {
			if (wait[j] != s.data[k]) {
				count++;//计数，若等于3，即物理块中没有页面号则缺页
			}
			else {//若存在，则提至栈顶
				judge(&s, wait[j]);
				if (j == 9) {
					printf("%d \n", wait[j]);
				}
				else {
					printf("%d ", wait[j]);
				}
				}
			if (count == 3) {
				sum++;
				s.data[caculate(wait, &s, j)]=wait[j];//缺页计算
				if (j == 9) {
					printf("%dx \n", wait[j]);
				}
				else {
					printf("%dx ",wait[j]);
				}
			}
		}
	}
	printf("LRU：f=%.1lf", (1.0 * sum) / (sizeof(wait) / sizeof(int)));//计算缺页率
}
int main() {
	lru();
}