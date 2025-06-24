package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "sync"
)

func ProcessLogs(inputFiles []string, outputFile string) error {
    errorChan := make(chan string, 100)
    var wg sync.WaitGroup

    for _, file := range inputFiles {
        wg.Add(1)
        go func(filename string) {
            defer wg.Done()
            if err := readLogFile(filename, errorChan); err != nil {
                log.Printf("Failed to read file %s: %v\n", filename, err)
            }
        }(file)
    }

    go func() {
        wg.Wait()
        close(errorChan)
    }()

    return writeErrorsToFile(outputFile, errorChan)
}

func readLogFile(filename string, ch chan<- string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "ERROR") {
            ch <- line
        }
    }

    return scanner.Err()
}

func writeErrorsToFile(outputFile string, ch <-chan string) error {
    outFile, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer outFile.Close()

    writer := bufio.NewWriter(outFile)
    for line := range ch {
        _, err := writer.WriteString(line + "\n")
        if err != nil {
            return err
        }
    }

    err = writer.Flush()
    if err != nil {
        return err
    }

    fmt.Println("Error is successfully written in the output file")
    return nil
}


func main() {
    inputFiles := []string{"server1.log", "server2.log", "server3.log"}
    err := ProcessLogs(inputFiles, "errors.log")
    if err != nil {
        log.Fatal(err)
    } else {
        log.Println("Log processing completed. Check errors.log")
    }
}
