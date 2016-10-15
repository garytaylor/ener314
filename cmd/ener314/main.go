package main

import (
	"fmt"
	"log"
	"time"

	"github.com/barnybug/ener314"
)

func fatalIfErr(err error) {
	if err != nil {
		panic(fmt.Sprint("Error:", err))
	}
}

func main() {
	dev := ener314.NewDevice()
	err := dev.Start()
	fatalIfErr(err)

	for {
		// poll receive
		msg := dev.Receive()
		if msg == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		record := msg.Records[0] // only examine first record
		switch t := record.(type) {
		case ener314.Join:
			log.Printf("%06x Join\n", msg.SensorId)
			dev.Join(msg.SensorId)
		case ener314.Temperature:
			log.Printf("%06x Temperature: %.2f°C\n", msg.SensorId, t.Value)
			if msg.SensorId == 0x00098b {
				dev.TargetTemperature(msg.SensorId, 25)
			}
			// dev.Voltage(msg.SensorId)
			// dev.Diagnostics(msg.SensorId)
		case ener314.Voltage:
			log.Printf("%06x Voltage: %.2fV\n", msg.SensorId, t.Value)
		case ener314.Diagnostics:
			log.Printf("%06x Diagnostic report: %s\n", msg.SensorId, t)
		}
	}
}
