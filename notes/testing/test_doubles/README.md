## 🧠 **Top-Down Overview: What Are Test Doubles?**

In software testing, especially in **unit testing**, we often replace real components (like databases, APIs, etc.) with **"test doubles"** to isolate the system under test (SUT).

Test doubles come in various forms:

| Term             | Purpose                                                     | Behavior                                                    |
| ---------------- | ----------------------------------------------------------- | ----------------------------------------------------------- |
| **Fake**         | Implements actual logic, just simplified or in-memory       | Used in testing without real dependencies                   |
| **Mock**         | Pre-programmed expectations and assertions                  | Checks that something _was called_ with _certain arguments_ |
| **Monkey Patch** | Replaces existing methods or objects dynamically at runtime | Used mainly in **dynamic languages** like Python            |

---

## 🛠️ 1. **Faking**

### 🔹 Concept:

A **fake** is a working implementation that mimics the real thing but is simpler (like using an in-memory DB instead of Postgres).

### ✅ Use Case:

When you want something that behaves like the real object but is fast and easy to control.

---

### ✅ Golang Example (Fake DB):

```go
type User struct {
    ID   int
    Name string
}

type UserRepository interface {
    GetUserByID(id int) (*User, error)
}

// Real implementation (e.g., querying SQL database)
type RealUserRepository struct{}

func (r *RealUserRepository) GetUserByID(id int) (*User, error) {
    // connect to database
    return &User{ID: id, Name: "Real User"}, nil
}

// Fake implementation (for testing)
type FakeUserRepository struct {
    Users map[int]*User
}

func (f *FakeUserRepository) GetUserByID(id int) (*User, error) {
    user, ok := f.Users[id]
    if !ok {
        return nil, fmt.Errorf("User not found")
    }
    return user, nil
}
```

---

### ✅ Python Example (Fake DB):

```python
class User:
    def __init__(self, id, name):
        self.id = id
        self.name = name

class UserRepository:
    def get_user_by_id(self, id):
        raise NotImplementedError()

class FakeUserRepository(UserRepository):
    def __init__(self):
        self.users = {1: User(1, "Alice")}

    def get_user_by_id(self, id):
        return self.users.get(id)
```

---

## 🧪 2. **Mocking**

### 🔹 Concept:

A **mock** doesn’t implement actual logic. It just records interactions and makes assertions about them.

### ✅ Use Case:

When you want to **verify behavior**, like "Was method `X` called with argument `Y`?"

---

### ✅ Golang Example (Manual Mock):

Golang doesn't have a built-in mocking library, so you manually record calls.

```go
type MockUserRepository struct {
    CalledWithID int
}

func (m *MockUserRepository) GetUserByID(id int) (*User, error) {
    m.CalledWithID = id
    return &User{ID: id, Name: "Mocked User"}, nil
}

// Test:
repo := &MockUserRepository{}
repo.GetUserByID(42)
if repo.CalledWithID != 42 {
    t.Errorf("Expected GetUserByID to be called with 42")
}
```

(Or you can use libraries like `github.com/stretchr/testify/mock`)

---

### ✅ Python Example (with `unittest.mock`):

```python
from unittest.mock import MagicMock

repo = MagicMock()
repo.get_user_by_id.return_value = {"id": 1, "name": "Mocked User"}

# Using mock
user = repo.get_user_by_id(1)
repo.get_user_by_id.assert_called_with(1)
```

---

## 🧩 3. **Monkey Patching**

### 🔹 Concept:

**Monkey patching** is replacing a method or function at runtime. It’s powerful and risky.

### ✅ Use Case:

Used in testing to override behavior of external libraries or methods.

---

### ⚠️ Not Available in Go (Statically Typed)

Go does **not** support monkey patching due to its strict type system and compile-time method resolution.

---

### ✅ Python Example (Monkey Patching):

```python
# real_module.py
def get_user():
    return "Real User"

# test_module.py
import real_module

def fake_get_user():
    return "Fake User"

# Monkey patching
real_module.get_user = fake_get_user

# Test
assert real_module.get_user() == "Fake User"
```

---

## 🔚 Summary Table

| Feature          | Go (Golang)                   | Python                           |
| ---------------- | ----------------------------- | -------------------------------- |
| **Fake**         | Struct implementing interface | Class implementing methods       |
| **Mock**         | Manual or via testify/mock    | `unittest.mock` or `pytest-mock` |
| **Monkey Patch** | ❌ Not supported              | ✅ Easy with dynamic assignment  |

## 🧾 References

- [ChatGPT](https://chatgpt.com/?model=auto&temporary-chat=true)
