# akora

```
akora send [filename]
```
(what we want it to do:)
lookup db
    if exist -> send() //whatsapp or email
    else 'ingest' -> 
        if exist -> send
        else retry until end

```
akora ingest
```

(what we want it to do:)
ingest()


```
akora pull sent
```

(what we want it to do:)
get sent email, dump to db, delete email