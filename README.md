# Event Package for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/parsilver/event.svg)](https://pkg.go.dev/github.com/parsilver/event)
[![Go Report Card](https://goreportcard.com/badge/github.com/parsilver/event)](https://goreportcard.com/report/github.com/parsilver/event)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A lightweight, flexible event management system for Go applications. This package provides a robust event-driven architecture with support for listeners, subscribers, prioritization, and propagation control.

## Features

- **Event Dispatching**: Dispatch events to registered listeners
- **Priority System**: Control the order of listener execution
- **Propagation Control**: Optionally stop event propagation at any point
- **Subscriber Interface**: Register multiple listeners at once using a declarative API
- **Middleware Support**: Add cross-cutting concerns like logging, timing, etc.
- **Thread-Safe**: Concurrent access to the dispatcher is properly handled
- **Zero Dependencies**: No external dependencies required

## Table of Contents

- [Event Package for Go](#event-package-for-go)
  - [Features](#features)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
  - [Core Concepts](#core-concepts)
    - [Events](#events)
    - [Listeners](#listeners)
    - [Priority](#priority)
    - [Stopping Propagation](#stopping-propagation)
    - [Subscribers](#subscribers)
  - [Advanced Usage](#advanced-usage)
  - [License](#license)

## Installation

```bash
go get github.com/parsilver/event
```

## Quick Start

```go
package main

import (
    "fmt"
    
    "github.com/parsilver/event"
)

func main() {
    // Create a dispatcher
    dispatcher := event.NewDispatcher()
    
    // Register a listener using a function
    dispatcher.AddListener("user.created", event.ListenerFunc(func(e event.Event) bool {
        fmt.Println("User created event received!")
        return true
    }))
    
    // Dispatch an event
    e := event.NewEvent("user.created")
    dispatcher.Dispatch(e)
}
```

## Core Concepts

### Events

Events represent something that has happened in your application. They have a name and can carry additional data.

```go
// Create a simple event
e := event.NewEvent("user.created")

// Create an event with data
e := event.NewEvent("user.created", map[string]interface{}{
    "user_id": 123,
    "email": "john@example.com",
})
```

You can also create custom event types for better type safety:

```go
type UserCreatedEvent struct {
    *event.BaseEvent
    UserID   int
    Username string
    Email    string
}

func NewUserCreatedEvent(userID int, username, email string) *UserCreatedEvent {
    return &UserCreatedEvent{
        BaseEvent: event.NewEvent("user.created"),
        UserID:    userID,
        Username:  username,
        Email:     email,
    }
}
```

### Listeners

Listeners handle events when they are dispatched. They implement the `Listener` interface:

```go
type MyListener struct{}

func (l *MyListener) Handle(e event.Event) bool {
    // Handle the event
    return true // Return true for successful handling
}

// Register the listener
dispatcher.AddListener("user.created", &MyListener{})
```

You can also use the `ListenerFunc` for simpler cases:

```go
dispatcher.AddListener("user.created", event.ListenerFunc(func(e event.Event) bool {
    // Handle the event
    return true
}))
```

### Priority

You can control the order of listener execution by assigning priorities:

```go
// Higher priority listeners execute first
dispatcher.AddListener("user.created", &FirstListener{}, 100)
dispatcher.AddListener("user.created", &SecondListener{}, 50)
dispatcher.AddListener("user.created", &ThirdListener{}, 10)
```

### Stopping Propagation

Listeners can stop event propagation to prevent other listeners from being called:

```go
func (l *MyListener) Handle(e event.Event) bool {
    // Do something
    
    // Stop propagation to other listeners
    e.StopPropagation()
    
    return true
}
```

### Subscribers

Subscribers can listen to multiple events with a single registration:

```go
type MySubscriber struct{}

func (s *MySubscriber) OnUserCreated(e event.Event) bool {
    // Handle user.created event
    return true
}

func (s *MySubscriber) OnUserUpdated(e event.Event) bool {
    // Handle user.updated event
    return true
}

func (s *MySubscriber) GetSubscribedEvents() map[string][]event.SubscriberConfig {
    return map[string][]event.SubscriberConfig{
        "user.created": {
            {
                Method:   "OnUserCreated",
                Priority: 100,
            },
        },
        "user.updated": {
            {
                Method:   "OnUserUpdated",
                Priority: 90,
            },
        },
    }
}

// Register the subscriber
event.RegisterSubscriber(dispatcher, &MySubscriber{})
```

## Advanced Usage

See the `examples` directory for more advanced usage, including:

- Custom events
- Middleware implementation
- Priority handling
- Propagation control

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.