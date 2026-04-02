#include<iostream>
#include<cstdlib>
#include<ctime>
using namespace std;
//定义双向链表结构体
typedef struct node {
    int data;
    struct node* last;
    struct node* next;
} node;
//定义头节点
struct node* head;
//初始化头节点
void init() {
    head = new node();
    head->next = NULL;
    head->last = NULL;
    head->data = NULL;
}
//尾插法插入值
node* insert(node* p,int t) {
    if (t == 1) {
        node* l = new node();
        l->data= rand() % 100 + 1;
        l->next = NULL;
        l->last = p;
        p->next = l;
        return l;
    }
    else if(t==2) {
        cout << "请输入值：";
        node* l = new node();
        cin >> l->data;
        l->next = NULL;
        l->last=p;
        p->next = l;
        return l;
    }
}
//输出链表元素
void display() {
    node* l = head->next;
     do{
        cout << l->data<<"\t";
        l = l->next;
     } while (l != NULL);
}
//插入排序
int insertSort(int* e) {
    int count = 0;
    // 从头节点的下一个节点开始遍历
    node* l = head->next;
    while (l != NULL) {
        // 保存当前节点
        node* current = l;
        l = l->next;
        (*e)++;
        // 如果当前节点的下一个节点存在并且当前节点的数据大于下一个节点的数据
        if (current->next != NULL && current->data > current->next->data) {
            count++;
            // 保存下一个节点
            node* p = current->next;
            // 重新连接当前节点的next指针
            current->next = p->next;
            // 如果下一个节点的下一个节点存在，更新其last指针
            if (p->next != NULL) {
                p->next->last = current;
            }
            // 找到合适的位置插入节点
            node* t = head;
            while (t->next != NULL && t->next->data < p->data) {
                t = t->next;
                count++;
                (*e)++;
            }
            // 插入节点
            p->next = t->next;
            if (t->next != NULL) {
                t->next->last = p;
            }
            t->next = p;
            p->last = t;
        }
    }
    return count;
}

//冒泡排序
int bubbleSort(int* t) {
    int count = 0; 
    node* l = head->next;   
    for (int i = 0; i < head->data; i++) {
        node* p = l->next;
        for (int j = i + 1; j < head->data; j++) {
            if (l->data > p->data) {
                count++;
                // 交换节点数据
                int temp = l->data;
                l->data = p->data;
                p->data = temp;
            }
            p = p->next;
            (*t)++;
        }
        l = l->next;
    }
    return count;  
}

//选择排序
int chooseSort(int*e) {
    int count = 0;
    node*l = head->next;
    node* min;
    for (int i = 0; i < head->data; i++) {
        min = l;
        node* p = l->next;
        for (int j = i + 1; j < head->data; j++) {
            if (min->data > p->data) {
                count++;
                min = p;
            }
            (*e)++;
            p = p->next;
        }
        if (l != min) {
            //交换数据
            int t = min->data;
            min->data = l->data;
            l->data = t;
        }
        
    }
    return count;
}
//创建排序表
void creat() {
    int n,t,h,z;
    cout << "----------------\n";
    cout << "| 1. 随机输入   |\n";
    cout << "| 2. 手动输入   |\n";
    cout << "| 0. 退出       |\n";
    cout << "----------------\n";
    cout << "请选择排序算法（输入0退出）：";
    cin >> z;
    if(z == 1){
        srand((unsigned)time(NULL));
        int r = rand() % 10 + 1;
        head->data = r;
        node* p = head;
        for (int i = 0; i < r; i++) {
            p = insert(p,z);
        }
    }
    else if(z == 2) {
        cout << "请输入要输入的个数：";
        cin >> n;    
        head->data = n;
        node* p = head;
        for (int i = 0; i < n; i++) {
            p = insert(p,z);
        }
    }
    else if (z == 0) {
        return;
    }
    else {
        cout << "错误";
        return;
    }
    //system("cls");
    cout << "排序前的顺序：" ;
    display();
    cout << endl;
    cout << "--------------------\n";
    cout << "| 1. 插入排序      |\n";
    cout << "| 2. 冒泡排序      |\n";
    cout << "| 3. 选择排序      |\n";
    cout << "| 0. 退出          |\n";
    cout << "--------------------\n";
    cout<<"请选择排序算法（输入0退出）：";
    cin>>t;int e = 0;
    switch (t) {
        case 1:
            h = insertSort(&e);
            cout << "排序后：";
            display();
            cout << "\n" << "总共移动" << h << "次" ;
            cout << "\n" << "总共比较" << e-1 << "次" ;
            return;
        case 2:
            h = bubbleSort(&e);
            cout << "排序后：";
            display();
            cout << "\n" << "总共移动" << h << "次" << endl;
            cout << "\n" << "总共比较" << e << "次" ;
            return;
        case 3:
            h = chooseSort(&e);
            cout << "排序后：";
            display();
            cout << "\n" << "总共移动" << h << "次" << endl;
            cout << "\n" << "总共比较" << e << "次" ;
            return;
        case 0:
            break;
        default:
            cout << "错误\n";
            return;
        }
}

int main(){
    init();
    creat();
}
