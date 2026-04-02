#include<stdio.h>
typedef struct lnode {
	int date;
	struct lnode *next;
}lnode,*linklist;
int initlist(linklist* l) {
	(*l)->next = NULL;
	return 1;
}
void creat(linklist* l, int n) {
	l = new lnode;
	l->next = NULL;
	linklist p;
	initlist(&l);
	for (int i = n; i > 0; --i) {
		p = new lnode;
		scanf("%d", &l->date[i]);
		p->next = l->next;
		l->next - p;
	}
}
void display(linklist* l) {
	linklist p = l->next;
	while (p) {
		printf("%d", p->date[p]);
		p = p->next;
	}
}
int main() {
	linklist l;
	creat(&l,4);
	display(&l);
	return 0;
}