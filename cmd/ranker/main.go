package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/nilsbu/ranker/pkg/rank"
)

func parse(path string) (vals []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return vals, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) > 0 {
			vals = append(vals, text)
		}
	}

	if err = scanner.Err(); err != nil {
		return vals, err
	}

	return vals, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("need to provide file name")
	}

	keyfile := os.Args[1]
	ranksfile := keyfile + "-ranks.json"

	var mtx *rank.Matrix

	if len(os.Args) >= 3 && os.Args[2] == "load" {
		data, err := ioutil.ReadFile(ranksfile)
		if err != nil {
			log.Fatal(err)
			return
		}

		mtx = rank.Deserialize(data)
	} else {
		keys, err := parse(os.Args[1])
		if err != nil {
			log.Fatal(err)
			return
		}

		mtx = rank.InitMatrix(keys)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("free spaces: %v\n", mtx.CountFree())

		target := rank.Position{}

		if pos, ok := mtx.FindFree(); ok {
			target = pos
			fmt.Printf("'%v' vs '%v':\n", target[0], target[1])

			txt, _ := reader.ReadString('\n')
			switch txt[:len(txt)-1] {
			case "k":
				mtx.Set(target, rank.A)
			case "l":
				mtx.Set(target, rank.B)
			case "i":
				mtx.Set(target, rank.AA)
			case "o":
				mtx.Set(target, rank.BB)
			}

			if filled, ok := mtx.SetImplied(); ok {
				mtx = filled
			} else {
				fmt.Println("conflict")
				mtx.Set(target, rank.X)
			}
		} else {
			for i, key := range mtx.Rank() {
				fmt.Printf("%v: %v\n", i, key)
			}
			return
		}

		if cycle, ok := mtx.FindCycle(); ok {
			for i := range cycle {
				fmt.Printf("%v: '%v' vs '%v'\n", i, cycle[i], cycle[(i+1)%len(cycle)])
			}

			txt, _ := reader.ReadString('\n')
			idx, err := strconv.Atoi(txt[:len(txt)-1])
			if err != nil {
				fmt.Println(err)
				continue
			}

			target[0] = cycle[idx]
			target[1] = cycle[(idx+1)%len(cycle)]

			mtx.Set(target, rank.B)
		}

		ioutil.WriteFile(ranksfile, mtx.Serialize(), 0644)
	}

}
