#!/bin/bash

# installing soda
if ! soda -v 2>/dev/null; then
    go get -v github.com/gobuffalo/pop/...
    go install github.com/gobuffalo/pop/soda
fi