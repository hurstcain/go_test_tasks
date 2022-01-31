package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Range struct {
	begin int
	end   int
}

func (self Range) String() string {
	return fmt.Sprintf("{begin: %d, end: %d}", self.begin, self.end)
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

// Функция поиска простых чисел.
// В качестве алгоритма взято Решето Эратосфена.
func FindPrimeNumbers(c chan<- int, wg *sync.WaitGroup, begin, end int) {
	defer wg.Done()

	nums := make([]bool, end+1)
	for i := 2; i <= end; i++ {
		if !nums[i] {
			if i >= begin {
				c <- i
			}
			for j := i; j <= end; j += i {
				nums[j] = true
			}
		}
	}
}

// Функция записи данных в файл.
func WriteToFile(ctx context.Context, c <-chan int, programEndingMsg chan<- string, wg *sync.WaitGroup, name string) {
	file, err := os.Create(name)
	CheckError(err)
	defer file.Close()

	// Канал, из которого можно считать данные только тогда, когда работа всех горутин, вычисляющих простые числа, завершена
	allWritersClosed := make(chan struct{})
	go func() {
		defer close(allWritersClosed)
		wg.Wait()
	}()

LOOP:
	for {
		select {
		case <-ctx.Done():
			programEndingMsg <- "The program ended by timeout"
			break LOOP
		case primeDigit := <-c:
			file.WriteString(strconv.Itoa(primeDigit) + " ")
		case <-allWritersClosed:
			programEndingMsg <- "The program ended because writing to the file is finished"
			break LOOP
		}
	}
}

func main() {
	timeout := flag.Int("timeout", 10, "Timeout")
	fileName := flag.String("file", "file.txt", "File's name")
	ranges := make([]Range, 0)

	flag.Func("range", "Range of prime numbers", func(s string) error {
		strArrRanges := strings.Split(s, ":")
		begin, err := strconv.Atoi(strArrRanges[0])
		CheckError(err)
		end, err := strconv.Atoi(strArrRanges[1])
		CheckError(err)

		ranges = append(ranges, Range{
			begin: begin,
			end:   end,
		})

		return nil
	})

	flag.Parse()

	// Канал, по которому передаются, а затем считываются простые числа.
	c := make(chan int)
	defer close(c)
	// Канал, в который записывается причина завершения программы
	programEndingMsg := make(chan string)
	defer close(programEndingMsg)
	wg := &sync.WaitGroup{}
	// Контекст, завершающийся по таймауту
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)

	for _, val := range ranges {
		wg.Add(1)
		go FindPrimeNumbers(c, wg, val.begin, val.end)
	}

	go WriteToFile(ctx, c, programEndingMsg, wg, *fileName)

	select {
	case msg := <-programEndingMsg:
		fmt.Println(msg)
	}
}
