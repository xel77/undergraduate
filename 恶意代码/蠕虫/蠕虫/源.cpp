#include <stdio.h>
#include <winsock2.h>

#define SERVER_IP "127.0.0.1"
#define PORT 12345
#define BUFFER_SIZE 1024

int main(int argc, char* argv[]) {
    if (argc != 2) {
        printf("Usage: %s <file_path>\n", argv[0]);
        return 1;
    }

    char* file_path = argv[1];
    char buffer[BUFFER_SIZE];

    // 初始化Winsock
    WSADATA wsaData;
    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0) {
        printf("WSAStartup failed\n");
        return 1;
    }

    // 创建套接字
    SOCKET client_socket = socket(AF_INET, SOCK_STREAM, 0);
    if (client_socket == INVALID_SOCKET) {
        printf("Error in socket creation\n");
        WSACleanup();
        return 1;
    }

    // 设置服务器地址结构
    struct sockaddr_in server_addr;
    server_addr.sin_family = AF_INET;
    if (inet_pton(AF_INET, SERVER_IP, &server_addr.sin_addr) <= 0) {
        printf("Invalid address: %s\n", SERVER_IP);
        closesocket(client_socket);
        WSACleanup();
        return 1;
    }
    server_addr.sin_port = htons(PORT);

    // 连接到服务器
    if (connect(client_socket, (struct sockaddr*)&server_addr, sizeof(server_addr)) == SOCKET_ERROR) {
        printf("Error in connecting to server\n");
        closesocket(client_socket);
        WSACleanup();
        return 1;
    }

    // 发送文件名
    char* filename = strrchr(file_path, '\\');
    if (filename == NULL) {
        filename = file_path;
    }
    else {
        filename++;
    }
    if (send(client_socket, filename, strlen(filename), 0) == SOCKET_ERROR) {
        printf("Error in sending filename\n");
        closesocket(client_socket);
        WSACleanup();
        return 1;
    }

    // 打开文件并发送文件内容
    FILE* file = fopen(file_path, "rb");
    if (file == NULL) {
        printf("Error in opening file\n");
        closesocket(client_socket);
        WSACleanup();
        return 1;
    }
    int bytes_read;
    while ((bytes_read = fread(buffer, 1, BUFFER_SIZE, file)) > 0) {
        if (send(client_socket, buffer, bytes_read, 0) == SOCKET_ERROR) {
            printf("Error in sending file\n");
            fclose(file);
            closesocket(client_socket);
            WSACleanup();
            return 1;
        }
    }
    if (bytes_read == 0) {
        printf("File sent successfully\n");
    }
    else {
        printf("Error in reading file\n");
    }

    fclose(file);
    closesocket(client_socket);
    WSACleanup();

    return 0;
}
