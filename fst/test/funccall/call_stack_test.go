package funccall

import (
	"runtime"
	"testing"
)

// +++++++++++++++++++++++++
// 设定中间件处理函数的数量
var handlerLen = 20
var Users []User

func init() {
	runtime.GOMAXPROCS(4)
	initData()
}

// A. 递归嵌套调用中间件模式
func BenchmarkCall(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	b.SetParallelism(20000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			callModelFunc()
		}
	})

	//for i := 0; i < b.N; i++ {
	//	callModelFunc()
	//}
}

// B. 循环调用中间件模式
func BenchmarkLoop(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	b.SetParallelism(20000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			loopModelFunc()
		}
	})

	//for i := 0; i < b.N; i++ {
	//	loopModelFunc()
	//}
}

// +++++++++++++++++++++++++
type User struct {
	name string
	age  int
}

// 准备测试数据
func initData() {
	for i := 0; i < handlerLen; i++ {
		Users = append(Users, User{name: "sdx", age: i + 1})
	}
}

// +++++++++++++++++++++++++
func callModelFunc() {
	index := 0
	next(index)
}

func next(index int) {
	if index < handlerLen {
		user := Users[index]
		index++
		handlerCall(user.name, user.age, index)
	}
}

func handlerCall(name string, age int, index int) int {
	arr := [100000]int{}
	ctLen := len(arr)
	for i := 0; i < ctLen; i++ {
		arr[i] = i * 10
	}

	//i := 0
	//i++
	// 栈空间申请
	//execCalc()
	next(index)
	//execCalc()
	//arr[0] = time.Now().Second() + i
	return arr[0]
}

// +++++++++++++++++++++++++
func loopModelFunc() {
	for _, user := range Users {
		handlerLoop(user.name, user.age)
	}
}

func handlerLoop(name string, age int) int {
	arr := [100000]int{}
	ctLen := len(arr)
	for i := 0; i < ctLen; i++ {
		arr[i] = i * 10
	}

	//i := 0
	//i++
	// 栈空间申请
	//execCalc()
	//execCalc()
	//arr[0] = time.Now().Second() + i
	return arr[0]
}

// +++++++++++++++++++++++++
// 重点改变这里的 数组大小，模拟业务代码的执行
func execCalc() {
	// 栈空间申请
	arr := [1]int{}
	ctLen := len(arr)
	for i := 0; i < ctLen; i++ {
		arr[i] = i * 10
	}
}

//
//func getRandomString(length int) string {
//	str := "0123456789abcdefghijklmnopqrstuvwxyz"
//	bytes := []byte(str)
//	result := []byte{}
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//	for i := 0; i < length; i++ {
//		result = append(result, bytes[r.Intn(len(bytes))])
//	}
//	return string(result)
//}
//
//func getRandomInt(max int) int {
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//	return r.Intn(max)
//}

//// ++++++++++++++++++++++++++++++++++
//// 打印 runtime 信息
//func printUsage() {
//	var m runtime.MemStats
//	runtime.ReadMemStats(&m)
//	log.Printf("CPU: %dm, MEMORY: Alloc=%.1fMi, TotalAlloc=%.1fMi, Sys=%.1fMi, NumGC=%d\n",
//		CpuUsage(), bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
//}
//
//var cpuUsage int64
//
//func CpuUsage() int64 {
//	return atomic.LoadInt64(&cpuUsage)
//}
//
//func bToMb(b uint64) float32 {
//	return float32(b) / 1024 / 1024
//}
