# web application of TODO list

## 

* backend: go
* frontend: polymer 

## local deployment

* go install ./cmd/...
* heroku local web

## heroku deployment

To set up:

* heroku login
* heroku create -b https://github.com/heroku/heroku-buildpack-go.git

To deploy

* godep save ./cmd/...
* git push heroku master

To maintain:

* heroku ps
* heroku open
* heroku logs
* heroku config
