# zditm-api

> A simple API for getting departure times from ZDiTM

Default port: `38126`

Endpoint: `GET /json/:id`  
Response:
```
[
    {
      line: string
      direction: string
      time: string
    }
]
```

Endpoint: `/text/:id`  
Response:
```
${line}: ${direction} ${time}
```

Best used with a HTTP widget like [this one](https://apkpure.com/http-widget/com.axgs.httpwidget).
