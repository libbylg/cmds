package cmds

import (
	"container/list"
	"fmt"
	"os"
)

//	Cmd	定义聊所有命令的统一接口
//	Name	函数用来获取本Cmd的主名称和别名，第一个字符串为主名称。比如，下面的这个Name()函数的实现
//	中，"check"为主名称，它还有两个别名"--check"和"-c"。
//
//		func (c MyCmd) Name []string {
//			return []string{"check", "--check", "-c"}
//		}
//
//	Help	函数用来获取本Cmd的帮助摘要和帮助详情，其中第一个字符串为帮助摘要。
//	Exec	函数用来直接执行命令行指令。
type Cmd interface {
	Name() []string
	Help() []string
	Exec(args []string) (int, error)
}

var (
	CmdList *list.List
	CmdsMap map[string]*list.Element
	HelpCmd Cmd
	UnspCmd Cmd
	MispCmd Cmd
)

func init() {
	Clear()
}

//	Clear	用来清除所有已经注册的指令。
//	执行完Clear函数之后，所有的全局变量将被重置。在Clear()函数调用之前的获取的指令将失效。
func Clear() {
	CmdList = list.New()
	CmdsMap = make(map[string]*list.Element)
	HelpCmd = new(cmdHelp)
	UnspCmd = new(cmdUnsp)
	MispCmd = new(cmdMisp)
}

//	Reg	注册一个指令对象
func Reg(c Cmd) {
	names := c.Name()
	if len(names) <= 0 {
		return
	}

	elem := CmdList.PushBack(c)
	for _, n := range names {
		CmdsMap[n] = elem
	}
}

func isHelp(name string) bool {
	for _, alias := range HelpCmd.Name() {
		if alias == name {
			return true
		}
	}

	return false
}

//	Exec	执行指令命令行
//	args 参数可以直接接受os.Args
func Exec(args []string) (int, error) {

	if len(args) <= 2 {
		return MispCmd.Exec(args)
	}

	//	检查一下是否帮助指令
	if isHelp(args[2]) {
		return HelpCmd.Exec(args)
	}

	//	从指令表中查找指令
	elem, exist := CmdsMap[args[2]]
	if exist {
		return UnspCmd.Exec(args)
	}

	//	执行找到的指令对象
	return elem.Value.(Cmd).Exec(args)
}

type cmdHelp struct {
}

func (h cmdHelp) Name() []string {
	return [] string{"help", "-h", "--help"}
}

func (h cmdHelp) Help() []string {
	return [] string{
		"Show this help",

		"help|-h|--help	   Show abstracts for all commands.",
		"help <COMMAND>    Show help detial for <COMMAND>.",
		"help help         Show this help.",
	}
}

func (h *cmdHelp) Exec(args []string) (int, error) {
	if len(args) <= 2 {
		for elem := CmdList.Front(); nil != elem; elem = elem.Next() {
			c := elem.Value.(Cmd)
			if "" != c.Name()[0] {
				fmt.Fprintf(os.Stderr, "%s\t%s\n", c.Name()[0], c.Help()[0])
			}
		}

		return 0, nil
	}

	//	查找到底是针对谁的帮助
	elem, exist := CmdsMap[args[2]]

	//	处理帮助指令不存在的场景：直接定位为不支持的指令
	c := elem.Value.(Cmd)
	if !exist {
		c = UnspCmd
	}

	//	获取帮助信息并打印
	helps := c.Help()
	for i := 1; i < len(helps); i++ {
		fmt.Fprintf(os.Stderr, "%s\n", helps[i])
	}

	return 0, nil
}

type cmdMisp struct {
}

func (h cmdMisp) Name() []string {
	return [] string{""}
}

func (h cmdMisp) Help() []string {
	return [] string{""}
}

func (h *cmdMisp) Exec(args []string) (int, error) {
	return 1, fmt.Errorf("Missing parameters, type -h for help")
}

type cmdUnsp struct {
}

func (h cmdUnsp) Name() []string {
	return [] string{""}
}

func (h cmdUnsp) Help() []string {
	return [] string{""}
}

func (h *cmdUnsp) Exec(args []string) (int, error) {
	return 2, fmt.Errorf("Unsupported command or help target('%s'), type -h for help", args[2])
}
