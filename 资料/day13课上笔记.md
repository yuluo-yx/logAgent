# day13课上笔记



# 内容回顾

### ginsession 的问题

1. sessionData中数据永不过期
   1. SessionData中保存过期时间（弊端是如果用户一直不登录，SessionData一直增长）
   2. 每次返回响应之前清除SessionData（弊端是会频繁连接Redis）
2. json反序列化之后整型和浮点型都成了`Float64`

### gob标准库

`gob`是Go语言中一个私有的编解码协议。

```go
func gobDemo() {
	var s1 = s{
		data: make(map[string]interface{}, 8),
	}
	s1.data["count"] = 1 // int
	// encode
	buf := new(bytes.Buffer) // 指针
	enc := gob.NewEncoder(buf) // 造一个编码器对象
	err := enc.Encode(s1.data)
	if err != nil {
		fmt.Println("gob encode failed, err:", err)
		return
	}
	b := buf.Bytes() // 拿到编码之后的字节数据
	fmt.Println(b)


	var s2 = s{
		data: make(map[string]interface{}, 8),
	}
	// decode
	dec := gob.NewDecoder(bytes.NewBuffer(b)) // 造一个新的解码器对象
	err = dec.Decode(&s2.data) // 把二进制数据解码到s2.data
	if err != nil {
		fmt.Println("gob decode failed, err", err)
		return
	}
	fmt.Println(s2.data)
	for _, v := range s2.data {
		fmt.Printf("value:%v, type:%T\n", v, v)
	}
}
```

# 今日内容

## make和new的区别

都是用来初始化内存

new多用来为基本数据类型（bool、string、int...）初始化内存，返回的是指针

make用来初始化`slice`、`map`、`channel`，返回的是对应类型。

## context

https://www.liwenzhou.com/posts/Go/go_context/

两个根函数

* `context.Background()`
* `context.TODO()`

四个with函数

* `context.WithCancel`:调用cancel看就会给`ctx.Done()`发信号
* `context.WithDeadline`：超时了也要调用`cancel()`
* `context.WithTimeOut`：超时了也要调用`cancel()`
* `context.WithValue`  ：注意key要使用自定义类型

## 性能调优（pprof）

[博客地址](https://www.liwenzhou.com/posts/Go/performance_optimisation/)

## 日志收集的项目（详见PDF）

1. 项目的架构
2. kafka
3. ZooKeeper

# 本周作业

1. 找一个内存优化的例子，使用pprof看一下效果