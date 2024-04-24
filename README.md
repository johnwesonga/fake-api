## Overview

fake-api is just what the title says, it's fake and does nothing.
I just wanted to play with Go and Docker.

To build the container to run on Apple silicon:

```
docker buildx build --platform linux/amd64 --load --tag fake-api .
```

To run the docker container:
```
docker run -p 8080:8080 fake-api
```