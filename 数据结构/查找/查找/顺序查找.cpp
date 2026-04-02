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
	(*st)->length = length+1;
	(*st)->r = new elem[length];
	cout << "输入表中的数据元素：" << endl;
	// 根据查找表中数据元素的总长度，在存储时，从数组下标为 1 的空间开始存储数据
	for (int i = 1; i < (*st)->length; i++) {
		cin >> (*st)->r[i].key;
	}
}

int search(table* st) {
	elem k;
	cout << "请输入你要查找的数：" << endl;
	cin >> k.key;
	st->r[0].key = k.key;//将关键字作为一个数据元素存放到查找表的第一个位置，起监视哨的作用
	int i = st->length;
	//从查找表的最后一个数据元素依次遍历，一直遍历到数组下标为0
	while (st->r[i].key != st->r[0].key) {
		i--;
	}
	//如果 i=0，说明查找失败；反之，返回的是含有关键字key的数据元素在查找表中的位置
	return i;
}
int main(){
	table* s; 
	
	int n;
	cout << "请输入表的长度："<<endl;
	cin >> n;
	create(&s, n);
	int i=search(s);
	if (i) {
		cout << "存在，在表第" << i<< "个位置" << endl;;
	}
	else {
		cout << "不存在" << endl;
	}
	delete[] s->r;
	delete s;
}
