# Cram your page into the url

## Usage

```
$ go run . url index.html
data=H4sIAAAAAAAA%2F7LJKMnNseNSsMlITUwB0fowRlJ%2BSqUdl4KCgk2GoZ1Hak5OvkJ4flFOCogLUgeRt9GHGAAIAAD%2F%2F7wb6LlIAAAA
$ go run . &
[1] 112938
$ curl localhost:8080?`go run . url index.html`
<html>
 <head>
 </head>
 <body>
   <h1>Hello World<h1>
 </body>
</html>

```
