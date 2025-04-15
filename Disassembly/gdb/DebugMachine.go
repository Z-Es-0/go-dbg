package gdb

import (
	"fmt"
	"syscall"
)

// Directive 反汇编指令结构体
type Directive struct {
	Length   uint32 // 指令长度
	Address  uint64 // 指令地址
	HexCodes []byte // 机器码
	Armcode  string // 汇编指令
	Comment  string // 注释
	//breakpoint bool // 是否设置断点
}

// Dbgbreak 增加原始字节字段
type Dbgbreak struct {
	address  uintptr
	original byte // 新增：保存被替换的原始字节
	rawcode  *Directive
}

type DbgMachine struct {
	process syscall.Handle // 进程句柄

	thread syscall.Handle // 线程句柄
	// 存储所有断点的切片
	breakpoints map[uintptr]*Dbgbreak

	textdata map[uintptr]*Directive // 存储所有指令的映射，key 为指令首地址，value 为指令结构体指针
}

// SetBreakpoint 为指定的指令设置断点。
// 该函数接收一个 uintptr 类型的参数，代表要设置断点的指令地址。
// 如果指令为空，则返回错误。
// 如果断点已经存在，则不做任何操作。
// 该函数会将断点信息添加到 DbgMachine 的 breakpoints 映射中，并将断点地址处的内存写入 0xCC（INT 3 指令）。
func (d *DbgMachine) SetBreakpoint(address uintptr) error {
	if d.process == 0 || d.thread == 0 {
		return fmt.Errorf("无效的调试器句柄")
	}

	// 检查断点是否已经存在于 breakpoints 映射中
	if _, exists := d.breakpoints[uintptr(address)]; exists {
		// 如果断点已经存在，直接返回 nil，表示不需要再次设置
		return nil
	}
	// 读取原始指令的机器码
	//oldCode := make([]byte, rawcode.Length)
	origBytes, err := ReadProcessMemory(d.process, uintptr(address), 1)

	if err != nil {
		return err
	}

	// 创建断点结构体
	breakpoint := &Dbgbreak{
		address:  address,
		original: origBytes[0],
		rawcode:  d.textdata[address],
	}

	// 写入 INT 3 指令 (0xCC)
	int3Code := []byte{0xCC}
	_, err = WriteProcessMemory(d.process, address, int3Code)
	if err == nil {
		// 将断点信息添加到 DbgMachine 的 breakpoints 映射中
		d.breakpoints[address] = breakpoint
		return nil
	}
	return err

}

// DeleteBreakpoint 为指定的指令删除断点。
// 该函数接收一个 *Directive 类型的参数，代表要删除断点的指令。
// 如果指令为空，则返回错误。
// 如果断点不存在，则不做任何操作。
// 该函数会将断点地址处的内存恢复为原始指令的机器码，并从 DbgMachine 的 breakpoints 映射中删除该断点信息。
func (d *DbgMachine) DeleteBreakpoint(address uintptr) error {

	// 检查断点是否已经存在于 breakpoints 映射中
	if _, exists := d.breakpoints[address]; !exists {
		// 如果断点不存在，直接返回 nil，表示不需要再次删除
		return nil
	}
	// 调用 WriteProcessMemory 函数将原始指令的机器码写入进程内存中的断点地址
	_, err := WriteProcessMemory(d.process, address, []byte{d.breakpoints[address].original})
	if err != nil {
		// 如果写入过程中出现错误，返回该错误
		return err
	}

	// 从 DbgMachine 的 breakpoints 映射中删除该断点信息
	delete(d.breakpoints, address)
	// 如果一切正常，返回 nil 表示删除断点成功
	return nil
}

// Maketextdata 读取线程上下文，获取指令指针 (RIP)，并从进程内存中读取指令。
// 它将读取的指令反汇编，并将结果存储在 DbgMachine 的 textdata 映射中。
// 如果获取线程上下文、读取内存或反汇编过程中出现错误，函数将返回相应的错误。
func (d *DbgMachine) Maketextdata() error {
	// 初始化 textdata 映射，如果还未初始化
	if d.textdata == nil {
		d.textdata = make(map[uintptr]*Directive)
	}

	// 获取线程上下文
	context, err := GetThreadContext(d.thread)
	if err != nil {
		fmt.Println("获取线程上下文失败:", err)
		return err
	}

	// 获取指令指针 (RIP)
	rip := context.Rip

	// 读取内存中的指令
	bytedata, err := ReadProcessMemory(d.process, uintptr(rip), 200)
	if err != nil {
		fmt.Println("读取内存失败:", err)
		return err
	}

	// 反汇编读取的指令
	data := DisassembleRange(bytedata, rip, 64) //  64 位架构

	// 遍历反汇编结果并存储到 textdata 映射中
	for _, directive := range *data {
		// 按16进制地址 - 16进制 - 汇编的格式打印
		hexCodesStr := ""
		for _, code := range directive.HexCodes {
			hexCodesStr += fmt.Sprintf("%02X ", code)
		}
		//fmt.Printf("%016X - %s - %s\n", directive.Address, hexCodesStr, directive.Armcode)

		d.textdata[uintptr(directive.Address)] = &directive
	}

	return nil
}

// FindBreakpoint 查找并处理当前线程指令指针处的断点。
// 该函数会获取当前线程的上下文，检查指令指针 (RIP) 前一个地址是否存在断点。
// 如果存在断点，它会修正线程上下文，将 RIP 恢复到断点前的位置，并删除该断点。
// 如果获取线程上下文或修改线程上下文时出现错误，函数将返回相应的错误。
func (d *DbgMachine) FindBreakpoint() (bool, error) {
	// 获取当前线程的上下文
	context, err := GetThreadContext(d.thread)
	if err != nil {
		// 若获取线程上下文失败，打印错误信息并返回错误
		fmt.Println("获取线程上下文失败:", err)
		return false, err
	}
	// 获取指令指针 (RIP)
	rip := context.Rip
	// 由于断点触发时 RIP 已经指向下一条指令，将 RIP 减 1 回到断点指令处
	rip = rip - 1
	// 检查断点是否存在于 breakpoints 映射中
	if _, exists := d.breakpoints[uintptr(rip)]; exists {
		// 若断点存在，修正线程上下文，将 RIP 恢复到断点前的位置
		err = ReviseThreadContext(d.thread, context, "Rip", rip)
		if err != nil {
			// 若修改线程上下文失败，打印错误信息并返回错误
			fmt.Println("修改线程上下文失败:", err)
			return false, err
		}
		// 删除该断点
		d.DeleteBreakpoint(uintptr(rip))
		return true, nil
	}
	// 若未找到断点，返回 nil 表示没有错误
	return false, nil
}

// Run 启动调试器，进入调试循环。
// 该函数会不断等待调试事件的发生，并根据事件类型进行相应的处理。
// 如果发生异常调试事件，则会打印异常信息，并停止调试器。
// 如果发生进程退出事件，则会打印进程退出信息，并停止调试器。
// 如果发生其他事件类型，则会打印事件类型，并继续等待下一个事件。
// 该函数会一直运行，直到遇到异常或进程退出事件。
func (d *DbgMachine) Run() {
	err := d.Maketextdata()
	if err != nil {
		fmt.Println("反汇编失败:", err)
	}

	threadID, err := GetThreadID(d.thread)
	if err != nil {
		fmt.Println("获取线程ID失败:", err)
	}
	processID, err := GetProcessID(d.process)
	if err != nil {
		fmt.Println("获取进程ID失败:", err)
	}

	debugEvent := &DEBUG_EVENT{
		DebugEventCode: 0,
		ThreadId:       threadID,
		ProcessId:      processID,
	}

	//time.Sleep(5 * time.Second)
	// 调试事件循环
	d.ShellMain()

	for {
		debugEvent, err = WaitForDebug(debugEvent)
		if err != nil {
			fmt.Println("等待调试事件失败:", err)
		}

		switch debugEvent.DebugEventCode {

		case EXCEPTION_DEBUG_EVENT:
			{
				_ = d.Maketextdata()

				switch (GetUnion[EXCEPTION_DEBUG_INFO](debugEvent)).ExceptionRecord.ExceptionCode {

				case EXCEPTION_ACCESS_VIOLATION:
					fmt.Println("内存访问冲突")

				case EXCEPTION_BREAKPOINT:
					//fmt.Println("断点触发")
					r, err := d.FindBreakpoint()
					if err != nil {
						fmt.Println("触发断点失败:", err)
					}
					if r {
						fmt.Println("断点触发")
					}

				case EXCEPTION_SINGLE_STEP:
					fmt.Println("单步执行异常") // 单步执行异常

				case EXCEPTION_GUARD_PAGE:
					fmt.Println("保护页异常") // 保护页异常

				case EXCEPTION_DATATYPE_MISALIGNMENT:
					fmt.Println("数据类型不匹配异常") // 数据类型不匹配异常

				case EXCEPTION_NONCONTINUABLE_EXCEPTION:
					fmt.Println("不可继续执行异常") // 不可继续执行异常

				case EXCEPTION_INT_DIVIDE_BY_ZERO:
					fmt.Println("整数除以零异常") // 整数除以零异常

				case EXCEPTION_INT_OVERFLOW:
					fmt.Println("整数溢出异常") // 整数溢出异常

				case EXCEPTION_PRIV_INSTRUCTION:
					fmt.Println("特权指令异常") // 特权指令异常

				case EXCEPTION_IN_PAGE_ERROR:
					fmt.Println("页面错误异常") // 页面错误异常

				case EXCEPTION_ILLEGAL_INSTRUCTION:
					fmt.Println("非法指令异常") // 非法指令异常

				default:
					fmt.Println((GetUnion[EXCEPTION_DEBUG_INFO](debugEvent)).ExceptionRecord.ExceptionCode)

					fmt.Println("其他异常")

				}

				d.ShellMain()
			}

		case CREATE_THREAD_DEBUG_EVENT:
			fmt.Println("线程创建")
			_ = d.Maketextdata()

		case LOAD_DLL_DEBUG_EVENT:
			fmt.Println("DLL加载")
			_ = d.Maketextdata()

		case EXIT_PROCESS_DEBUG_EVENT:
			fmt.Println("目标进程正常退出")
			return

		}

		ContinueDebugEvent(
			debugEvent.ProcessId,
			debugEvent.ThreadId,
			DBG_CONTINUE,
		)

	}

}
