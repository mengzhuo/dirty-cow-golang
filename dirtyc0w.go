// Not work, why?

package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
)

const (
	TryTimes = 100000000
)

var (
	filePath = flag.String("f", "foo", "root file path")
	content  = flag.String("c", "mooooooooo", "write content")
	MAP      uintptr
)

func main() {
	flag.Parse()
	fmt.Println(">>>", *filePath, "with", *content)
	file, err := os.OpenFile(*filePath, os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	size := stat.Size()
	/*
		void *mmap(void *addr, size_t length, int prot, int flags,
				   int fd, off_t offset);
	*/

	MAP, _, _ = syscall.Syscall6(
		syscall.SYS_MMAP,
		uintptr(1),
		uintptr(stat.Size()),
		uintptr(syscall.PROT_READ),
		uintptr(syscall.MAP_PRIVATE),
		file.Fd(),
		0)
	go madvise(int(size))
	selfMem()
}

func madvise(size int) {
	var err error
	sl := struct {
		addr uintptr
		len  int
		cap  int
	}{MAP, size, size}
	/*
		for i := 0; i < TryTimes; i++ {
			err = syscall.Madvise(*(*[]byte)(unsafe.Pointer(&sl)), syscall.MADV_DONTNEED)
		}
	*/

	r1, r2, eo := syscall.Syscall(syscall.SYS_MADVISE, MAP, uintptr(100), syscall.MADV_DONTNEED)

	fmt.Println("madvise", r1, r2, eo)
}

func selfMem() {
	f, err := os.OpenFile("/proc/self/mem", syscall.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	con := []byte(*content)
	c := 0
	for i := 0; i < TryTimes; i++ {
		syscall.Syscall(syscall.SYS_LSEEK, f.Fd(), MAP, uintptr(os.SEEK_SET))
		n, _ := f.Write(con)
		c += n
	}
	fmt.Printf("Self Mem:%d", c)
}
