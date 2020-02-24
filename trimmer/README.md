# mb-trimmer

## Usage

Start server.

```sh
$ make run
```

Get video.

```sh
$ http -v \
      --download \
      :3011/video \
      url=='https://www.youtube.com/watch?v=p5BzZNH2mkU'
```

Get sound.

```sh
$ http -v \
      --download \
      :3011/sound \
      url=='https://www.youtube.com/watch?v=p5BzZNH2mkU' \
      start-ms==10000 \
      duration-ms==5000
```
