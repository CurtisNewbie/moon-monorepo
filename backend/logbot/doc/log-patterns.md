# Log Patterns

Logs are parsed using regex, there are two log patterns already available:

```yaml
log:
  pattern:
    go: '^([0-9]{4}\-[0-9]{2}\-[0-9]{2} [0-9:\.]+) +(\w+) +\[([\w ]+),([\w ]+)\] +([\w\.\(\)\.\*_\-]+) +: *((?s).*)'
    java: '^([0-9]{4}\-[0-9]{2}\-[0-9]{2} [0-9:\.]+) +(\w+) +\[[\w \-]+,([\w ]*),([\w ]*),[\w ]*\] [\w\.]+ \-\-\- \[[\w\- ]+\] ([\w\-\.]+) +: *((?s).*)'
```

With `go` pattern, the logs looks like this:

```log
2023-06-13 22:16:13.746 ERROR [v2geq7340pbfxcc9,k1gsschfgarpc7no] main.registerWebEndpoints.func2 : Oh on!
continue on a new line :D
```

With `java` pattern, the logs looks like this:

```log
2023-06-17 17:34:48.762  INFO [auth-service,,,] 78446 --- [           main] .c.m.r.c.YamlBasedRedissonClientProvider : Loading RedissonClient from yaml config file, reading environment property: redisson-config
```

The pattern provided will need to match following groups:

1. time that matches `2006-01-02 15:04:05.000`
2. log level
3. trace id
4. span id
5. name of the method
6. log content