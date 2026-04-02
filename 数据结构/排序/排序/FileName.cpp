#include <stdio.h>
#include <stdlib.h>
//节点
typedef struct {
    int data;
} elem;
//表
typedef struct {
    elem* r;
    int length;
} table;
//输出序列
void display(table *a) {
    static int t = 1;
    printf("第%d躺排序结果：", t);
    for (int i = 0; i < a->length; i++) {
        printf("%d\t", a->r[i].data);
    }
    t++;
}
//插入排序
void insertSort(table* a) {
    int t;
    printf("\n请输入元素个数: ");
    scanf_s("%d", &t);
    a->r = (elem*)malloc(t+1 * sizeof(elem));
    a->length = t+1;
    printf("请输入元素: \n");
    for (int i = 1; i < a->length; i++) {
        scanf_s("%d", &(a->r[i].data));
    }
    int k;
    for (int j = 2; j < a->length; j++) {
        if (a->r[j].data < a->r[j -1].data) {
            a->r[0].data = a->r[j].data;
            for ( k= j - 1; a->r[k].data > a->r[0].data; k--) {
                a->r[k + 1].data = a->r[k].data;
            }
            a->r[k + 1].data = a->r[0].data;
        }
        printf("第%d趟排序结果：", j - 1);
        for (int l = 1; l < a->length; l++) {
            printf("%d\t", a->r[l].data);
        }
        printf("\n");
    }
}
//冒泡排序
void bubbleSort(table*a) {
    int t;
    printf("\n请输入元素个数: ");
    scanf_s("%d", &t);
    a->r = (elem*)malloc(t * sizeof(elem));
    a->length = t ;
    printf("请输入元素: \n");
    for (int i = 0; i < a->length; i++) {
        scanf_s("%d", &(a->r[i].data));
    }
    int k, temp;
    for (int j = 0; j <= a->length; j++) {
        for (k = j + 1; k < a->length; k++) {
            if (a->r[j].data > a->r[k].data) {
                temp = a->r[j].data;
                a->r[j].data = a->r[k].data;
                a->r[k].data = temp;
            }
        }
        printf("第%d趟排序结果：", j+1);
        for (int l = 0; l < a->length; l++) {
            printf("%d\t", a->r[l].data);
            }
        printf("\n");
    }
}
//快速排序
int partition(table *a, int low, int high) {
    int pivot = a->r[high].data;
    int i = low - 1;

    for (int j = low; j <= high - 1; j++) {
        if (a->r[j].data < pivot) {
            i++;
            int temp = a->r[i].data;
            a->r[i].data = a->r[j].data;
            a->r[j].data = temp;
        }
    }
    int temp = a->r[i + 1].data;
    a->r[i + 1].data = a->r[high].data;
    a->r[high].data = temp;
    return (i + 1);
}
int count = 0;
void quick_sort(table* a,int low,int high) {
    
    if (low < high) {
        count++;
        int pi = partition(a, low, high);
        printf("第%d趟排序结果：", count);
        for (int i = 0; i < a->length; i++) {
            printf("%d ", a->r[i].data);
        }
        printf("\n");
        quick_sort(a,low,pi-1);
        quick_sort(a,pi+1,high);
    }
}
//快排元素输入
void quickSort(table* a) {
    int t;
    printf("\n请输入元素个数: ");
    scanf_s("%d", &t);
    a->r = (elem*)malloc(t * sizeof(elem));
    a->length = t;
    printf("请输入元素: \n");
    for (int i = 0; i < a->length; i++) {
        scanf_s("%d", &(a->r[i].data));
    }
    int low = 0, high = a->length - 1;
    quick_sort(a, low, high);
    free(a->r);
}
//选择排序
void chooseSort(table* a) {
    int t;
    printf("\n请输入元素个数: ");
    scanf_s("%d", &t);
    a->r = (elem*)malloc(t * sizeof(elem));
    a->length = t;
    printf("请输入元素: \n");
    for (int i = 0; i < a->length; i++) {
        scanf_s("%d", &(a->r[i].data));
    }
    int min,temp;
    for (int j = 0; j < a->length; j++) {
        min = j;
        for (int k = j + 1; k < a->length; k++) {
            if (a->r[min].data > a->r[k].data) {
                min = k;
            }
        }
        if (min != j) {
            temp = a->r[min].data;
            a->r[min].data = a->r[j].data;
            a->r[j].data = temp;
        }
        printf("第%d趟排序结果：", j + 1);
        for (int l = 0; l < a->length; l++) {
            printf("%d\t", a->r[l].data);
        }
        printf("\n");
    }
}
//堆排序
void adjustHeap(table* a, int i, int length) {
    int temp = a->r[i].data;
    for (int k = i * 2 + 1; k < length; k = k * 2 + 1) {
        if (k + 1 < length && a->r[k].data < a->r[k + 1].data) {
            k++;
        }
        if (temp > a->r[k].data) { 
            break;
        }
        a->r[i].data = a->r[k].data;
        i = k;
    }
    a->r[i].data = temp;
}
void heapSort(table* a) {
    int t;
    printf("\n请输入元素个数: ");
    scanf_s("%d", &t);
    a->r = (elem*)malloc(t * sizeof(elem));
    a->length = t;
    printf("请输入元素: \n");
    for (int i = 0; i < a->length; i++) {
        scanf_s("%d", &(a->r[i].data));
    }
    for (int i = t / 2 - 1; i >= 0; i--) {
        adjustHeap(a, i, t);
    }
    for (int j = t - 1; j > 0; j--) {
        int temp = a->r[0].data;
        a->r[0].data = a->r[j].data;
        a->r[j].data = temp;
        adjustHeap(a, 0, j);
        printf("第%d趟排序结果：", t - j);
        for (int l = 0; l < a->length; l++) {
            printf("%d\t", a->r[l].data);
        }
        printf("\n");
    }
}
//归并排序左右部分比较
void merge(table* result, table* left, table* right) {
    int i = 0, j = 0, k = 0;
    while (i < left->length && j < right->length) {
        if (left->r[i].data <= right->r[j].data) {
            result->r[k++] = left->r[i++];
        }
        else {
            result->r[k++] = right->r[j++];
        }
    }
    while (i < left->length) {
        result->r[k++] = left->r[i++];
    }
    while (j < right->length) {
        result->r[k++] = right->r[j++];
    }
    display(result);
    printf("\n");
}
//归并排序表分块
void mergesort(table* a) {
    if (a->length <= 1) {
        return;
    }
    int mid = a->length / 2;
    table left, right;
    left.r = (elem*)(malloc(mid * sizeof(elem)));
    if (left.r == NULL) {
        printf("内存分配失败");
        return;
    }
    left.length = mid;
    right.r = (elem*)(malloc((a->length - mid) * sizeof(elem)));
    if (right.r == NULL) {
        printf("内存分配失败");
        free(left.r);
        return;
    }
    right.length = a->length - mid;
    for (int i = 0; i < mid; i++) {
        left.r[i] = a->r[i];
    }
    for (int i = mid; i < a->length; i++) {
        right.r[i - mid] = a->r[i];
    }
    mergesort(&left);
    mergesort(&right);
    merge(a, &left, &right);
    free(left.r);
    free(right.r);
}
//归并排序数据输入
void mergeSort(table* a) {
    int t;
    printf("\n请输入元素个数: ");
    scanf_s("%d", &t);
    a->r = (elem*)malloc(t * sizeof(elem));
    if (a->r == NULL) {
        printf("内存分配失败");
        return ;
    }
    a->length = t;
    printf("请输入元素: \n");
    for (int i = 0; i < a->length; i++) {
        scanf_s("%d", &(a->r[i].data));
    }
    mergesort(a);
    free(a->r);
}

int main() {
    int t;
    table a;
    do {
        printf("--------------------\n");
        printf("| 1. 插入排序      |\n");
        printf("| 2. 冒泡排序      |\n");
        printf("| 3. 快速排序      |\n");
        printf("| 4. 选择排序      |\n");
        printf("| 5. 堆排序        |\n");
        printf("| 6. 归并排序      |\n");
        printf("| 0. 退出          |\n");
        printf("--------------------\n");
        printf("请选择排序算法（输入0退出）：");
        scanf_s("%d", &t);
        switch (t) {
        case 1:
            insertSort(&a);
            break;
        case 2:
            bubbleSort(&a);
            break;
        case 3:
            quickSort(&a);
            break;
        case 4:
            chooseSort(&a);
            break;
        case 5:
            heapSort(&a);
            break;
        case 6:
            mergeSort(&a);
            break;
        case 0:
            break;
        default:
            printf("错误\n");
            break;
        }
    } while (t != 0);
    return 0;
}
