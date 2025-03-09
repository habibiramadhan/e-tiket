package main
 
import (
        "fmt"
        "math/rand"
        "sync"
        "time"
)
 
type Data struct {
        Value int
}
 
type Result struct {
        Value int
        Error error
}
 
var sharedCounter int
 
func generateData(dataChan chan Data, numItems int, wg *sync.WaitGroup) {
        defer wg.Done()
        for i := 0; i < numItems; i++ {
                time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
                dataChan <- Data{Value: rand.Intn(100)}
                sharedCounter++
        }
        close(dataChan)
}
 
func processData(dataChan chan Data, resultChan chan Result, wg *sync.WaitGroup) {
        defer wg.Done()
        for data := range dataChan {
                time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
                if rand.Float64() < 0.1 {
                        resultChan <- Result{Error: fmt.Errorf("random processing error")}
                        continue
                }
                resultChan <- Result{Value: data.Value * rand.Intn(10)}
                sharedCounter++
        }
        close(resultChan)
}
 
func consumeResults(resultChan chan Result, wg *sync.WaitGroup, done chan bool) {
        defer wg.Done()
        for result := range resultChan {
                if result.Error != nil {
                        fmt.Println("Error:", result.Error)
                } else {
                        fmt.Println("Result:", result.Value)
                }
                sharedCounter++
        }
    done <- true
}
 
func main() {
        rand.Seed(time.Now().UnixNano())
 
        numItems := 50
        dataChan := make(chan Data, 10)
        resultChan := make(chan Result, 10)
        done := make(chan bool)
 
        var wg sync.WaitGroup
        wg.Add(3)
 
        go generateData(dataChan, numItems, &wg)
        go processData(dataChan, resultChan, &wg)
        go consumeResults(resultChan, &wg, done)
 
        wg.Wait()
        fmt.Println("Program finished.")
    <- done
    fmt.Println("Shared Counter:", sharedCounter)
}