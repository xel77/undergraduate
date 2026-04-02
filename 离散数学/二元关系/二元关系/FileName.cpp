#include<iostream>
using namespace std;
int flag = 1;
int judge(int a[][4]) {
	for (int i = 0; i < 4; i++) {
		if (a[i][i]!=1) {
			return 0;
		}
	}
	return 1;
}
int judge(int a[][4], int b) {
	for (int i = 0; i < 4; i++) {
		for (int j = 0; j < 4; j++) {
			if (a[i][j] != a[j][i]) {
				return 0;
			}
		}
	}
	return 1;
}
int judge(int a[][4], char b) {
	for (int i = 0; i < 4; i++) {
		for (int j = 0; j < 4; j++) {
			for (int k = 0; k < 4; k++) {
				if (a[i][j]==1 && a[j][k]==1&&a[i][k]!=1) {
					return 0;
				}
			}
			
		}
	}
	return 1;
}
int main() {
	int r[4][4] = { 1,0,0,0,
					0,1,1,0,
					0,0,1,1,
					0,1,0,1 };
	if (judge(r)) {
		cout << "R具有自反性" << endl;
	}
	else {
		cout << "R不具有自反性" << endl;
		flag = 0;
	}
	if (judge(r,1)) {
		cout << "R具有对称性" << endl;
	}
	else {
		cout << "R不具有对称性" << endl;
		flag = 0;
	}
	if (judge(r, '1')) {
		cout << "R具有传递性" << endl;
	}
	else {
		cout << "R不具有传递性" << endl;
		flag = 0;
	}
	if (flag) {
		cout << "是等价关系" << endl;
		cout<<"商集为:\n";
		int B[4][4] = { 0 };
		int num = 0;
		for (int i = 0; i < 8; i++) {
			int flag = 0;
			for (int j = 0; j < num; j++) {
				if (B[j][0] == r[i][0]) {
					B[j][1 + B[j][1]] = r[i][1];
					B[j][1] = B[j][1] + 1;
					flag = 1;
					break;
				}
			}
			if (!flag) {
				B[num][0] = r[i][0];
				B[num][1] = 1;
				B[num][2] = r[i][1];
				num++;
			}
		}
		for (int i = 0; i < num; i++) {
			printf("等价类 %d 为: {", i + 1);
			for (int j = 0; j < B[i][1]; j++) {
				printf("%d ", B[i][2 + j]);
			}
			printf("}\n");
		}
	}
	else {
		cout << "不是等价关系" << endl;
	}
}