package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	DefaultDownloadSize = 100 * 1024 * 1024 // 100 MB
	DefaultHost         = "0.0.0.0"
	DefaultPort         = "80"
)

var tmpl = template.Must(template.New("home").Parse(`
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>HTTP Speed Test</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #e9ecef;
            display: flex;
            justify-content: center;
            align-items: center;
            flex-direction: column;
            height: 100vh;
        }

        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 20px;
        }

        form {
            margin-top: 20px;
            max-width: 400px;
            width: 90%;
            background-color: #fff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }

        label {
            display: block;
            margin-bottom: 10px;
            font-weight: bold;
            color: #555;
        }

        input[type="number"],
        select {
            width: 100%;
            padding: 12px;
            margin-bottom: 15px;
            border: 1px solid #ccc;
            border-radius: 6px;
            box-sizing: border-box;
            transition: border-color 0.3s;
        }

        input[type="number"]:focus,
        select:focus {
            border-color: #007bff;
            outline: none;
        }

        input[type="submit"] {
            background-color: #007bff;
            color: white;
            padding: 12px 20px;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            width: 100%;
            transition: background-color 0.3s;
        }

        input[type="submit"]:hover {
            background-color: #0056b3;
        }

        .footer {
            margin-top: 20px;
            text-align: center;
        }

        .footer a {
            color: #007bff;
            text-decoration: none;
            font-weight: bold;
        }

        .footer a:hover {
            text-decoration: underline;
        }

        @media (max-width: 600px) {
            form {
                padding: 15px;
            }

            input[type="number"],
            select,
            input[type="submit"] {
                padding: 14px;
            }
        }
    </style>
</head>

<body>
    <div>
        <h1>HTTP Speed Test</h1>
        <form action="/download" method="get">
            <label for="size">Enter size in bytes:</label>
            <input type="number" id="size" name="size" value="104857600" min="1" required>
            <label for="chunk_size">Select chunk size:</label>
            <select id="chunk_size" name="chunk_size">
                <option value="1048576">1 MB</option>
                <option value="5242880">5 MB</option>
                <option value="10485760" selected>10 MB</option>
                <option value="20971520">20 MB</option>
            </select>
            <input type="submit" value="Start Download">
        </form>
    </div>
    <div class="footer">
        <a href="https://github.com/DiheChen/speedtest" target="_blank">View on GitHub</a>
    </div>
</body>

</html>
`))

func generateRandomData(size int, chunkSize int, w http.ResponseWriter) {
	buf := make([]byte, chunkSize)
	for size > 0 {
		n := chunkSize
		if size < chunkSize {
			n = size
		}
		if _, err := rand.Read(buf[:n]); err != nil {
			http.Error(w, "Error generating random data", http.StatusInternalServerError)
			log.Printf("Error generating random data: %v", err)
			return
		}
		if _, err := w.Write(buf[:n]); err != nil {
			log.Printf("Error writing data to response: %v", err)
			return
		}
		size -= n
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil || size <= 0 {
		size = DefaultDownloadSize
	}

	chunkSize, err := strconv.Atoi(r.URL.Query().Get("chunk_size"))
	if err != nil || chunkSize <= 0 {
		chunkSize = 10 * 1024 * 1024 // Default to 10 MB
	}

	w.Header().Set("Content-Length", strconv.Itoa(size))
	w.Header().Set("Content-Type", "application/octet-stream")

	generateRandomData(size, chunkSize, w)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func main() {
	host := getEnv("HOST", DefaultHost)
	port := getEnv("PORT", DefaultPort)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/download", downloadHandler)

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
