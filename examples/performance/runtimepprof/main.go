package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()

	// --- 1. 开启 CPU 分析 ---
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatalf("could not create cpu profile: %v", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalf("could not start cpu profile: %v", err)
		}
		defer pprof.StopCPUProfile()
	}

	// ... begin of the program

	// --- 2. 中间部分的业务逻辑 ---
	fmt.Println("程序开始运行...")
	// 模拟任务 A：产生了大量临时垃圾（运行完就应该被回收）
	generateTransientGarbage()
	// 模拟任务 B：产生了一些长期驻留的数据（比如最终结果，直到程序结束都在）
	data := generatePersistentData()
	fmt.Printf("业务逻辑结束，生成了 %d 条持久数据。\n", len(data))

	// ... end of the program

	// --- 3. 结束前记录内存快照 ---
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatalf("could not create mem profile: %v", err)
		}
		defer f.Close()
		// 关键点：手动 GC。
		// 如果不 GC，任务 A 的垃圾可能还在堆里，会干扰你判断“谁泄露了”。
		// 我们想看的是：程序跑完后，到底是谁还在占着内存不肯走？
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatalf("could not write mem profile: %v", err)
		}
	}
}

// 模拟：在这个函数里申请了 50MB 内存，但函数结束后这些内存就没用了
func generateTransientGarbage() {
	fmt.Println("执行任务 A (产生临时垃圾)...")
	for i := 0; i < 10; i++ {
		// 每次申请 5MB，这一瞬间内存会飙高
		_ = make([]byte, 5*1024*1024)
		// 假装做了一些计算
		time.Sleep(time.Millisecond * 50)
	}
}

// 模拟：这个函数返回的数据，会在 main 函数里一直存活
func generatePersistentData() [][]int {
	fmt.Println("执行任务 B (产生持久数据)...")
	var result [][]int
	for i := 0; i < 10000; i++ {
		// 这些数据被 append 到了 result，并被返回了，无法被回收
		row := make([]int, 100)
		for j := 0; j < 100; j++ {
			row[j] = rand.Int()
		}
		result = append(result, row)
	}
	return result
}
