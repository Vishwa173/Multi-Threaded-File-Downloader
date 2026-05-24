# Multi-Threaded Adaptive File Downloader

A modular high-performance file downloader built in Go, designed to optimize large file transfers using adaptive chunk scheduling, concurrent worker pools, resumable downloads, and integrity-aware transfer validation.

The system combines concurrent segmented downloading, centralized scheduling, fault-tolerant retries, and persistent checkpointing to improve throughput and reliability under unstable network conditions.

---

# Overview

The downloader is designed as a systems-focused networking project exploring:

- Concurrent network programming
- Adaptive scheduling systems
- Fault-tolerant distributed workflows
- Persistent state management
- Reliability engineering
- Binary integrity verification

The architecture implements:

- Concurrent segmented downloads using goroutines
- Adaptive throughput-aware chunk scheduling
- Centralized retry orchestration
- Persistent resumable downloads
- Source-aware scheduling foundations
- SHA256 integrity validation

---

# System Architecture

## Download Pipeline

```text
File Request
      ↓
Central Scheduler
      ↓
Dynamic Chunk Allocation
      ↓
Worker Pool Execution
      ↓
Concurrent HTTP Range Downloads
      ↓
Persistent Chunk Storage
      ↓
Chunk Merge
      ↓
SHA256 Integrity Verification
      ↓
Final Output
```

---

# Features

## Concurrent Segmented Downloading

- Parallel chunk-based downloading using goroutines
- Worker-pool architecture with centralized scheduling
- Dynamic runtime chunk allocation
- Efficient large-file transfer optimization

## Adaptive Scheduling

- Throughput-aware chunk sizing
- Dynamic worker load balancing
- Exponential moving average (EMA) throughput tracking
- Runtime-generated chunk scheduling

## Fault Tolerance

- Centralized retry orchestration
- Failed chunk requeueing
- Exponential backoff retry strategy
- Retry caps to prevent infinite retry loops

## Resumable Downloads

- Persistent checkpoint metadata
- Crash-safe recovery
- Persistent chunk storage
- Selective chunk reuse during resume

## Integrity Validation

- SHA256 file integrity verification
- Corruption detection
- Optional expected-hash validation
- Deterministic binary reconstruction

## Multi-Source Scheduling Foundations

- Source-aware scheduling architecture
- Source health tracking
- Throughput-aware source selection
- Dynamic source load balancing

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
│   ├── logger.go
│   ├── metrics.go
│   ├── scheduler.go
│   ├── storage.go
│   ├── types.go
│   └── worker.go
│
├── utils/
│   ├── download.go
│   ├── merge.go
│   └── temp.go
│
└── README.md
```

---

# Core Components

## scheduler.go

Implements:
- centralized chunk scheduling
- adaptive chunk sizing
- retry orchestration
- source-aware scheduling
- load balancing logic

## worker.go

Implements:
- worker pool execution
- concurrent chunk downloads
- throughput measurement
- retry-aware execution

## checkpoint.go

Implements:
- persistent checkpoint storage
- scheduler state persistence
- resumable download restoration

## integrity.go

Implements:
- SHA256 computation
- integrity validation
- corruption detection

## storage.go

Implements:
- chunk persistence
- chunk existence validation
- resumable chunk reuse

---

# Design Decisions

## Why Adaptive Chunk Scheduling?

Static chunk sizes can lead to inefficient resource utilization when workers operate under varying network conditions. Adaptive scheduling dynamically adjusts chunk sizes based on observed throughput to improve overall download efficiency.

## Why Centralized Scheduling?

A centralized scheduler simplifies:
- load balancing
- retry orchestration
- chunk ownership management
- adaptive scheduling decisions

while avoiding duplicate chunk assignments.

## Why Persistent Checkpointing?

Large downloads are vulnerable to:
- crashes
- interruptions
- unstable network conditions

Persistent checkpointing enables resumable downloads without restarting transfers from scratch.

## Why SHA256 Validation?

Segmented downloads may silently produce corrupted outputs due to:
- incomplete writes
- network instability
- corrupted responses

SHA256 verification guarantees deterministic binary reconstruction.

---

# Technologies Used

- Go (Golang)
- Goroutines
- Mutex Synchronization
- HTTP Range Requests
- JSON Checkpointing
- SHA256 Hashing

---

# Installation

## Clone Repository

```bash
git clone https://github.com/<your-username>/multi-threaded-downloader.git
cd multi-threaded-downloader
```

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

# Example Workflow

## Query

```text
Download large binary file
```

## Pipeline

```text
Chunk generation
→ concurrent segmented downloads
→ adaptive scheduling
→ retry orchestration
→ persistent chunk storage
→ resumable recovery
→ chunk merge
→ SHA256 verification
→ final output generation
```

---

# Future Improvements

Potential future extensions include:

- HTTP/2 multiplexing
- Real multi-mirror downloading
- Bandwidth throttling
- CLI flag support
- Per-chunk hashing
- Peer-to-peer transfer extensions
- Web dashboard and monitoring

---

# Learning Outcomes

This project explores core concepts in:

- Concurrent Systems Programming
- Networking and HTTP Internals
- Adaptive Scheduling Systems
- Fault Tolerance and Recovery
- Persistent State Management
- Reliability Engineering
- Binary Integrity Verification

---

# Author

Vishwa
