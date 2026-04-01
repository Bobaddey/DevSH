# devsh (Developer Shell)

`devsh` is a production-grade CLI tool written in Go that allows developers to execute terminal-based operations using natural language. It safely translates plain English into terminal commands across Linux, Kubernetes, AWS, Terraform, Docker, and Git.

## Architecture Highlights
- **Context Engine**: Automatically sniffs environment markers (e.g. `.git`, `.terraform`, `KUBECONFIG`, `.minikube`) and the host **Operating System** (macOS, Linux) to prioritize contextual, platform-correct commands.
- **Rule Engine**: Fast fallback for common commands (like `docker ps` or `ls`) using static pattern matching for instant, deterministic results.
- **LLM Engine**: Pluggable architecture supporting OpenAI (GPT-4), Anthropic (Claude), Gemini, and local models via Ollama.
- **Safety Engine**: Validates commands before execution against blocklists, risk levels, and user-configured safety preferences.
- **Plugin System**: Modular generators for domain-specific tools.

## Installation

```bash
git clone https://github.com/devsh/devsh.git
cd devsh
go build -o devsh main.go
sudo mv devsh /usr/local/bin/
```

## Configuration

The easiest way to configure `devsh` is via the interactive setup wizard:
```bash
devsh config
```
This command allows you to seamlessly configure:
- **LLM Provider**: Choose between OpenAI, Anthropic, Gemini, or Ollama.
- **API Keys / Base URLs**: Effortlessly set your API secrets or configure OpenAI-compatible providers (Groq, vLLM, DeepSeek, etc.) by passing your custom Remote Base URL.
- **Model Selection**: Switch between gpt-4o, claude-3-5-sonnet-20240620, gemini-1.5-pro, or llama3.
- **Safety Levels**: Control devsh's aggressiveness (Low, Medium, High).

Alternatively, you can manually modify `~/.devsh/config.yaml` or set environment variables like `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, or `GEMINI_API_KEY`.

## Usage Examples

**Single Command Mode**
```bash
$ devsh create a folder called logs
🤖 Explained: Creates a new directory
💻 Command:  mkdir -p logs (Confidence: 1.00) (Risk: low)
🚀 Executing...

$ devsh list all pods in kube-system namespace
🤖 Explained: Direct Kubernetes command execution
💻 Command:  kubectl get pods -n kube-system (Confidence: 1.00) (Risk: medium)
⚠️ This command requires confirmation. Execute? (y/N): y
🚀 Executing...
```

**Context Awareness**
If you are currently inside a directory with `.terraform` files:
```bash
$ devsh "deploy the infrastructure"
# devsh will detect Terraform context and favor `terraform apply` over `kubectl apply`.
```

**Interactive REPL Mode**
```bash
$ devsh --interactive
🚀 Welcome to devsh interactive mode (type 'exit' or 'quit' to leave)

devsh> show pods
...

devsh> delete that pod
...
```

**Dry Run Mode**
```bash
$ devsh --dry-run "scale my payment service to 3 replicas"
🤖 Explained: Scales the payment deployment to 3 replicas in Kubernetes
💻 Command:  kubectl scale deployment payment --replicas=3 (Confidence: 0.95) (Risk: medium)
```

## Security
By default, `devsh` runs with `safety_level` set to high, requiring explicit y/N confirmation for almost any command it generates using the LLM. It hard blocks objectively dangerous commands like fork bombs or `rm -rf /`. You can change this behavior in your `~/.devsh/config.yaml`.
