# Go Bayes

Mainly a POC package to explain how it works. But the package is strong enough to be used.

A command line is implemented in cmd/bayes

# Package Installation

The easy way

```
go get -u github.com/metal3d/go-bayes
```

The package is named "bayes"

# POC tool installation


```
go get -u github.com/metal3d/go-bayes/cmd/bayes
```

And use:

```
bayes -h
  -category string
    	category name
  -check
    	check result
  -check-all
    	check result in each categories
  -content string
    	message to train
  -lang string
    	language for stemming, if not set the content will not be stemmed
  -save string
    	file base (default "bayes.json")
  -save-if float
    	if > 0, save the result in found categorie if bayes result is greater than this value

```

Example:

```
bayes -category animal -content "My cat has a long tail , yellow eyes and meows a lot"
bayes -category animal -content "My dog has a no tail , he's really sweet"
bayes -category object -content "My table is yellow too, but it's made with wood"
bayes -category object -content "My chair is red, made with plastic"
bayes -check animal -content "What is made with wood ?"
```

