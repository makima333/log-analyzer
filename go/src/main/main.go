package main

import (
	"fmt"
	"github.com/urfave/cli"
	"logparse"
	"os"
	"time"
	"sort"
	"log"
)

type Maptype struct {
	RemoteHost map[string]int
	Time       map[int]int
}

const dateformat = "2006-01-02"

func callparse(c *cli.Context, start string, end string, flag string) Maptype {
	lfactor := make(chan logparse.Line, 1000000)
	fin := make(chan int, c.NArg())
	if c.NArg() == 0 {
		log.Fatal("need log file as argument")
	}
	for i := 0; i < c.NArg(); i++ {
		go func(i int) {
			logparse.ReadLines(c.Args().Get(i), lfactor, fin)
		}(i)
	}
	var sum Maptype
	sum.Time = map[int]int{}
	sum.RemoteHost = map[string]int{}
	start_date, err := time.Parse(dateformat, start)
	if err != nil && start != "" {
		log.Fatal("error:start date format is wrong!", "example:", dateformat, "\n", err)
	}
	end_date, err := time.Parse(dateformat, end)
	if err != nil && end != "" {
		log.Fatal("error:end date format is wrong!", "example:", dateformat, "\n", err)
	}

	for i := range lfactor {
		switch {
		case start_date.Before(i.Time) && !end_date.After(i.Time) && start != "" && end != "":
		case !start_date.Before(i.Time) && end == "" && start != "":
		case !end_date.After(i.Time) && start == "" && end != "":
		case flag == "t":
			sum.Time[i.Time.Hour()]++
			//fmt.Println(start_date.After(i.Time))
		case flag == "r":
			sum.RemoteHost[i.RemoteHost]++
		}
		if len(lfactor) == 0 { //deadlock prevention
			time.Sleep(1 * time.Second)
			if len(lfactor) == 0 && len(fin) == c.NArg() {
				close(lfactor)
			}
		}
	}
	return sum

}

// sort a Map[string]int by its values++++++++++++++++++++++
func rankByWordCount(wordFrequencies map[string]int) PairList{
  pl := make(PairList, len(wordFrequencies))
  i := 0
  for k, v := range wordFrequencies {
    pl[i] = Pair{k, v}
    i++
  }
  sort.Sort(sort.Reverse(pl))
  return pl
}

type Pair struct {
  Key string
  Value int
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }
//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++


func main() {
	var start string
	var end string

	app := cli.NewApp()
	app.Name = "parse apache log"
	app.Version = "1.0"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:        "start, s",
			Destination: &start,
		},
		cli.StringFlag{
			Name:				 "end, e",
			Destination: &end,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "time",
			Aliases: []string{"t"},
			Usage:   "echo time_map",
			Action: func(c *cli.Context) error {
				time_map := callparse(c, start, end, "t")
				fmt.Println("===============================")
				fmt.Println("term :" ,start , "~", end)
				fmt.Println("===============================")
				fmt.Println("[time] : [the number of access]")
				for i := 0; i < 24; i++ {
					fmt.Printf("%02d     : %d\n", i, time_map.Time[i])
				}
				return nil
			},
		},
		{
			Name:    "host",
			Aliases: []string{"h"},
			Usage:   "echo host_rank",
			Action: func(c *cli.Context) error {
				host_map := callparse(c, start, end, "r")
				//fmt.Println(host_map.RemoteHost)
			  rank := rankByWordCount(host_map.RemoteHost)
				fmt.Println("===============================")
				fmt.Println("term :" ,start , "~", end)
				fmt.Println("===============================")
				fmt.Println("[IP]         : [the number of access]")
				for i := 0; i < len(rank); i++ {
					fmt.Printf("%-12s : %d\n" ,rank[i].Key ,rank[i].Value)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
