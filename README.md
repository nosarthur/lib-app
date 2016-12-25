# web application of TODO list

This is an attempt to combine Go backend, Polymer frontend, and Slack robot into a todo list app.

The basic data type is called `Ticket` and each `Ticket` contains multiple `Todo`s.

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

To maintain:

* heroku ps
* heroku open
* heroku logs
* heroku config[:set|unset]
