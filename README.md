# log4go
a logkit like as log4j with go language

# Introduce
This project is come from an unmaintained fork, the branch can see http://log4go.googlecode.com

- left only so it doesn't break imports.
- some enhancement if to need 


Usage:
- Add the following import:
import l4g "github.com/gohutool/log4go"
```
LoggerManager.InitWithXML("./examples/example.xml")
var logger = log4go.LoggerManager.GetLogger("com.hello")
logger.Info("hello")
logger.Info("hello")
```

Acknowledgements:
- pomack
  For providing awesome patches to bring log4go up to the latest Go spec