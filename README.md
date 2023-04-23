# cookiecutter-coldbrew

Powered by [Cookiecutter](https://github.com/audreyr/cookiecutter), Cookiecutter Coldbrew is a framework for jumpstarting production-ready go projects quickly.

## Features

- Generous `Makefile` with management commands
- injects build time and git hash at build time.
- Powered by [ColdBrew](https://docs.coldbrew.cloud)

## Constraints

- Uses `mod` for dependency management
- Only maintained 3rd party libraries are used.
- Use multistage docker builds for super small docker images
- Make sure '$GOBIN' is set in PATH

## Docker

This template uses docker multistage builds to make images slimmer and containers only the final project binary and assets with no source code whatsoever.

## Usage

Let's pretend you want to create a project called "echoserver".

Rather than starting from scratch maybe copying some files and then editing the results to include your name, email, and various configuration issues that always get forgotten until the worst possible moment, get cookiecutter to do all the work.

### Prerequisites
First, get Cookiecutter. Trust me, it's awesome:

```shell
$ pip install cookiecutter
```

Alternatively, you can install `cookiecutter` with homebrew:

```shell
$ brew install cookiecutter
```
### Using the ColdBrew Cookiecutter Template

To run it based on this template, type:

```shell
$ cookiecutter gh:go-coldbrew/cookiecutter-coldbrew
```

You will be asked about your basic info \(name, project name, app name, etc.\). This info will be used to customise your new project.

### Providing your app information to the cookiecutter

Warning: After this point, change 'github.com/ankurs', 'MyApp', etc to your own information.

Answer the prompts with your own desired options. For example:

```shell
source_path [github.com/ankurs]: github.com/ankurs
app_name [MyApp]: MyApp
grpc_package [github.com.ankurs]: github.com.ankurs
service_name [MySvc]: MySvc
project_short_description [A Golang project.]: A Golang project
docker_image [alpine:latest]:
docker_build_image [golang]:
Select docker_build_image_version:
1 - 1.19
2 - 1.20
Choose from 1, 2 [1]: 2
```

### Checkout your new project

Enter the project and take a look around:

```shell
$ cd MyApp/
$ ls
```

Run `make help` to see the available management commands, or just run `make build` to build your project.

```shell
$ make run
```
### Working with your new project

Your project is now ready to be worked on. You can find the generated `README.md` file in the project root directory. It contains a lot of useful information about the project. You can also find the generated `Dockerfile` in the project root directory. It contains a lot of useful commands to build, test, and run your project. You can also find the generated `Makefile` in the project root directory. It contains a lot of useful commands to build, test, and run your project. You can run `make help` to see the available management commands.

