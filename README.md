# krakend-cached-router
A simple plugin able to cache request at router level.

![KrakenD router-proxy](https://www.krakend.io/images/documentation/krakend-plugins.png)

Out of the box KrakenD already has a caching mechanism,with some limitations:
 - caching only at http-client level (blue dot on the right)
 - caching is only in-memory (same memory of running container)
Also the kind of caching provided is only in-memory.

The goal of this plugin is to provide a more complex caching strategy:
 - [x] http-handler caching
 - [x] caching on redis store
 - [ ] some tests
 - [ ] caching on memcache store
---
## How to?
### Build
A basic Makefile is provided for build the plugin:\
```make build```    

To change the golang version to use simply pass a makefile argument:\
```make GO_VERSION=1.1 build```

> :warning: **Note**: during compile is mandatory to use the same golang version of your KrakenD instace, you can find it [here](https://plugin-tools.krakend.io/).

### Use
Add plugin load to `krakend.json`
```json
{
  "version": 2,
  "plugin": {
    "folder": "./plugins/build/",
    "pattern": ".so"
  },
  "extra_config": {
    "github_com/devopsfaith/krakend/transport/http/server/handler": {
      "name": [
        "cached-router"
      ],
      "type": "redis",
      "host": "redis:6379"
    }
  },
}
```
