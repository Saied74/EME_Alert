package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"gopkg.in/yaml.v2"

	// "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	secondAccuracy = 3 //seconds to half a minute time
	midMinute      = 30
)

type configType struct {
	WSJTLog   string    `yaml:"wsjtLog"`
	Map65Log  string    `yaml:"map65Log"`
	TestLog   string    `yaml:"testLog"`
	StartTime time.Time `yaml:"startTime"`
	EndTime   time.Time `yaml:"endTime"`
}

type pastLog struct {
	callSign  string
	date      string
	time      string
	grid      string
	frequency string
}

type pastLogs map[string]*pastLog

// var wsjtxCnt = 0

// var map65Cnt = 0

func main() {
	//check for test flag
	testOption := flag.Bool("t", false, "true runs in test mode, see README.md")
	flag.Parse()
	test := *testOption
	//get configuration data
	configData, err := getConfigData()
	if err != nil {
		log.Fatal(err)
	}
	//run the test go routine if testing is called for
	if test {
		go configData.testFeeder()
		time.Sleep(120 * time.Second)
	}
	configData.display()
}

//read config.yaml file and return configuration data.
func getConfigData() (*configType, error) {
	config := &configType{}

	goPath := os.Getenv("GOPATH")
	configPath := filepath.Join(goPath, "EME_Alert/config.yaml")

	configData, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, syscall.ENOENT) { //if no yaml file, keep the previous configuration numbers
			fmt.Println("No adjust.yaml file found", err)
		}
		return &configType{}, err
	}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		return &configType{}, err
	}
	return config, nil
}

//return the number of secconds between now and a half a second.
func halfMinute() int {
	t := time.Now()
	d := t.Second() - midMinute
	return d
}

//build a history of the wsjt-x logs between the minimum and maximum dates
//in the config.yanl file.
func (c *configType) buildHistory() (pastLogs, error) {
	history := pastLogs{}
	f, err := os.ReadFile(c.WSJTLog)
	if err != nil {
		return history, err
	}
	data := string(f)
	lines := strings.Split(data, "\n")

	if len(lines) == 0 {
		return history, fmt.Errorf("zero length data from %s\n", c.WSJTLog)
	}

	for _, line := range lines {
		h := pastLog{}
		items := strings.Split(line, ",")
		if len(items) < 8 {
			fmt.Printf("malformed line %v in %s\n", line, c.WSJTLog)
			continue
		}

		h.callSign = items[4]
		h.date = items[0]
		h.time = items[1]
		h.grid = items[5]
		h.frequency = items[6]

		tS := h.date + "T" + h.time + "Z"
		tt, err := time.Parse(time.RFC3339, tS)
		if err != nil {
			log.Fatal(err)
		}
		start := c.StartTime.Sub(tt) <= 0
		end := c.EndTime.Sub(tt) > 0
		if start && end {
			history[items[4]] = &h
		}
	}

	return history, nil
}

func (c *configType) makeDisplayData(h pastLogs) (pastLogs, error) {
	l := pastLogs{}
	mD, err := os.ReadFile(c.Map65Log)
	if err != nil {
		return l, err
	}
	mapData := string(mD)
	if len(mapData) == 0 {
		return l, fmt.Errorf("map65 log was empty")
	}
	mapLines := strings.Split(mapData, "\n")
	for i, lineVal := range mapLines {
		line := strings.Split(lineVal, ",")
		if len(line) < 3 {
			fmt.Printf("malformed map65 line %d, as %v\n", i, line)
			continue
		}

		_, ok := h[line[2]]
		if !ok {
			l[line[2]] = &pastLog{
				callSign: line[2],
				date:     line[0],
				time:     line[1] + ":00",
			}
		}
	}

	return l, nil
}

func (c *configType) display() {
	a := app.New()
	w := a.NewWindow("Stations Not Worked")
	dvra, err := fyne.LoadResourceFromPath("ui/static/img/dvra.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	w.SetIcon(dvra)
	row1Col1 := widget.NewLabel("Call Sign")
	row1Col2 := widget.NewLabel("Last Decoded")
	row1 := container.New(layout.NewGridLayout(2), row1Col1, row1Col2)
	grid := container.New(layout.NewVBoxLayout(), row1)
	w.SetContent(grid)
	w.Resize(fyne.NewSize(180, 350))

	go func() {
		delT := halfMinute()
		time.Sleep(time.Duration(60-delT) * time.Second)
		for {

			h, err := c.buildHistory()
			d, err := c.makeDisplayData(h)
			if err != nil {
				log.Fatal(err)
			}
			rows := []fyne.CanvasObject{row1}
			for _, val := range d {
				rowData := []fyne.CanvasObject{widget.NewLabel(val.callSign),
					widget.NewLabel(fmt.Sprintf("%v %v", val.date, val.time))}
				row := container.New(layout.NewGridLayout(2), rowData...)
				rows = append(rows, row)
			}
			grid := container.New(layout.NewVBoxLayout(), rows...)
			w.SetContent(grid)
			w.Show()
			time.Sleep(60 * time.Second)
		}
	}()
	w.Show()
	a.Run()
	fmt.Println("Exited")
}
