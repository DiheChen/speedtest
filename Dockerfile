FROM python:3.13-alpine

WORKDIR /app

RUN pip install fastapi uvicorn

COPY app.py /app

ENTRYPOINT [ "python", "app.py" ]
