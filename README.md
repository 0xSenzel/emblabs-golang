# Assignments Answer

## Q1: Payment Service - Idempotency
- **File:** [assignments/internal/service/payment.go](./assignments/internal/service/payment.go)
- Follow the **"Run API Server"** section in [assignments/README.md](./assignments/README.md#1️⃣-run-api-server)

---
## Q2: Concurrency - Worker Pool
- **File:** [assignments/workerpool/workerpool.go](./assignments/workerpool/workerpool.go)
- Follow the **"Run Worker Pool"** section in [assignments/README.md](./assignments/README.md#3️⃣-run-worker-pool)

---
## Q3: Code Review - Bad Go Code
You are given the following code:
```go
var data = ""

func handler(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		data = string(body)
		fmt.Fprintf(w, "Saved: %s", data)
		defer r.Body.Close()
}
```

---

### ❌ Issues Found:

#### **1. Global Variable**
**Problem:** `var data = ""` is declared as a global variable, but it's only used locally within the handler function

**Issue:** multiple goroutines can access and modify the same global variable simultaneously, causing unpredictable behavior

**Fix:** Remove the global variable declaration and use a local variable inside the handler function instead

---

#### **2. Ignored Error Handling**
**Problem:** `body, _ := ioutil.ReadAll(r.Body)` ignores errors using blank identifier `_`

**Issue:** errors are suppressed and not handled, which can lead to potential panics or unexpected behavior

**Fix:**
```go
body, err := io.ReadAll(r.Body)
if err != nil {
    http.Error(w, "Failed to read body", http.StatusBadRequest)
    return
}
```

---

#### **3. Deprecated Package**
**Problem:** `ioutil.ReadAll` is deprecated since Go 1.16

**Issue:** using deprecated functions makes code outdated and harder to maintain

**Fix:**
```go
import "io"

body, err := io.ReadAll(r.Body) // Use io.ReadAll instead
```

---

#### **4. Defer Placement**
**Problem:** `defer r.Body.Close()` is placed after reading the body instead of immediately after function starts

**Issue:** If the code crashes before `defer` is reached, the connection stays open forever, wasting memory (similar to not calling `stream.close()` in a try-finally block)

**Fix:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close() // Close immediately after function starts
    
    body, err := io.ReadAll(r.Body)
    // ... rest of code
}
```

---

#### **5. Missing HTTP Status Code**
**Problem:** No explicit status code is set for successful response

**Issue:** clients may not know if the request was successful, defaults to 200 but not explicit

**Fix:**
```go
w.WriteHeader(http.StatusOK) // Explicit 200 status
fmt.Fprintf(w, "Saved: %s", data)
```

---
## Q4: SQL
### 4a) Join
You have two tables:

users:
| id | name |
| --- | --- |
| 101 | Alice |
| 102 | Bob |

orders:
| id | user_id | amount |
| --- | --- | --- |
| 1 | 101 | 50.00 |
| 2 | 101 | 75.00 |
| 3 | 102 | 30.00 |

Task: Write a query to list all users and their total order amount, including users with no orders.

**Answer:**
```sql
SELECT
    u.id,
    u.name,
    COALESCE(SUM(o.amount), 0) AS total_amount
FROM
    users u
LEFT JOIN
    orders o ON u.id = o.user_id
GROUP BY
    u.id, u.name;
```
**Expected Output:**
|id	  |name	|total_amount|
| --- | --- | --- |
|101  |Alice   |125.00|
|102  |Bob	   |30.00|
|103  |Charlie |0.00|

---
### 4b) Optimization / Indexing
Suppose you have a table with millions of rows:

```sql
CREATE TABLE transactions (
	id SERIAL PRIMARY KEY,
	user_id INT,
	amount DECIMAL,
	created_at TIMESTAMP
);
```

Task:

- Write a query to get the last 10 transactions of a given user (user_id = 123).
- Suggest one index that would improve query performance.

**Answer:**

#### 📝 **Query:**
```sql
SELECT 
    id,
    user_id,
    amount,
    created_at
FROM 
    transactions
WHERE 
    user_id = 123
ORDER BY 
    created_at DESC
LIMIT 10;
```

#### **Recommended Index:**
```sql
CREATE INDEX idx_user_created ON transactions(user_id, created_at DESC);
```
- Index reads **left to right**: `user_id` → `created_at`
- Query filters by `user_id` first (WHERE clause), then sorts by `created_at` (ORDER BY)
- Single index finds user_id quickly but still requires to sort the created_at

---
# Q5: Code Review Exercise
Given a poorly structured Go code and suggest way to refactor and improve the code. 

```go
package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
)

var result string

func handler(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadAll(r.Body) 
    result = string(body)             
    fmt.Fprintf(w, "Saved: %s", result)
    defer r.Body.Close()              
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

---

### ❌ Issues Found:

#### **1. Handler Function Issues**
**Note:** The `handler` function has multiple issues already covered in [**Q3**](#q3-code-review---bad-go-code)

Please refer to [**Q3**](#q3-code-review---bad-go-code) for detailed fixes.

---

#### **2. Missing Error Handling in main()**
**Problem:** `http.ListenAndServe(":8080", nil)` ignores the returned error

**Issue:** If the server fails to start, the program exits silently without any error message, making debugging difficult

**Fix:**
```go
func main() {
    http.HandleFunc("/", handler)
    
    fmt.Println("Server starting on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Printf("Server failed to start: %v\n", err)
    }
}
```

---

#### **3. No HTTP Method Validation**
**Problem:** Handler accepts all HTTP methods (GET, POST, PUT, DELETE, etc.) without validation

**Issue:** A GET request would try to read an empty body, and the handler would process it incorrectly. This can cause unexpected behavior

**Fix:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Only accept POST requests
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    defer r.Body.Close()
    // ... rest of code
}
```

---

#### **4. No Graceful Shutdown**
**Problem:** Server doesn't handle shutdown signals (CTRL+C), connections are abruptly terminated

**Issue:** potentially causing data loss or incomplete operations

**Fix:**
```go
import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // 1. Create server configuration
    server := &http.Server{
        Addr:    ":8080",
        Handler: http.DefaultServeMux,
    }
    
    http.HandleFunc("/", handler)
    
    // 2. Start server in background (non-blocking)
    go func() {
        fmt.Println("Server starting on :8080...")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            fmt.Printf("Server error: %v\n", err)
        }
    }()
    
    // 3. Create channel to receive OS signals (CTRL+C)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit  // Block here until signal received
    
    // 4. Begin graceful shutdown
    fmt.Println("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // 5. Wait for existing requests to complete (max 5 seconds)
    if err := server.Shutdown(ctx); err != nil {
        fmt.Printf("Server forced to shutdown: %v\n", err)
    }
    
    fmt.Println("Server stopped")
}
```

---

### ✅ **Final Refactored Code:**

```go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // Validate HTTP method
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Close body immediately
    defer r.Body.Close()
    
    // Read body with error handling
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    
    // Use local variable instead of global
    data := string(body)
    
    // Send response with explicit status code
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Saved: %s", data)
}

func main() {
    // Configure server
    server := &http.Server{
        Addr:         ":8080",
        Handler:      http.DefaultServeMux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    http.HandleFunc("/", handler)
    
    // Start server in goroutine
    go func() {
        fmt.Println("Server starting on :8080...")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            fmt.Printf("Server error: %v\n", err)
        }
    }()
    
    // Wait for interrupt signal for graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    fmt.Println("Shutting down server...")
    
    // Graceful shutdown with 5 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        fmt.Printf("Server forced to shutdown: %v\n", err)
    }
    
    fmt.Println("Server stopped gracefully")
}
```
---

# Bonus:
- unit test refers to [README.md](./assignments/README.md#2️⃣-run-handler-tests)
- Dockerfile refers to [README.md](./assignments/README.md#4️⃣-run-with-docker)
- CI/CD pipeline refers to [assignments/.github/workflows/ci-cd.yml](./assignments/.github/workflows/ci-cd.yml)
