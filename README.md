[![CircleCI](https://circleci.com/gh/marccarre/todo.txt-googletasks/tree/master.svg?style=shield)](https://circleci.com/gh/marccarre/todo.txt-googletasks/tree/master)
[![Docker Repository on Quay](https://quay.io/repository/marccarre/todo.txt-googletasks/status)](https://quay.io/repository/marccarre/todo.txt-googletasks)
[![Go Report Card](https://goreportcard.com/badge/github.com/marccarre/todo.txt-googletasks)](https://goreportcard.com/report/github.com/marccarre/todo.txt-googletasks)

# todo.txt-googletasks

## Features

- Batteries included: no need to install any 3rd party dependency since the plugin is compiled as a static Go binary, with all dependencies inside it already.
- Caching of Google OAuth token: once authenticated you can run things in a head-less way in scripts, via `cron`, etc.
- Supported operations:
  - Delete all tasks in all lists.

## Installation

### Enable Google Tasks API

- Go to [Google APIs' console](http://code.google.com/apis/console).
- Click on "_Create project_" and give it whichever name you like, e.g.: `todotxt-googletasks`
- Click on "[_Dashboard_](https://console.developers.google.com/apis/dashboard?supportedpurview=project)", then [_Enable APIs and services_](https://console.developers.google.com/apis/library?supportedpurview=project). Filter APIs by typing "_task_" in the search box, and click on "[_Tasks API_](https://console.developers.google.com/apis/library/tasks.googleapis.com)". Click "_Enable_".
- After a few seconds, you should see the message "_To use this API, you may need credentials. Click 'Create credentials' to get started._". Click on [_Create credentials_](https://console.developers.google.com/apis/credentials/wizard).
- Under "_Which API are you using?_", select "_Tasks API_".
  Under "_Where will you be calling the API from?_", select "_Other UI (e.g. Windows, CLI tool).
  Under "_What data will you be accessing?_", select "_User data_".
  Click "_What credentials do I need?_". You should arrive on a page saying "_Create an OAuth 2.0 client ID_".
- Enter `todotxt-googletasks` under "_Name_", click "_Create OAuth client ID_".
- Select your email address, enter `todotxt-googletasks` under "_Product name_" , click "_Continue_".
- Click "_Download_". This should download a JSON file with your client ID and client secret in it.
  Place this file under your home directory (`~`), and rename it to `.todo.txt-googletasks_credentials.json`.
  This is where the addon will look for your credentials.

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

### Testing

```console
make test
```

Note that the above will not run integration tests.
To do so, pass your credentials via the `CLIENT_ID` and `CLIENT_SECRET` environment variables.

```console
make CLIENT_ID=... CLIENT_SECRET=... test
```

N.B.: using [`direnv`](https://direnv.net/) / an `.envrc` file to automatically provide your credentials may be convenient.
