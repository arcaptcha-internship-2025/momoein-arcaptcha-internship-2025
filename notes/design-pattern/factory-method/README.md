## Factory Method

The Factory Method is a creational design pattern that provides an interface for creating objects in a superclass but allows subclasses to alter the type of objects that will be created.

It helps promote loose coupling and adheres to the Open/Closed Principle by allowing new product types without changing existing code.

## ✅ Concept:

- Instead of calling a constructor directly, you call a factory method.
- Subclasses can override this method to change the class of objects that will be created.

## 🐍 Python Example:

```python
from abc import ABC, abstractmethod

# Product
class Button(ABC):
    @abstractmethod
    def render(self):
        pass

# Concrete Products
class WindowsButton(Button):
    def render(self):
        print("Rendering a Windows button.")

class MacButton(Button):
    def render(self):
        print("Rendering a Mac button.")

# Creator
class Dialog(ABC):
    @abstractmethod
    def create_button(self) -> Button:
        pass

    def render_window(self):
        button = self.create_button()
        button.render()

# Concrete Creators
class WindowsDialog(Dialog):
    def create_button(self) -> Button:
        return WindowsButton()

class MacDialog(Dialog):
    def create_button(self) -> Button:
        return MacButton()

# Client Code
def client_code(dialog: Dialog):
    dialog.render_window()

# Usage
client_code(WindowsDialog())  # Rendering a Windows button.
client_code(MacDialog())      # Rendering a Mac button.
```

## 🦫 Go Example:

```go
package main

import "fmt"

// Product
type Button interface {
	Render()
}

// Concrete Products
type WindowsButton struct{}
func (w WindowsButton) Render() {
	fmt.Println("Rendering a Windows button.")
}

type MacButton struct{}
func (m MacButton) Render() {
	fmt.Println("Rendering a Mac button.")
}

// Creator
type Dialog interface {
	CreateButton() Button
}

// Concrete Creators
type WindowsDialog struct{}
func (WindowsDialog) CreateButton() Button {
	return WindowsButton{}
}

type MacDialog struct{}
func (MacDialog) CreateButton() Button {
	return MacButton{}
}

// Client Code
func renderWindow(d Dialog) {
	button := d.CreateButton()
	button.Render()
}

func main() {
	renderWindow(WindowsDialog{}) // Rendering a Windows button.
	renderWindow(MacDialog{})     // Rendering a Mac button.
}
```

## 📌 Summary:

Use it when: You need to delegate the instantiation logic to subclasses.

Benefit: Avoid tight coupling between the creator and concrete products.

Drawback: More classes and complexity due to subclassing.

## 🗒️ References

- [refactoring: factory method](https://refactoring.guru/design-patterns/factory-method)
