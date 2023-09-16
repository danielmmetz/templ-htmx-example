This is a quick demo of templ w/ htmx. It represents a shopping-list esque app.

As a quick demo, it may not respect best practices or be properly secure.


```
❯ fd | grep -v "templ.go$" | entr -r -s 'templ generate && go run .'
```
