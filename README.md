# lite
Simple Golang library with basic tools

### Content
- [Get started](#get-started)
- [Config](#config)
- [Logging](#logging)
- [Error](#error)
- [Router](#router)
- [Worker](#worker)
- [Broker](#broker)
- [Storage](#storage)
- - [SQL](#sql)
- - [Mongo](#mongo)
- [Tools](#tools)
- - [types](#types)
- - [system](#system)
- - [collections](#collections)
- - [fs](#fs)

### Get started
```go
package main

import (
	"github.com/boostgo/lite"
	"github.com/boostgo/lite/app/api"
	"github.com/boostgo/lite/config"
	"github.com/labstack/echo/v4"
)

type Config struct {
	IsDebug bool `env:"IS_DEBUG" envDefault:"false"`
}

func main() {
	// set debug mode
	lite.SetDebug(true)
	
	// parse custom config structure from env
	var cfg Config
	config.MustRead(&cfg) // throws panic if could not parse config
	
	setupRoutes()
	
	// run http server
	lite.Run("localhost:8080")
}

func setupRoutes() {
	lite.GET("/health", func(ctx echo.Context) error {
		return api.Ok(ctx, map[string]any{
			"Message": "OK",
		})
    })

	users := lite.Group("/api/v1/user")
	users.GET("/get", getUser)
}

type User struct {
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Age       int    `json:"age"`
}

func getUser(ctx echo.Context) error {
	return api.Ok(ctx, []User{
		{
			LastName:  "John",
			FirstName: "Doe",
			Age:       42,
		},
		{
			LastName:  "Jane",
			FirstName: "Doe",
			Age:       42,
		}
    })
}
```

### Config
Config package parses provided structure and fill env variables with tag "env" and "envDefault" <br>
Ref: [go-envconfig](https://github.com/sethvargo/go-envconfig) <br>
Import:
```go
import "github.com/boostgo/lite/config"
```
Example:
```go
package main

import (
	"fmt"
	"github.com/boostgo/lite/config"
	"os"
)

type Config struct {
	IsDebug    bool   `env:"IS_DEBUG" envDefault:"false"`
	ServiceURL string `env:"SERVICE_URL"`
}

func main() {
	_ = os.Setenv("SERVICE_URL", "http://localhost:8000/api/v1/service")

	// parse custom config structure from env
	var cfg Config
	config.MustRead(&cfg)       // throws panic if could not parse config
	fmt.Println(cfg.IsDebug)    // print value of IS_DEBUG env variable
	fmt.Println(cfg.ServiceURL) // print "http://localhost:8000/api/v1/service"
}
```

### Logging
Logger uses "zerolog". If use errs.Error logger will print all custom error fields<br>
Ref: [zerolog](https://github.com/rs/zerolog) <br>
Import:
```go
import "github.com/boostgo/lite/log"
```

Example:
```go
package main

import (
	"errors"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/log"
	"net/http"
)

func main() {
	log.Info("main").Msg("Hello world")
	// print:
	// {"level":"info","namespace":"main","time":"2024-06-26T15:01:24+05:00","message":"Hello world"}

	log.Error().Err(errs.New("Custom error").SetHttpCode(http.StatusNotFound).SetError(errors.New("inner error")))
	// print:
	// {"level":"error","innerError":"inner error","statusCode":404,"time":"2024-06-26T15:05:29+05:00","message":"Custom error"}
}
```

### Error
Custom error with fields:
- message
- inner error (builtin error)
- http code
- error type (string)
<br<br>

Also supports **Is**, **As**, **Unwrap** functions 

Example:
```go
package main

import (
	"errors"
	"fmt"
	"github.com/boostgo/lite/errs"
	"net/http"
)

var (
	ErrNotFound = errors.New("route not found")
)

func main() {
	fmt.Println(errs.
		New("Not found").
		SetHttpCode(http.StatusNotFound).
		SetError(ErrNotFound).
		SetType("HTTP"))
	
	// print:
	// [HTTP] Not found: route not found
	
	err := getNotFoundError()
	fmt.Println(errors.Is(err, ErrNotFound))
	// print:
	// true
}

func getNotFoundError() error {
    return errs.
		New("Not found").
		SetHttpCode(http.StatusNotFound).
		SetError(ErrNotFound).
		SetType("HTTP")
}
```

### Router


### Worker
Example #1:
```go
package main

import (
	"fmt"
	"github.com/boostgo/lite/app/worker"
	"github.com/boostgo/lite/system/life"
	"time"
)

func main() {
	worker.Run("Event", time.Second*1, eventWorker)
	// or another option run with fromStart flag
	// worker.Run("Event", time.Second*1, eventWorker, true)
	life.Wait() // wait "interrupt" or "kill" signal
	
	// when app off, print:
	// {"level":"info","namespace":"workers","worker":"Event","time":"2024-06-26T15:26:27+05:00","message":"Stop worker by context"}
}

func eventWorker() error {
	fmt.Println("worker action")
	return nil
}
```

Example #2:
```go
package main

import (
	"errors"
	"fmt"
	"github.com/boostgo/lite/app/worker"
	"github.com/boostgo/lite/system/life"
	"time"
)

func main() {
	worker.
		New("Event", time.Second*1, eventWorker).
		FromStart(). // run at start once
		ErrorHandler(workerErrorHandler). // set custom error handler
		Run() // run (async)

	life.Wait() // wait "interrupt" or "kill" signal

	// when app off, print:
	// {"level":"info","namespace":"workers","worker":"Event","time":"2024-06-26T15:26:27+05:00","message":"Stop worker by context"}
}

var (
	ErrWorkerStop = errors.New("worker stop")
	
	cnt = 1
)

func eventWorker() error {
	if cnt%3 == 0 {
		fmt.Println("cnt case:", cnt)
		return ErrWorkerStop
	}

	fmt.Println("action:", cnt)
	cnt++
	return nil
}

func workerErrorHandler(err error) bool {
	fmt.Println("handled error:", err, ". Is worker stop:", errors.Is(err, ErrWorkerStop))
	return !errors.Is(err, ErrWorkerStop)
}
```

### Storage
#### SQL
Example:
```go
package main

func main() {
	//
}
```
#### Mongo
Example:
```go
package main

func main() {
	//
}
```

### Broker
#### Kafka
Consumer group example:
```go
package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/broker/kafka"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/life"
	"github.com/boostgo/lite/types/to"
)

func main() {
	cfg := kafka.Config{
		Brokers: []string{"localhost:19092"},
		Topics:  []string{"some_topic"},
		GroupID: "test.some_topic_group_id",
	}

	consumer, err := kafka.NewConsumerGroup(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("create kafka consumer")
	}

	fmt.Println("consuming started")
	consumer.Consume("event", kafka.ConsumerGroupHandler(
		func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, message *sarama.ConsumerMessage) {
			fmt.Println(to.String(message.Value))
			session.MarkMessage(message, "")
		},
		nil, nil,
	))

	life.Wait()
}
```

## Tools
### types
Convert provided value to types <br>
**Examples:**
```go
package main

import (
	"fmt"
	"github.com/boostgo/lite/types/to"
)

type Person struct {
    LastName string `json:"last_name"`
	FirstName string `json:"first_name"`
}

func main() {
	// string
	fmt.Println(to.String("hello world")) // "hello world"
	fmt.Println(to.String(123)) // "123"
	fmt.Println(to.String(true)) // "true"
	fmt.Println(to.String(Person{
		LastName: "John",
		FirstName: "Smith",
    })) // {"last_name":"John","first_name":"Smith"}
	
	// int
	fmt.Println(to.Int("123")) // 123
	fmt.Println(to.Int("123asd")) // 0
	
	// float
	fmt.Println(to.Float32("123.23"))
	fmt.Println(to.Float64("123.23"))
	
	// bool
	fmt.Println(to.Bool("true")) // true
	
	var oneInt int = 1
	var oneInt8 = 1
	var oneFloat32 = 1
	fmt.Println(to.Bool(oneInt)) // true
	fmt.Println(to.Bool(oneInt8)) // true
	fmt.Println(to.Bool(oneFloat32)) // true
	
	// bytes
	fmt.Println(to.Bytes(map[string]any{
		"field1": "value1",
		"field2": "value2",
    }))
}
```

Simple param based on string which convert to many types <br>
Example:

```go
package main

import (
	"fmt"
	"github.com/boostgo/lite/types/param"
)

func main() {
	intParam := param.New("123asd")
	fmt.Println(intParam.MustInt(-1)) // -1
	
	uuidParam := param.New("c967fd9f-5e73-4b1c-8d3d-e4c527f7da12")
	fmt.Println(uuidParam.MustUUID()) // uuid object
}
```

### system
Example:

```go
package main

import (
	"context"
	"fmt"
	"github.com/boostgo/lite/system/life"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/system/try"
	"github.com/google/uuid"
)

func main() {
	// trace
	trace.IAmMaster()
	if trace.AmIMaster() {
		fmt.Println("I am trace master service")
    }
	ctx := context.Background()
	trace.Set(ctx, uuid.New().String())
	traceID := trace.Get(ctx)
	fmt.Println("trace id:", traceID)

	// try
	defer func() {
		err := try.CatchPanic(recover()) // return error object
		if err != nil {
			fmt.Println("Catch panic:", err)	
        }
    }()
	if err := try.Try(func() error {
		return connectToResource()
	}); err != nil {
		panic(err)
	}
	
	// life
	life.Tear(disconnectResource) // call after wait
	life.Context() // get app context
	life.Cancel() // cancel app context func
	life.Wait() // wait for signals and call teardown functions
}

func connectToResource() error {
	return nil
}

func disconnectResource() error {
    return nil
}
```

### collections
**Iterator:**
Example:

```go
package main

import (
	"fmt"
	"github.com/boostgo/lite/collections/iterator"
)

func main() {
	// simple example
	params := []string{"param1", "param2", "param3"}
	it := iterator.New(params)
	for it.Next() {
		fmt.Println(it.Get())
	}
	
	// skip
	it = iterator.New(params)
	it.Skip(2)
	fmt.Println(it.TryGet()) // "param3"
}
```

**list:**
Example:

```go
package main

import (
	"fmt"
	"github.com/boostgo/lite/collections/list"
	"strings"
)

func main() {
	params := list.Of([]string{"param1", "param2", "param3", "test"})
	fmt.Println(params) // ["param1", "param2", "param3", "test"]

	param1 := params.Single(func(s string) bool {
		return s == "param1"
	})
	if param1 == nil {
		fmt.Println("param1 not found")
	}

	fmt.Println("found:", *param1)

	onlyParams := params.Filter(func(s string) bool {
		return strings.Contains(s, "param")
	})
	fmt.Println(onlyParams.Slice()) // ["param1", "param2", "param3"]
}
```

### fs
Example:
```go
package main

func main() {
	//
}
```

