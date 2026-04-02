#include <stdio.h>
#include <stdlib.h>

typedef struct {
    int maxKey;
    int start;
} Index;

int compare(const void *a, const void *b) {
    return (*(int*)a - *(int*)b);
}

int BlockSearch(int *num, Index *index, int numLen, int indexLen, int key) {
    int low = 0, high = indexLen - 1;
    int mid;
    while (low <= high) {
        mid = low + (high - low) / 2;
        if (key <= index[mid].maxKey) {
            break;
        } else {
            low = mid + 1;
        }
    }

    int start = index[mid].start;
    int end = (mid == indexLen - 1) ? numLen : index[mid + 1].start;

    for (int i = start; i < end; i++) {
        if (key == num[i]) {
            return i;
        }
    }
    return -1;
}

int main() {
    int size;
    printf("请输入表的长度：\n");
    scanf("%d", &size);
    int *a = (int*)malloc(size * sizeof(int));
    printf("请输入数据：\n");
    for (int j = 0; j < size; j++) {
        scanf("%d", &a[j]);
    }

    qsort(a, size, sizeof(int), compare);

    Index index[3];
    index[0].maxKey = a[size / 3];
    index[0].start = 0;
    index[1].maxKey = a[size * 2 / 3];
    index[1].start = size / 3;
    index[2].maxKey = a[size - 1];
    index[2].start = size * 2 / 3;

    int key;
    printf("请输入要查找的数：\n");
    scanf("%d", &key);
    int pos = BlockSearch(a, index, size, 3, key);
    if (pos != -1) {
        printf("查找成功，位置 = %d\n", pos);
    } else {
        printf("未找到该数\n");
    }

    free(a);
    return 0;
}

