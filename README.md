# todo.txt-googletasks

## Features

- Batteries included: no need to install any 3rd party dependency since the plugin is compiled as a static Go binary, with all dependencies inside it already.
- Caching of Google OAuth token: once authenticated you can run things in a head-less way in scripts, via `cron`, etc.
- Supported operations:
  - Delete all tasks in all lists.

## Development

### Setup

- Install [`docker`](https://store.docker.com/search?type=edition&offering=community)
- Install `make`

That's all folks!
All other tools are packaged in build Docker images (see `Dockerfile`) to ensure any machine can build easily, hence avoiding the "[_it works on my machine_](http://www.codinghorror.com/blog/2007/03/the-works-on-my-machine-certification-program.html)" syndrome.

### Build

```console
make
```
