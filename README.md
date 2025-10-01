# System Monitor

A simple, real-time system monitoring dashboard built with Go, [raylib-go](https://github.com/gen2brain/raylib-go) for rendering, and [gopsutil](https://github.com/shirou/gopsutil) for metrics collection.

This application provides a clean, visual overview of key system resources, including CPU, memory, disk, and network activity.

**This project was created for learning purposes, as a way to explore Go programming and practice working with system metrics and graphics.**

![preview](https://github.com/user-attachments/assets/d884d213-50ec-41ba-922b-350bccde9f59)

## Features

  - **CPU Monitoring:**
      - Real-time usage percentage displayed with a visual bar graph.
      - Current CPU frequency.
  - **Memory Monitoring:**
      - Real-time memory usage percentage with a visual bar graph.
      - Displays used vs. total available RAM.
      - Displays used vs. total available Swap memory.
  - **Disk Monitoring:**
      - Primary disk usage percentage with a visual bar graph.
      - Displays used vs. total available disk space.
  - **Network Monitoring:**
      - Live upload and download speed (KB/s).
      - **Real-time line graphs** showing the last 60 seconds of network activity.
      - Displays current speed and peak speed within the historical window.
      - Total data sent and received since system startup.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

  - **Go**: Version 1.18 or later.
  - **A C Compiler**: `raylib` is a C library, so a compiler like `GCC` or `Clang` is required for the bindings to work.
      - On **Ubuntu/Debian**: `sudo apt-get install build-essential`
      - On **Fedora/RHEL**: `sudo dnf groupinstall "Development Tools"`
      - On **macOS**: Install Xcode Command Line Tools.
      - On **Windows**: Install MinGW-w64 or TDM-GCC.
  - **raylib System Dependencies**: Follow the installation instructions for your OS on the [raylib-go page](https://www.google.com/search?q=https://github.com/gen2brain/raylib-go%23installation).

## Installation & Usage

1.  **Clone the repository:**

    ```sh
    git clone [https://github.com/your-username/system-monitor.git](https://github.com/DiegoAndradeD/system-monitor)
    cd system-monitor
    ```

2.  **Tidy the dependencies:**
    Go will automatically fetch the required packages (`raylib-go` and `gopsutil`).

    ```sh
    go mod tidy
    ```

3.  **Run the application:**

    ```sh
    go run cmd/monitor/main.go  
    ```

    The application window should appear, displaying your system's metrics in real-time.

## Built With

  - [Go](https://golang.org/) - The programming language used.
  - [raylib-go](https://github.com/gen2brain/raylib-go) - For creating the graphical user interface.
  - [gopsutil](https://github.com/shirou/gopsutil) - For cross-platform system metrics collection.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.
