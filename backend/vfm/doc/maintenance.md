# Maintenance

Calculate size of all directories recursively, bubbling up to the root:

```sh
curl -X POST "http://localhost:8086/compensate/dir/calculate-size"
```

Compensate thumbnail generations, those that are images/videos (guessed by names) are processed to generate thumbnails:

```sh
curl -X POST "http://localhost:8086/compensate/thumbnail"
```
