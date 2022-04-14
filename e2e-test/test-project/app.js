const express = require("express")
const { expressCspHeader, INLINE, SELF } = require('express-csp-header');
const app = express()

app.use(expressCspHeader({
  directives: {
    'default-src': [SELF],
    'script-src': [SELF, INLINE, 'somehost.com'],
  }
}))

app.use(express.static("static"))

const port = process.env.PORT ? PORT : 8080
app.listen(port, () => console.log(`Listening on ${port}`))
