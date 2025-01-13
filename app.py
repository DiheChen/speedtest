import os

import uvicorn
from fastapi import FastAPI
from fastapi.responses import HTMLResponse, StreamingResponse

app = FastAPI()

CHUNK_SIZE = 1024 * 1024 * 10  # 10 MB
DEFAULT_DOWNLOAD_SIZE = 104857600  # 100 MB


@app.get("/download")
async def download(size: int = DEFAULT_DOWNLOAD_SIZE) -> StreamingResponse:
    async def generate():
        full_chunks = size // CHUNK_SIZE
        remainder = size % CHUNK_SIZE

        for _ in range(full_chunks):
            yield os.urandom(CHUNK_SIZE)
        if remainder:
            yield os.urandom(remainder)

    return StreamingResponse(generate(), media_type="application/octet-stream")


@app.get("/", response_class=HTMLResponse)
async def home() -> str:
    return """
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
                background-color: #f4f4f9;
                display: flex;
                justify-content: center;
                align-items: center;
                height: 100vh;
            }
            h1 {
                color: #333;
                text-align: center;
            }
            form {
                margin-top: 20px;
                max-width: 400px;
                width: 90%;
                background-color: #fff;
                padding: 20px;
                border-radius: 8px;
                box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            }
            label {
                display: block;
                margin-bottom: 8px;
                font-weight: bold;
            }
            input[type="number"] {
                width: 100%;
                padding: 10px;
                margin-bottom: 10px;
                border: 1px solid #ccc;
                border-radius: 4px;
                box-sizing: border-box;
            }
            input[type="submit"] {
                background-color: #007bff;
                color: white;
                padding: 10px 15px;
                border: none;
                border-radius: 4px;
                cursor: pointer;
                width: 100%;
            }
            input[type="submit"]:hover {
                background-color: #0056b3;
            }
            @media (max-width: 600px) {
                form {
                    padding: 15px;
                }
                input[type="number"],
                input[type="submit"] {
                    padding: 12px;
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
                <input type="submit" value="Start Download">
            </form>
        </div>
    </body>
    </html>
    """


if __name__ == "__main__":
    host = os.getenv("HOST", "0.0.0.0")
    port = int(os.getenv("PORT", 80))
    uvicorn.run(app, host=host, port=port)
