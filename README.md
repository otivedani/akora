# akora

mail-puller: 
flow : user fills google form with emails -> (this repo) read gmails from google spreadsheet -> send files to the gmails.

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