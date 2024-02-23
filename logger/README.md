# Logger
Default logger for your service.

# Examples
## Default
```go
import (
  "github.com/Moranilt/http-utils/logger"
)

func main() {
  log := logger.New(os.Stdout, logger.TYPE_JSON)
  log.Info("Hello World")
  // Output: {"level":"INFO","message":"Hello World","time":"2020-07-20T17:22:54+03:00"}
}
```

## Global
```go
import (
  "github.com/Moranilt/http-utils/logger"
)

func main(){
  log := logger.New(os.Stdout, logger.TYPE_JSON)
  logger.SetDefault(log)

  logger.Default().Info("Hello World")
  // Output: {"level":"INFO","message":"Hello World","time":"2020-07-20T17:22:54+03:00"}
}
```