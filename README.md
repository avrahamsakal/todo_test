# TODO Test (deleteme)

## Setup

To install Go on a Mac, run the following in the terminal. (Project uses Go version >= 1.19)
~~~
brew install go
~~~
then run
~~~
make godownload gotidy build
~~~

Ostensibly it would be necessary for a dev to copy .env.example to .env so server can load env var overrides, but it's currently not

## Running

### Run locally
~~~
make runlocal
~~~

### Run tests
~~~
make test
~~~

### Docker image

#### Build
~~~
make build-image
~~~

#### Run
~~~
make rundocker
~~~
