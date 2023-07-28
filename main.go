package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func logUsage() {
	for {
		// Get current time
        currentTime := time.Now()

        // Create a folder for logs if it doesn't exist
        if _, err := os.Stat("logs"); os.IsNotExist(err) {
            os.Mkdir("logs", 0755)
        }

		// Create a new log file for the current date
        logFileName := fmt.Sprintf("logs/log_%s.txt", currentTime.Format("2006-01-02"))
		logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening log file:", err)
			os.Exit(1)
		}
		defer logFile.Close()

		// Get CPU usage
		percent, err := cpu.Percent(time.Second, false)
		if err != nil {
			fmt.Printf("Error getting CPU usage: %s\n", err)
		} 

		// Get memory usage
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			fmt.Printf("Error getting memory usage: %s\n", err)
		} 
		
        // Log the data
        logLine := fmt.Sprintf("%s | CPU Usage: %.2f%% | Memory Usage: %.2f%% \n", currentTime.Format("15:04:05"), percent[0], memInfo.UsedPercent)
        fmt.Print(logLine)
        logFile.WriteString(logLine)

		// Wait for 1 seconds before logging again
		time.Sleep(1 * time.Second)
	}
}

func main() {
	fmt.Println("Monitoring client started.")

	// Handle Ctrl+C signal to exit the program gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Monitoring client stopped.")
		os.Exit(0)
	}()

	// Start logging CPU and memory usage
	logUsage()
}
