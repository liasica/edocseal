FROM --platform=linux/amd64 liasica/edocseal:latest

COPY ./build/release/edocseal /app/
COPY ./signer.py /app/

RUN chmod +x /app/signer.py

WORKDIR /app

ENTRYPOINT ["/app/edocseal", "server", "-c", "/app/config/config.yaml"]
