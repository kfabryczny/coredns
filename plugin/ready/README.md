# ready

## Name

*ready* - enables a readiness check endpoint.

## Description

By enabling *ready* CoreDNS will wait until all plugins signal readiness. When they all called in,
the startup is continued. This plugin is global in the sense that any plugin, regardless of Server
Block, will report readiness even when *ready* is only specific in one server block - as this
setting affect the overall process readiness.

## Syntax

~~~
ready [ADDRESS]
~~~

*ready* optionally takes an address; the default is `:8181`. The path is fixed to `/ready`. The
readiness endpoint returns a 200 response code and the word "OK" when this server is ready. It
returns a 503 otherwise.

## Plugins

Any plugin wanting to signal readiness will need to add the following code in their setup function:

~~~ go
ready.RegisterPlugin("erratic")
~~~

And then at some point in the future call `ready.Signal` with the same name.

~~~ go
ready.Signal("erratic")
~~~

## Examples

Run another ready endpoint on <http://localhost:8091/ready>.

~~~ txt
. {
    ready localhost:8091
}
~~~
