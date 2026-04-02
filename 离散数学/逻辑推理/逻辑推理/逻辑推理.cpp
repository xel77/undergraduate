#include <iostream>
using namespace std;
//逻辑非
int neg(int a) {
    if (a==1) {
        return 0;
    }
    else {
        return  1;
    }
}
//合取
int con(int a, int b) {
    switch (a)
    {
    case 0:
        return 0;
    case 1:
        if (b == 0) {
            return 0;
        }
        return 1;
    }
}
//析取
int ext(int a, int b) {
    switch (a)
    {
    case 1:
        return 1;
    case 0:
        if (b == 0) {
            return 0;
        }
        return 1;
    }
}
//蕴含
int contain(int a, int b) {
    switch (a)
    {
    case 0:
        return 1;
    case 1:
        if (b == 0) {
            return 0;
        }
        return 1;
    }
}
//等价
int equ(int a, int b) {
    if (a == b) {
        return 1;
    }
    return 0;
}
int main() {
    char A, B, C;
    A = 'A'; // 副班长
    B = 'B'; // 正班长
    C = 'C'; // 团支书

    // 根据同学的意见和队长的决策规则进行判断
    if ((A == 'A' && B != 'B') || (A != 'A' && B == 'B')) {
        printf("M应该担任团支书（%c）\n", C);
    }
    else if ((A == 'A' && C != 'C') || (A != 'A' && C == 'C')) {
        printf("M应该担任正班长（%c）\n", B);
    }
    else if ((B == 'B' && C != 'C') || (B != 'B' && C == 'C')) {
        printf("M应该担任副班长（%c）\n", A);
    }
    else {
        printf("无法确定M应该担任哪个班干部职位\n");
    }

    return 0;
}
