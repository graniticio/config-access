# Config Access

Config Access provides tools for loading, layering and accessing configuration files. It works with any configuration
format that can be parsed into a ```map[string]interface{}```, which includes Go's built in JSON parser and the most
[widely used YAML parsers](https://pkg.go.dev/gopkg.in/yaml.v3).

This package provides a small subset of the functionality of larger configuration libraries like [Viper](https://github.com/spf13/viper),
but has a fraction of the upstream dependencies. This package provides the configuration handling functionality for the
[Granitic microservices framework](https://granitic.io/).

GoDoc for this package can be [found here](https://pkg.go.dev/github.com/graniticio/config-access)

## Loading configuration

Config Access is agnostic of how your configuration is stored and parsed, as long as the result of parsing is a
```map[string]interface{}```. For brevity, Config Access defines a type alias ```ConfigNode``` for 
```map[string]interface{}```.

For example, a JSON file could be loaded and parsed as:

```go
  var config config_access.ConfigNode

  f, _ := os.Open("/your/config.json")
  b, _ := io.ReadAll(f)

  json.Unmarshal(bytes, &config)
```

## Accessing configuration

Once loaded, Config Access provides a way of recovering individual configuration values while converting them to 
various Go builtin types through an interface called a ```Selector```. Individual items are accessed using
a dot delimited 'path' syntax.

```go
  selector := config_access.NewDefaultSelector(config, true, true)
	
  if cs.PathExists("my.string") {
    s := cs.StringVal("my.string")		
  }
```

Methods exist to try and interpret configuration values as ```string```, ```int```, ```float64```, ```bool```, slices
```[]interface{}``` and objects ```map[string]interface{}```.

## 'Quiet' access

If you do not want to handle errors whenever you attempt to access a configuration value, you can use a `QuietSelector`
that instead of returning an error when encountering a problem or missing value, executes a function that you provide.

This can be used to implement 'fail fast' behaviour by exiting the application on an error or logging the problem. If
the function you provide does not exit the application, configuration value methods (e.g. `StringVal()`) will return 
the zero value relevant to the type you asked for (`nil`, `""`, 0 etc")

```go
  quiet := config_access.NewDeferredErrorQuietSelector(config,
	  func(path string, err error){
		  fmt.printf("Error while trying to load config at %s: %s", path, err.Error())
      })
```


## Layering and merging

Configuration sources can be merged together. The most common use case is to combine some base common configuration with
configuration that is specific to a deployment environment.

```go
  var base config_access.ConfigNode
  var prod config_access.ConfigNode

  f, _ := os.Open("/your/base-config.json")
  b, _ := io.ReadAll(f)

  json.Unmarshal(bytes, &base)

  f, _ = os.Open("/your/prod-config.json")
  b, _ = io.ReadAll(f)

  json.Unmarshal(bytes, &prod)
  
  combined := config_access.Merge(base, prod, false)
```

## Injecting configuration

Configuration loaded into a ```ConfigNode``` can be used to populate the fields of a struct in one call:

```go
  //Assume config previously loaded into variable config
  type Name struct {
    First string
    Middle []string
    Last string
  }
  
  var name Name

  ca.Populate("orderDetails.recipient", &name, config)
```