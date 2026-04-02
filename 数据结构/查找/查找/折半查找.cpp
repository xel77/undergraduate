#include<iostream>
using namespace std;
#define MAX 100
typedef struct {
	int key;
} elem;
typedef struct {
	elem* r;
	int length;
} table;

void create(table** st, int length) {
	*st = new table;
	(*st)->length = length ;
	(*st)->r = new elem[length];
	cout << "输入表中的数据元素：" << endl;
	for (int i = 0; i < (*st)->length; i++) {
		cin >> (*st)->r[i].key;
	}
}
int middlesearch(table* st) {
	int key;
	cout << "请输入要查找的数：" << endl;
	cin >> key;
	int low = 0;
	int high = st->length - 1;
	int m = (low + high) / 2;
	while (low <= high) {
		if (key > st->r[m].key) {
			low = m+1;
			m = (low + high) / 2;
		}
		else if (key < st->r[m].key) {
			high = m-1;
			m = (low + high) / 2;
		}
		else {
			return m+1;
		}
	}
	cout << "不存在";
	return 0;
}
int main() {
	table* s;

	int n;
	cout << "请输入表的长度：" << endl;
	cin >> n;
	create(&s, n);
	int i = middlesearch(s);
	if (i != 0) {
		cout << "存在，在第" << i << "个位置" << endl;
	}
	delete[] s->r;
	delete s;
}
