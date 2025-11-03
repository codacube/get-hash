# Use official Golang image as base
FROM docker.io/library/golang:1.25 AS dev

# Set working directory inside container
WORKDIR /app

# Install Air (live reload tool) for hot-reloading dev workflow
RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# For delve
EXPOSE 43000

# Jump into bash
# dlv debug --headless --listen=127.0.0.1:43000 --api-version=2 -- main.go -- "file_to_hash.ext"
CMD ["/bin/bash"]

# OR use this to fire up remote debugging using delve
# Arguments must be passed here when delve launches
# FROM https://github.com/charmbracelet/bubbletea
# dlv debug --headless --api-version=2 --listen=127.0.0.1:43000 .
# CMD ["dlv", "debug", "--headless", "--listen=:43000", "--api-version=2", ".", "--", "file_to_hash.ext"] 
# OR No file args
# CMD ["dlv", "debug", "--headless", "--listen=:43000", "--api-version=2", "."] 

# To build the image
# docker build -t get-hash .

# To run it - mount the $(pwd) to the image /app folder (volume mount)
# docker run -rm -it -v $(pwd):/app get-hash
# docker run -rm -it -v $(pwd):/app get-hash <file_to_hash.exe> - pass args directly to executable running in image. Haven't tested this.