#include<iostream>
#include <windows.h>
using namespace std;
int rejudge(int *a,int n,int t) {
	for(int i=0;i<n;i++){
		if (a[i] == t) {
			return 0;
		}
	}
	return 1;
}
void ziji(int *a,int* b,int n) {
	int m;
re:
	cout << "\t请输入子集的元素个数：" << endl;
	cin >> m;
	if (m > n) {
		cout << "子集元素数不能大于全集，请重新输入";
		Sleep(3000);
		system("cls");
		goto re;
	}
	cout << "请输入元素：" << endl;
	for (int i = 0; i < m; i++) {
		cin >> b[i];
		int count = 0;
		for (int j = 0; j < n; j++) {
			if (a[j] != b[i]) {
				count++;
			}
			if (!rejudge(b, i, b[i])) {
				cout << "元素重复，请重新输入子集";
				Sleep(3000);
				system("cls");
				goto re;
			}
			if (count == n) {
				cout << "全集中不存在该元素，请重新输入子集";
				Sleep(3000);
				system("cls");
				goto re;
			}

		}
	}
	b[m] = NULL;
}
void display(int* a) {
	int i = 0;
	while (a[i] != NULL) {
		cout << a[i];
		i++;
	}
}
void con(int *a, int *b,int c[]) {
	int k = 0; 
	for (int m=0; a[m] != NULL; m++) {
		for (int j = 0; b[j] != NULL; j++) {
			if (a[m] == b[j]) {
				c[k] = a[m];
				k++;
			}
		}
	}
	c[k++] = NULL;
	cout << "A和B集合的交集为：" << endl;
	display(c);
}
void uni(int a[], int b[], int c[]) {
	int i;
	for (i = 0; a[i] != NULL; i++) {
		c[i] = a[i];
	}
	for (int j = 0; b[j] != NULL; j++) {
		int flag = 1;
		for (int k = 0; a[k] != NULL; k++) {
			if (b[j] == a[k]) {
				flag = 0;
				break;
			}
		}
		if (flag) {
			c[i++] = b[j];
		}
	}
	c[i++] = NULL;
	cout << "A和B集合的并集为：" << endl;
	display(c);
}
void diff(int a[], int b[], int c[]) {
	int i,d[100];
	for (i = 0; a[i] != NULL; i++) {
		c[i] = a[i];
	}
	c[i++] = NULL;
	for (int j = 0; b[j] != NULL; j++) {
		int flag = 0;
		for (int k = 0; a[k] != NULL; k++) {
			if (b[j] == a[k]) {
				flag = 1;
				break;
			}
		}
		if (flag) {
			int t = 0;
			while (c[t] != NULL) {
				if (c[t] == b[j]) {
					c[t] = -1;
				}
				t++;
			}
		}
	}
	int q = 0;
	for (int r = 0; c[r] != NULL; r++) {
		if (c[r] != -1) {
			d[q] = c[r];
			q++;
		}
	}
	d[q++] = NULL;
	display(d);
}
int main() {
	int e[100], b[100], a[100], c[100];
	int n;
	cout << "请输入全集的元素个数："<<endl;
	cin >> n;
rem:
	cout << "请输入全集元素：" << endl;
	for (int i = 0; i < n; i++) {
		cin >> e[i];
		if (!rejudge(e, i, e[i])) {
			cout << "元素重复，请重新输入";
			Sleep(5000);
			system("cls");
			goto rem;
		}
	}
	e[n] = NULL;
	cout << "请输入A子集："<<endl;
	ziji(e, a, n);
	cout << "请输入B子集" << endl;
	ziji(e, b, n);
	con(a, b, c);
	cout << endl;
	uni(a, b, c); 
	cout << endl;
	cout << "A和B集合的差集为：" << endl;
	diff(a, b, c);
	cout << endl;
	cout << "A集合的补集为：" << endl;
	diff(e, a, c);
}