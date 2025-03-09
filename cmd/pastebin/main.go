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
        // goroutine processor >=  proses pembacaan hasil nybabin error
        // Masalahnya: goroutine processor bisa selesai lebih cepat dari consumer,
		// tapi consumer masih perlu bac data dari channel ini.
		// kalo ditutup di sini, consumer bisa kebingungan ya karna channelnya
		// udah ditutup padahal masih ada data yg harus diproses.
        // close(resultChan)
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
 
		// kasih dua wg terpisah biar nggak deadlock
		var wg sync.WaitGroup
		// wg.Add(3) // salah nih bang! Jangan campur consumer dengan generator di satu wg
		wg.Add(2) // Yang benar: wg ini khusus buat generator dan processor aja
        
        var consumerWg sync.WaitGroup
        consumerWg.Add(1)
 
        go generateData(dataChan, numItems, &wg)
        go processData(dataChan, resultChan, &wg)
		// pake wg terpisah buat consumer biar nggak deadlock
		// go consumeResults(resultChan, &wg, done)
		go consumeResults(resultChan, &consumerWg, done)
 
		// deadlockbang klo gini! 
		// wg.Wait() nunggu consumer, consumer nunggu channel ditutup,
		// tapi channel ditutup sama processor yg juga ditunggu sama wg.Wait()
		// jadinya saling tunggu terus gak selesai-selesai kayak nunggu gebetan yang gak pernah notice bang
		// wg.Wait()
		// fmt.Println("Program finished.")
		// <- done
        
        // Urutan yang benar:
        wg.Wait() // Tunggu generator dan processor selesai
        close(resultChan) // Tutup resultChan setelah processor selesai
        <- done // Tunggu consumer mengirim sinyal selesai
        fmt.Println("Program finished.")
        fmt.Println("Shared Counter:", sharedCounter)
}