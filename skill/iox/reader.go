package iox

import "io"

// Copy from io/io.go 638行的函数，用最有可能的[]byte长度申请内存空间，防止动态扩容
// 而标准库默认字节数组512字节，内容超过了会发生slice自动grow
// NOTE：传入的size，会在申请内存时自动+1
// ReadAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.
func ReadAll(r io.Reader, size int64) ([]byte, error) {
	if size <= 0 {
		size = 512 // 默认初始空间512字节
	}
	// 防止刚好读取了所有字符，但是没有得到EOF标记，还会循环读取一次，才会得到EOF标记
	// 此时 append(b, 0) 这个逻辑会造成 b 自动扩容；但是扩容却是无用功
	size += 1
	b := make([]byte, 0, size) // 内存空间尽量一次性分配到位
	for {
		if len(b) == cap(b) {
			b = append(b, 0)[:len(b)] // Add more capacity (let append pick how much).
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
	}
}
