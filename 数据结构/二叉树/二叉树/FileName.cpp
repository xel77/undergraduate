
#include <stdio.h>
#include <stdlib.h>
typedef struct BiNode {
	int data;
	struct BiNode *lchild, *rchild; 
} BiNode, *BiTree;

BiTree CreateLink() {
	int data;
	scanf("%d", &data);
	if (data == -1) {
		return NULL;
	} else {
		BiTree node = (BiTree)malloc(sizeof(BiNode));
		node->data = data;
		printf("请输入%d的左子树(输入-1结束)： ", data);
		node->lchild = CreateLink();
		printf("请输入%d的右子树： ", data);
		node->rchild = CreateLink();
		return node;
	}
}

int countnodes(BiNode *root) {
	if (root == NULL) {
		return 0;
	}
	if (root->lchild == NULL && root->rchild == NULL) {
		return 1;
	}
	return countnodes(root->lchild) + countnodes(root->rchild);
}

void Createtree(BiTree T, BiTree *L) {
	if (T == NULL) return;
	Createtree(T->lchild, L);
	*L = T;
	(*L)->lchild = NULL;
	(*L)->rchild = NULL;
	Createtree(T->rchild, L);
}

// 先序遍历二叉链表
void PreOrderTraverse(BiTree T) {
	if (T == NULL) {
		return;
	}
	printf("%d ", T->data);
	PreOrderTraverse(T->lchild);
	PreOrderTraverse(T->rchild);
}

// 中序遍历二叉链表
void InOrderTraverse(BiTree T) {
	if (T == NULL) {
		return;
	}
	InOrderTraverse(T->lchild);
	printf("%d ", T->data);
	InOrderTraverse(T->rchild);
}

// 后序遍历二叉链表
void PostOrderTraverse(BiTree T) {
	if (T == NULL) {
		return;
	}
	PostOrderTraverse(T->lchild);
	PostOrderTraverse(T->rchild);
	printf("%d ", T->data);
}

int main() {
	BiTree L;
	printf("输入节点的数据：\n");
	L = CreateLink();
	BiTree *root = (BiTree *)malloc(sizeof(BiTree));
	Createtree(L, root);
	free(root); // 释放分配的内存
	// 分别输出先序、中序和后序遍历的结果
	printf("先序遍历：");
	PreOrderTraverse(L);
	printf("\n中序遍历：");
	InOrderTraverse(L);
	printf("\n后序遍历：");
	PostOrderTraverse(L);
	printf("叶结点个数： %d", countnodes(L)); // 计算叶结点个数
	return 0;
}
