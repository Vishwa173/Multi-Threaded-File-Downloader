# Multi-Threaded Distributed Download Accelerator

A high-performance distributed file downloader built in Go, designed to optimize large file transfers through adaptive chunk scheduling, concurrent worker pools, fault-tolerant retries, resumable downloads, and integrity-aware transfer validation.

The system leverages goroutines, mutex synchronization, HTTP range requests, and centralized orchestration to maximize throughput while maintaining reliability under unstable network conditions.

---

# Features

## Concurrent Multi-Threaded Downloading
- Parallel segmented downloads using goroutines
- Worker-pool architecture with centralized task scheduling
- Adaptive chunk allocation for efficient resource utilization

---

## Adaptive Scheduling Engine
- Runtime-generated chunk scheduling
- Throughput-aware chunk sizing
- Dynamic worker load balancing
- Exponential moving average (EMA) throughput tracking

---

## Fault-Tolerant Retry System
- Centralized retry orchestration
- Failed chunk requeueing
- Exponential backoff strategy
- Retry caps to prevent infinite loops

---

## Multi-Source Download Architecture
- Source-aware scheduling foundations
- Source health tracking
- Throughput-aware source selection
- Dynamic load balancing across sources

---

## Resumable Downloads
- Persistent checkpoint metadata
- Crash-safe recovery
- Chunk persistence across restarts
- Selective chunk reuse during resume

---

## Integrity Validation
- SHA256 file integrity verification
- Corruption detection
- Optional expected-hash validation
- Deterministic binary reconstruction

---

# Architecture Overview

```text
                    +----------------------+
                    |   Central Scheduler  |
                    +----------------------+
                     |    |    |    |
         +-----------+    |    |    +-----------+
         |                |    |                |
         v                v    v                v

    +---------+      +---------+          +---------+
    | Worker  |      | Worker  |   ...    | Worker  |
    +---------+      +---------+          +---------+
         |                |                    |
         +-------- Concurrent Downloads -------+
                          |
                          v
                +------------------+
                | Chunk Storage    |
                | (.part files)    |
                +------------------+
                          |
                          v
                +------------------+
                | Final Merge      |
                +------------------+
                          |
                          v
                +------------------+
                | SHA256 Verify    |
                +------------------+
```

---

# Tech Stack

- **Language:** Go (Golang)
- **Concurrency:** Goroutines, Channels, Mutexes
- **Networking:** HTTP Range Requests
- **Synchronization:** sync.Mutex, sync.RWMutex
- **Persistence:** JSON Checkpointing
- **Integrity:** SHA256 Hashing

---

# Project Structure

```text
multi-threaded-downloader/
│
├── cmd/
│   └── main.go
│
├── downloader/
│   ├── checkpoint.go
│   ├── downloader.go
│   ├── integrity.go
│   ├── metrics.go
│   ├── scheduler.go
│   ├── storage.go
│   ├── types.go
│   ├── worker.go
│   └── logger.go
│
├── utils/
│   ├── download.go
│   ├── merge.go
│   └── temp.go
│
└── README.md
```

---

# How It Works

## 1. File Segmentation

The scheduler dynamically generates byte-range chunks using HTTP range requests.

Example:

```text
Chunk 0 → bytes 0-8MB
Chunk 1 → bytes 8MB-16MB
Chunk 2 → bytes 16MB-24MB
```

---

## 2. Worker Pool Execution

Persistent workers continuously pull chunks from the centralized scheduler.

- Fast workers naturally process more chunks
- Slow workers receive smaller workloads
- Dynamic balancing improves throughput

---

## 3. Fault Recovery

Failed chunks are:
- requeued
- retried with exponential backoff
- redistributed to available workers

---

## 4. Resume Support

During downloads:
- chunk files persist on disk
- checkpoint metadata is periodically saved

On restart:
- completed chunks are reused
- unfinished chunks resume automatically

---

## 5. Integrity Verification

After merge:
- SHA256 hash is computed
- optional expected-hash verification performed

This guarantees deterministic binary reconstruction.

---

# Installation

## Clone Repository

```bash
git clone https://github.com/<your-username>/multi-threaded-downloader.git
cd multi-threaded-downloader
```

---

## Build

```bash
go build ./cmd/main.go
```

---

# Usage

## Basic Download

```bash
go run ./cmd/main.go <url> <output-file> <workers>
```

Example:

```bash
go run ./cmd/main.go https://example.com/file.zip file.zip 8
```

---

## Download With SHA256 Verification

```bash
go run ./cmd/main.go <url> <output-file> <workers> <sha256>
```

Example:

```bash
go run ./cmd/main.go \
https://example.com/file.zip \
file.zip \
8 \
d2d2d2d2d2d2d2d2d2d2d2d2d2d2d2d2
```

---

# Resume Support

If the downloader is interrupted:

```text
file.zip.parts/
file.zip.meta.json
```

remain on disk.

Re-running the same command resumes the download automatically.

---

# Example Features Demonstrated

- Adaptive scheduling
- Distributed task orchestration
- Throughput-aware chunk allocation
- Fault-tolerant retry systems
- Persistent state management
- Binary integrity validation
- Concurrent network programming

---

# Future Improvements

- HTTP/2 multiplexing
- Real multi-mirror downloading
- Bandwidth throttling
- CLI flags
- Per-chunk hashing
- Peer-to-peer transfer extensions
- Web dashboard / monitoring

---

# Benchmarks

| Workers | Avg Speed |
|---|---|
| 1 | Baseline |
| 4 | ~3-4× improvement |
| 8 | Higher throughput under stable networks |

*(Results depend on network bandwidth and server range-request support.)*

---

# Key Concepts Demonstrated

- Concurrent systems programming
- Networking and HTTP internals
- Distributed scheduling concepts
- Fault tolerance and recovery
- Adaptive load balancing
- Persistent checkpointing
- Systems reliability engineering

---

---

# Author

Vishwa
