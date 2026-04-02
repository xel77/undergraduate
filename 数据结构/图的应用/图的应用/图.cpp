#include<iostream>
#include<queue>
using namespace std;

#define MVNum 100              // 最大顶点数 

typedef struct {
    string pointname[MVNum];    // 顶点表
    int side[MVNum][MVNum];     // 邻接矩阵
    int vexnum, arcnum;         // 图的当前点数和边数
}graph;

void init(graph* a) {
    for (int i = 0; i < MVNum; i++) {
        for (int j = 0; j < MVNum; j++) {
            a->side[i][j] = 0;  // 0表示没有边相连
        }
    }
    a->arcnum = 0;
    a->vexnum = 0;
}

void creat(graph* a) {
    int i, j, m, n;
    cout << "请输入顶点数和边数\n";
    cin >> a->vexnum >> a->arcnum;

    for (i = 0; i < a->vexnum; i++) {
        cout << "请输入顶点" << i << "的名称：" << endl;
        cin >> a->pointname[i];
    }

    for (j = 0; j < a->arcnum; j++) {
        cout << "请输入vi-vj的边的下标i，j：" << endl;
        cin >> m >> n;
        a->side[m][n] = 1;
        a->side[n][m] = 1;  // 无向图，需要对称
    }
}

void BFS(graph* a, int start) {
    bool visited[MVNum] = { false };
    queue<int> q;

    visited[start] = true;
    q.push(start);

    while (!q.empty()) {
        int node = q.front();
        q.pop();
        cout << a->pointname[node] << " ";

        for (int i = 0; i < a->vexnum; i++) {
            if (a->side[node][i] == 1 && !visited[i]) {
                visited[i] = true;
                q.push(i);
            }
        }
    }
}

int main() {
    graph g;
    init(&g);
    creat(&g);
    cout << "广度优先搜索结果：" << endl;
    BFS(&g, 0);
    return 0;
}
