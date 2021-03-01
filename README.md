# cookiecutter-coldbrew

Powered by [Cookiecutter](https://github.com/audreyr/cookiecutter), Cookiecutter Coldbrew is a framework for jumpstarting production-ready go projects quickly.

## Features

- Generous `Makefile` with management commands
- injects build time and git hash at build time.

## Constraints

- Uses `dep` or `mod` for dependency management
- Only maintained 3rd party libraries are used.
- Use multistage docker builds for super small docker images
- Make sure '$GOBIN' is set in PATH

## Docker

This template uses docker multistage builds to make images slimmer and containers only the final project binary and assets with no source code whatsoever.

You can find the image dokcer file in this [repo](https://github.com/lacion/alpine-golang-buildimage) and more information about docker multistage builds in this [blog post](https://www.critiqus.com/post/multi-stage-docker-builds/).

## Usage

Let's pretend you want to create a project called "echoserver". Rather than starting from scratch maybe copying 
some files and then editing the results to include your name, email, and various configuration issues that always 
get forgotten until the worst possible moment, get cookiecutter to do all the work.

First, get Cookiecutter. Trust me, it's awesome:
```console
$ pip install cookiecutter
```

Alternatively, you can install `cookiecutter` with homebrew:
```console
$ brew install cookiecutter
```

Finally, to run it based on this template, type:
```console
$ cookiecutter gh:go-coldbrew/cookiecutter-coldbrew
```

You will be asked about your basic info (name, project name, app name, etc.). This info will be used to customize your new project.

Warning: After this point, change 'github.com/ankurs', 'MyApp', etc to your own information.

Answer the prompts with your own desired [options](). For example:
```console
source_path [github.com/ankurs]: github.com/ankurs
app_name [MyApp]: MyApp
grpc_package [github.com.ankurs]: github.com.ankurs
service_name [MySvc]: MySvc
project_short_description [A Golang project.]: A Golang project
docker_image [alpine:latest]:
docker_build_image [golang]:
Select docker_build_image_version:
1 - 1.15
2 - 1.16
Choose from 1, 2 [1]: 2
```

Enter the project and take a look around:
```console
$ cd MyApp/
$ ls
```

Run `make help` to see the available management commands, or just run `make build` to build your project.
```console
$ make run
```


