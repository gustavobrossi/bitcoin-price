# bitcoin-price
## Overview
Bitcoin-Price is a high-performance solution built in Go, designed to fetch cryptocurrency values concurrently from various REST APIs. It's engineered for efficiency and speed, utilizing Go's robust concurrency model to make multiple endpoint calls simultaneously. This project stands out for its minimalistic approach, resulting in an extremely small Docker image based on scratch, making it lightweight and fast.

## Features
* Concurrent API Calls: Leverages Go's concurrency for efficient API querying.
* Flexible Configuration: Easily configurable for fetching values of various assets, not limited to Bitcoin.
* Minimal Docker Image: Built on scratch, ensuring a small footprint.
* Robust Error Handling: Handles API rate limits and networking issues with a retry mechanism.

## Usage
### Local Setup
```bash
go build -o main .
./main
```

### Docker Setup
Build the Docker image:
```bash
docker build -t bitcoin-price .
```

Run the Docker container (ensure config.json is in your current directory):
```bash
docker run -w $(pwd) -v $(pwd):$(pwd) bitcoin-price
```

### Configuration
config.json allows adding new exchanges or assets easily. Example:

```json
[
    {
        "exchange": "NewExchange",
        "url": "https://api.newexchange.com/price",
        "jsonPath": "price"
    }
    // ... other configurations ...
]
```

### Contribution
Contributions are welcome! Feel free to submit pull requests or open issues for enhancements, bug fixes, or suggestions.