
<a name="v1.4.0"></a>
## [v1.4.0](https://github.com/kilianc/pretender/compare/v1.3.0...v1.4.0)

> 2024-04-25

### Feat

* add healthcheck ([#50](https://github.com/kilianc/pretender/issues/50))


<a name="v1.3.0"></a>
## [v1.3.0](https://github.com/kilianc/pretender/compare/v1.2.0...v1.3.0)

> 2024-04-24

### Feat

* add support for repeatable responses ([#47](https://github.com/kilianc/pretender/issues/47))


<a name="v1.2.0"></a>
## [v1.2.0](https://github.com/kilianc/pretender/compare/v1.1.0...v1.2.0)

> 2024-04-24

### Feat

* add native support for json body ([#43](https://github.com/kilianc/pretender/issues/43))


<a name="v1.1.0"></a>
## [v1.1.0](https://github.com/kilianc/pretender/compare/v1.0.3...v1.1.0)

> 2024-04-14

### Feat

* add support for http status code, headers and response delay ([#11](https://github.com/kilianc/pretender/issues/11))
* add `--version` flag ([#14](https://github.com/kilianc/pretender/issues/14))


<a name="v1.0.3"></a>
## [v1.0.3](https://github.com/kilianc/pretender/compare/v1.0.2...v1.0.3)

> 2024-04-14

### Feat

* gracefully shut down on SIGINT and SIGTERM ([#12](https://github.com/kilianc/pretender/issues/12))

### Fix

* address bug in test for hh.LoadResponsesFile ([#10](https://github.com/kilianc/pretender/issues/10))
* change path for make run


<a name="v1.0.2"></a>
## [v1.0.2](https://github.com/kilianc/pretender/compare/v1.0.1...v1.0.2)

> 2024-04-12

### Fix

* change executable name when using go install


<a name="v1.0.1"></a>
## [v1.0.1](https://github.com/kilianc/pretender/compare/v1.0.0...v1.0.1)

> 2024-04-12

### Fix

* change package name to allow install


<a name="v1.0.0"></a>
## v1.0.0

> 2024-04-12

### Fix

* improve test suite approach by using osFileReader
* release lock on error
* make test catch path errors
* concurrency by embedding Mutex
* use an atomic counter instead of mutating the responses slice

