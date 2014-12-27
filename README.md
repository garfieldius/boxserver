# Boxserver

This is a lightweight HTTP server that emulates the behaviour of [Atlas](https://atlas.hashicorp.com/), formerly known as Vagrant Cloud. You can use boxes with version numbers, use vagrants upgrade functionality and add or delete boxes and versions using PUT and DELETE HTTP requests.

**_Important_: This is _NOT_ a replacement for Vagrant Cloud / Atlas, as it only contains a subset of its features. Especially the access permissions or other control options are not included. Everyone who can send a HTTP request to it, can delete all boxes**

## Installation

Have [go](https://golang.org/dl/) installed, clone the repo, build and start it

```bash
mkdir -p $GOPATH/github.com/trenker/boxserver
cd $GOPATH/github.com/trenker/boxserver

git clone https://github.com/trenker/boxserver.git .
go build

./boxserver
```

## Configuration

### Parameters

There is only one command switch: `-c`. Use it to set a custom configuration file. An example configuration is included in the file `config.json`.

### JSON

As the file suffix suggests, the actual configuration is a simple JSON struct offering the following options

#### data

This is the base directory where all boxes will be stored. It needs to be writeable for the boxserver and readable for a webserver, because the webserver needs to deliver those, the boxserver only tells vagrant about their location.

#### baseUrl

The URL prepended to box URLs, including scheme, domain and ending slash.

eg.: If you have a document root of `/var/www/html` and the *data* option set to `/var/www/html/boxes`. The *baseUrl* setting should be `http://domain.tld/boxes/`

#### port

The port to listen at.

#### proxy

You can block the configured port for access from outside and have your webserver proxy the requests to the boxserver. In that case, you need to set the proxy location.

eg.: If you configure nginx like this:

```nginx
location /vagrant {
	proxy_pass http://127.0.0.1:8001;
}
```

You should set *proxy* to `/vagrant`

#### cors

If you want the boxserver to send CORS headers, set this to the desired value for `Access-Control-Allow-Origin`. In this case it will also set a `Access-Control-Allow-Methods` with a hardcoded list of all supported methods.

This should be set when using the [boxserver frontend](https://github.com/trenker/boxserver-frontend)

## Useage

If you have a box, send it to the boxserver with curl. The URL path must have this structure:

`[PROJECT|USER]/[BOX]/[VERSION]/[PROVIDER]`

eg.: to create a new box called `test/demo` for virtualbox, run this:

```bash
curl -X "PUT" -F="box=@my-box-file.box" http://localhost:8001/test/demo/1.0.0/virtualbox
```

Now you need to tell vagrant about this in your vagrant file:

```ruby
Vagrant.configure("2") do |config|
	config.vm.box_server_url = "http://localhost:8001"
	config.vm.box = "test/demo"
end
```

Now run `vagrant up` and it will import version 1.0.0 of your box.

## Reference

Once it is running, all operations are done using HTTP. This can be done using a command line tool like cURL or the [boxserver frontend](https://github.com/trenker/boxserver-frontend)

The response is always JSON

Here is a list of commands and some examples.

### GET / HEAD

HEAD and GET requests work the same, with the (obvious?) difference, that HEAD requests have no response body.

`/`: Get all projects:

```bash
curl 'http://localhost:8001'
[
	"project1",
	"project2"
]
```

`/PROJECT`: Get all boxes of a project / user / group:

```bash
curl 'http://localhost:8001/project1'
[
	"box-a",
	"box-b"
]
```

`/PROJECT/BOX`: Get versions of box

```bash
curl 'http://localhost:8001/project1/box-a'
{
	"name": "project1/box-a",
	"description": "",
	...
	"versions": [
		{
			"version": "1.0.0",
			...
			"providers": [
				{
					"provider": "virtualbox",
					"url": "http://.../virtualbox.box"
				}
			]
		}
	]
}
```

### POST / PUT

POST and PUT requests can be used synonymous.

`/PROJECT/BOX/VERSION/PROVIDER`: Add a box for the given project and provider, giving it the specified version

```bash
curl -F="box=@my-box-file.box" http://localhost:8001/project/box-a/1.0.0/virtualbox
```

### DELETE

Removes boxes. All DELETEs are recursive, meaning, if you delete the last provider, the version will be deleted, if you delete the last version, the box will be deleted and so on.

Delete provider
```bash
curl -X "DELETE" 'http://localhost:8001/project1/box-a/1.0.0/virtualbox'
```

Delete version and all providers
```bash
curl -X "DELETE" 'http://localhost:8001/project1/box-a/1.0.0'
```

Delete box and all versions
```bash
curl -X "DELETE" 'http://localhost:8001/project1/box-a'
```

Delete project and all boxes
```bash
curl -X "DELETE" 'http://localhost:8001/project1'
```

## LICENSE

The MIT License (MIT)

Copyright (c) 2014 Georg Großberger

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
