#include<iostream>
using namespace std;
void dis(int a[][100],int n) {
	for (int i = 0; i < n; i++) {
		for (int j = 0; j < n; j++) {
			cout << a[i][j];
		}
		cout << endl;
	}

}
void zifan(int a[][100], int n) {
	int b[100][100];
	for (int i = 0; i < n; i++)
	{
		for (int j = 0; j < n; j++)
		{
			b[i][j] = a[i][j];
		}
	}
	for (int i = 0; i < n; i++) {
		b[i][i] = 1;
	}
	cout << "自反闭包是：" << endl;
	dis(b, n);
}
void duichen(int a[][100], int n) {
	int b[100][100];
	for (int i = 0; i < n; i++)
	{
		for (int j = 0; j < n; j++)
		{
			b[i][j] = a[i][j];
		}
	}
	for (int i = 0; i < n; i++)
	{
		for (int j = 0; j < n; j++)
		{
			if (b[i][j] == 1)
				b[j][i] = 1;
		}
	}
	cout << "对称闭包是：" << endl;
	dis(b, n);
}
void chuandi(int a[][100], int n) {
	int b[100][100];
	for (int i = 0; i < n; i++)
	{
		for (int j = 0; j < n; j++)
		{
			b[i][j] = a[i][j];
		}
	}
	for (int i = 0; i < n; i++)
	{
		for (int j = 0; j < n; j++)
		{
			if (b[i][j] == 1) {
				if (i != j) {
					for (int k = 0; k < n; k++) {
						if (b[j][k] == 1) {
							if (!(j==k)) {
								b[i][k] =1;
							}
						}
					}
				}
			}
		}
	}
	cout << "传递闭包是：" << endl;
	dis(b, n);
}
int main() {
	int a[100][100];
	int n;
	cout << "请输入关系矩阵阶数：";
	cin >> n; cout << "请输入关系矩阵：" << endl;;
	for (int i = 0; i < n; i++) {
		for (int j = 0; j < n; j++) {
			cin >> a[i][j];
		}
	}
	cout << "原关系矩阵是：" << endl;
	dis(a, n);
	zifan(a, n);
	duichen(a, n);
	chuandi(a, n);
}