# hej 🐚

> A zero-dependency AI-powered CLI assistant that helps you discover and run shell commands from natural language.
> Built in **Go**, shipped as a single static binary.

---

## ✨ Features

* 🔹 No dependencies (single binary, no Python/Node required).
* 🔹 Cross-platform (Linux, macOS, Windows).
* 🔹 Works in **offline mode** (basic commands help).
* 🔹 Works in **online mode** with AI backends (via OpenRouter API).
* 🔹 Plug-and-play: `curl → chmod → run`.

---

## 📦 Installation

### 🛠 Local Development Setup

If you want to modify `hej` and test it locally:

```bash
# Clone the repo
git clone https://github.com/balajimalathi/hej.git
cd hej

# Build the binary (Linux/macOS)
go build -o hej .

# Run locally
./hej --help

# Example test
./hej "list all running docker containers"
```

For **Windows (PowerShell):**

```powershell
git clone https://github.com/balajimalathi/hej.git
cd hej

# Build
go build -o hej.exe .

# Run
.\hej.exe --help
```

---

### Server (Ubuntu/Debian)

```bash
# Download nightly build for Linux ARM64/AMD64
curl -sSL https://github.com/balajimalathi/hej/releases/latest/download/hej-linux-arm64 -o /usr/local/bin/hej

# Set permissions
sudo chmod +x /usr/local/bin/hej

# Verify
hej --help
```

---

## 🔑 Configuration

For **online mode** (AI suggestions), set your OpenRouter API key:

```bash
export OPENROUTER_API_KEY="your_api_key_here"
```

To persist across sessions, add it to `~/.bashrc` or `~/.zshrc`.

Without this key, `hej` will run in **offline mode** (basic predefined commands only).

---

## 🚀 Usage

```bash
hej "find all files larger than 500MB"
```

Example output:

```bash
Command: find . -type f -size +500M
Description: Lists all files larger than 500MB in the current directory.
```

---

## 🛠 Troubleshooting

### ❌ `Permission denied`

You forgot to make the binary executable:

```bash
chmod +x /usr/local/bin/hej
```

### ❌ `command not found`

The binary is not in your `PATH`.
Check:

```bash
which hej
```

If empty, move it into `/usr/local/bin/`:

```bash
sudo mv hej /usr/local/bin/hej
```

### ❌ `cannot execute binary file: Exec format error`

You downloaded the wrong binary for your architecture.
Check your system:

```bash
uname -m
```

* `x86_64` → use `hej-linux-amd64`
* `aarch64` → use `hej-linux-arm64`

### ❌ `Set OPENROUTER_API_KEY env var to enable online mode.`

You need to export your API key:

```bash
export OPENROUTER_API_KEY="your_api_key_here"
```

---

## 📅 Nightly Builds

Nightly builds are automatically published from the `main` branch with the format:

```
hej-<os>-<arch>-nightly-YYYYMMDD-HHMM
```

Example:

```
hej-linux-arm64-nightly-20250816-0730
```

---

## 📜 License

MIT License © 2025 [Balaji Malathi](https://github.com/balajimalathi)

---