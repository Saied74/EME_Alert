package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func (c *configType) testFeeder() {
	d := halfMinute()
	lC := 0
	time.Sleep(time.Duration(90-d) * time.Second)
	for {
		tD, err := os.ReadFile(c.TestLog)
		if err != nil {
			log.Fatal(err)
		}
		testData := string(tD)
		if len(testData) == 0 {
			log.Fatalf("empty test file\n`")
		}
		lines := strings.Split(testData, "\n")
		if len(lines) == 1 {
			log.Fatalf("test data did not split with newline symbol \n")
		}
		// lines = lines[lC:]
		for i := 0; i < 2; i++ {
			line := lines[2*lC+i]
			items := strings.Split(line, ",")
			if len(items) < 8 {
				log.Fatalf("malformed test line %v", line)
			}
			// if lC%3 == 0 {
			newLine := ""
			tS := items[0] + "T" + items[1] + "Z"
			t, err := time.Parse(time.RFC3339, tS)
			if err != nil {
				log.Fatal(err)
			}
			month := fmt.Sprintf("%v", t.Month())[0:3]
			newLine += items[0][0:5]
			newLine += month + "-"
			newLine += items[0][8:10] + ","
			newLine += items[1][0:5] + "," + items[4] + "\n"
			f, err := os.OpenFile(c.Map65Log, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.WriteString(newLine)
			if err != nil {
				f.Close()
				log.Fatal(err)
			}
			f.Close()
		}

		oldLine := lines[lC] + "\n"
		lC++
		// if lC%3 == 0 {
		g, err := os.OpenFile(c.WSJTLog, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
		}

		_, err = g.WriteString(oldLine)
		if err != nil {
			g.Close()
			log.Fatal(err)
		}
		g.Close()
		// }

		time.Sleep(60 * time.Second)
	}
}
