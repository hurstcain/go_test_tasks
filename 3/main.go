package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
func FindPrimeNumbers(c chan<- int, begin, end int) {
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
func WriteToFile(c <-chan int, name string) {
	file, err := os.Create(name)
	CheckError(err)
	defer file.Close()

	for {
		file.WriteString(strconv.Itoa(<-c) + " ")
	}
}

func main() {
	timeout := flag.Int("timeout", 10, "Timeout")
	file_name := flag.String("file", "file.txt", "File's name")
	ranges := make([]Range, 0)

	flag.Func("range", "Range of prime numbers", func(s string) error {
		str_arr_ranges := strings.Split(s, ":")
		begin, err := strconv.Atoi(str_arr_ranges[0])
		CheckError(err)
		end, err := strconv.Atoi(str_arr_ranges[1])
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
	// Таймер продолжительностью timeout секунд
	timer := time.NewTimer(time.Duration(*timeout) * time.Second)

	go WriteToFile(c, *file_name)

	for _, val := range ranges {
		go FindPrimeNumbers(c, val.begin, val.end)
	}

	select {
	case <-timer.C:
		// Когда таймер завершается, он записывает в свой канал текущее время.
		// Таким образом, после получения данных из канала таймера программа завершает свою работу.
		fmt.Println("The program has completed")
	}
}
