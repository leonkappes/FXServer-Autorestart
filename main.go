package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("---------------------")
	fmt.Println("FXServer-Restarter")
	fmt.Println("© Leon Kappes")
	fmt.Println("---------------------")

	restartProcess(true)

	times := make([]int, 0)

	f, err := os.Open("restart.csv")
	if err != nil {

		log.Fatal(err)
	}
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		for value := range record {

			i, err := strconv.Atoi(record[value])
			if err != nil {
				log.Fatal(err)
			}
			times = append(times, i)
		}
	}

	c := cron.New()
	c.AddFunc("1 * * * *", func() {
		t := time.Now()
		h := t.Hour()
		if contains(times, h) {
			time.Sleep(10 * time.Second)
			restartProcess(false)
		}
	})
	c.Start()

	fmt.Println("Press r to restart the server")
	for {
		fmt.Print("-> ")
		scanner.Scan()
		var text = scanner.Text()
		if text == "r" {
			restartProcess(false)
		}

	}
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func restartProcess(firstrun bool) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if runs, err := isProcRunning("FXServer.exe"); runs == true && err == nil && firstrun == false {
		fmt.Println("Stopping FXServer")
		cmd := exec.Command("taskkill.exe", "/F", "/IM", "FXServer.exe")
		err = cmd.Start()
		if err != nil {
			fmt.Println("Server was not running?!?")
		}
		err = os.RemoveAll(dir + "\\cache")
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(15 * time.Second)
		cmd = exec.Command("taskkill.exe", "/F", "/IM", "cmd.exe")
		err = cmd.Start()
		if err != nil {
			fmt.Println("Server was not running?!?")
		}
	}
	fmt.Println("Starting FXServer")
	cmd := exec.Command("cmd", "/C", "start", "cmd.exe", "@cmd", "/k", dir+"\\run.cmd")
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func isProcRunning(names ...string) (bool, error) {
	if len(names) == 0 {
		return false, nil
	}

	cmd := exec.Command("tasklist.exe", "/fo", "csv", "/nh")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	for _, name := range names {
		if bytes.Contains(out, []byte(name)) {
			return true, nil
		}
	}
	return false, nil
}
