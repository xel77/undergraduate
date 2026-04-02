#include<iostream>
using namespace std;

struct graph {
    int vexnum, arcnum;
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
    while (a->linjie[i][i] == 0) {
        a->vexname[i] = 0;
        i++;
    }
}

int judge(graph* a) {
    for (int i = 0; i < a->vexnum; i++) {
        int indegree = 0;
        int outdegree = 0;
        for (int j = 0; j < a->vexnum; j++) {
            if (a->linjie[i][j] == 1) {
                outdegree++;
            }
            if (a->linjie[j][i] == 1) {
                indegree++;
            }
        }
        if (indegree != outdegree) {
            return 0;
        }
    }
    return 1;
}

void creat(graph* a) {
    cout << "请输入顶点数：";
    cin >> a->vexnum;
    for (int i = 0; i < a->vexnum; i++) {
        cout << "请输入第" << i + 1 << "个顶点名：";
        cin >> a->vexname[i];
    }
    cout << "请输入边数：";
    cin >> a->arcnum;

    int m, n;
    cout << "请选择图的类型：1.有向图 2.无向图" << endl;
    int type;
    cin >> type;

    for (int j = 0; j < a->arcnum; j++) {
        cout << "请输入vi-vj：";
        cin >> m >> n;
        if (type == 1) {
            a->linjie[m - 1][n - 1] = 1;
        }
        else if (type == 2) {
            a->linjie[m - 1][n - 1] = a->linjie[n - 1][m - 1] = 1;
        }
    }

    if (judge(a)) {
        cout << "是欧拉图" << endl;
    }
    else {
        cout << "不是欧拉图" << endl;
    }
}

int main() {
    graph a;
    init(&a);
    creat(&a);
}
