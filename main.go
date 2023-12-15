package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strconv"
    "sync"
    "time"
    "github.com/oliveagle/jsonpath"
)

// Global variables for calculating the average price
var (
    sum   float64
    count int
    mutex sync.Mutex
)

// APIConfig defines the structure for API configuration
type APIConfig struct {
    Exchange string `json:"exchange"`
    URL      string `json:"url"`
    JSONPath string `json:"jsonPath"`
}

// fetchAPI fetches data from the API and updates global price sum and count
func fetchAPI(cfg APIConfig, wg *sync.WaitGroup) {
    defer wg.Done()

    maxRetries := 3 // Maximum number of retries for the HTTP request
    var resp *http.Response
    var err error

    // Retry loop for the HTTP request
    for attempt := 1; attempt <= maxRetries; attempt++ {
        resp, err = http.Get(cfg.URL)
        if err == nil && resp.StatusCode != http.StatusTooManyRequests {
            break // Successful request
        }
        if attempt < maxRetries {
            time.Sleep(time.Duration(attempt) * time.Second) // Exponential backoff
        }
    }

    if err != nil {
        fmt.Printf("HTTP request failed after %d attempts: %s\n", maxRetries, err)
        return
    }
    defer resp.Body.Close()

    // Parse the JSON response
    value, err := parseJSONResponse(resp.Body, cfg.JSONPath)
    if err != nil {
        fmt.Println("Error parsing JSON:", err)
        return
    }

    fmt.Printf("%s: %s\n", cfg.Exchange, value)

    // Safely update the global sum and count
    mutex.Lock()
    price, err := strconv.ParseFloat(value, 64)
    if err == nil {
        sum += price
        count++
    }
    mutex.Unlock()
}

// parseJSONResponse parses the JSON response from the API
func parseJSONResponse(body io.ReadCloser, jsonPath string) (string, error) {
    var jsonData interface{}
    err := json.NewDecoder(body).Decode(&jsonData)
    if err != nil {
        return "", err
    }

    var jsonPathExpr = fmt.Sprintf("$.%s", jsonPath)
    res, err := jsonpath.JsonPathLookup(jsonData, jsonPathExpr)
    if err != nil {
        return "", err
    }

    value, ok := res.(string)
    if !ok {
        return "", fmt.Errorf("value at path %s is not a string", jsonPath)
    }
    return value, nil
}

func main() {
    var wg sync.WaitGroup

    // Read and parse the config.json file
    file, err := os.ReadFile("config.json")
    if err != nil {
        fmt.Println("Error reading config file:", err)
        return
    }

    var configs []APIConfig
    err = json.Unmarshal(file, &configs)
    if err != nil {
        fmt.Println("Error unmarshalling JSON:", err)
        return
    }

    // Fetch data concurrently for each API configuration
    for _, config := range configs {
        wg.Add(1)
        go fetchAPI(config, &wg)
    }

    wg.Wait()

    // Calculate and print the average price
    if count > 0 {
        avg := sum / float64(count)
        fmt.Println(int(avg))
    }
}
