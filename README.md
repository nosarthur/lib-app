# web application of TODO list

This is an attempt to combine Go backend, Polymer frontend, and Slack robot into a todo list app.

The basic data type is called `Ticket` and each `Ticket` contains multiple `Todo`s.

## test

* go test ./storage -v
* go test ./server -v

## local deployment

* go install ./cmd/...
* heroku local web

## heroku deployment

To set up:

* heroku login
* heroku create -b https://github.com/heroku/heroku-buildpack-go.git

To deploy:

* godep save ./cmd/...
* git push heroku master
* heroku open

To create database

* heroku addons:create heroku-postgresql:hobby-dev
* heroku pg:info
* heroku pg:psql
* heroku run bash
    * initTodobotDB

To maintain:

* heroku apps:info
* heroku ps
* heroku logs --tail
* heroku config[:set|unset]

## install polymer and initialize the project

* install [node.js](https://nodejs.org/en/)
* sudo npm install npm@latest -g
* sudo npm install -g polymer-cli
* sudo npm install -g bower
* polymer init (optional)
* bower init
* bower install iron-ajax
* polymer build (optional)
