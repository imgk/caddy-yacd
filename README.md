# Caddy-PProf

## Build with [xcaddy](https://github.com/caddyserver/xcaddy)

```
$ xcaddy build \
    --with github.com/imgk/caddy-yacd
```

## Config

```
{
    "apps": {
        "http": {
            "servers": {
                "": {
                    "routes": [
                        {
                            "handle": [
                                {
                                    "handler": "yacd"
                                }
                            ]
                        }
                    ]
                }
            }
        }
    }
}

```
