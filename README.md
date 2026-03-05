# gtr 

[English](README.md) | [中文](README-CN.md)

gtr is a go terminal tools for quick conversions and transformations in the CLI.

## Installation

```bash
go build -o gtr
```

## Usage

### 1. Coordinate Convert

Convert between different coordinate systems: WGS84, GCJ02 (火星坐标系), and BD09 (百度坐标系).

**Syntax:**
```bash
gtr coordinate [type] <longitude,latitude>
gtr coordinate [type] <longitude|latitude>
gtr coordinate [type] <longitude> <latitude>
```

**Types:** `wgs`, `wgs84`, `gcj`, `gcj02`, `gd`, `bd`, `bd09`, `gg`

**Examples:**
```bash
gtr coordinate 113.901495,22.499501
gtr coordinate gcj 113.901495,22.499501
gtr coordinate wgs 113.901495 22.499501
```

**Output:** Displays conversions to all three coordinate systems with 6 decimal precision.

---

### 2. HTTP Request

Send GET or POST HTTP requests and view detailed response information.

**Syntax:**
```bash
gtr http <URL>
gtr http <URL> '<JSON_DATA>'
```

**Examples:**
```bash
gtr http https://api.example.com/endpoint
gtr http https://api.example.com/post '{"key":"value","name":"test"}'
```

**Output:** Shows request URL, POST data, response headers, timing info, status code, and response body.

---

### 3. Time Convert

Convert between timestamps (10-digit and 13-digit) and human-readable date formats.

**Syntax:**
```bash
gtr time <timestamp|date>
```

**Supported Formats:**
- 10-digit Unix timestamp: `1727087511`
- 13-digit millisecond timestamp: `1727087511000`
- Standard format: `2024-09-23 10:31:51`
- Compact format: `20240923103151`

**Examples:**
```bash
gtr time 1727087511
gtr time 1727087511000
gtr time "2024-09-23 10:31:51"
gtr time 20240923103151
```

**Output:** Displays all format conversions including both 10-digit and 13-digit timestamps.

---

### 4. Text Encode/Decode

Encode and decode text using Base64, URL encoding, and generate MD5 hashes.

#### Base64 Encode/Decode
```bash
gtr text base64 encode "<text>"
gtr text base64 decode "<base64_string>"
```

**Examples:**
```bash
gtr text base64 encode "hello world"
gtr text base64 decode "aGVsbG8gd29ybGQ="
```

#### URL Encode/Decode
```bash
gtr text url encode "<text>"
gtr text url decode "<url_encoded_string>"
```

**Examples:**
```bash
gtr text url encode "hello world"
gtr text url decode "hello%20world"
```

#### MD5 Hash
```bash
gtr text md5 "<text>"
```

**Examples:**
```bash
gtr text md5 "password"
```

---

## Command Aliases

All commands support multiple aliases for convenience:

- `time`: `t`, `-t`, `--t`, `-time`, `--time`
- `coordinate`: `coor`, `-coor`, `--coor`, `c`, `-c`, `--c`
- `http`: `-http`, `--http`, `h`, `-h`, `--h`
- `text`: `t`, `-t`, `--t`, `-text`, `--text`

**Examples:**
```bash
gtr t 1727087511           # Same as: gtr time 1727087511
gtr coor 113.901495,22.499501  # Same as: gtr coordinate 113.901495,22.499501
gtr h https://example.com  # Same as: gtr http https://example.com
```

---

### 5. Todo Manager

Interactive task management with TUI (Text User Interface) and command-line modes.

**TUI Mode:**
```bash
gtr todo              # Launch interactive TUI with table view
```

**Command-line Mode:**
```bash
gtr todo add "Task title" --priority high --due 2025-03-10
gtr todo list [--status pending] [--priority high]
gtr todo delete <task-id>
```

**TUI Features:**
- Table view with: Title, Priority, Status, Due Date
- Navigation: (j/k) or (↑/↓) to move cursor
- (e) - Edit task: cycle priority with (p), cycle status with (c)
- (d) - Delete task
- Filters: (0)clear (1)pending (2)inprogress (3)done (h)igh (m)edium (n)low
- (?) - Help
- (q) - Quit

**Task States:**
- Pending → InProgress → Done (cycle with 'c' in TUI)

**Priorities:**
- High, Medium, Low (cycle with 'p' in TUI)

---

## Features

✅ Multiple coordinate system conversions  
✅ HTTP request testing with detailed response info  
✅ Flexible timestamp and date format conversion  
✅ Base64, URL, and MD5 text transformations  
✅ Interactive Todo manager with TUI  
✅ Command aliases for quick access  
✅ Clean, formatted output