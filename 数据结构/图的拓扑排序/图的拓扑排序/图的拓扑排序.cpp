#include<stdio.h>
#include<iostream>
#include<string.h>
using namespace std;
//边节点
typedef struct node {
	int date;
	struct node* next;
};
//顶点节点
typedef struct vnode {
	int data;
	char name[100];
	node* first;
	int rudushu;
}vnode,list[100];
//定义图
typedef struct graph {
	int vexnum, arcnum;
	list vexname;
};
//定义栈
typedef struct stack{
	int* top;
	int data[100];
	int length;
};
//初始化栈
void init(stack *s) {
	s->length = 0;
	s->top = s->data;
}
//压栈
void push(stack* s, int e) {
	if (s->length > 11) {
		cout << "栈已满";
		return;
	}
	*(s->top) = e;
	(s->top)++;
	s->length++;
}
//出栈
int pop(stack* p) {
	int t;
	if (p->top != p->data) {//栈非空
		p->top--;
		t = *(p->top);
		p->length--;
		return t;
	}
}
//判空
int empty(stack* s) {
	return s->length;
}
//查找顶点下标
int locate(graph g,char a[]) {
	for (int i = 0; i < g.vexnum; i++) {
		if (strcmp(g.vexname[i].name,a)==0) {
			return i;
		}
	}
	return -1;
}
//判断是否结束
int judge(graph g) {
	for (int i = 0; i < g.vexnum; i++) {
		if (g.vexname[i].rudushu > -1) {
			return 1;
		}
	}
	return 0;
}
//判断是否存在环
int judgeloop(graph g) {
	int count = 0;
	for (int i = 0; i < g.vexnum; i++) {
		if (g.vexname[i].rudushu > 0) {
			count++;
		}
	}
	if (count >= g.vexnum) {
		return 1;
	}
	return 0;
}
//拓扑排序
void toposort(graph* g, stack* s) {
	// 遍历图中的所有顶点
	for (int i = 0; i < g->vexnum; i++) {
		// 如果当前顶点的入度为0，将其压入栈中
		if (g->vexname[i].rudushu == 0) {
			push(s, g->vexname[i].data);
		}
	}
	// 初始化栈顶指针j为栈的长度减1
	int j = s->length - 1;
	// 当栈不为空时，执行循环
	while (!empty(s)) {
		// 弹出栈顶元素t
		int t = pop(s);
		// 将弹出元素的入度减1
		g->vexname[t].rudushu--;
		// 输出弹出元素的顶点名
		cout << g->vexname[t].name << "\t";
		// 获取弹出元素指向的邻接表头节点k
		node* k = g->vexname[s->data[j]].first;
		// 遍历邻接表中的节点
		while (k) {
			// 将邻接节点的入度减1
			g->vexname[k->date].rudushu--;
			// 移动邻接节点指针
			k = k->next;
		}
		// 更新栈顶指针j
		j--;
	}
}

//创建邻接表
void creat()
{
	graph g;
	stack* s;
	s = new stack;
	init(s);
	cout << "请输入顶点数和边数：";
	cin >> g.vexnum >> g.arcnum;               //输入总顶点数，总边数
	for (int i = 0; i < g.vexnum; i++)         //输入各点，构造表头(顶点)节点表
	{
		cout << "请输入课程名称：";
		cin >> g.vexname[i].name;		//输	入顶点值
		g.vexname[i].data = i;			//前缀
		g.vexname[i].rudushu = 0;		//初始化入度数
		g.vexname[i].first = NULL;         //初始化表头结点的指针域
	}
	char v1[10], v2[10];
	for (int j = 0; j < g.arcnum; j++) {
		cout << "请输入课程：";
		cin >> v1 ;
		cout << "请输入之前要学习的课程：";
		cin >> v2;
		int a = locate(g,v1);
		int b = locate(g,v2);
		if (a == -1 || b == -1) {
			if (a == -1) {
				cout << "没有"<<v1<<"这门课";
			}
			else{
				cout << "没有" << v2 << "这门课";
			}
			return;
		}
		g.vexname[a].rudushu++;
		//前插法
		node* p = new node();
		p->date = a;
		p->next = g.vexname[b].first;
		g.vexname[b].first = p;
	}
	while (judge(g)) {
		if (judgeloop(g)) {
			cout << "存在环";
			return;
		}
		toposort(&g, s);
	}
}
int main() {
	creat();
}
