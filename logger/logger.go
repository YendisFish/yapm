package logger

import (
	"fmt"

	"github.com/fatih/color"
)

func Error(entry *LogEntry) {
	c := color.New(color.FgRed, color.Bold)

	c.Print("Error: ")
	fmt.Println(entry.Message)

	for k, v := range entry.Info {
		c.Print("    " + k + ": ")
		fmt.Println(v)
	}
}

func Info(entry *LogEntry) {
	c := color.New(color.FgGreen, color.Bold)

	c.Print("Info: ")
	fmt.Println(entry.Message)

	for k, v := range entry.Info {
		c.Print("    " + k + ": ")
		fmt.Println(v)
	}
}

func Warn(entry *LogEntry) {
	c := color.New(color.FgYellow, color.Bold)

	c.Print("Warning: ")
	fmt.Println(entry.Message)

	for k, v := range entry.Info {
		c.Print("    " + k + ": ")
		fmt.Println(v)
	}
}

func Log(entry *LogEntry, key string, col ...color.Attribute) {
	c := color.New(col...)

	c.Print(key + ": ")
	fmt.Println(entry.Message)

	for k, v := range entry.Info {
		c.Print("    " + k + ": ")
		fmt.Println(v)
	}
}

func LogString(message string, key string, col ...color.Attribute) {
	c := color.New(col...)

	c.Print(key + ": ")
	fmt.Println(message)
}

func LogRaw(message string, col ...color.Attribute) {
	c := color.New(col...)
	c.Print(message)
}

func LogRawln(message string, col ...color.Attribute) {
	c := color.New(col...)
	c.Println(message)
}

type LogEntry struct {
	Message string
	Info    map[string]string
}

func CreateLogEntry(message string, info ...[]string) *LogEntry {
	ret := &LogEntry{Message: message, Info: make(map[string]string)}

	for _, v := range info {
		if len(v) != 2 {
			Error(CreateLogEntry("Info array length should be of size 2!"))
			return nil
		}

		ret.Info["INFO | "+v[0]] = v[1]
	}

	return ret
}
