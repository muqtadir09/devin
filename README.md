# Devin

[![Build-Check](https://github.com/muqtadir09/devin/actions/workflows/build.yaml/badge.svg?branch=main)](https://github.com/muqtadir09/devin/actions/workflows/build.yaml)
[![Registry](https://github.com/muqtadir09/devin/actions/workflows/registry.yaml/badge.svg?branch=main)](https://github.com/muqtadir09/devin/actions/workflows/registry.yaml)

## About

**Devin** is an autonomous AI software engineer designed to execute the complete software development lifecycle inside isolated execution runtimes. You submit tasks, and Devin analyzes the problem, provisions an isolated sandbox (or microVM), plans the solution, writes clean code, runs terminal commands, executes tests, and persists results to git.

### Key Highlights
- **Autonomous Reasoning & Planning**: Powered by the Brain harness (`@devin/brain`), Devin analyzes prompts, breaks down tasks, and selects appropriate tech stacks (`nextjs`, `node`, `go`, `python`, `rust`).
- **Isolated Execution Plane (Devbox)**: Code execution is sandboxed inside lightweight Firecracker microVMs or containerized runtimes with custom supervisors, ensuring security and zero pollution of the host machine.
- **Interactive Multi-Modal Interface**: Real-time streaming terminal, file tree explorer, interactive diff reviewer, and live desktop browser view via VNC.
- **Enterprise-Grade Security**: Strict prompt injection and jailbreak defenses (`<untrusted>` wrappers) treating all external user inputs, tool outputs, and repo files as unprivileged data.
- **Hybrid Deployment Models**: Can be run locally with zero-overhead dry-run mode, or scaled in production on AWS EKS with dedicated EC2 hardware-accelerated KVM execution hosts.

---

## System Architecture

Devin adopts a clean separation between the **Control Plane** (Web, API Server, Database, and LLM Brain) and the **Execution Plane** (Scheduler, Orchestrator, Tool Gateway, and MicroVM Sandboxes).

```
                      ┌───────────────────────────────────────────────┐
                      │              Web UI (Next.js 16)              │
                      │               http://localhost:3000           │
                      └──────────────────────┬────────────────────────┘
                                             │ HTTP / WebSockets
                                             ▼
                      ┌───────────────────────────────────────────────┐
                      │             API Server (@devin/server)        │
                      │               http://localhost:8080           │
                      │  • Better-Auth Authentication                 │
                      │  • WebSocket Event / VNC Tunneling            │
                      └──────────────┬─────────────────┬──────────────┘
                                     │                 │
                   Drizzle ORM       │                 │ HTTP Proxy
                                     ▼                 ▼
             ┌─────────────────────────────┐    ┌──────────────────────────────────┐
             │       PostgreSQL 16         │    │      Brain Engine (@devin/brain) │
             │   Users, Tasks, Sessions    │    │        http://localhost:9092     │
             └─────────────────────────────┘    │  • OpenAI Reasoning Harness      │
                                                │  • Intent & Stack Classification │
                                                │  • Tool Call Coordination        │
                                                └──────────────┬───────────────────┘
                                                               │
                                       POST /internal/v1/jobs  │
                                                               ▼
             ┌─────────────────────────────────────────────────────────────────────┐
             │                          Execution Plane                            │
             │                                                                     │
             │  ┌───────────────────────┐           ┌───────────────────────────┐  │
             │  │ Scheduler (:9091)     │──────────►│ Orchestrator (:9090)      │  │
             │  │ • Job Queues          │           │ • K8s Controller / Dry-run│  │
             │  │ • Devbox Standby/Wake │           └─────────────┬─────────────┘  │
             │  └───────────┬───────────┘                         │                │
             │              │                                     ▼                │
             │              │ Tool Execution        ┌───────────────────────────┐  │
             │              │ Requests              │ Firecracker MicroVM Host  │  │
             │              ▼                       │ • Golden Snapshot Restore │  │
             │  ┌───────────────────────┐           │ • Isolated /workspace     │  │
             │  │ Tool Gateway (:9095)  │──────────►│ • Supervisor API (:8081) │  │
             │  │ • Go gRPC Bridge      │  gRPC     └───────────────────────────┘  │
             │  └───────────────────────┘                                          │
             └─────────────────────────────────────────────────────────────────────┘
```

### Components & Responsibilities
- **`apps/web`** (Port `3000`): Modern Next.js 16 interface with real-time task progress, file diff viewer, terminal emulator, and desktop VNC stream.
- **`apps/server`** (Port `8080`): Central backend gateway managing user sessions, Better-Auth authentication, database schema migrations, and event streaming.
- **`apps/brain`** (Port `9092`): Autonomous LLM reasoning loop. Dispatches prompts to OpenAI (`gpt-4o-mini` / `gpt-4o`), determines if code execution is required, infers runtime requirements, and invokes coding tools.
- **`apps/scheduler`** (Port `9091`): Manages task scheduling, queues (in-memory locally or AWS SQS in cloud), and microVM session persistence.
- **`apps/orchestrator`** (Port `9090`): Go controller that provisions and reconciles sandbox instances either via Kubernetes CRDs or local dry-run HTTP.
- **`apps/tool-gateway`** (Port `9095`): High-performance Go gRPC daemon translating high-level tool calls into low-level runtime commands.
- **`apps/runtime`** (Port `8081`): In-guest Go supervisor acting as PID 1 inside the sandbox, exposing `/terminal`, `/files/*`, and `/git/*` APIs.
- **`apps/firecracker`**: Host microVM manager utilizing Linux KVM and CNI networking for sub-second VM snapshot restoration.

---

## Monorepo Layout

```
devin/
├── apps/
│   ├── web/               # Next.js 16 frontend application
│   ├── server/            # Express / Hono API server & Better-Auth
│   ├── brain/             # LLM reasoning harness & OpenAI loop
│   ├── scheduler/         # Task queue worker & session manager
│   ├── orchestrator/      # Go Kubernetes / dry-run sandbox controller
│   ├── tool-gateway/      # Go gRPC DevboxTools RPC gateway
│   ├── runtime/           # In-sandbox supervisor daemon (PID 1)
│   └── firecracker/       # Firecracker microVM daemon & CNI setup
├── packages/
│   ├── api/               # API route definitions (v1 and v2)
│   ├── drizzle/           # PostgreSQL schema, Drizzle ORM, & migrations
│   ├── types/             # Shared TypeScript types & runtime inference
│   ├── secrets/           # Envelope encryption for stored credentials
│   ├── ui/                # Shared UI design system components
│   └── events/            # Internal event bus & streaming primitives
├── docker/                # Dockerfiles & docker-compose configurations
├── infra/                 # Terraform modules for AWS EKS & execution hosts
└── tests/                 # Unit, integration, and E2E test suites
```

---

## Features

- **End-to-End Autonomous Coding**: Submits prompts, analyzes repos, modifies source trees, runs builds, and verifies output.
- **Dynamic Stack Selection**: Intelligently detects and selects golden snapshot images for Next.js, Node.js, Go, Python, and Rust.
- **Zero-Friction Local Development**: Full mock / dry-run modes enabling development on macOS/Windows without Linux KVM or AWS accounts.
- **Interactive VNC & Browser**: Real-time browser automation and virtual desktop capture via Xvfb and noVNC.
- **Hardened Security Boundaries**:
  - Prompt injection prevention via delimiter enforcement (`<untrusted>` blocks).
  - Secret sanitization stripping tokens from persisted jobs and logs.
  - MicroVM isolation preventing guest access to host resources.
- **1-Click Local Dev Login**: Instant local developer login bypassing external OAuth/email configuration.

---

## Prerequisites

- **Bun**: `^1.2.3` (or `v1.4+`)
- **Node.js**: `>= 18` (Node 20 or 22 recommended)
- **Go**: `>= 1.22` (for orchestrator, tool gateway, and runtime supervisor)
- **Docker**: For running the PostgreSQL database or containerized stack
- **OpenAI API Key**: For Brain LLM task planning (`OPENAI_API_KEY`)

---

## Quick Start (Local Setup)

### 1. Clone & Install Dependencies
```bash
git clone https://github.com/muqtadir09/devin.git
cd devin
bun install
```

### 2. Start PostgreSQL Database
```bash
docker compose -f docker/compose-dev.yaml up postgres -d
```

### 3. Initialize Environment Files
Copy sample environment configurations:
```bash
cp packages/drizzle/.env.sample packages/drizzle/.env
cp apps/server/.env.sample apps/server/.env
cp apps/brain/.env.sample apps/brain/.env
cp apps/orchestrator/.env.sample apps/orchestrator/.env
cp apps/firecracker/.env.sample apps/firecracker/.env
cp apps/runtime/.env.sample apps/runtime/.env
```

Generate a local session secret and enable local dev login:
```bash
sed -i '' "s/BETTER_AUTH_SECRET=/BETTER_AUTH_SECRET=$(openssl rand -hex 32)/" apps/server/.env
echo "ALLOW_LOCAL_DEV_LOGIN=true" >> apps/server/.env
```

Add your `OPENAI_API_KEY` in `apps/brain/.env`:
```env
OPENAI_API_KEY=sk-...
```

Enable the local developer login button in the web UI:
```bash
cp apps/web/components/auth/local-dev-login.local.example.tsx apps/web/components/auth/local-dev-login.local.tsx
```

### 4. Run Database Migrations
```bash
bun run migrate
```

### 5. Launch the Development Stack
```bash
bun run dev
```

Open your browser:
- **Web Interface**: `http://localhost:3000` (Click **"Local Dev Login"** to sign in instantly)
- **API Server**: `http://localhost:8080`
- **Brain Engine**: `http://localhost:9092`

---

## Testing & Quality Assurance

Run the test suite and quality verification:

```bash
# Run all unit and integration tests (159 passing)
bun run test

# Run TypeScript typechecks across all 26 packages
bun run check-types

# Run ESLint across web and UI packages
bun run lint

# Build all packages and binaries
bun run build
```

---

## Production Deployment (AWS / Kubernetes)

For cloud production environments:
1. **Control Plane on EKS**: Deploy `web`, `server`, and `brain` into an Amazon EKS cluster using GitOps manifests.
2. **Dedicated KVM Execution Hosts**: Provision EC2 compute-optimized or bare-metal instances (e.g. `c7g` / `c6i`) with `/dev/kvm` hardware acceleration enabled for real Firecracker microVMs.
3. **Database**: Use managed PostgreSQL (such as Neon, Supabase, or AWS RDS).
4. See [infra/README.md](file:///Users/abdulmuqtadir/dev/devin/infra/README.md) and [infra/deployment.md](file:///Users/abdulmuqtadir/dev/devin/infra/deployment.md) for full Terraform and GitOps runbooks.

---

## License

This project is licensed under the MIT License - see the [LICENSE](file:///Users/abdulmuqtadir/dev/devin/LICENSE) file for details.
