# 🦥 Vengo
*A lightweight Python virtual environment manager built in Go*

I made this because Python venvs suck and hopefully this makes them easier to manage.

---

## 🚀 Features

- ✅ **Create** Python virtual environments instantly.
- 📌 **List** existing environments along with their last modification dates.
- 🔑 **Activate** environments seamlessly from your shell.
- Automatic shell integration for effortless activation (`vengo activate myenv`).
- 🗑 **Delete** virtual environments easily.
---

## 💻 Installation

### Using Go (recommended):

```bash
go install github.com/glancing/vengo@latest
```
Make sure your Go bin directory (`$HOME/go/bin`) is in your `PATH`. If it's not, add this line to your shell config (`~/.bashrc` or `~/.zshrc`

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Then reload your shell config:

```bash
source ~/.bashrc  # or ~/.zshrc
```
---

## 📖 Usage

### 🔨 Create a new environment

```bash
vengo create myenv
```

### 📋 List your environments

```bash
vengo list
```

**Example output**:
```
📌 Your virtual environments:
 - myenv (Last Modified: 2025-03-10 14:22:11)
 - test (Last Modified: 2025-03-08 10:01:05)
```

### Activate an environment

Activate a virtual environment directly in your terminal:

```bash
vengo activate myenv
```

Your shell prompt will now indicate that `myenv` is active.

To **deactivate** the environment, run:

```bash
deactivate
```

### Delete an environment

Safely remove an environment when it's no longer needed:

```bash
vengo delete myenv
```

---

## Shell Integration

Vengo automatically sets up shell integration when you first run it. It adds a function to your `.bashrc` or `.zshrc`, enabling seamless activation of environments.

If activation isn't working immediately after installation, reload your shell configuration:

```bash
source ~/.bashrc
# or
source ~/.zshrc
```

---

## 🖥 Example Workflow

```bash
# Create and activate an environment
vengo create myenv
vengo activate myenv

# Install packages
pip install requests flask

# Deactivate the environment
deactivate

# Remove the environment when done
vengo delete myenv
```

---

## 🐛 Troubleshooting

- **`vengo command not found`?**  
  Check that `$HOME/go/bin` is in your PATH.
- **Shell function not working?**  
  Run `source ~/.bashrc` or `source ~/.zshrc`.

---

## 🤝 Contributing

Issues, feature requests, and pull requests are welcome.  
Please open an issue first to discuss changes.

---

## 📄 License

Vengo is licensed under the MIT License.
