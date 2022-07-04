<h1 align="center">httpx</h1>

<p align="center">
A simple wrapper around <code>net/http</code> which makes
HTTP handlers in Go more convenient to use.
</p>

---

**httpx** wraps `http.Request` into `httpx.Request` and
defines a `httpx.Responder` that implements `http.ResponseWriter`
to provide handy functions to send variety of HTTP responses.

The library also provides an idiomatic way to handle errors inside
HTTP handlers by defining a `httpx.HandlerFunc` (which implements `http.Handler`)
that has a better signature than `http.HandlerFunc`.

> This is *not* a web framework, but rather an unopinionated library that
> extends `net/http`, as it does not include a HTTP router nor any mechanism
> that handles middlewares and the database layer.
>
> In fact, you can leverage your favourite router e.g. `gorilla/mux` or `chi` to
> provide routing and middleware ability, and **httpx** works _well_ with them.

Additionally, `httpx.HandlerFunc` is essentially a drop-in replacement for `http.HandlerFunc`.
Most of the methods of `httpx.Request` and `httpx.Responder` are extracted from
[Echo](https://github.com/labstack/echo)'s [`Context`](https://github.com/labstack/echo/blob/c0c00e6241/context.go).

Consider this library as a _lite_ version of Echo.

## Why?

I want to use my favourite router along with some convenient methods
from Echo's `Context` at the same time, but without other bloat and dependencies.

## Roadmap

- [ ] ~~Request Data Binding~~
  - [x] exposed a `RequestBinder` interface instead
- [ ] Real IP Extractor
- [ ] Graceful Shutdown
- [ ] Better TLS Support

## Credits

Most of the code is adapted and modified from
[labstack/echo](https://github.com/labstack/echo),
a high performance, minimalist web framework.

---

<p align="center">
  <sub><strong>Made with ♥︎ by tnychn</strong></sub>
  <br>
  <sub><strong>MIT © 2022 Tony Chan</strong></sub>
</p>
