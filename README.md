<table align="center"><tr><td align="center" width="9999">
<img src="images/gophemeral.png" align="center" width="300" alt="Gophemeral">

# Gophemeral

Easy Secrets Sharing 

</td></tr></table>



## Overview
Gophemeral is a temporary secret sharing tool. You can input a string and a number of views and Gophemeral will keep the string secret until the number of views runs out.

Secrets must be under 100 characters. 

Site is [https://gophemeral.com](https://gophemeral.com).

## Usage

Gophemeral also has an API. 

## Create Secret

To create a secret, send a POST request with this payload to `https://gophemeral.com/api/secret`:

```
{
	"text": "this is a test",
	"views": 1
}
```

## Lookup Secret

To retrieve a secret, send a GET request to `https://gophemeral.com/api/secret?id={message-id}` and the password in the header `X-Header`.


## Technologies

The backend for Gophemeral right now is BoltDB.

The site is embedded and uses [HTMX](https://htmx.org/).

Everything is deployed on [Fly.io](https://fly.io)
