# :warning: ABANDONED :warning:

With the advancements of Docker for Mac I see no reason to keep switching between Docker for Mac and Docker Toolbox is no longer needed.

# docker-set
docker-set is a simple tool to switch between docker environments, virtual machines and docker for mac.

## Installation

```sh
$ go get -u github.com/fmenezes/docker-set
```

## Usage

### 1. List all machines

```sh
$ docker-set list

ACTIVE NAME           DRIVER         STATE
*      docker-for-mac docker-for-mac Unknown
       default        docker-machine Running
       test           vagrant        running
```

### 2. Sets the environment

```sh
eval $(docker-set env default)
```

All further docker commands will run in the selected machine

### 3. Adds a vagrant box to the list

```sh
$ docker-set add test vagrant /path/to/Vagrantfile

Done
```

### 4. Removes a vagrant box from the list

```sh
$ docker-set rm test

Done
```

### 5. Starts a vm from the list

```sh
$ docker-set start test

Done
```


### 6. Stops a vm from the list

```sh
$ docker-set stop test

Done
```

## Notes
- When adding or removing a vagrant machine a file named `$HOME/.docker-set` will be stored.

