# mb-trimmer

## Usage

Start server.

```sh
$ make run
```

Send request.

```sh
$ http -v \
    -o kirinuki.mp3 \
    :3011/kirinuki \
    url=='https://www.youtube.com/watch?v=p5BzZNH2mkU' \
    start-ms==10000 \
    duration-ms==5000
```
