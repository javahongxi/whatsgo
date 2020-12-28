package main

import (
	"bufio"
	"fmt"
	"github.com/javahongxi/golab/pipeline"
	"os"
	"strconv"
)

func main() {
	p := createNetworkPipeline("small.in", 800000000, 4)
	writeToFile(p, "small.out")
	printFile("small.out")
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(file, -1)
	count := 0
	for v := range p {
		fmt.Println(v)
		count++
		if count >= 100 {
			break
		}
	}
}

func writeToFile(p <-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriterSink(writer, p)
}

func createNetworkPipeline(filename string, filesize, chunkCount int) <-chan int {
	chunkSize := filesize / chunkCount
	pipeline.Init()
	var sortAddrs []string
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		file.Seek(int64(chunkSize*i), 0)

		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)

		addr := ":" + strconv.Itoa(7000+i)
		pipeline.NetworkSink(addr, pipeline.InMemSort(source))
		sortAddrs = append(sortAddrs, addr)
	}

	var sortResults []<-chan int
	for _, addr := range sortAddrs {
		sortResults = append(sortResults, pipeline.NetworkSource(addr))
	}
	return pipeline.MergeN(sortResults...)
}
