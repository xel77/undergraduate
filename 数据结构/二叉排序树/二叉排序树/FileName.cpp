#include <stdio.h>
#include <stdlib.h>

// 定义二叉排序树节点
struct TreeNode {
    int value;
    struct TreeNode* left;
    struct TreeNode* right;
};

// 插入元素
struct TreeNode* insert(struct TreeNode* node, int value) {
    if (node == NULL) {
        struct TreeNode* new_node = (struct TreeNode*)malloc(sizeof(struct TreeNode));
        new_node->value = value;
        new_node->left = NULL;
        new_node->right = NULL;
        return new_node;
    }
    if (value < node->value) {
        node->left = insert(node->left, value);
    }
    else {
        node->right = insert(node->right, value);
    }
    return node;
}

// 显示二叉排序树（中序遍历）
void display(struct TreeNode* node) {
    if (node) {
        display(node->left);
        printf("%d ", node->value);
        display(node->right);
    }
}

// 查找元素
struct TreeNode* search(struct TreeNode* node, int value) {
    if (node == NULL || node->value == value) {
        return node;
    }
    if (value < node->value) {
        return search(node->left, value);
    }
    else {
        return search(node->right, value);
    }
}

// 查找最小值节点
struct TreeNode* findmin(struct TreeNode* node) {
    while (node->left) {
        node = node->left;
    }
    return node;
}

// 删除元素
struct TreeNode* del(struct TreeNode* node, int value) {
    if (node == NULL) {
        return node;
    }
    if (value < node->value) {
        node->left = del(node->left, value);
    }
    else if (value > node->value) {
        node->right = del(node->right, value);
    }
    else {
        if (node->left == NULL) {
            struct TreeNode* temp = node->right;
            free(node);
            return temp;
        }
        else if (node->right == NULL) {
            struct TreeNode* temp = node->left;
            free(node);
            return temp;
        }
        struct TreeNode* temp = findmin(node->right);
        node->value = temp->value;
        node->right = del(node->right, temp->value);
    }
    return node;
}

int main() {
    struct TreeNode* root = NULL;
    int value;

    do {
        printf("请输入要插入的值(0结束): ");
        scanf_s("%d", &value);
        if (value == 0) {
            break;
        }
        root = insert(root, value);
    } while (1);

    printf("\n二叉排序树:\n");
    display(root);
    printf("\n");

    int search_value;
    printf("请输入要查找的数: ");
    scanf_s("%d", &search_value);
    struct TreeNode* result = search(root, search_value);
    if (result) {
        printf("存在\n");
    }
    else {
        printf("不存在\n");
    }

    int delete_value;
    printf("请输入要删除的元素: ");
    scanf_s("%d", &delete_value);
    root = del(root, delete_value);

    printf("删除后的二叉树:\n");
    display(root);

    return 0;
}
