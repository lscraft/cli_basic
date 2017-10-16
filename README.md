# learing_go_homework

## selpg程序
基本功能是读取文件，并根据-d参数是否有被使用来输出想印的页，指定了就将信息pipe到命令 "lp -d**"中，没指定就送入标准输出流  
使用方法：
```
    ./cli [-options] filename
    -s=1     开始页数，默认为1  
    -e=1     终止页数，默认为1，输出包括终止页  
    -l=72    一页包含的行数，默认72，可以指定  
    -d=lp1   默认为空，如果指定了就送到指定的打印机里  
    -f       布尔值，默认为false
```
不指定文件名就会从标准输入流中读数据  
送到标准输出中的信息可以使用linux的管线命令来进一步处理 
实现思路是使用arginit函数来parse程序运行时的参数，参数数据游全局变量存储。
然后judge方法进行基本的判断，比如给出数据数字不能小于1，-f和-l不能同时使用。
之后根据有没有给出文件名来决定从标准输出还是从文件里读取。
将这个reader交给processInfo函数处理，这个函数中使用两个函数，SplitWithoutF和SplitWithF,分别处理没f和有f参数的情况
没有-f时，使用bufio.Scanner.Scan()来读取，有-f时一个字节一个字节地读取，并比对是否是"\f"，如果是就加一个换行到输出中。
最后如果有-d参数，就把读取的信息送入lp指令中


一些实例：
test sample: test.txt 
```
num:1
num:2
num:3
num:4
num:16
num:1
num:1
num:1
num:1
num:1
num:1
num:1
num:10end

```
1 input:./cli -s 1 -e 1 test.txt
```
XXXX@ubuntu:~/Desktop/go$ ./cli -s 1 -e 1 test.txt
num:1
num:2
num:3
num:4
num:16
num:1
num:1
num:1
num:1
num:1
num:1
num:1
num:10end

```
2 input:./cli -s 1 -e 1 test.txt
```
XXXX@ubuntu:~/Desktop/go$ ./cli -s 1 -e 1 < test.txt
num:1
num:2
num:3
num:4
num:16
num:1
num:1
num:1
num:1
num:1
num:1
num:1
num:10end

```
3. input: cat test.txt | ./cli -s 1 -e 1
```
XXXX@ubuntu:~/Desktop/go$ cat test.txt | ./cli -s 1 -e 1
num:1
num:2
num:3
num:4
num:16
num:1
num:1
num:1
num:1
num:1
num:1
num:1
num:10end

```
4 input: ./cli -s 1 -e 1 test.txt > ./hehe.txt
```
XXXX@ubuntu:~/Desktop/go$ ./cli -s 1 -e 1 test.txt > ./hehe.txt
```
同一目录下生成hehe.txt
打开一看，哇：
```
num:1
num:2
num:3
num:4
num:16
num:1
num:1
num:1
num:1
num:1
num:1
num:1
num:10end
```
因为各种管道，重定向都是有linux实现，程序所做的不过是把数据输入到标准输出中而已，所以略过。
5.-d参数实验
因为没有打印机，所以为了测试-d参数是否有效，将lp命令换成了cat命令，并且将标准输出和cmd的输出连接起来。
```
	//cmdArg := "-d" + d
	cmd := exec.Command("cat")
	cmdin, cmderr := cmd.StdinPipe()
	if cmderr != nil {
		cmdin.Close()
		os.Exit(1)
	}
 cmd.Stdout = os.Stdout
	cmd.Start()
	io.WriteString(cmdin, cmdinfo)
	cmdin.Close()
	cmd.Wait()
```
结果:
```
XXXX@ubuntu:~/Desktop/go$ ./cli -s 1 -e -d lp1 test.txt
num:1
num:2
num:3
num:4
num:16
num:1
num:1
num:1
num:1
num:1
num:1
num:1
num:10end

```
6. -f 参数
这次用比较特殊的测试文件，1~9，每行一个数字，数字后有一个"n"字符，暂时用n来分页符
```
XXXX@ubuntu:~/Desktop/go$ ./cli -s 3 -e 5 -f t2.txt

3
4
5

```
可以看到分页符没了，并把需要的信息分割出来，当然也可以保留不去掉分页符，稍微修改下代码就能做到。

7 -l 参数  该参数用于指定一页有多少行，你不能指定了-l参数还用-f参数，不然会报错（不过你填的是默认的72又用了-f,会直接按照-f执行）
这次我们用1行作为1页
input:./cli -s 3 -e 10 -l 1 test.txt
```
XXXX@ubuntu:~/Desktop/go$ ./cli -s 3 -e 10 -l 1 test.txt
num:3
num:4
num:16
num:1
num:1
num:1
num:1
num:1

```


# 总结
golang真的很奇怪的语言，语法上让人觉得很奇怪，特别是声明变量类型我经常转不过弯。
写这个程序的大部分时间都在查找golang相关的文档，以及实践操作，一个比较欣慰的是go的很多函数用法你看了文档后试着去用，然后是正确，但是C++中你看着文档然后试着去用，总是出现奇怪的错误LUL
