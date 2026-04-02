// autostart.h

#ifndef AUTOSTART_H
#define AUTOSTART_H

#include <windows.h>
#include <stdio.h>

// 添加自启动项
// 参数：
//   szPath: 要添加到自启动项的程序路径
//   szName: 在注册表中的键名
// 返回值：
//   0 表示成功，其他值表示失败
int AddToAutoStart(const char* szPath, const char* szName) {
    HKEY hKey;
    wchar_t wPath[MAX_PATH];
    wchar_t wName[MAX_PATH];

    // 将 char* 转换为 wchar_t*
    MultiByteToWideChar(CP_ACP, 0, szPath, -1, wPath, MAX_PATH);
    MultiByteToWideChar(CP_ACP, 0, szName, -1, wName, MAX_PATH);

    // 检查注册表中是否已存在相同的自启动项
    if (RegOpenKeyExW(HKEY_LOCAL_MACHINE,
        L"Software\\Microsoft\\Windows\\CurrentVersion\\Run",
        0, KEY_READ, &hKey) == ERROR_SUCCESS) {
        // 检查注册表中是否已存在相同的键名
        DWORD valueType;
        DWORD dataSize = 0;
        if (RegQueryValueExW(hKey, wName, NULL, &valueType, NULL, &dataSize) == ERROR_SUCCESS) {
            // 如果键名已存在，则关闭注册表项，返回0
            printf("自启动项已存在，无需创建。\n");
            RegCloseKey(hKey);
            return 0;
        }
        RegCloseKey(hKey);
    }

    // 创建自启动项
    if (RegOpenKeyExW(HKEY_LOCAL_MACHINE,
        L"Software\\Microsoft\\Windows\\CurrentVersion\\Run",
        0, KEY_WRITE, &hKey) == ERROR_SUCCESS) {
        // 设置注册表值
        if (RegSetValueExW(hKey, wName, 0, REG_SZ, (BYTE*)wPath, (wcslen(wPath) + 1) * sizeof(wchar_t)) == ERROR_SUCCESS) {
            printf("程序添加到自启动项成功。\n");
            RegCloseKey(hKey);
            return 1; // 创建成功返回1
        }
        else {
            printf("设置注册表值失败。\n");
            RegCloseKey(hKey);
            return -1; // 创建失败返回-1
        }
    }
    else {
        printf("打开注册表项失败。\n");
        return -1; // 创建失败返回-1
    }
}


#endif // AUTOSTART_H
