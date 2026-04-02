#include<iostream>
using namespace std;
struct graph {
	int vexnunm, arcnum;
	int vexname[100];
	int linjie[100][100];
};
void init(graph* a) {
	for (int i = 0; i < 100; i++) {
		for (int j = 0; j < 100; j++) {
			a->linjie[i][j] = 0;
		}
	}
	int i = 0;
	while (a->linjie[i][i]==0) {
		a->vexname[i] = 0;
		i++;
	}
}
int judge(graph* a) {
	for (int i = 0; i < a->vexnunm; i++) {
		for (int j = 0; j < a->vexnunm; j++) {
			if (a->linjie[i][j] == 1) {
				return 1;
			}
		}
	}
	return 0;
}
int fun(graph* a,int i) {
	int flag = 0, count = 0;
	for (int k = 0; k < a->vexnunm; k++) {
		if (a->linjie[i][k] == 1) {
			count++;
		}
		if (k == 0) {
			flag = 1;
		}
	}
	if (flag) {
		if (count==1) {
			return 0;
		}
		else {
			return 1;
		}
	}	
	return 0;
}
void suanfa(graph* a) {
	int i = 0, flag=0,l=0;
	int t[100] = {};
	do{
		int k = i;
		flag = 0;
		if (fun(a, i)) {
			for (int j = 1; j < a->vexnunm; j++) {
				if (a->linjie[i][j] == 1) {
					a->linjie[i][j] = -1;
					a->linjie[j][i] = -1;
					t[l++] = a->vexname[i];
					i = j;
					break;
				}
			}
		}
		else {
			for (int j = 0; j < a->vexnunm; j++) {
				if (a->linjie[i][j] == 1) {
					a->linjie[i][j] = -1;
					a->linjie[j][i] = -1;
					t[l++] = a->vexname[i];
					i = j;
					break;
				}
			}
		}
		
		if (i != 0) {
			for (int m = i; m < a->vexnunm; m++) {
				for (int n = 0; n < a->vexnunm; n++) {
					if (a->linjie[m][n] == 1) {
						flag = 1;
					}
				}

			}
		}
		else {
			flag = 1;
		}
		if (i == k||flag==0) {
			cout << "NULL";
			return;
		}
	} while (judge(a));
	int r = 0;
	cout << "欧拉回路结果：";
	while(t[r]){
		cout << t[r]<<"\t";
		r++;
	}
	cout << a->vexname[0];
}
void creat(graph* a) {
	cout << "请输入顶点数：";
	cin >> a->vexnunm;
	for (int i = 0; i < a->vexnunm; i++) {
		cout << "请输入第"<<i+1<<"个顶点名：";
		cin >> a->vexname[i];
	}
	cout << "请输入边数：";
	cin >> a->arcnum;
	
	int m, n;
	for (int j = 0; j < a->arcnum; j++) {
		cout << "请输入vi-vj：";
		cin >> m >> n;
		a->linjie[m-1][n-1] =a->linjie[n-1][m-1]= 1;
	}
	suanfa(a);
}
int main() {
	graph a;
	init(&a);
	creat(&a);
}