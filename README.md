# muellkalender-oldenburg_enhancer
Tool zum Anreichern der ICS-Daten des oldenburger Müllkalenders um Erinnerungen zum Müllrausbringen ;)
Tool for enhancing the ICS file from the [online collection calendar](https://services.oldenburg.de/index.php?id=45)
for Oldenburg (GER) with alerts for the calendar events. Additionally a reminder for downloading the new ICS file if the 
current dates will end.

## Requirements

This tool is build with Go 1.14

## Usage

```shell script
make build-osx        # builds an executable for OSX under /build/ics_enhancer
make build-linux      # builds an executable for Linux under /build/ics_enhancer_linux
make build-windows    # builds an executable for Windows under /build/ics_enhancer.exe
make clean            # cleans Go and deletes /build
make all              # builds all artifacts
```

## Weblinks
* [How to Write Go Code](https://golang.org/doc/code.html)

## Author

**Stephan Zerhusen**

* https://twitter.com/stzerhus
* https://github.com/szerhusenBC

## Copyright and license

The code is released under the [MIT license](LICENSE?raw=true).
