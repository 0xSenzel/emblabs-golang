## 📁 Project Structure

```
assignments/
├── cmd/
│   └── api/
│       └── main.go              # API server entry point
├── internal/
│   ├── handler/
│   │   ├── handler.go           # HTTP request handlers
│   │   └── handler_test.go      # Handler unit tests
│   └── service/
│       ├── model.go             # Transaction data model
│       └── payment.go           # Payment processing logic
├── go.mod                       # Go module dependencies
└── README.md                    # This file

../workerpool/
└── workerpool.go                # Concurrent worker pool implementation
```

---

## Getting Started

### Prerequisites
- Go 1.21 or higher

### Setup Dependencies

This project uses Go modules for dependency management. Although it currently only uses the standard library, you should still run these commands:

```bash
cd assignments

# Download dependencies (if any)
go mod download

# Clean up and sync dependencies
go mod tidy
```
---

## 1️⃣ Run API Server

### Start the server:
```bash
cd assignments
go run ./cmd/api
```

**Expected output:**
```
Server is running on :8080
```

### Test the endpoint:
```bash
# Using curl
curl -X POST http://localhost:8080/pay \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "amount": 99.99,
    "transaction_id": "txn_001",
    "status": "PENDING"
  }'
```

**Response (201 Created):**
```json
{
  "user_id": "user123",
  "amount": 99.99,
  "transaction_id": "txn_001",
  "status": "SUCCESS"
}
```

---

## 2️⃣ Run Handler Tests

### Run all tests with verbose output:
```bash
cd assignments
go test ./internal/handler/... -v
```

### Run specific test:
```bash
go test ./internal/handler -run TestPayHandler_NewTransaction -v
```

**Expected output:**
```
=== RUN   TestPayHandler_NewTransaction
--- PASS: TestPayHandler_NewTransaction (0.00s)
=== RUN   TestPayHandler_IdempotentTransaction
--- PASS: TestPayHandler_IdempotentTransaction (0.00s)
PASS
ok      github.com/0xsenzel/emblabs-golang/internal/handler
```

---

## 3️⃣ Run Worker Pool

### Execute the worker pool:
```bash
cd workerpool
go run workerpool/workerpool.go
```

**What it does:**
- Creates 5 concurrent workers
- Processes 100 jobs (calculates square of each number)
- Each job takes ~2 seconds
- Total time: ~40 seconds (vs 200 seconds sequential)

**Sample output:**
```
Worker 1 starting job 1
Worker 2 starting job 2
Worker 3 starting job 3
...
Worker 1 finished job 1 (Result: 1)
Worker 2 finished job 2 (Result: 4)
...
Result: 1
Result: 4
Result: 9
...
```

---

## 4️⃣ Run with Docker

### Prerequisites:
- Docker installed on your system

### Build the Docker image:
```bash
cd assignments
docker build -t payment-api .
```

### Run the container:
```bash
docker run -d -p 8080:8080 --name payment-service payment-api
```

**If port 8080 is already in use:**
```bash
docker run -d -p 8081:8080 --name payment-service payment-api
```
---
## 5️⃣ Setup CI/CD

### Prerequisites:
- Docker Hub account

### Step 1: Create Docker Hub Personal Access Token

1. Log in to [Docker Hub](https://hub.docker.com)
2. Navigate to **Account Settings** > **Security** > **Personal access tokens**
3. Click **Create new token**
4. Enter a token name (e.g., `github-actions`)
5. Ensure **Read & Write** permissions are selected (required for pushing images)
6. Click **Generate** and copy the token immediately (you won't see it again)

### Step 2: Create Docker Hub Repository

1. Go to your Docker Hub account and click **Create repository**
2. Name the repository matching the `IMAGE_NAME` in [../.github/workflows/ci-cd.yml](../.github/workflows/ci-cd.yml)
3. Set visibility to **Public** (or Private if preferred)
4. Click **Create**

### Step 3: Add GitHub Secrets

1. Go to your GitHub repository
2. Navigate to **Settings** > **Secrets and variables** > **Actions**
3. Click **New repository secret** and add:
   - **Name:** `DOCKERHUB_USERNAME` | **Value:** Your Docker Hub username
   - **Name:** `DOCKERHUB_TOKEN` | **Value:** The token you created in Step 1

### Verification

Once configured, commits to your repository will automatically:
- Build the Docker image
- Run tests
- Push the image to Docker Hub (if all checks pass)

You can monitor progress in the **Actions** tab of your GitHub repository.