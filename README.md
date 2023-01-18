# TODO Test (deleteme)

## Setup

To install Go on a Mac, run the following in the terminal. Project uses go version >= 1.19
~~~
brew install go
~~~
then run
~~~
make build
~~~

Make sure to update your environment with the values in .env.example (not really necessary, but ostensibly it would be necessary for a dev to copy .env.example to .env so server can load env var overrides)

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
