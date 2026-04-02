#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>
#include <direct.h> // 用于创建目录
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>
// 定义一个函数，用于执行命令并将结果写入文件
void run_command(const char* command, const char* filename) {
    FILE* fp;
    errno_t err = fopen_s(&fp, filename, "a");
    if (err != 0) {
        printf("无法打开文件 %s\n", filename);
        return;
    }

    fprintf(fp, "---------- %s ----------\n", command); // 在文件中写入命令名称
    fflush(fp);

    // 使用 _popen() 函数捕获命令输出
    FILE* pipe = _popen(command, "r");
    if (pipe == NULL) {
        fclose(fp);
        return;
    }

    char buffer[1024];
    while (fgets(buffer, sizeof(buffer), pipe) != NULL) { // 逐行读取命令输出
        fprintf(fp, "%s", buffer); // 将命令输出写入文件
    }

    _pclose(pipe); // 关闭管道
    fclose(fp); // 关闭文件
}

// 定义一个函数，用于执行一系列收集系统信息的指令
void collect_system_info() {
  
    char current_dir[MAX_PATH];
    DWORD result = GetCurrentDirectoryA(MAX_PATH, current_dir);
    // 自动创建一个文件来保存系统信息
     // 拼接文件路径
    const char* filename = "info.txt";
    char filepath[MAX_PATH];
    strcpy_s(filepath, current_dir); // 将当前目录复制到 filepath
    strcat_s(filepath, "\\dist\\fileupload\\"); // 添加路径分隔符
    strcat_s(filepath, filename); // 添加文件名

    // 创建文件
    FILE* fp;
    errno_t err = fopen_s(&fp, filepath, "w");
    if (err != 0) {
        printf("无法创建文件 %s\n", filepath);
    }
    fclose(fp);

    printf("文件 %s 创建成功\n", filepath);

    // 执行命令并将信息写入文件
    run_command("ipconfig", filename); // 执行 ipconfig 命令
    run_command("systeminfo", filename); // 执行 systeminfo 命令
    run_command("netstat -ano", filename); // 执行 netstat -ano 命令
    run_command("tasklist /svc", filename); // 执行 tasklist /svc 命令
    run_command("net user", filename); // 执行 net user 命令
    run_command("net localgroup", filename); // 执行 net localgroup 命令
    run_command("net share", filename); // 执行 net share 命令
    //for (int i = 0; current_dir[i] != '\0'; i++) {
    //    if (current_dir[i] == '\\') {
    //        current_dir[i] = '/';
    //    }
    //}

    //// 拼接路径
    //char command[MAX_PATH + 50];  // 路径长度 + "/dist/fileupload/fileupload.exe" 字符串长度
    //sprintf_s(command, "%s/dist/fileupload/fileupload.exe", current_dir);

    //// 执行拼接后的命令
    //int status = system(command);
    //if (status == -1) {
    //    printf("无法执行命令\n");
    //}
    //else {
    //    printf("命令执行返回码: %d\n", status);
    //}

}
