/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-26 14:20:50
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-03-26 14:31:26
 * @FilePath: \ZesOJ\Disassembly\gdb\debugger_win.cpp
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
#include <windows.h>
#include <iostream>
#include <vector>

/**
 * @brief ，用于创建一个新进程并以调试模式启动，使其在启动时暂停。
 * 
 * 此函数使用Windows API的CreateProcessW来创建一个新进程，并将其置于调试模式下。
 * 调试模式允许调用者在进程启动时暂停它，以便进行调试操作。
 * 
 * @param exePath 要执行的可执行文件的路径。
 * @param cmdLine 传递给可执行文件的命令行参数。
 * @param hProcess 用于接收新进程句柄的指针。
 * @param hThread 用于接收新进程主线程句柄的指针。
 * @return 如果进程创建成功，则返回true；否则返回false，并输出错误信息。
 */
bool CreateAndBlockProcess(
    LPCWSTR exePath, 
    LPCWSTR cmdLine, 
    HANDLE* hProcess, 
    HANDLE* hThread) {

    // 初始化STARTUPINFOW结构体，该结构体包含新进程的窗口信息
    STARTUPINFOW si = { sizeof(STARTUPINFOW) };
    // 初始化PROCESS_INFORMATION结构体，该结构体包含新进程和其主线程的信息
    PROCESS_INFORMATION pi = {0};

    // 调用CreateProcessW函数创建新进程
    if (!CreateProcessW(
        exePath,
        const_cast<LPWSTR>(cmdLine), // 命令行参数
        NULL,    // 进程安全属性
        NULL,    // 线程安全属性
        FALSE,   // 不继承句柄
        DEBUG_PROCESS, // 调试模式标志，使进程在启动时暂停
        NULL,    // 环境变量块
        NULL,    // 当前目录
        &si,
        &pi)) 
    {
        // 如果CreateProcessW失败，输出错误信息
        std::cerr << "CreateProcess failed. Error: " << GetLastError() << std::endl;
        return false;
    }

    // 将新进程的句柄赋值给传入的指针
    *hProcess = pi.hProcess;
    // 将新进程主线程的句柄赋值给传入的指针
    *hThread = pi.hThread;
    return true;
}




/**
 * @brief 从指定进程的内存中读取数据，对应Go的ReadProcessMemory。
 * 
 * 此函数使用Windows API的ReadProcessMemory从指定进程的内存中读取指定大小的数据。
 * 读取的数据存储在一个std::vector<BYTE>中，并返回该向量。
 * 
 * @param hProcess 要读取内存的进程的句柄。
 * @param address 要读取的内存地址。
 * @param size 要读取的字节数。
 * @return 包含读取数据的std::vector<BYTE>。如果读取失败，返回空向量。
 */
std::vector<BYTE> ReadProcessMemory(HANDLE hProcess, LPVOID address, SIZE_T size) {
    // 创建一个大小为size的BYTE向量，用于存储读取的数据
    std::vector<BYTE> buffer(size);
    // 用于存储实际读取的字节数
    SIZE_T bytesRead = 0;

    // 调用ReadProcessMemory函数从指定进程的内存中读取数据
    if (!::ReadProcessMemory(
        hProcess,
        address,
        buffer.data(),
        size,
        &bytesRead)) 
    {
        // 如果ReadProcessMemory失败，输出错误信息
        std::cerr << "ReadProcessMemory failed. Error: " << GetLastError() << std::endl;
        // 返回空向量
        return {};
    }

    // 调整向量的大小为实际读取的字节数
    buffer.resize(bytesRead);
    // 返回包含读取数据的向量
    return buffer;
}
