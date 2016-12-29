# web application of TODO list

[![Build Status](https://travis-ci.org/nosarthur/todobot.svg?branch=master)](https://travis-ci.org/nosarthur/todobot)
![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)

This is an attempt to combine Go backend, Polymer frontend, and Slack slash command into a todo list app.
A working example can be seen [here](http://tiny.cc/rfn0hy).

The basic data type is called `Ticket`, and each `Ticket` contains multiple `Todo`s.

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
