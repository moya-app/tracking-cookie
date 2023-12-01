# Tracking cookies

This is a simple module which adds a non-expiring tracking cookie to requests which enables following a user session, even if the backend service does not do this.

## Configuration

```
experimental:
  plugins:
    tracking-cookie:
      moduleName: github.com/moya-app/tracking-cookie
      version: v0.0.1
```

and set up a middleware like:

```
http:
  middlewares:
    add-tracking-cookie:
      plugin:
        tracking-cookie:
          # The domain that this should be issued on. Defaults to the domain of
          # the request, but you could modify it to be eg the parent domain so
          # that it is persisted over all subdomains
          domain: test.com

          # The name of the cookie (defaults to reqId)
          name: moyaReqId

          # Cookie Expiry in seconds (default 100 years). Eg below 1 year:
          expires: 31536000
```

## Extracting cookies

You can later use `accesslog` to extract all cookies and parse out the cookie field via config line like:

```
accessLog:
  fields:
    headers:
      names:
        Cookie: keep
        Set-Cookie: keep
```
