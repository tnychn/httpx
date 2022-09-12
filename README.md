# httpx

[![Go Reference](https://pkg.go.dev/badge/github.com/tnychn/httpx.svg)](https://pkg.go.dev/github.com/tnychn/httpx)
[![Tag](https://img.shields.io/github/v/tag/tnychn/httpx)](https://github.com/tnychn/httpx/tags)
[![License](https://img.shields.io/github/license/tnychn/httpx)](./LICENSE.txt)

_A simple and convenient `net/http` wrapper with an [Echo](https://github.com/labstack/echo)-like interface._

---

**httpx** wraps `http.Request` into `httpx.Request` and
defines a `httpx.Responder` that implements `http.ResponseWriter`
to provide handy functions to send variety of HTTP responses.

The library also provides an idiomatic way to handle errors inside
HTTP handlers by defining a `httpx.HandlerFunc` (which implements `http.Handler`)
that has a better signature than `http.HandlerFunc`.

Additionally, `httpx.HandlerFunc` is essentially a drop-in replacement for `http.HandlerFunc`.
Most of the methods of `httpx.Request` and `httpx.Responder` are extracted from
[Echo](https://github.com/labstack/echo)'s [`Context`](https://github.com/labstack/echo/blob/c0c00e6241/context.go).

> This is *not* a web framework, but rather an unopinionated library that
> extends `net/http`, as it does not include a HTTP router nor any mechanism
> that handles middlewares and the database layer.
>
> In fact, you can leverage your favourite router e.g. `gorilla/mux` or `chi` to
> provide routing and middleware ability, and **httpx** works _well_ with them.

In short, consider this library a _lite_ version of Echo, but compatible with `net/http`,
without any third party dependencies.

## Why?

I want to use my favourite router along with some convenient methods
from Echo's `Context` at the same time, but without other bloat and dependencies.

## Roadmap

- [x] Request Data Binding
- [ ] File Responder Methods
- [ ] Real IP Extractor
- [ ] Graceful Shutdown
- [ ] Better TLS Support

## Credits

Most of the code is adapted and modified from
[labstack/echo](https://github.com/labstack/echo)
@[v5_alpha](https://github.com/labstack/echo/tree/v5_alpha),
a high performance, minimalist web framework.
