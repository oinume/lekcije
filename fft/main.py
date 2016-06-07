# coding: utf-8
# An entry point of web application

import os

from bottle import run
from fft.config import config
from fft.web import app
import fft.view.index
#import md2ameblo.view.blogger

if __name__ == "__main__":
    app.log.debug("config = %s" % app.config)
    #app.log.debug("routes = %s" % app.routes)
    if app.config["debug"]:
        routes_debug = ""
        # TODO: Print a file name of function
        for route in app.routes:
            routes_debug += "%6s %s\n" % (route.method, route.rule)
        #self.router.add(route.rule, route.method, route, name=route.name)
        #inspect.getargspec(func)
        app.log.debug("===== routes =====\n" + routes_debug)

    port = os.environ.get("PORT")
    if port:
        config["port"] = port
    run(app, **config)
